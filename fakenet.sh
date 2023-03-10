#!/bin/bash

./build/arthera-node --fakenet 1/1 --http --http.corsdomain "*" \
    --ws --ws.origins "*" \
    --http.api=eth,web3,net,txpool,art,abft \
    --ws.api=eth,web3,net,txpool,art,abft
