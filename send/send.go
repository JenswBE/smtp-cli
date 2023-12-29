package send

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/mail"
	"net/smtp"

	"github.com/rs/zerolog/log"
)

type EmailConfig struct {
	Host             string
	Port             uint
	Username         string
	Password         string
	FromName         string
	FromAddress      string
	ToName           string
	ToAddress        string
	Subject          string
	BodyReader       io.Reader
	Security         string
	AllowInsecureTLS bool
}

const (
	EmailSecurityForceTLS = "FORCE_TLS"
	EmailSecuritySTARTTLS = "STARTTLS"
)

func SendEmail(c EmailConfig) (err error) {
	// Create client
	var client *smtp.Client
	switch c.Security {
	case EmailSecurityForceTLS:
		client, err = createForceTLSClient(c.Host, c.Port, c.AllowInsecureTLS)
	case EmailSecuritySTARTTLS:
		client, err = createSTARTTLSClient(c.Host, c.Port, c.AllowInsecureTLS)
	default:
		return fmt.Errorf("unknown email security option: %s", c.Security)
	}
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}

	// Defer server connection close
	defer func() {
		quitErr := client.Quit()
		if quitErr != nil {
			if err == nil {
				err = fmt.Errorf("failed to close connection with server: %w", quitErr)
			} else {
				log.Error().Err(err).Msg("Failed to close connection with server")
			}
		}
	}()

	// Authenticate to server
	if c.Username != "" || c.Password != "" {
		err = client.Auth(smtp.PlainAuth("", c.Username, c.Password, c.Host))
		if err != nil {
			log.Error().Err(err).Msg("Failed to authenticate to SMTP server")
			return fmt.Errorf("failed to authenticate to SMTP server: %w", err)
		}
	} else {
		log.Debug().Msg("Authentication skipped as both username and password are empty")
	}

	// Set the sender and recipient
	if c.FromAddress == "" {
		// From address not set => Default to username
		c.FromAddress = c.Username
	}
	from := mail.Address{Name: c.FromName, Address: c.FromAddress}
	if err := client.Mail(from.String()); err != nil {
		log.Error().Err(err).Str("from_name", c.FromName).Str("from_address", c.FromAddress).Msg("Failed to send MAIL command to server and set sender")
		return fmt.Errorf("failed to send MAIL command to server and set sender: %w", err)
	}
	to := mail.Address{Name: c.ToName, Address: c.ToAddress}
	if err := client.Rcpt(to.String()); err != nil {
		log.Error().Err(err).Str("to_name", c.ToName).Str("to_address", c.ToAddress).Msg("Failed to send RCPT command to server and set receiver")
		return fmt.Errorf("failed to send RCPT command to server and set receiver: %w", err)
	}

	// Send the email body from stdin
	bodyWriter, err := client.Data()
	if err != nil {
		log.Error().Err(err).Msg("Failed to send DATA command to server")
		return fmt.Errorf("failed to send DATA command to server: %w", err)
	}
	_, err = io.WriteString(bodyWriter, fmt.Sprintf("Subject: %s\r\n\r\n", c.Subject))
	if err != nil {
		log.Error().Err(err).Msg("Failed to write subject to email body")
		return fmt.Errorf("failed to write subject to email body: %w", err)
	}
	_, err = io.Copy(bodyWriter, c.BodyReader)
	if err != nil {
		log.Error().Err(err).Msg("Failed to write message from stdin to email body")
		return fmt.Errorf("failed to write message from stdin to email body: %w", err)
	}
	_, err = io.WriteString(bodyWriter, "\r\n")
	if err != nil {
		log.Error().Err(err).Msg("Failed to write final new line to email body")
		return fmt.Errorf("failed to write final new line to email body: %w", err)
	}
	err = bodyWriter.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close email body")
		return fmt.Errorf("failed to close email body: %w", err)
	}

	// Email sent successfully
	return nil
}

func createForceTLSClient(host string, port uint, allowInsecureTLS bool) (*smtp.Client, error) {
	// Connect to server
	hostPort := fmt.Sprintf("%s:%d", host, port)
	conn, err := tls.Dial("tcp", hostPort, &tls.Config{InsecureSkipVerify: allowInsecureTLS}) // #nosec G402
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SMTP server at %s over TLS: %w", hostPort, err)
	}

	// Start SMTP session
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil, fmt.Errorf("failed to create SMTP client from TLS connection at %s: %w", hostPort, err)
	}
	return client, nil
}

func createSTARTTLSClient(host string, port uint, allowInsecureTLS bool) (*smtp.Client, error) {
	// Connect to server
	hostPort := fmt.Sprintf("%s:%d", host, port)
	client, err := smtp.Dial(hostPort)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SMTP server at %s: %w", hostPort, err)
	}

	// Switch to STARTTLS
	err = client.StartTLS(&tls.Config{InsecureSkipVerify: allowInsecureTLS}) // #nosec G402
	if err != nil {
		return nil, fmt.Errorf("failed to switch to TLS using STARTTLS on server %s: %w", hostPort, err)
	}
	return client, err
}
