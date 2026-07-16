package caddyfile_editor

import (
	"testing"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func TestUnmarshalCaddyfile(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantAuth  string
		wantHash  string
	}{
		{
			name:      "valid no_password",
			input:     "admin_panel no_password",
			wantErr:   false,
			wantAuth:  "no_password",
			wantHash:  "",
		},
		{
			name:      "valid bcrypt",
			input:     `admin_panel bcrypt $2a$12$Wv9hQoMf3AIa5qEdwd/95uq0oyJacFTD03/cMKnBAQ0zm54ovS/9K`,
			wantErr:   false,
			wantAuth:  "bcrypt",
			wantHash:  "$2a$12$Wv9hQoMf3AIa5qEdwd/95uq0oyJacFTD03/cMKnBAQ0zm54ovS/9K",
		},
		{
			name:      "invalid auth method",
			input:     "admin_panel invalid_auth",
			wantErr:   true,
			wantAuth:  "invalid_auth",
			wantHash:  "",
		},
		{
			name:      "missing auth method",
			input:     "admin_panel",
			wantErr:   true,
			wantAuth:  "",
			wantHash:  "",
		},
		{
			name:      "missing bcrypt hash",
			input:     "admin_panel bcrypt",
			wantErr:   true,
			wantAuth:  "bcrypt",
			wantHash:  "",
		},
		{
			name:      "invalid bcrypt hash",
			input:     "admin_panel bcrypt not_a_hash",
			wantErr:   true,
			wantAuth:  "bcrypt",
			wantHash:  "not_a_hash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Middleware{}
			d := caddyfile.NewTestDispenser(tt.input)
			err := m.UnmarshalCaddyfile(d)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalCaddyfile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if m.AuthMethod != tt.wantAuth {
					t.Errorf("UnmarshalCaddyfile() AuthMethod = %v, want %v", m.AuthMethod, tt.wantAuth)
				}
				if m.AdminPasswordHash != tt.wantHash {
					t.Errorf("UnmarshalCaddyfile() AdminPasswordHash = %v, want %v", m.AdminPasswordHash, tt.wantHash)
				}
			}
		})
	}
}
