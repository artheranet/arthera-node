GOPROXY?="https://proxy.golang.org,direct"
GIT_COMMIT?=$(shell git rev-list -1 HEAD | xargs git rev-parse --short)
GIT_DATE?=$(shell git log -1 --date=short --pretty=format:%ct)
VERSION=1.1.0-rc.3

.PHONY: all
all: arthera

.PHONY: arthera
arthera:
	@echo "Building version: $(VERSION)"
	go build \
	    -ldflags "-s -w -X github.com/artheranet/arthera-node/version.GitCommit=$(GIT_COMMIT) -X github.com/artheranet/arthera-node/version.GitDate=$(GIT_DATE)" \
	    -o build/arthera-node \
	    ./cmd/arthera

.PHONY: clean
clean:
	rm -fr ./build/*

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

## <RPC node>
rpc_release: check_changes tag_release push_changes rpc_docker
	docker login && \
	docker image push arthera/arthera-rpc:$(VERSION) && \
	docker image push arthera/arthera-rpc:latest && \
	docker logout

rpc_docker: rpc_docker_build rpc_docker_tag

rpc_docker_build:
	docker build -f Dockerfile.rpc --network host --build-arg "GIT_COMMIT=$(GIT_COMMIT)" --build-arg "GIT_DATE=$(GIT_DATE)" . -t arthera/arthera-rpc:$(VERSION)

rpc_docker_tag:
	docker tag arthera/arthera-rpc:$(VERSION) arthera/arthera-rpc:latest
## </RPC node>

## <RPC trace node>
rpc_trace_release: check_changes tag_release push_changes rpc_trace_docker
	docker login && \
	docker image push arthera/arthera-rpc-trace:$(VERSION) && \
	docker image push arthera/arthera-rpc-trace:latest && \
	docker logout

rpc_trace_docker: rpc_trace_docker_build rpc_trace_docker_tag

rpc_trace_docker_build:
	docker build -f Dockerfile.rpc.tracenode --network host --build-arg "GIT_COMMIT=$(GIT_COMMIT)" --build-arg "GIT_DATE=$(GIT_DATE)" . -t arthera/arthera-rpc-trace:$(VERSION)

rpc_trace_docker_tag:
	docker tag arthera/arthera-rpc-trace:$(VERSION) arthera/arthera-rpc-trace:latest
## </RPC trace node>

## <Validator node>
node_release: check_changes tag_release push_changes node_docker
	docker login && \
	docker image push arthera/arthera-node:$(VERSION) && \
	docker image push arthera/arthera-node:latest && \
	docker logout

node_docker: node_docker_build node_docker_tag

node_docker_build:
	docker build -f Dockerfile.node --network host --build-arg "GIT_COMMIT=$(GIT_COMMIT)" --build-arg "GIT_DATE=$(GIT_DATE)" . -t arthera/arthera-node:$(VERSION)

node_docker_tag:
	docker tag arthera/arthera-node:$(VERSION) arthera/arthera-node:latest
## </Validator node>
