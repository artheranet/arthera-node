#!/bin/bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)
NODE_DIR=$SCRIPT_DIR/node-api

"$SCRIPT_DIR/../../build/arthera-node" --devnet --genesis devnet.g --genesis.allowUnknown --port 6537 --netrestrict 127.0.0.1/8 --nousb \
  --datadir "$NODE_DIR" --verbosity=2 --cache 8000 --nat=none \
  --bootnodes "enode://4dbc94a60d0d5c91b0fcafd8dd931bb77a2de8b269c80a58da676af3a74fcf9fa5457c536aea40544080780a99b0dcf6629f34f0974d21da7a4c2f62a0074eec@127.0.0.1:6534" \
  --http.addr 0.0.0.0 --http.port 18548 --http --http.corsdomain "*" --http.vhosts="*" \
  --http.api=eth,web3,net,txpool,art,abft \
  --ws --ws.addr 0.0.0.0 --ws.port 18558 --ws.origins "*" --ws.api=eth,web3,net,txpool,art,abft
