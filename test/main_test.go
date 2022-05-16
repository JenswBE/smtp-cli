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
	time.Sleep(2 * time.Second)
}

func (s *E2ETestSuite) SetupTest() {
	err := clearMessages()
	require.NoError(s.T(), err)
}
