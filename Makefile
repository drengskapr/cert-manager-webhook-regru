IMAGE_NAME := "webhook"
IMAGE_TAG := "latest"

ENVTEST_VERSION ?= latest
ENVTEST_K8S_VERSION ?= 1.30.0
ENVTEST_BIN := $(shell go env GOPATH)/bin/setup-envtest
ENVTEST_ASSETS ?= .envtest-assets

.PHONY: build test

build:
	docker build -t "$(IMAGE_NAME):$(IMAGE_TAG)" .

$(ENVTEST_BIN):
	go install sigs.k8s.io/controller-runtime/tools/setup-envtest@$(ENVTEST_VERSION)

test: $(ENVTEST_BIN)
	$(eval ASSETS := $(shell $(ENVTEST_BIN) use $(ENVTEST_K8S_VERSION) -p path --bin-dir $(ENVTEST_ASSETS)))
	TEST_ASSET_ETCD=$(ASSETS)/etcd \
	TEST_ASSET_KUBE_APISERVER=$(ASSETS)/kube-apiserver \
	go test -tags integration -v -timeout 15m ./...
