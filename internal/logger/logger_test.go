package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewDefaultLogger(t *testing.T) {
	tests := []struct {
		name      string
		service   string
		wantField string
	}{
		{"empty service", "", ""},
		{"non-empty service", "test-service", "test-service"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewDefaultLogger(tt.service)
			got := logger.With().Logger()

			out := &bytes.Buffer{}
			got = got.Output(out)
			got.Info().Msg("test")
			logOutput := out.String()

			if tt.wantField != "" && !strings.Contains(logOutput, tt.wantField) {
				t.Errorf("expected log to contain service %q, got %q", tt.wantField, logOutput)
			}
		})
	}
}

func TestNewStdErr(t *testing.T) {
	tests := []struct {
		name           string
		level          string
		service        string
		expectedLevel  string
		expectedPanic  bool
		expectedOutput string
	}{
		{"valid level, empty service", "info", "", "info", false, "test"},
		{"valid level, non-empty service", "debug", "test-service", "debug", false, "test-service"},
		{"invalid level", "invalid", "", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectedPanic {
						t.Errorf("unexpected panic: %v", r)
					}
				}
			}()

			if tt.expectedPanic {
				NewStdErr(tt.level, tt.service)
				return
			}

			logger := NewStdErr(tt.level, tt.service)
			out := &bytes.Buffer{}
			logger = logger.Output(out)
			logger.Info().Msg("test")

			logOutput := out.String()
			if tt.service != "" && !strings.Contains(logOutput, tt.service) {
				t.Errorf("expected log to contain service %q, got %q", tt.service, logOutput)
			}
			if !strings.Contains(logOutput, tt.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tt.expectedOutput, logOutput)
			}
		})
	}
}
