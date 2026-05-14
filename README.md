# cert-manager-webhook-regru

cert-manager DNS01 webhook solver for the [REG.RU](https://www.reg.ru/) DNS provider.

## Prerequisites

- [cert-manager](https://cert-manager.io/docs/installation/) installed in your cluster
- `kubectl` and `helm` available locally
- A REG.RU account with API access enabled from your cluster's outbound IP address
  - Configure API access: https://www.reg.ru/user/account/#/settings/api/

## REG.RU API access

REG.RU restricts API access to allowlisted IP addresses. Before deploying the webhook, add your cluster's outbound IP to the allowlist in your REG.RU account settings.

API documentation: https://www.reg.ru/reseller/api2doc#common

If access is not configured, calls to the REG.RU API will fail with:

```json
{
   "charset" : "utf-8",
   "error_code" : "ACCESS_DENIED_FROM_IP",
   "error_params" : {
      "command_name" : "zone/get_resource_records"
   },
   "error_text" : "Access to API from this IP denied",
   "messagestore" : null,
   "result" : "error"
}
```

## Installation

### 1. Create the credentials secret

Create a Kubernetes secret with your REG.RU login and password in the `cert-manager` namespace:

```bash
kubectl --namespace cert-manager create secret generic regru-api-creds \
  --from-literal=login='<your-username>' \
  --from-literal=password='<your-password>'
```

For alternative ways to create secrets, see the [Kubernetes documentation](https://kubernetes.io/docs/tasks/configmap-secret/).

### 2. Deploy the webhook

```bash
git clone https://github.com/drengskapr/cert-manager-webhook-regru.git

helm --namespace cert-manager upgrade --install regru-webhook \
  ./cert-manager-webhook-regru/deploy/helm/regru-webhook/ \
```

## Configuration

Create a `ClusterIssuer` referencing the webhook:

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: regru-dns
spec:
  acme:
    email: username@example.com
    privateKeySecretRef:
      name: letsencrypt-private-key
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    solvers:
    - dns01:
        webhook:
          config:
            apiLoginRef:
              key: login
              name: regru-api-creds
            apiPasswordRef:
              key: password
              name: regru-api-creds
          groupName: acme.regru.ru
          solverName: regru
```

> **Note:** The example above uses the Let's Encrypt **staging** server. For production certificates, replace the `server` value with:
> ```
> https://acme-v02.api.letsencrypt.org/directory
> ```

Once the issuer is ready, create a `Certificate` resource to request a certificate:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-com
  namespace: default
spec:
  secretName: example-com-tls
  issuerRef:
    name: regru-dns
    kind: ClusterIssuer
  dnsNames:
  - example.com
  - "*.example.com"
```

## Running the test suite

### Unit tests

`subdomain_test.go` contains fast, dependency-free unit tests and runs without any environment setup:

```bash
go test ./...
```

### Integration / conformance tests

The integration test (`main_test.go`) runs the cert-manager DNS01 conformance suite against a real REG.RU zone. It requires:

- A domain you control that is hosted on REG.RU
- REG.RU API credentials with access allowed from your machine's IP

**Credentials** — choose one of:

1. **Environment variables (recommended):** the test writes `testdata/regru/manifests/secret.yaml` automatically:

   ```bash
   export TEST_ZONE_NAME=example.com.
   export TEST_REGRU_LOGIN=your-login
   export TEST_REGRU_PASSWORD=your-password
   make test
   ```

> **Note:** trailing dot in `TEST_ZONE_NAME` is required

2. **Manual secret file:** copy the example and fill in your credentials, then run without the login/password vars:

   ```bash
   cp testdata/regru/manifests/secret.yaml.example \
      testdata/regru/manifests/secret.yaml
   # edit secret.yaml and set login / password
   TEST_ZONE_NAME=example.com. make test
   ```

The conformance fixture uses `1.1.1.1:53` for DNS propagation checks by default. Override with `TEST_DNS_SERVER=<host:port>` if needed.

The solver config lives in [`testdata/regru/config.json`](testdata/regru/config.json).
