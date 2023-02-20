#!/bin/bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)
DEVNET_DIR=$SCRIPT_DIR/../testnet
GENESIS_FILE="$DEVNET_DIR/testnet.genesis"

./build/opera --port 5050 --http --graphql --ws --http.port 18545  --ws.port 18546 \
              --genesis "$GENESIS_FILE" \
              --identity "node1" --nodekey $DEVNET_DIR/node1.testnet.key \
              --datadir $DEVNET_DIR/node1 --verbosity=3 \
              --bootnodes ""

