# Arthera Beacon chain 

This is the initial chain to bootstrap the Arthera Network, which we named the "beacon" chain. It's an EVM-compatible chain, secured by Fantom's Lachesis consensus algorithm.
The purpose of the Beacon chain is to gather initial support around the Arthera Network, prior to the launch of the new Main chain that will have all the features of Arthera. 

## Building the source

Building the `arthera-node` requires both a Go (version 1.14 or later) and a C compiler. You can install
them using your favourite package manager. Once the dependencies are installed, run

```shell
make arthera
```
The build output is ```build/arthera-node``` executable.

## Building the Docker image
docker build . -t arthera/arthera-node:1.0
docker tag arthera/arthera-node:1.0 arthera/arthera-node:latest
docker login
docker image push arthera/arthera-node:1.0
docker image push arthera/arthera-node:latest
docker logout
