package main

import (
	"errors"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid https", "https://example.com", false},
		{"valid with spaces around", "  http://example.com/a  ", false},
		{"empty", "", true},
		{"whitespace only", "   ", true},
		{"no scheme", "example.com", true},
		{"unsupported scheme", "ftp://example.com", true},
		{"missing host", "http:///path", true},
		{"port too big", "http://example.com:70000", true},
		{"port is not a number", "http://example.com:abc", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := URLAnalyzer{}.Validate(tc.input)
			if tc.wantErr != (err != nil) {
				t.Fatalf("err = %v, wantErr = %v", err, tc.wantErr)
			}
			if err != nil && !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("err = %v, want ErrInvalidInput", err)
			}
		})
	}
}

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Result
	}{
		{
			name:  "full url",
			input: "https://Example.com:8443/a/b/?x=1&x=2&y=3#top",
			want: Result{
				Scheme: "https", Host: "example.com", Port: 8443, IsSecure: true,
				Path: "/a/b/", PathDepth: 2, QueryParams: 3, HasFragment: true,
			},
		},
		{
			name:  "http defaults",
			input: "http://example.com",
			want: Result{
				Scheme: "http", Host: "example.com", Port: 80,
				Path: "/", PathDepth: 0,
			},
		},
		{
			name:  "https default port",
			input: "https://example.com/docs",
			want: Result{
				Scheme: "https", Host: "example.com", Port: 443, IsSecure: true,
				Path: "/docs", PathDepth: 1,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := URLAnalyzer{}.Analyze(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got  %+v\nwant %+v", got, tc.want)
			}
		})
	}
}

func TestAnalyzeInvalid(t *testing.T) {
	_, err := URLAnalyzer{}.Analyze("")
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("err = %v, want ErrInvalidInput", err)
	}
}
