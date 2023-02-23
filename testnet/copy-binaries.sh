#!/bin/bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)

scp -P22022 -i ~/digitalocean $SCRIPT_DIR/../build/arthera-node root@testnet1:/opt/arthera/bin/
scp -P22022 -i ~/digitalocean $SCRIPT_DIR/../build/arthera-node root@testnet2:/opt/arthera/bin/
scp -P22022 -i ~/digitalocean $SCRIPT_DIR/../build/arthera-node root@testnet3:/opt/arthera/bin/
