package curlbuilder

import (
	"net/http"
	"strings"
	"testing"
)

func TestBuildCurlCommand(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "GET request",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "https://example.com/api", nil)
					req.Header.Set("Authorization", "Bearer token")
					req.Header.Set("Accept", "application/json")
					return req
				}(),
			},
			want:    "curl -X GET https://example.com/api -H 'Authorization: Bearer token' -H 'Accept: application/json' ",
			wantErr: false,
		},
		{
			name: "POST request with body",
			args: args{
				req: func() *http.Request {
					body := strings.NewReader("foo=bar")
					req, _ := http.NewRequest("POST", "https://example.com/post", body)
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					return req
				}(),
			},
			want:    "curl -X POST https://example.com/post -H 'Content-Type: application/x-www-form-urlencoded' -d 'foo=bar' ",
			wantErr: false,
		},
		{
			name: "nil request",
			args: args{
				req: nil,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildCurlCommand(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("BuildCurlCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if tt.name == "GET request" {
				if !strings.Contains(got, "-H 'Authorization: Bearer token'") || !strings.Contains(got, "-H 'Accept: application/json'") {
					t.Errorf("BuildCurlCommand() missing expected headers: %v", got)
				}
				if !strings.Contains(got, "curl -X GET https://example.com/api") {
					t.Errorf("BuildCurlCommand() missing expected command: %v", got)
				}
			} else if got != tt.want {
				t.Errorf("BuildCurlCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
