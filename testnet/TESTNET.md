1. Generate node keys:

go-ethereum/build/bin/devp2p key generate node1.testnet.key
go-ethereum/build/bin/devp2p key generate node2.testnet.key
go-ethereum/build/bin/devp2p key generate node3.testnet.key

2. Create 3 validators
./build/arthera validator new

ID: 0
Stake: 1_000_000
Public key:  0xc0041d7405a8bc7dabf1e397e6689ff09482466aea9d3a716bf1dd4fd971c22d035d8d939c88764136a3213106282887f9005b5addf23af781302a0119400706996e
Address:     0x7a97E50436a074ADDB9A51D50Fbd35ADAFE88442

ID: 1
Stake: 1_000_000
Public key:  0xc004a61ec5eb3cf8d6b399ff56682b95277337b601fb31e1a254dd451101b8aafb0218d428fc814faee132aabcc17b3dd39fa35dfce2d5ce29d6bd05615bbd571016
Address:     0xfE8301b91A8Eb4734ed954f8E2FB84c2F72Cef8a

ID: 2
Stake: 1_000_000
Public key:  0xc004c39c38dc49cc4c9b64ea9d817545e713635f808d692f2f500ad801e002c50987e15cf4d9419731adf4cd83edf2207a806685cb2b75c3027d2dcdd78ec126f430
Address:     0xF51e935061731a129765ff63b3Af0Adb5e4486aC

3. Create 3 accounts
./build/arthera account new

Balance: 10_000_000

0x40bd65cfc4D95844704F4b2a2c46a60f6d6CE766 / b16a37247cc832aa1859d20a3849debbaa800833580a1ccb66ec0b78326f1c01
0x35E58946b74fDbD9032aed876FC58629A6e65E79 / d4115a3665416057d3d5f76a36f662e7c7b16562b46f89b7ea392e001f2a5036
0x846032c611697818a31cC090D436664b263C6E54 / 623cf3e5f616abe95cbdf1507063eefe98aac6641f4d7054f65d6942def07cbb

4. Update validators, accounts and genesis hashes in:
cmd/arthera/launcher/creategen.go

5. Run ./build/arthera-node --genesis.type=testnet creategen testnet.genesis

6. Change genesis hashes in:
cmd/arthera/launcher/params.go
