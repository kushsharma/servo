package mailrelay

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	gosmtp "github.com/emersion/go-smtp"
	"github.com/kushsharma/servo/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//StartSMTP starts smtp server on provided ip and port
func StartSMTP(config internal.SMTPConfig, be gosmtp.Backend) (err error) {

	hostname, err := os.Hostname()
	if err != nil {
		hostname = viper.GetString("appname")
	}

	s := gosmtp.NewServer(be)
	s.Addr = fmt.Sprintf("%s:%d", config.LocalListenIP, config.LocalListenPort)
	s.Domain = hostname
	s.ReadTimeout = 30 * time.Second
	s.WriteTimeout = 30 * time.Second
	s.MaxMessageBytes = 1024 * 1024 * 5
	s.MaxRecipients = 10000
	s.AllowInsecureAuth = false
	cert, err := tls.LoadX509KeyPair(config.TLSCert, config.TLSPrivateKey)
	if err != nil {
		return fmt.Errorf("error while loading the certificate: %s", err)
	}
	s.TLSConfig = &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
	}

	log.Infof("starting smtp server at: %s", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Error(err)
	}
	return err
}
