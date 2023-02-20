.PHONY: all
all: arthera

GOPROXY ?= "https://proxy.golang.org,direct"
.PHONY: arthera
arthera:
	GIT_COMMIT=`git rev-list -1 HEAD 2>/dev/null || echo ""` && \
	GIT_DATE=`git log -1 --date=short --pretty=format:%ct 2>/dev/null || echo ""` && \
	GOPROXY=$(GOPROXY) \
	go build \
	    -ldflags "-s -w -X github.com/artheranet/arthera-node/cmd/arthera/launcher.gitCommit=$${GIT_COMMIT} -X github.com/artheranet/arthera-node/cmd/arthera/launcher.gitDate=$${GIT_DATE}" \
	    -o build/arthera-node \
	    ./cmd/arthera

.PHONY: clean
clean:
	rm -fr ./build/*
