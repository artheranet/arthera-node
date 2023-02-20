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

### Launching a network

You will need a genesis file to join a network, which may be found in https://github.com/Fantom-foundation/lachesis_launch

Launching `arthera-node` readonly (non-validator) node for network specified by the genesis file:

```shell
$ arthera-node --genesis file.g
```

### Configuration

As an alternative to passing the numerous flags to the `arthera-node` binary, you can also pass a
configuration file via:

```shell
$ arthera-node --config /path/to/your_config.toml
```

To get an idea how the file should look like you can use the `dumpconfig` subcommand to
export your existing configuration:

```shell
$ arthera-node --your-favourite-flags dumpconfig
```

#### Validator

New validator private key may be created with `arthera-node validator new` command.

To launch a validator, you have to use `--validator.id` and `--validator.pubkey` flags to enable events emitter.

```shell
$ arthera-node --nousb --validator.id YOUR_ID --validator.pubkey 0xYOUR_PUBKEY
```

`arthera-node` will prompt you for a password to decrypt your validator private key. Optionally, you can
specify password with a file using `--validator.password` flag.

#### Participation in discovery

Optionally you can specify your public IP to straighten connectivity of the network.
Ensure your TCP/UDP p2p port (5050 by default) isn't blocked by your firewall.

```shell
$ arthera-node --nat extip:1.2.3.4
```
