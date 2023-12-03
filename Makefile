GOPROXY?="https://proxy.golang.org,direct"
GIT_COMMIT?=$(shell git rev-list -1 HEAD | xargs git rev-parse --short)
GIT_DATE?=$(shell git log -1 --date=short --pretty=format:%ct)
VERSION_MAJOR=1
VERSION_MINOR=2
VERSION_PATCH=0
VERSION_META=$$ENV_VERSION_META
VERSION="$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_PATCH)-$(VERSION_META)"
VERSION_FLAGS=-X github.com/artheranet/arthera-node/version.VersionMajor=$(VERSION_MAJOR) -X github.com/artheranet/arthera-node/version.VersionMinor=$(VERSION_MINOR) -X github.com/artheranet/arthera-node/version.VersionPatch=$(VERSION_PATCH) -X github.com/artheranet/arthera-node/version.VersionMeta=$(VERSION_META) -X github.com/artheranet/arthera-node/version.GitCommit=$(GIT_COMMIT) -X github.com/artheranet/arthera-node/version.GitDate=$(GIT_DATE)

.PHONY: all
all: arthera

all_release: rpc_docker rpc_trace_docker node_docker
	docker login && \
	docker image push arthera/arthera-rpc:$(VERSION) && \
	docker image push arthera/arthera-rpc-trace:$(VERSION) && \
	docker image push arthera/arthera-node:$(VERSION) && \
	docker logout

.PHONY: arthera
arthera:
	@echo "Building version: $(VERSION)"
	go build \
	    -ldflags "-s -w $(VERSION_FLAGS)" \
	    -o build/arthera-node \
	    ./cmd/arthera

.PHONY: clean
clean:
	rm -fr ./build/*

## <RPC node>
rpc_release: rpc_docker
	docker login && \
	docker image push arthera/arthera-rpc:$(VERSION) && \
	docker logout

rpc_docker:
	docker build -f Dockerfile.rpc --network host --build-arg "GIT_COMMIT=$(GIT_COMMIT)" --build-arg "GIT_DATE=$(GIT_DATE)" . -t arthera/arthera-rpc:$(VERSION)
## </RPC node>

## <RPC trace node>
rpc_trace_release: rpc_trace_docker
	docker login && \
	docker image push arthera/arthera-rpc-trace:$(VERSION) && \
	docker logout

rpc_trace_docker:
	docker build -f Dockerfile.rpc.tracenode --network host --build-arg "GIT_COMMIT=$(GIT_COMMIT)" --build-arg "GIT_DATE=$(GIT_DATE)" . -t arthera/arthera-rpc-trace:$(VERSION)
## </RPC trace node>

## <Validator node>
node_release: node_docker
	docker login && \
	docker image push arthera/arthera-node:$(VERSION) && \
	docker logout

node_docker:
	docker build -f Dockerfile.node --network host --build-arg "GIT_COMMIT=$(GIT_COMMIT)" --build-arg "GIT_DATE=$(GIT_DATE)" . -t arthera/arthera-node:$(VERSION)
## </Validator node>

check_changes:
	@if ! git diff-index --quiet HEAD --; then \
		echo "You have uncommitted changes. Please commit or stash them before making a release."; \
		exit 1; \
	fi

push_changes:
	git push --tags

tag_release:
	@echo "Checking tag $(VERSION)"
	@if ! git tag -l | grep "$(VERSION)"; then \
		echo "Creating git tag $(VERSION)"; \
		git tag -a $(VERSION) -m "Release $(VERSION)"; \
	fi

