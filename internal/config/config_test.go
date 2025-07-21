package config

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

func TestBootStrapLooger(t *testing.T) {
	tests := []struct {
		name      string
		envLevel  string
		wantDebug bool
		wantInfo  bool
		wantWarn  bool
		wantError bool
	}{
		{"DEBUG level", "DEBUG", true, true, true, true},
		{"INFO level", "INFO", false, true, true, true},
		{"WARN level", "WARN", false, false, true, true},
		{"ERROR level", "ERROR", false, false, false, true},
		{"Default level", "", false, true, true, true},
		{"Invalid level", "INVALID", false, true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envLevel != "" {
				os.Setenv("LOG_LEVEL", tt.envLevel)
				defer os.Unsetenv("LOG_LEVEL")
			}

			logger := BootStrapLooger()

			// Test each level
			if got := logger.Enabled(context.TODO(), slog.LevelDebug); got != tt.wantDebug {
				t.Errorf("Debug level enabled = %v, want %v", got, tt.wantDebug)
			}
			if got := logger.Enabled(context.TODO(), slog.LevelInfo); got != tt.wantInfo {
				t.Errorf("Info level enabled = %v, want %v", got, tt.wantInfo)
			}
			if got := logger.Enabled(context.TODO(), slog.LevelWarn); got != tt.wantWarn {
				t.Errorf("Warn level enabled = %v, want %v", got, tt.wantWarn)
			}
			if got := logger.Enabled(context.TODO(), slog.LevelError); got != tt.wantError {
				t.Errorf("Error level enabled = %v, want %v", got, tt.wantError)
			}
		})
	}
}

func TestNew(t *testing.T) {
	// Setup
	logger := BootStrapLooger()

	tests := []struct {
		name        string
		appPort     string
		databaseURL string
		wantPort    int
		wantDSN     string
		wantErr     bool
	}{
		{"Valid config", "8080", "postgres://user:pass@localhost:5432/db", 8080, "postgres://user:pass@localhost:5432/db", false},
		{"Default port", "", "postgres://user:pass@localhost:5432/db", 8888, "postgres://user:pass@localhost:5432/db", false},
		{"Missing database URL", "8080", "", 0, "", true},
		{"Invalid port", "invalid", "postgres://user:pass@localhost:5432/db", 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.appPort != "" {
				os.Setenv("APP_PORT", tt.appPort)
			} else {
				os.Unsetenv("APP_PORT")
			}

			if tt.databaseURL != "" {
				os.Setenv("DATABASE_URL", tt.databaseURL)
			} else {
				os.Unsetenv("DATABASE_URL")
			}

			got, err := New(logger)

			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Server.Port != tt.wantPort {
					t.Errorf("New() got port = %v, want %v", got.Server.Port, tt.wantPort)
				}

				if got.Database.DSN != tt.wantDSN {
					t.Errorf("New() got DSN = %v, want %v", got.Database.DSN, tt.wantDSN)
				}
			}

			// Cleanup
			os.Unsetenv("APP_PORT")
			os.Unsetenv("DATABASE_URL")
		})
	}
}
