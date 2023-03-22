package abis

const NodeDriverStr = `[
	{
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "num",
          "type": "uint256"
        }
      ],
      "name": "AdvanceEpochs",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint8",
          "name": "version",
          "type": "uint8"
        }
      ],
      "name": "Initialized",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "bytes",
          "name": "diff",
          "type": "bytes"
        }
      ],
      "name": "UpdateNetworkRules",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "version",
          "type": "uint256"
        }
      ],
      "name": "UpdateNetworkVersion",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "validatorID",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "bytes",
          "name": "pubkey",
          "type": "bytes"
        }
      ],
      "name": "UpdateValidatorPubkey",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "validatorID",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "weight",
          "type": "uint256"
        }
      ],
      "name": "UpdateValidatorWeight",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "backend",
          "type": "address"
        }
      ],
      "name": "UpdatedBackend",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "num",
          "type": "uint256"
        }
      ],
      "name": "advanceEpochs",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "acc",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "from",
          "type": "address"
        }
      ],
      "name": "copyCode",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "validatorID",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "status",
          "type": "uint256"
        }
      ],
      "name": "deactivateValidator",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "acc",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "diff",
          "type": "uint256"
        }
      ],
      "name": "incNonce",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_backend",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_evmWriterAddress",
          "type": "address"
        }
      ],
      "name": "initialize",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256[]",
          "name": "offlineTimes",
          "type": "uint256[]"
        },
        {
          "internalType": "uint256[]",
          "name": "offlineBlocks",
          "type": "uint256[]"
        },
        {
          "internalType": "uint256[]",
          "name": "uptimes",
          "type": "uint256[]"
        },
        {
          "internalType": "uint256[]",
          "name": "originatedTxsFee",
          "type": "uint256[]"
        }
      ],
      "name": "sealEpoch",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256[]",
          "name": "nextValidatorIDs",
          "type": "uint256[]"
        }
      ],
      "name": "sealEpochValidators",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_backend",
          "type": "address"
        }
      ],
      "name": "setBackend",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "acc",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "setBalance",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "delegator",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "toValidatorID",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "stake",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "lockedStake",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "lockupFromEpoch",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "lockupEndTime",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "lockupDuration",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "earlyUnlockPenalty",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "rewards",
          "type": "uint256"
        }
      ],
      "name": "setGenesisDelegation",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_auth",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "validatorID",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "pubkey",
          "type": "bytes"
        },
        {
          "internalType": "uint256",
          "name": "status",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "createdEpoch",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "createdTime",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "deactivatedEpoch",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "deactivatedTime",
          "type": "uint256"
        }
      ],
      "name": "setGenesisValidator",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "acc",
          "type": "address"
        },
        {
          "internalType": "bytes32",
          "name": "key",
          "type": "bytes32"
        },
        {
          "internalType": "bytes32",
          "name": "value",
          "type": "bytes32"
        }
      ],
      "name": "setStorage",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "acc",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "with",
          "type": "address"
        }
      ],
      "name": "swapCode",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes",
          "name": "diff",
          "type": "bytes"
        }
      ],
      "name": "updateNetworkRules",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "version",
          "type": "uint256"
        }
      ],
      "name": "updateNetworkVersion",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "validatorID",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "pubkey",
          "type": "bytes"
        }
      ],
      "name": "updateValidatorPubkey",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "validatorID",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "updateValidatorWeight",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
]`

const EVMWriterStr = `[
	{
      "inputs": [
        {
          "internalType": "address",
          "name": "acc",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "from",
          "type": "address"
        }
      ],
      "name": "copyCode",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "acc",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "diff",
          "type": "uint256"
        }
      ],
      "name": "incNonce",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "acc",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "setBalance",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "acc",
          "type": "address"
        },
        {
          "internalType": "bytes32",
          "name": "key",
          "type": "bytes32"
        },
        {
          "internalType": "bytes32",
          "name": "value",
          "type": "bytes32"
        }
      ],
      "name": "setStorage",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "acc",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "with",
          "type": "address"
        }
      ],
      "name": "swapCode",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
]`

const NetworkInitializerStr = `[
	{
      "inputs": [],
      "name": "EVM_WRITER",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "NODE_DRIVER",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "NODE_DRIVER_AUTH",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "PYAG_REBATES",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "STAKING",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "SUBSCRIBERS",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "VALIDATOR_INFO",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "sealedEpoch",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "totalSupply",
          "type": "uint256"
        },
        {
          "internalType": "address",
          "name": "_owner",
          "type": "address"
        }
      ],
      "name": "initializeAll",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
]`

const SubscribersStr = `[
	{
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint8",
          "name": "version",
          "type": "uint8"
        }
      ],
      "name": "Initialized",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "previousOwner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "OwnershipTransferred",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "planId",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "duration",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "units",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "usdPrice",
          "type": "uint256"
        },
        {
          "internalType": "enum SubscriptionPlan.CapFrequency",
          "name": "capFrequency",
          "type": "uint8"
        },
        {
          "internalType": "uint256",
          "name": "capUnits",
          "type": "uint256"
        },
        {
          "internalType": "bool",
          "name": "dapp",
          "type": "bool"
        }
      ],
      "name": "createPlan",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "subscriber",
          "type": "address"
        },
        {
          "internalType": "uint56",
          "name": "units",
          "type": "uint56"
        }
      ],
      "name": "credit",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "subscriber",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "units",
          "type": "uint256"
        }
      ],
      "name": "debit",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "subscriber",
          "type": "address"
        }
      ],
      "name": "getBalance",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "subscriber",
          "type": "address"
        }
      ],
      "name": "getCapRemaining",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "subscriber",
          "type": "address"
        }
      ],
      "name": "getCapWindow",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "subscriber",
          "type": "address"
        }
      ],
      "name": "getEndTime",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "getPlan",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "planId",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "duration",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "units",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "price",
          "type": "uint256"
        },
        {
          "internalType": "enum SubscriptionPlan.CapFrequency",
          "name": "capFrequency",
          "type": "uint8"
        },
        {
          "internalType": "uint256",
          "name": "capUnits",
          "type": "uint256"
        },
        {
          "internalType": "bool",
          "name": "dapp",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "subscriber",
          "type": "address"
        }
      ],
      "name": "getStartTime",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "name": "getSub",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "planId",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "balance",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "startTime",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "endTime",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "lastCapReset",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "periodUsage",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "subscriber",
          "type": "address"
        }
      ],
      "name": "hasActiveSubscription",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "subscriber",
          "type": "address"
        }
      ],
      "name": "hasSubscription",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "initialize",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "isOwner",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "owner",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "renewSubscription",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "renounceOwnership",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_priceProvider",
          "type": "address"
        }
      ],
      "name": "setPriceProvider",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "dapp",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "planId",
          "type": "uint256"
        }
      ],
      "name": "subscribeToDappPlan",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "planId",
          "type": "uint256"
        }
      ],
      "name": "subscribeToUserPlan",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "transferOwnership",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "planId",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "duration",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "units",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "usdPrice",
          "type": "uint256"
        },
        {
          "internalType": "enum SubscriptionPlan.CapFrequency",
          "name": "capFrequency",
          "type": "uint8"
        },
        {
          "internalType": "uint256",
          "name": "capUnits",
          "type": "uint256"
        }
      ],
      "name": "updatePlan",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
]`
