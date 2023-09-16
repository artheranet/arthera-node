#!/bin/bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)
NODE_DIR=$SCRIPT_DIR/node1

"$SCRIPT_DIR/../../build/arthera-node" --devnet --genesis devnet.g --genesis.allowUnknown --port 6534 --netrestrict 127.0.0.1/8 \
  --identity "node1" --nodekey "$NODE_DIR/node.key" --nousb --nat=none --cache 8000 \
  --datadir "$NODE_DIR" --verbosity=2 \
  --bootnodes "" \
  --validator.id 1 \
  --validator.pubkey "0xc004bf689e0aa508fc18c9820348cea64cc8b3b3dff85af513fef6309a514c21b33d96e6113904e21c49e012cb73c46d1e5b8ab7cad64131b27a8578d9f87a298f49" \
  --validator.password "$NODE_DIR/keystore/validator/password" \
  --http.addr 0.0.0.0 --http.port 18545 --http --http.corsdomain "*" --http.vhosts="*" \
  --http.api=eth,web3,net,txpool,art,abft \
  --ws --ws.addr 0.0.0.0 --ws.port 18555 --ws.origins "*" --ws.api=eth,web3,net,txpool,art,abft

