#!/bin/bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)
NODE_DIR=$SCRIPT_DIR/../testnet/node3

"$SCRIPT_DIR/../../build/arthera-node" --testnet --genesis testnet2.g --genesis.allowUnknown --port 6536 --netrestrict 127.0.0.1/8 --cache 5000 \
  --identity "node3" --nodekey "$NODE_DIR/node.key" --nousb \
  --datadir "$NODE_DIR" --verbosity=4 \
  --bootnodes "enode://4dbc94a60d0d5c91b0fcafd8dd931bb77a2de8b269c80a58da676af3a74fcf9fa5457c536aea40544080780a99b0dcf6629f34f0974d21da7a4c2f62a0074eec@127.0.0.1:6534" \
  --validator.id 3 \
  --validator.pubkey "0xc004c39c38dc49cc4c9b64ea9d817545e713635f808d692f2f500ad801e002c50987e15cf4d9419731adf4cd83edf2207a806685cb2b75c3027d2dcdd78ec126f430" \
  --validator.password "$NODE_DIR/keystore/validator/password"
