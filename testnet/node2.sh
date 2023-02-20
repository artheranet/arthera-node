#!/bin/bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)
DEVNET_DIR=$SCRIPT_DIR/../devnet

GENESIS_FILE="genesis-1676823010098"
BN1="enode://acf68d98f9ba6a7f6eab0d9d04fca859804c48fcb08c2d466532bc8cff9ae97c90d7242deba5dc158fd103d363f9633e2fedd4529286a1e5ff423536185f36dd@127.0.0.1:5050"

./build/opera --port 6050 --http --graphql --ws --http.port 19545  --ws.port 19546 \
              --genesis.allowExperimental --genesis "$GENESIS_FILE" \
              --identity "node2" --nodekey $DEVNET_DIR/node2.key \
              --datadir $DEVNET_DIR/node2 --verbosity=3 \
              --bootnodes "$BN1"

