FROM golang:1.14-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR /arthera

COPY ../arthera-go-ethereum .
COPY ../arthera-node .

RUN cd arthera-node
ARG GOPROXY
RUN go mod download
RUN make

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /arthera/arthera-node/build/arthera-node /

EXPOSE 5050 18545 18546 18547 19090

ENTRYPOINT ["/arthera-node"]
