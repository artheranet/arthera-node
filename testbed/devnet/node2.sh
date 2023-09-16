#!/bin/bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)
NODE_DIR=$SCRIPT_DIR/node2

"$SCRIPT_DIR/../../build/arthera-node" --devnet --genesis devnet.g --genesis.allowUnknown --port 6535 --netrestrict 127.0.0.1/8 \
  --identity "node2" --nodekey "$NODE_DIR/node.key" --nousb --nat=none --cache 8000 \
  --datadir "$NODE_DIR" --verbosity=2 \
  --bootnodes "enode://2195b2c0ca9695cff2fe84851e6213f2a231d21d056572f6e9bc52b800299beea846ac47ab37256200638aa574fb0387df006638eb0026eaac8ef8a5a3f7b604@127.0.0.1:6534" \
  --validator.id 2 \
  --validator.pubkey "0xc0046218198298ade0acaecde7816c1513c40c359673b516449f4e383d87fa53b54c245a75ed98629f2a35eae140306d7d81b6ba33feec81f0b98ddb0b529c48db32" \
  --validator.password "$NODE_DIR/keystore/validator/password" \
  --http.addr 0.0.0.0 --http.port 18546 --http --http.corsdomain "*" --http.vhosts="*" \
  --http.api=eth,web3,net,txpool,art,abft \
  --ws --ws.addr 0.0.0.0 --ws.port 18556 --ws.origins "*" --ws.api=eth,web3,net,txpool,art,abft
