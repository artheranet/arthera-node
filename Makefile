GOPROXY?="https://proxy.golang.org,direct"
GIT_COMMIT=$(shell git rev-list -1 HEAD | xargs git rev-parse --short)
GIT_DATE=$(shell git log -1 --date=short --pretty=format:%ct)
VERSION=1.0.0-rc.2-$(GIT_COMMIT)-$(GIT_DATE)
DOCKER_IMAGE=arthera/arthera-node:$(VERSION)

.PHONY: all
all: arthera

.PHONY: arthera
arthera:
	@echo "Building version: $(VERSION)"
	go build \
	    -ldflags "-s -w -X github.com/artheranet/arthera-node/cmd/arthera/launcher.gitCommit=$(GIT_COMMIT) -X github.com/artheranet/arthera-node/cmd/arthera/launcher.gitDate=$(GIT_DATE)" \
	    -o build/arthera-node \
	    ./cmd/arthera

.PHONY: clean
clean:
	rm -fr ./build/*

docker: docker_build docker_tag

docker_build:
	docker build . -t $(DOCKER_IMAGE)

docker_tag:
	docker tag $(DOCKER_IMAGE) arthera/arthera-node:latest

release:
	docker login && \
	git commit -am "Release $(VERSION)"
#	git tag -a $(VERSION) -m "Release $(VERSION)" && \
#	git push --tags && \
#	docker image push $(DOCKER_IMAGE) && \
#	docker image push arthera/arthera-node:latest && \
#	docker logout

