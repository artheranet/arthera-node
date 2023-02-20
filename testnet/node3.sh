#!/bin/bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)
DEVNET_DIR=$SCRIPT_DIR/../devnet

GENESIS_FILE="genesis-1676823010098"
BN1="enode://acf68d98f9ba6a7f6eab0d9d04fca859804c48fcb08c2d466532bc8cff9ae97c90d7242deba5dc158fd103d363f9633e2fedd4529286a1e5ff423536185f36dd@127.0.0.1:5050"
BN2="enode://d83a16757082952a84660aebf421c7a57ab80f18650c3c4f4dee9dde33aac72ccb5813eb7fd082711817cd80f481d856b74aa91ca14ff84c16c1f96a3d7084d2@127.0.0.1:6050"

./build/opera --port 7050 --http --graphql --ws --http.port 20545  --ws.port 20546 \
              --genesis.allowExperimental --genesis "$GENESIS_FILE" \
              --identity "node3" --nodekey $DEVNET_DIR/node3.key \
              --datadir $DEVNET_DIR/node3 --verbosity=3 \
              --bootnodes "$BN1,$BN2"

