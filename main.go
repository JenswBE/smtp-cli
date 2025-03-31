package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/jenswbe/smtp-cli/email"
)

var Version = "unknown"

func main() {
	// Parse flags
	var (
		printVersion     = flag.Bool("version", false, "Print version and exit")
		host             = flag.String("host", "localhost", "Hostname of the server")
		port             = flag.Uint("port", 465, "Port of the server")
		username         = flag.String("username", "", "Username for authentication")
		password         = flag.String("password", "", "Password for authentication")
		fromName         = flag.String("from-name", "", "Name of the sender")
		fromAddress      = flag.String("from-address", "", "Address of the sender. Defaults to username.")
		toName           = flag.String("to-name", "", "Name of the receiver")
		toAddress        = flag.String("to-address", "", "Address of the receiver")
		subject          = flag.String("subject", "", "Subject of the email")
		security         = flag.String("security", "FORCE_TLS", "Supported options: FORCE_TLS (= implicit TLS), STARTTLS")
		allowInsecureTLS = flag.Bool("allow-insecure-tls", false, "Skip TLS certificate verification. Should only be used for testing!")
	)
	flag.Parse()

	if *printVersion {
		fmt.Println(Version) //nolint
		return
	}

	// Send email
	err := email.Send(email.Config{
		Host:             *host,
		Port:             *port,
		Username:         *username,
		Password:         *password,
		FromName:         *fromName,
		FromAddress:      *fromAddress,
		ToName:           *toName,
		ToAddress:        *toAddress,
		Subject:          *subject,
		BodyReader:       os.Stdin,
		Security:         *security,
		AllowInsecureTLS: *allowInsecureTLS,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to send email")
	}
	log.Info().Str("subject", *subject).Str("to_name", *toName).Str("to_address", *toAddress).Msg("Email successfully sent")
}
