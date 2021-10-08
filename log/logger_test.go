package log

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	t.Run("Given a msg should print correctly the info msg", func(t *testing.T) {
		// Given
		app := "test"
		msg := "Test info msg"
		file := fmt.Sprintf("./full_%s.log", app)
		os.Remove(file)
		l := NewLogger("./", "test", false)

		// When
		l.Info(msg)

		// Then
		data, err := os.ReadFile(file)
		assert.NoError(t, err)
		line := string(data)
		assert.Contains(t, line, "[INFO]")
		assert.Contains(t, line, msg)
		os.Remove(fmt.Sprintf("./full_%s.log", app))
	})

	t.Run("Given a msg should print correctly the error msg", func(t *testing.T) {
		// Given
		app := "test"
		msg := "Test error msg"
		fullFile := fmt.Sprintf("./full_%s.log", app)
		errorFile := fmt.Sprintf("./error_%s.log", app)
		os.Remove(fullFile)
		os.Remove(errorFile)
		l := NewLogger("./", "test", false)

		// When
		l.Error(msg)

		// Then
		dataFull, err := os.ReadFile(fullFile)
		assert.NoError(t, err)
		lineFull := string(dataFull)
		assert.Contains(t, lineFull, "[ERROR]")
		assert.Contains(t, lineFull, msg)
		os.Remove(fullFile)

		dataError, err := os.ReadFile(errorFile)
		assert.NoError(t, err)
		lineError := string(dataError)
		assert.Contains(t, lineError, "[ERROR]")
		assert.Contains(t, lineError, msg)
		os.Remove(errorFile)
	})

	t.Run("Given a msg should print correctly the fatal msg", func(t *testing.T) {
		// Given
		app := "test"
		msg := "Test fatal msg"
		fullFile := fmt.Sprintf("./full_%s.log", app)
		errorFile := fmt.Sprintf("./error_%s.log", app)
		os.Remove(fullFile)
		os.Remove(errorFile)
		l := NewLogger("./", "test", false)

		// When
		l.Fatal(msg)

		// Then
		dataFull, err := os.ReadFile(fullFile)
		assert.NoError(t, err)
		lineFull := string(dataFull)
		assert.Contains(t, lineFull, "[FATAL]")
		assert.Contains(t, lineFull, msg)
		os.Remove(fullFile)

		dataError, err := os.ReadFile(errorFile)
		assert.NoError(t, err)
		lineError := string(dataError)
		assert.Contains(t, lineError, "[FATAL]")
		assert.Contains(t, lineError, msg)
		os.Remove(errorFile)
	})

	t.Run("Given a msg should print correctly the warning msg", func(t *testing.T) {
		// Given
		app := "test"
		msg := "Test warning msg"
		file := fmt.Sprintf("./full_%s.log", app)
		os.Remove(file)
		l := NewLogger("./", "test", false)

		// When
		l.Warning(msg)

		// Then
		data, err := os.ReadFile(file)
		assert.NoError(t, err)
		line := string(data)
		assert.Contains(t, line, "[WARN]")
		assert.Contains(t, line, msg)
		os.Remove(fmt.Sprintf("./full_%s.log", app))
	})

	t.Run("Given a debug mode and a msg should print correctly the debug msg", func(t *testing.T) {
		// Given
		app := "test"
		msg := "Test debug msg"
		file := fmt.Sprintf("./full_%s.log", app)
		os.Remove(file)
		l := NewLogger("./", "test", true)

		// When
		l.Debug(msg)

		// Then
		data, err := os.ReadFile(file)
		assert.NoError(t, err)
		line := string(data)
		assert.Contains(t, line, "[DEBUG]")
		assert.Contains(t, line, msg)
		os.Remove(fmt.Sprintf("./full_%s.log", app))
	})

	t.Run("Given a debug mode disable should not print debug msg", func(t *testing.T) {
		// Given
		app := "test"
		msg := "Test debug msg"
		file := fmt.Sprintf("./full_%s.log", app)
		os.Remove(file)
		l := NewLogger("./", "test", false)

		// When
		l.Debug(msg)

		// Then
		_, err := os.ReadFile(file)
		assert.Error(t, err)
	})
}
