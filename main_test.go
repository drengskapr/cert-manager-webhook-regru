//go:build integration

package main

import (
	"fmt"
	"os"
	"testing"

	dns "github.com/cert-manager/cert-manager/test/acme"
)

func TestConformance(t *testing.T) {
	zone := os.Getenv("TEST_ZONE_NAME")
	if zone == "" {
		t.Skip("TEST_ZONE_NAME not set")
	}

	if login := os.Getenv("TEST_REGRU_LOGIN"); login != "" {
		if err := writeSecretManifest(login, os.Getenv("TEST_REGRU_PASSWORD")); err != nil {
			t.Fatalf("failed to write credentials manifest: %v", err)
		}
	}

	fixture := dns.NewFixture(
		&regruDNSProviderSolver{},
		dns.SetResolvedZone(zone),
		dns.SetManifestPath("testdata/regru"),
		dns.SetStrict(false),
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
	return os.WriteFile(
		"testdata/regru/manifests/secret.yaml",
		[]byte(fmt.Sprintf(tmpl, login, password)),
		0600,
	)
}
