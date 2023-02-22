#!/bin/bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)
NODE_DIR=$SCRIPT_DIR/../testnet/node1
GENESIS_FILE="$SCRIPT_DIR/../testnet/testnet.genesis"

./build/arthera-node --testnet --port 6060 --http --ws --http.port 18545  --ws.port 18546 \
              --genesis "$GENESIS_FILE" --genesis.allowExperimental  \
              --identity "node1" --nodekey $NODE_DIR/node.key \
              --datadir $NODE_DIR --verbosity=3 \
              --bootnodes "" --allow-insecure-unlock \
              --validator.id 0 \
              --validator.pubkey "0xc0041d7405a8bc7dabf1e397e6689ff09482466aea9d3a716bf1dd4fd971c22d035d8d939c88764136a3213106282887f9005b5addf23af781302a0119400706996e" \
              --validator.password $NODE_DIR/keystore/validator/password

