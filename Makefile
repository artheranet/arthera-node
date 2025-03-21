GOPROXY?="https://proxy.golang.org,direct"
GIT_COMMIT?=$(shell git rev-list -1 HEAD | xargs git rev-parse --short)
GIT_DATE?=$(shell git log -1 --date=short --pretty=format:%ct)
VERSION_MAJOR=1
VERSION_MINOR=2
VERSION_PATCH=2
VERSION_META=$$ENV_VERSION_META
VERSION="$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_PATCH)-$(VERSION_META)"
VERSION_FLAGS=-X github.com/artheranet/arthera-node/version.VersionMajor=$(VERSION_MAJOR) -X github.com/artheranet/arthera-node/version.VersionMinor=$(VERSION_MINOR) -X github.com/artheranet/arthera-node/version.VersionPatch=$(VERSION_PATCH) -X github.com/artheranet/arthera-node/version.VersionMeta=$(VERSION_META) -X github.com/artheranet/arthera-node/version.GitCommit=$(GIT_COMMIT) -X github.com/artheranet/arthera-node/version.GitDate=$(GIT_DATE)

.PHONY: all
all: arthera

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

## <Validator node>
node_release: check_changes node_docker
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
