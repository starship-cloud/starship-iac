package logging_test

import (
	"testing"

	"github.com/starship-cloud/starship-iac/server/logging"
	"github.com/stretchr/testify/assert"
)

func TestStructuredLoggerSavesHistory(t *testing.T) {
	logger := logging.NewNoopLogger(t)

	historyLogger := logger.WithHistory()

	expectedStr := "[DBUG] Hello World\n[INFO] foo bar\n"

	historyLogger.Debug("Hello World")
	historyLogger.Info("foo bar")

	assert.Equal(t, expectedStr, historyLogger.GetHistory())
}
