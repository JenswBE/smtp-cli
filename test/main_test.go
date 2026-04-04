package test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	smtpPortImplicitTLS = 8465
	smtpPortSTARTTLS    = 8587
)

type E2ETestSuite struct {
	suite.Suite
}

func Test_E2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

func (s *E2ETestSuite) SetupSuite() {
	// Poll SMTP mock server for completed start
	slog.Info("Checking if SMTP mock servers are reachable ...")
	checkSMTPMockRunning(s, smtpMockImplictTLSBaseURL)
	checkSMTPMockRunning(s, smtpMockSTARTTLSBaseURL)
	slog.Info("SMTP mock servers up and running")
}

func checkSMTPMockRunning(s *E2ETestSuite, baseURL string) {
	for i := 0; true; i++ {
		err := pingServer(baseURL)
		if err == nil {
			// Server started
			break
		}
		// Server still down
		if i >= 10 {
			// Exceeded 10 tries => Fail unit tests
			s.Require().FailNow("Unable to contact API of SMTP mock server after 20 seconds")
		}
		// Retry in 2 seconds
		slog.Info("Polling API of SMTP mock server failed, retrying in 2 seconds", "error", err)
		time.Sleep(2 * time.Second)
	}
}

func (s *E2ETestSuite) SetupTest() {
	err := clearMessages(smtpMockImplictTLSBaseURL)
	s.Require().NoError(err)
	err = clearMessages(smtpMockSTARTTLSBaseURL)
	s.Require().NoError(err)
}
