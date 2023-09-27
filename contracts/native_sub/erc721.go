package native_sub

import (
	"bytes"
	"github.com/artheranet/arthera-node/contracts"
	"github.com/artheranet/arthera-node/contracts/abis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

const (
	OwnersSlot = 1024 + iota
	BalancesSlot
	TokenApprovalsSlot
)

var (
	transferEvent = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	approvalEvent = crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))

	balanceOfMethodID        []byte
	ownerOfMethodID          []byte
	safeTransferFromMethodID []byte
	transferFromMethodID     []byte
	approveMethodID          []byte
	getApprovedMethodID      []byte
	nameMethodID             []byte
	symbolMethodID           []byte
	tokenURIMethodID         []byte
)

type ERC721 struct {
	evm                 *vm.EVM
	name                string
	symbol              string
	tokenUri            string
	beforeTokenTransfer func(vm.StateDB, common.Address, common.Address, *big.Int) []byte
	afterTokenTransfer  func(vm.StateDB, common.Address, common.Address, *big.Int) []byte
	contractAddress     common.Address
}

func (erc721 *ERC721) ProcessMethod(methodId []byte, caller common.Address, args []byte, suppliedGas uint64) ([]byte, uint64, error) {
	if bytes.Equal(methodId, balanceOfMethodID) {
		return erc721.balanceOf(caller, args, suppliedGas)
	} else if bytes.Equal(methodId, ownerOfMethodID) {
		return erc721.ownerOf(caller, args, suppliedGas)
	} else if bytes.Equal(methodId, safeTransferFromMethodID) {
		return erc721.safeTransferFrom(caller, args, suppliedGas)
	} else if bytes.Equal(methodId, transferFromMethodID) {
		return erc721.transferFrom(caller, args, suppliedGas)
	} else if bytes.Equal(methodId, approveMethodID) {
		return erc721.approve(caller, args, suppliedGas)
	} else if bytes.Equal(methodId, getApprovedMethodID) {
		return erc721.getApproved(caller, args, suppliedGas)
	} else if bytes.Equal(methodId, nameMethodID) {
		return erc721.Name(caller, args, suppliedGas)
	} else if bytes.Equal(methodId, symbolMethodID) {
		return erc721.Symbol(caller, args, suppliedGas)
	} else if bytes.Equal(methodId, tokenURIMethodID) {
		return erc721.TokenURI(caller, args, suppliedGas)
	}

	return nil, suppliedGas, nil
}

// function balanceOf(address owner) returns (uint256 balance);
func (erc721 *ERC721) balanceOf(_ common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	if !abis.HasNumArgs(input, 1) {
		return nil, suppliedGas, vm.ErrExecutionReverted
	}
	owner := abis.GetAddressArg(input, 0)
	if owner == abis.ZeroAddress {
		return abis.PackRevert("ERC721: address zero is not a valid owner"), suppliedGas, vm.ErrExecutionReverted
	}
	return abis.PackAbiUint256(erc721.getBalance(owner)), suppliedGas, nil
}

// function ownerOf(uint256 tokenId) returns (address owner);
func (erc721 *ERC721) ownerOf(_ common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	if !abis.HasNumArgs(input, 1) {
		return nil, suppliedGas, vm.ErrExecutionReverted
	}
	tokenId := abis.GetUint256Arg(input, 0)
	owner := erc721.getOwner(tokenId)
	return abis.PackAbiAddress(owner), suppliedGas, nil
}

// function safeTransferFrom(address from, address to, uint256 tokenId)
func (erc721 *ERC721) safeTransferFrom(caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return erc721.transferFrom(caller, input, suppliedGas)

}

// function transferFrom(address from, address to, uint256 tokenId)
func (erc721 *ERC721) transferFrom(caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	if !abis.HasNumArgs(input, 3) {
		return nil, suppliedGas, vm.ErrExecutionReverted
	}
	tokenId := abis.GetUint256Arg(input, 2)
	approved, err := erc721._isApprovedOrOwner(caller, tokenId)
	if err != nil {
		return err, suppliedGas, vm.ErrExecutionReverted
	}
	if !approved {
		return abis.PackRevert("ERC721: caller is not token owner or approved"), suppliedGas, vm.ErrExecutionReverted
	}

	err = erc721._transfer(abis.GetAddressArg(input, 0), abis.GetAddressArg(input, 1), tokenId)
	if err != nil {
		return err, suppliedGas, vm.ErrExecutionReverted
	}

	return nil, suppliedGas, vm.ErrExecutionReverted
}

