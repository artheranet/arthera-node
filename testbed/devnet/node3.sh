#!/bin/bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)
NODE_DIR=$SCRIPT_DIR//node3

"$SCRIPT_DIR/../../build/arthera-node" --devnet --genesis devnet.g --genesis.allowUnknown --port 6536 --netrestrict 127.0.0.1/8 \
  --identity "node3" --nodekey "$NODE_DIR/node.key" --nousb \
  --datadir "$NODE_DIR" --verbosity=4 \
  --bootnodes "enode://2195b2c0ca9695cff2fe84851e6213f2a231d21d056572f6e9bc52b800299beea846ac47ab37256200638aa574fb0387df006638eb0026eaac8ef8a5a3f7b604@127.0.0.1:6534" \
  --validator.id 3 \
  --validator.pubkey "0xc004c0cc8ddc257ed1aadd6aec58c40592c00ce653c1e96c5856f3cbb57371b01a7e8c0f2b76255271004673930ebcbc798c360f8df39c08e831bc815ee9f13dc6a6" \
  --validator.password "$NODE_DIR/keystore/validator/password"
