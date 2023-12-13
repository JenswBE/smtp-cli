package test

import (
	"bytes"

	"github.com/jenswbe/smtp-cli/send"
)

func (s *E2ETestSuite) TestSendEmail() {
	// Send email
	err := send.SendEmail(send.EmailConfig{
		Host:             "localhost",
		Port:             smtpPort,
		Username:         "TestUsername",
		Password:         "TestPassword",
		FromName:         "TestFromName",
		FromAddress:      "TestFromAddress@example.com",
		ToName:           "TestToName",
		ToAddress:        "TestToAddress@example.com",
		Subject:          "TestSubject",
		BodyReader:       bytes.NewBufferString("TestBody"),
		AllowInsecureTLS: true,
	})
	s.Require().NoError(err)

	// Validate email
	messages, err := getMessages()
	s.Require().NoError(err)
	s.Require().Len(messages, 1, "Server should have received a single message")
	s.Require().Equal(`"TestFromName" <TestFromAddress@example.com>`, messages[0].From)
	s.Require().Equal(`"TestToName" <TestToAddress@example.com>`, messages[0].To)
	s.Require().Equal("TestSubject", messages[0].Subject)
	messageBody, err := getMessageBody(messages[0].ID)
	s.Require().NoError(err)
	s.Require().Contains(messageBody, "TestBody")
}
