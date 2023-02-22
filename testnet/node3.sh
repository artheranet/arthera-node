#!/bin/bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)
NODE_DIR=$SCRIPT_DIR/../testnet/node3
GENESIS_FILE="$SCRIPT_DIR/../testnet/testnet.genesis"
BN1="enode://4dbc94a60d0d5c91b0fcafd8dd931bb77a2de8b269c80a58da676af3a74fcf9fa5457c536aea40544080780a99b0dcf6629f34f0974d21da7a4c2f62a0074eec@127.0.0.1:6060"

./build/arthera-node --testnet --port 6062 \
              --genesis "$GENESIS_FILE" --genesis.allowExperimental  \
              --identity "node2" --nodekey "$NODE_DIR/node.key" \
              --datadir "$NODE_DIR" --verbosity=3 \
              --bootnodes "$BN1" --allow-insecure-unlock \
              --validator.id 2 \
              --validator.pubkey "0xc004c39c38dc49cc4c9b64ea9d817545e713635f808d692f2f500ad801e002c50987e15cf4d9419731adf4cd83edf2207a806685cb2b75c3027d2dcdd78ec126f430" \
              --validator.password "$NODE_DIR/keystore/validator/password"
