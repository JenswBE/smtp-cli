package test

import (
	"os/exec"
	"testing"
	"time"

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
	exec.Command("docker", "compose", "up", "-d")
	for i := 0; true; i++ {
		_, err := getMessages()
		if err == nil {
			// Server started
			break
		} else {
			// Server still down
			if i >= 10 {
				// Exceeded 10 tries => Fail unit tests
				require.FailNow(s.T(), "Unable to contact SMTP mock server after 20 seconds")
			}
			// Retry in 2 seconds
			time.Sleep(2 * time.Second)
		}
	}
}

func (s *E2ETestSuite) SetupTest() {
	err := clearMessages()
	require.NoError(s.T(), err)
}
