package main

import (
	"testing"
)

func TestSubdomainFromFQDN(t *testing.T) {
	tests := []struct {
		name    string
		fqdn    string
		zone    string
		want    string
		wantErr bool
	}{
		{
			name: "normal subdomain",
			fqdn: "_acme-challenge.example.com.",
			zone: "example.com.",
			want: "_acme-challenge",
		},
		{
			name: "multi-level subdomain",
			fqdn: "_acme-challenge.sub.example.com.",
			zone: "example.com.",
			want: "_acme-challenge.sub",
		},
		{
			name:    "apex — fqdn equals zone",
			fqdn:    "example.com.",
			zone:    "example.com.",
			wantErr: true,
		},
		{
			name:    "fqdn not in zone",
			fqdn:    "_acme-challenge.other.com.",
			zone:    "example.com.",
			wantErr: true,
		},
		{
			name:    "spoofed suffix — fakeexample.com does not match example.com",
			fqdn:    "_acme-challenge.fakeexample.com.",
			zone:    "example.com.",
			wantErr: true,
		},
		{
			name: "no trailing dots",
			fqdn: "_acme-challenge.example.com",
			zone: "example.com",
			want: "_acme-challenge",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := subdomainFromFQDN(tc.fqdn, tc.zone)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
