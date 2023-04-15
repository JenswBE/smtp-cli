package test

import (
	"os/exec"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
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
	// Start SMTP mock server
	log.Info().Msg("Starting SMTP mock server ...")
	cmd := exec.Command("docker", "compose", "up", "-d")
	result, err := cmd.CombinedOutput()
	require.NoErrorf(s.T(), err, "Failed to start Docker Compose: %s", string(result))

	// Poll SMTP mock server for completed start
	log.Info().Msg("Polling SMTP mock server ...")
	for i := 0; true; i++ {
		_, err := getMessages()
		if err == nil {
			// Server started
			break
		}
		// Server still down
		if i >= 10 {
			// Exceeded 10 tries => Fail unit tests
			require.FailNow(s.T(), "Unable to contact SMTP mock server after 20 seconds")
		}
		// Retry in 2 seconds
		log.Info().Msg("Polling SMTP mock server failed, retrying in 2 seconds")
		time.Sleep(2 * time.Second)
	}
	log.Info().Msg("SMTP mock server up and running")
}

func (s *E2ETestSuite) SetupTest() {
	err := clearMessages()
	require.NoError(s.T(), err)
}
