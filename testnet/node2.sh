#!/bin/bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)
NODE_DIR=$SCRIPT_DIR/../testnet/node2

"$SCRIPT_DIR/../build/arthera-node" --testnet --port 6535 --netrestrict 127.0.0.1/8 \
  --identity "node2" --nodekey "$NODE_DIR/node.key" \
  --datadir "$NODE_DIR" --verbosity=3 --allow-insecure-unlock \
  --bootnodes "enode://4dbc94a60d0d5c91b0fcafd8dd931bb77a2de8b269c80a58da676af3a74fcf9fa5457c536aea40544080780a99b0dcf6629f34f0974d21da7a4c2f62a0074eec@127.0.0.1:6534" \
  --validator.id 2 \
  --validator.pubkey "0xc004a61ec5eb3cf8d6b399ff56682b95277337b601fb31e1a254dd451101b8aafb0218d428fc814faee132aabcc17b3dd39fa35dfce2d5ce29d6bd05615bbd571016" \
  --validator.password "$NODE_DIR/keystore/validator/password"
