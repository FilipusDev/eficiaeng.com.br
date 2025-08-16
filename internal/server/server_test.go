package server

import (
	"log/slog"
	"testing"
)

// TestNewServer tests the newServer constructor
func TestNewServer(t *testing.T) {
	tests := []struct {
		name          string
		logger        *slog.Logger
		assetsVersion string
		wantVersion   string
	}{
		{
			name:          "with version",
			logger:        slog.Default(),
			assetsVersion: "v1.0.0",
			wantVersion:   "v1.0.0",
		},
		{
			name:          "empty version",
			logger:        slog.Default(),
			assetsVersion: "",
			wantVersion:   "dev-", // We'll check prefix since timestamp varies
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := newServer(tt.logger, tt.assetsVersion)

			if srv.logger != tt.logger {
				t.Errorf("got logger %v, want %v", srv.logger, tt.logger)
			}

			if tt.assetsVersion == "" {
				if srv.assetsVersion[:4] != "dev-" {
					t.Errorf("got version %s, want prefix 'dev-'", srv.assetsVersion)
				}
			} else {
				if srv.assetsVersion != tt.wantVersion {
					t.Errorf("got version %s, want %s", srv.assetsVersion, tt.wantVersion)
				}
			}
		})
	}
}
