package mailrelay

import (
	"fmt"

	guerrilla "github.com/flashmob/go-guerrilla"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/log"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/kushsharma/servo/internal"
)

// StartGuerrillaSMTP starts the smtp server
// use `openssl req -new -newkey rsa:4096 -x509 -sha256 -days 3650 -nodes -out smtpCert.crt -keyout smtpKey.key`
// to generate private key and signed certificate for TLS
func StartGuerrillaSMTP(config internal.SMTPConfig) (err error) {

	listen := fmt.Sprintf("%s:%d", config.LocalListenIP, config.LocalListenPort)

	cfg := &guerrilla.AppConfig{LogFile: log.OutputStdout.String(), AllowedHosts: []string{"*"}}
	sc := guerrilla.ServerConfig{
		ListenInterface: listen,
		IsEnabled:       true,
		TLS: guerrilla.ServerTLSConfig{
			StartTLSOn:     true,
			PrivateKeyFile: config.TLSPrivateKey,
			PublicKeyFile:  config.TLSCert,
		},
		MaxClients: 1,
		Timeout:    360,
	}
	cfg.Servers = append(cfg.Servers, sc)

	bcfg := backends.BackendConfig{
		"save_workers_size":  1, //3
		"save_process":       "HeadersParser|Header|Hasher|Debugger|MailRelay",
		"log_received_mails": true,
	}
	cfg.BackendConfig = bcfg
	d := guerrilla.Daemon{Config: cfg}
	d.AddProcessor("MailRelay", mailRelayProcessor)

	return d.Start()
}

// mailRelayProcessor decorator relays emails to another SMTP server.
var mailRelayProcessor = func() backends.Decorator {
	return func(p backends.Processor) backends.Processor {
		return backends.ProcessWith(
			func(e *mail.Envelope, task backends.SelectTask) (backends.Result, error) {
				if task == backends.TaskSaveMail {

					fmt.Print(e)
					// err := sendMail(e)
					// if err != nil {
					// 	fmt.Printf("!!! %v\n", err)
					// 	return backends.NewResult(fmt.Sprintf("554 Error: %s", err)), err
					// }

					return p.Process(e, task)
				}
				return p.Process(e, task)
			},
		)
	}
}