// function approve(address to, uint256 tokenId)
func (erc721 *ERC721) approve(caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	if !abis.HasNumArgs(input, 2) {
		return nil, suppliedGas, vm.ErrExecutionReverted
	}
	to := abis.GetAddressArg(input, 0)
	tokenId := abis.GetUint256Arg(input, 1)

	owner := erc721.getOwner(tokenId)
	if to == owner {
		return abis.PackRevert("ERC721: approval to current owner"), suppliedGas, vm.ErrExecutionReverted
	}

	if caller != owner {
		return abis.PackRevert("ERC721: approve caller is not token owner"), suppliedGas, vm.ErrExecutionReverted
	}

	erc721.setTokenApproval(tokenId, to)
	erc721.emitApprovalEvent(erc721.getOwner(tokenId), to, tokenId)
	return nil, suppliedGas, vm.ErrExecutionReverted
}

// function getApproved(uint256 tokenId) returns (address operator);
func (erc721 *ERC721) getApproved(_ common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	if !abis.HasNumArgs(input, 1) {
		return nil, suppliedGas, vm.ErrExecutionReverted
	}
	tokenId := abis.GetUint256Arg(input, 0)
	addr, err := erc721._getApproved(tokenId)
	if err != nil {
		return err, suppliedGas, vm.ErrExecutionReverted
	}
	return abis.PackAbiAddress(*addr), suppliedGas, nil
}

func (erc721 *ERC721) Name(_ common.Address, _ []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return abis.PackAbiString(erc721.name), suppliedGas, nil
}

func (erc721 *ERC721) Symbol(_ common.Address, _ []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return abis.PackAbiString(erc721.symbol), suppliedGas, nil
}

func (erc721 *ERC721) TokenURI(_ common.Address, _ []byte, suppliedGas uint64) ([]byte, uint64, error) {
	return abis.PackAbiString(erc721.tokenUri), suppliedGas, nil
}

func (erc721 *ERC721) _isApprovedOrOwner(spender common.Address, tokenId *big.Int) (bool, []byte) {
	owner := erc721.getOwner(tokenId)
	approvedAddr, err := erc721._getApproved(tokenId)
	if err != nil {
		return false, err
	}
	return spender == owner || *approvedAddr == spender, nil
}

func (erc721 *ERC721) _exists(tokenId *big.Int) bool {
	owner := erc721.getOwner(tokenId)
	return owner != abis.ZeroAddress
}

func (erc721 *ERC721) _requireMinted(tokenId *big.Int) []byte {
	exists := erc721._exists(tokenId)
	if !exists {
		return abis.PackRevert("ERC721: token does not exist")
	}
	return nil
}

func (erc721 *ERC721) _getApproved(tokenId *big.Int) (*common.Address, []byte) {
	err := erc721._requireMinted(tokenId)
	if err != nil {
		return nil, err
	}
	addr := common.BytesToAddress(
		abis.GetMappingData(erc721.evm.StateDB, contracts.SubscribersSmartContractAddress, TokenApprovalsSlot, tokenId.Bytes()).Bytes(),
	)
	return &addr, nil
}

