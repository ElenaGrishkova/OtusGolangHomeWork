package logger

import (
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	logger, err := New("wrongLevel")
	require.Error(t, err)
	require.Nil(t, logger)

	logger, err = New("info")
	require.NoError(t, err)
	require.NotNil(t, logger)
	logger.Info("Message to log")
	logger.Debug("Message to log not show")
}
