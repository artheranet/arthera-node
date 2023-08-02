# Arthera Beacon chain
This is the initial chain to bootstrap the Arthera Network, which we named the "beacon" chain. It's an EVM-compatible chain, secured by the Lachesis DAG consensus algorithm.
The purpose of the Beacon chain is to gather initial support around the Arthera Network, prior to the launch of the new Main chain that will have all the features of Arthera. 

## Building the source
Building the `arthera-node` requires both a Go version 1.14 or later and a C compiler.

To build the node binary just run:

```shell
make arthera
```
The build output is ```build/arthera-node``` executable.

## Building the Docker image
First, commit your changes to the repository, then run:

```shell
make docker
```

## Creating a new release
First, commit your changes to the repository, then run:

## Versioning
The versioning the Semver-compatible, with the following format:

MAJOR.MINOR.PATCH-META-GITCOMMIT-GITDATE

Where:
- MAJOR - is the major version, incremented when there are incompatible changes
- MINOR - is the minor version, incremented when there are new features in a backward compatible manner
- PATCH - is the patch version, incremented when there backward compatible are bug fixes
- META - is the metadata, which is either "alpha", "beta" or "rc" (release candidate)
- GITCOMMIT - is the short git commit hash
- GITDATE - is the date of the commit in Unix epoch time

## License
The Arthera Beacon chain is licensed under the [GPLv3](https://www.gnu.org/licenses/gpl-3.0.en.html) license.