func (erc721 *ERC721) _transfer(from common.Address, to common.Address, tokenId *big.Int) []byte {
	owner := erc721.getOwner(tokenId)
	if owner != from {
		return abis.PackRevert("ERC721: transfer from incorrect owner")
	}
	if to == abis.ZeroAddress {
		return abis.PackRevert("ERC721: transfer to the zero address")
	}

	if erc721.beforeTokenTransfer != nil {
		err := erc721.beforeTokenTransfer(erc721.evm.StateDB, from, to, tokenId)
		if err != nil {
			return err
		}
	}

	// Check that tokenId was not transferred by 'beforeTokenTransfer` hook
	owner = erc721.getOwner(tokenId)
	if owner != from {
		return abis.PackRevert("ERC721: transfer from incorrect owner")
	}

	// Clear approvals from the previous owner
	erc721.setTokenApproval(tokenId, abis.ZeroAddress)

	// Update balances
	balanceFrom := new(big.Int).Sub(erc721.getBalance(from), big.NewInt(1))
	erc721.setBalance(from, balanceFrom)
	balanceTo := new(big.Int).Add(erc721.getBalance(to), big.NewInt(1))
	erc721.setBalance(to, balanceTo)

	erc721.setOwner(tokenId, to)
	erc721.emitTransferEvent(from, to, tokenId)

	if erc721.afterTokenTransfer != nil {
		err := erc721.afterTokenTransfer(erc721.evm.StateDB, from, to, tokenId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (erc721 *ERC721) _mint(to common.Address, tokenId *big.Int) []byte {
	if to == abis.ZeroAddress {
		return abis.PackRevert("ERC721: mint to the zero address")
	}

	if erc721._exists(tokenId) {
		return abis.PackRevert("ERC721: token already minted")
	}

	if erc721.beforeTokenTransfer != nil {
		err := erc721.beforeTokenTransfer(erc721.evm.StateDB, abis.ZeroAddress, to, tokenId)
		if err != nil {
			return err
		}
	}

	// Check that tokenId was not minted by `beforeTokenTransfer` hook
	if erc721._exists(tokenId) {
		return abis.PackRevert("ERC721: token already minted")
	}

	balance := new(big.Int).Add(erc721.getBalance(to), big.NewInt(1))
	erc721.setBalance(to, balance)
	erc721.setOwner(tokenId, to)
	erc721.emitTransferEvent(abis.ZeroAddress, to, tokenId)

	if erc721.afterTokenTransfer != nil {
		err := erc721.afterTokenTransfer(erc721.evm.StateDB, abis.ZeroAddress, to, tokenId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (erc721 *ERC721) _burn(tokenId *big.Int) []byte {
	owner := erc721.getOwner(tokenId)
	if erc721.beforeTokenTransfer != nil {
		err := erc721.beforeTokenTransfer(erc721.evm.StateDB, owner, abis.ZeroAddress, tokenId)
		if err != nil {
			return err
		}
	}
	// Update ownership in case tokenId was transferred by `beforeTokenTransfer` hook
	owner = erc721.getOwner(tokenId)

	// Clear approvals
	erc721.setTokenApproval(tokenId, abis.ZeroAddress)

	balance := new(big.Int).Add(erc721.getBalance(owner), big.NewInt(1))
	erc721.setBalance(owner, balance)
	erc721.setOwner(tokenId, owner)
	erc721.emitTransferEvent(owner, abis.ZeroAddress, tokenId)

	if erc721.afterTokenTransfer != nil {
		err := erc721.afterTokenTransfer(erc721.evm.StateDB, owner, abis.ZeroAddress, tokenId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (erc721 *ERC721) getBalance(owner common.Address) *big.Int {
	return abis.GetMappingData(erc721.evm.StateDB, erc721.contractAddress, BalancesSlot, owner.Bytes()).Big()
}

func (erc721 *ERC721) setBalance(owner common.Address, balance *big.Int) {
	abis.SetMappingData(
		erc721.evm.StateDB, erc721.contractAddress, BalancesSlot,
		owner.Bytes(), common.BigToHash(balance),
	)
}

func (erc721 *ERC721) setTokenApproval(tokenId *big.Int, addr common.Address) {
	abis.SetMappingData(
		erc721.evm.StateDB, erc721.contractAddress, TokenApprovalsSlot,
		tokenId.Bytes(), addr.Hash(),
	)
}

func (erc721 *ERC721) getTokenApproval(tokenId *big.Int) common.Address {
	return common.BytesToAddress(
		abis.GetMappingData(erc721.evm.StateDB, erc721.contractAddress, TokenApprovalsSlot, tokenId.Bytes()).Bytes(),
	)
}

func (erc721 *ERC721) setOwner(tokenId *big.Int, addr common.Address) {
	abis.SetMappingData(
		erc721.evm.StateDB, erc721.contractAddress, OwnersSlot,
		tokenId.Bytes(), addr.Hash(),
	)
}

func (erc721 *ERC721) getOwner(tokenId *big.Int) common.Address {
	return common.BytesToAddress(
		abis.GetMappingData(erc721.evm.StateDB, contracts.SubscribersSmartContractAddress, OwnersSlot, tokenId.Bytes()).Bytes(),
	)
}

func (erc721 *ERC721) emitTransferEvent(from common.Address, to common.Address, tokenId *big.Int) {
	var topics = []common.Hash{
		transferEvent,
		common.BytesToHash(common.LeftPadBytes(from.Bytes(), 32)),
		common.BytesToHash(common.LeftPadBytes(to.Bytes(), 32)),
	}
	event := types.Log{
		Address:     erc721.contractAddress,
		Topics:      topics,
		Data:        abis.PackAbiUint256(tokenId),
		BlockNumber: erc721.evm.Context.BlockNumber.Uint64(),
	}
	erc721.evm.StateDB.AddLog(&event)
}

func (erc721 *ERC721) emitApprovalEvent(owner common.Address, spender common.Address, tokenId *big.Int) {
	var topics = []common.Hash{
		approvalEvent,
		common.BytesToHash(common.LeftPadBytes(owner.Bytes(), 32)),
		common.BytesToHash(common.LeftPadBytes(spender.Bytes(), 32)),
	}
	event := types.Log{
		Address:     erc721.contractAddress,
		Topics:      topics,
		Data:        abis.PackAbiUint256(tokenId),
		BlockNumber: erc721.evm.Context.BlockNumber.Uint64(),
	}
	erc721.evm.StateDB.AddLog(&event)
}
