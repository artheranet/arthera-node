#!/bin/bash

./build/arthera-node --fakenet 1/1 --http.addr 0.0.0.0 --http --http.corsdomain "*" \
    --http.addr 0.0.0.0 --http.vhosts="*" --ws --ws.origins "*" \
    --http.api=eth,web3,net,txpool,art,abft,debug \
    --ws.api=eth,web3,net,txpool,art,abft \
    --verbosity=4
