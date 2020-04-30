package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/kushsharma/servo/backup"
	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/logtool"
	"github.com/kushsharma/servo/routes"
	"github.com/kushsharma/servo/tunnel"
	cron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	termChan     = make(chan os.Signal, 1)
	shutdownWait = 30 * time.Second
	serverPort   = 9090
)

func initService() *cobra.Command {
	srvCmd := &cobra.Command{
		Use: "service",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			log.Infof("starting servo service on %d...", serverPort)
			CronManager.Start()

			appConfig, ok := viper.Get("app").(internal.ApplicationConfig)
			if !ok {
				return errors.New("unable to find application config")
			}

			//setup crons of machine
			// Funcs are invoked in their own goroutine, asynchronously.
			// e.g.: "@every 0h0m5s" "* * * * *" "1 * * * *" "@daily" "@hourly"
			for _, machine := range appConfig.Machines {
				if _, err := cron.ParseStandard(machine.Schedule); err == nil {

					//backup schedule
					CronManager.AddFunc(machine.Schedule, func() {

						//fs backup scheduled job
						fsService := backup.NewFSService(machine.Backup.FS)
						if err := backupFS(fsService); err != nil {
							log.Error(err)
							internal.AppStats.BackupFSError()
						} else {
							log.Infof("fs backup completed successfully for %s\n", machine.Name)
						}

						//db backup scheduled job
						localTnl := tunnel.NewLocalTunnel()
						defer localTnl.Close()
						dbService := backup.NewDBService(localTnl, machine.Backup.DB)
						if err := backupDB(dbService); err != nil {
							log.Error(err)
							internal.AppStats.BackupDBError()
						} else {
							log.Infof("db backup completed successfully for %s\n", machine.Name)
						}
						internal.AppStats.Backedup()

						//log clean up scheduled job
						logToolService := logtool.NewRcloneService(machine.Clean)
						if err := logClean(logToolService, machine.Clean); err != nil {
							log.Error(err)
							internal.AppStats.LogCleanError()
						} else {
							log.Infof("logs cleaned successfully for %s\n", machine.Name)
							internal.AppStats.LogCleaned()
						}
					})
				} else {
					log.Panic(err)
				}
			}

			//log application stats every few hours
			CronManager.AddFunc("@every 6h", func() {
				if stats, err := json.Marshal(internal.AppStats); err == nil {
					log.Infof("logging application stats: %s", string(stats))
				} else {
					log.Error(err)
				}
			})

			// Inspect the cron job entries' next and previous run times.
			log.Infof("scheduling %d job", len(CronManager.Entries()))

			router := mux.NewRouter()
			router.HandleFunc("/ping", routes.PingHandler)
			router.HandleFunc("/stats", routes.StatsHandler)
			http.Handle("/", router)

			srv := &http.Server{
				Handler:      router,
				Addr:         fmt.Sprintf("127.0.0.1:%d", serverPort),
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			// Run our server in a goroutine so that it doesn't block.
			go func() {
				if err := srv.ListenAndServe(); err != nil {
					log.Warn(err)
				}
			}()

			// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
			// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
			signal.Notify(termChan, os.Interrupt)
			signal.Notify(termChan, os.Kill)

			// Block until we receive our signal.
			<-termChan
			log.Info("termination request received")
			log.Info("waiting for few seconds to clean up scheduled job if any")
			if stats, err := json.Marshal(internal.AppStats); err == nil {
				log.Info(string(stats))
			} else {
				log.Error(err)
			}

			//termination request received, stop futher scheduling and wait for running ones to finish
			<-CronManager.Stop().Done()

			// Create a deadline to wait for server
			ctx, cancel := context.WithTimeout(context.Background(), shutdownWait)
			defer cancel()

			// Doesn't block if no connections, but will otherwise wait
			// until the timeout deadline.
			if err := srv.Shutdown(ctx); err != nil {
				log.Warn(err)
			}

			// Optionally, you could run srv.Shutdown in a goroutine and block on
			// <-ctx.Done() if your application should wait for other services
			// to finalize based on context cancellation.
			log.Print("bye")
			return nil
		},
	}

	srvCmd.Flags().IntVarP(&serverPort, "port", "p", serverPort, "port on which server needs to listen on")
	return srvCmd
}
