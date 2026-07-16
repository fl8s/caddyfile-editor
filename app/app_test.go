package app

import (
	"testing"

	_ "github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
    "github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
    "github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
    httpcaddyfile.RegisterHandlerDirective("admin_panel", func(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
		return nil, nil
	})
}

func TestAdaptCaddyfile(t *testing.T) {
	appStruct := &App{}

	tests := []struct {
		name       string
		content    string
		wantErr    bool
		wantWarn   bool
	}{
		{
			name: "valid syntax with admin_panel",
			content: `localhost:8080 {
				admin_panel no_password
			}`,
			wantErr:  false,
			wantWarn: false,
		},
		{
			name: "valid syntax without admin_panel",
			content: `localhost:8080 {
				respond "Hello World"
			}`,
			wantErr:  false,
			wantWarn: true,
		},
		{
			name: "invalid syntax",
			content: `localhost:8080 {
				invalid_directive
			}`,
			wantErr:  true,
			wantWarn: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := appStruct.AdaptCaddyfile(nil, tt.content)

			if (res.AdaptError != "") != tt.wantErr {
				t.Errorf("AdaptCaddyfile() error = %v, wantErr %v", res.AdaptError, tt.wantErr)
			}

			hasWarn := false
			for _, w := range res.Warnings {
				if w.Directive == "HACK_WHOLEFILE" {
					hasWarn = true
					break
				}
			}

			if hasWarn != tt.wantWarn {
				t.Errorf("AdaptCaddyfile() warning = %v, wantWarn %v", hasWarn, tt.wantWarn)
			}
		})
	}
}
