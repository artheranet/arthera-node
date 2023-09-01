#!/bin/bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)
rm -rf $SCRIPT_DIR/node1/arthera-node
rm -rf $SCRIPT_DIR/node1/chaindata
rm -rf $SCRIPT_DIR/node2/arthera-node
rm -rf $SCRIPT_DIR/node2/chaindata
rm -rf $SCRIPT_DIR/node3/arthera-node
rm -rf $SCRIPT_DIR/node3/chaindata
rm -rf $SCRIPT_DIR/node-api/arthera-node
rm -rf $SCRIPT_DIR/node-api/chaindata
