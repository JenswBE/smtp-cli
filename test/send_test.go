package test

import (
	"bytes"

	"github.com/jenswbe/smtp-cli/email"
)

func (s *E2ETestSuite) TestSendEmailImplicitTLS() {
	config := getEmailConfig(smtpPortImplicitTLS, email.SecurityForceTLS)
	err := email.Send(config)
	s.Require().NoError(err)
	validateEmailMessages(s, smtpMockImplictTLSBaseURL)
}

func (s *E2ETestSuite) TestSendEmailSTARTTLS() {
	config := getEmailConfig(smtpPortSTARTTLS, email.SecuritySTARTTLS)
	err := email.Send(config)
	s.Require().NoError(err)
	validateEmailMessages(s, smtpMockSTARTTLSBaseURL)
}

func getEmailConfig(port uint, security string) email.Config {
	return email.Config{
		Host:             "localhost",
		Port:             port,
		Username:         "TestUsername",
		Password:         "TestPassword",
		FromName:         "TestFromName",
		FromAddress:      "TestFromAddress@example.com",
		ToName:           "TestToName",
		ToAddress:        "TestToAddress@example.com",
		Subject:          "TestSubject",
		BodyReader:       bytes.NewBufferString("TestBody"),
		Security:         security,
		AllowInsecureTLS: true,
	}
}

func validateEmailMessages(s *E2ETestSuite, baseURL string) {
	messages, err := getMessages(baseURL)
	s.Require().NoError(err)
	s.Require().Len(messages, 1, "Server should have received a single message")
	s.Require().Equal(`"TestFromName" <TestFromAddress@example.com>`, messages[0].From)
	s.Require().Equal(`"TestToName" <TestToAddress@example.com>`, messages[0].To[0])
	s.Require().Equal("TestSubject", messages[0].Subject)
	messageBody, err := getMessageBody(baseURL, messages[0].ID)
	s.Require().NoError(err)
	s.Require().Contains(messageBody, "TestBody")
}
