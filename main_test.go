//go:build integration

package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	dns "github.com/cert-manager/cert-manager/test/acme"
)

func TestConformance(t *testing.T) {
	zone := os.Getenv("TEST_ZONE_NAME")
	if zone == "" {
		t.Skip("TEST_ZONE_NAME not set")
	}

	if login := os.Getenv("TEST_REGRU_LOGIN"); login != "" {
		password := os.Getenv("TEST_REGRU_PASSWORD")
		if password == "" {
			t.Fatal("TEST_REGRU_LOGIN is set but TEST_REGRU_PASSWORD is empty")
		}
		if err := writeSecretManifest(login, password); err != nil {
			t.Fatalf("failed to write credentials manifest: %v", err)
		}
	}

	dnsServer := os.Getenv("TEST_DNS_SERVER")
	if dnsServer == "" {
		dnsServer = "127.0.0.53:53"
	}

	fixture := dns.NewFixture(
		&regruDNSProviderSolver{},
		dns.SetResolvedZone(zone),
		dns.SetManifestPath("testdata/regru"),
		dns.SetStrict(false),
		dns.SetDNSServer(dnsServer),
		dns.SetUseAuthoritative(false),
		dns.SetPropagationLimit(4*time.Minute),
	)
	fixture.RunConformance(t)
}

func writeSecretManifest(login, password string) error {
	const tmpl = `apiVersion: v1
kind: Secret
metadata:
  name: regru-api-creds
stringData:
  login: %q
  password: %q
`
	if err := os.MkdirAll("testdata/regru/manifests", 0700); err != nil {
		return err
	}
	return os.WriteFile(
		"testdata/regru/manifests/secret.yaml",
		[]byte(fmt.Sprintf(tmpl, login, password)),
		0600,
	)
}
