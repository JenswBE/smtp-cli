package test

import (
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

const smtpPort = 8465

type E2ETestSuite struct {
	suite.Suite
}

func Test_E2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

func (s *E2ETestSuite) SetupSuite() {
	// Poll SMTP mock server for completed start
	log.Info().Msg("Checking if SMTP mock server is reachable ...")
	for i := 0; true; i++ {
		_, err := getMessages()
		if err == nil {
			// Server started
			break
		}
		// Server still down
		if i >= 10 {
			// Exceeded 10 tries => Fail unit tests
			s.Require().FailNow("Unable to contact SMTP mock server after 20 seconds")
		}
		// Retry in 2 seconds
		log.Info().Err(err).Msg("Polling SMTP mock server failed, retrying in 2 seconds")
		time.Sleep(2 * time.Second)
	}
	log.Info().Msg("SMTP mock server up and running")
}

func (s *E2ETestSuite) SetupTest() {
	err := clearMessages()
	s.Require().NoError(err)
}
