package test

import (
	"bytes"

	"github.com/jenswbe/smtp-cli/send"
	"github.com/stretchr/testify/require"
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
	require.NoError(s.T(), err)

	// Validate email
	messages, err := getMessages()
	require.NoError(s.T(), err)
	require.Len(s.T(), messages, 1, "Server should have received a single message")
	require.Equal(s.T(), `"TestFromName" <TestFromAddress@example.com>`, messages[0].From)
	require.Equal(s.T(), `"TestToName" <TestToAddress@example.com>`, messages[0].To)
	require.Equal(s.T(), "TestSubject", messages[0].Subject)
	messageBody, err := getMessageBody(messages[0].ID)
	require.NoError(s.T(), err)
	require.Contains(s.T(), messageBody, "TestBody")
}
