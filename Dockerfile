FROM golang:1.19.11-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR /arthera

RUN git clone --depth 1 https://github.com/artheranet/lachesis.git
RUN git clone --depth 1 https://github.com/artheranet/arthera-go-ethereum.git

COPY . arthera-node

WORKDIR /arthera/arthera-node

RUN make

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /arthera/arthera-node/build/arthera-node /usr/local/bin/

EXPOSE 6060 18545 18546

ENTRYPOINT ["arthera-node", "--datadir", "/data", "--verbosity", "3", "--cache", "4000"]