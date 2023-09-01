#!/bin/bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)
NODE_DIR=$SCRIPT_DIR/../testnet/node1

"$SCRIPT_DIR/../build/arthera-node" --testnet --genesis testnet.g --genesis.allowUnknown --port 6534 --netrestrict 127.0.0.1/8 --cache 5000 \
  --identity "node1" --nodekey "$NODE_DIR/node.key" --nousb \
  --datadir "$NODE_DIR" --verbosity=4 \
  --bootnodes "" \
  --validator.id 1 \
  --validator.pubkey "0xc0041d7405a8bc7dabf1e397e6689ff09482466aea9d3a716bf1dd4fd971c22d035d8d939c88764136a3213106282887f9005b5addf23af781302a0119400706996e" \
  --validator.password "$NODE_DIR/keystore/validator/password"

