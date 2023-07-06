// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
//
// StateProcessor implements Processor.
type StateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for block rewards
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(config *params.ChainConfig, bc *BlockChain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

// Process processes the state changes according to the Ethereum rules by running
// the transaction messages using the statedb and applying any rewards to both
// the processor (coinbase) and any included uncles.
//
// Process returns the receipts and logs accumulated during the process and
// returns the amount of gas that was used in the process. If any of the
// transactions failed to execute due to insufficient gas it will return an error.
func (p *StateProcessor) Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, uint64, error) {
	var (
		receipts    types.Receipts
		usedGas     = new(uint64)
		header      = block.Header()
		blockHash   = block.Hash()
		blockNumber = block.Number()
		allLogs     []*types.Log
		gp          = new(GasPool).AddGas(block.GasLimit())
	)
	// Mutate the block and state according to any hard-fork specs
	if p.config.DAOForkSupport && p.config.DAOForkBlock != nil && p.config.DAOForkBlock.Cmp(block.Number()) == 0 {
		misc.ApplyDAOHardFork(statedb)
	}
	var (
		context = NewEVMBlockContext(header, p.bc, nil)
		vmenv   = vm.NewEVM(context, vm.TxContext{}, statedb, p.config, cfg)
		signer  = types.MakeSigner(p.config, header.Number, header.Time)
	)
	if len(block.Transactions()) > 0 {
		types.InitFile(blockNumber)
	}
	// Iterate over and process the individual transactions
	for i, tx := range block.Transactions() {
		msg, err := TransactionToMessage(tx, signer, header.BaseFee)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("could not apply tx %d [%v]: %w", i, tx.Hash().Hex(), err)
		}
		statedb.SetTxContext(tx.Hash(), i)

		// minerExtra
		types.WriteTxHash(blockNumber, tx.Hash())

		types.WritePreRoot(blockNumber, statedb.IntermediateRoot(p.config.IsEIP158(blockNumber)))
		l := types.ReadTxFile(blockNumber, tx.Hash())
		for _, s := range l {
			sArr := strings.Split(s, "\t")
			t := sArr[0]
			switch t {
			case "0":
				addr := common.HexToAddress(sArr[1])
				if enc, error := statedb.GetAccountEnc(addr); error == nil {
					//proof, err := statedb.GetProof(addr)
					//if err != nil {
					//	log.Info("GetProof Err", "addr", addr, "GetProof", err)
					//	return nil, nil, 0, err
					//}
					types.WritePreState(blockNumber, addr, enc)
					//types.WritePreStateProof(blockNumber, addr, proof)
					//log.Info("账户信息", "blockNumber", blockNumber, "addr", addr)
					//roothash := statedb.IntermediateRoot(p.config.IsEIP158(blockNumber))
					//tempProof := types.NewExtraProof(proof)
					//log.Info("Proof信息", "accountProof", roothash, "roothash[:]", roothash[:], "tempProof", tempProof, "proof", proof)
					//value2, err2 := trie.VerifyProof(roothash, crypto.Keccak256Hash(addr.Bytes()).Bytes(), tempProof)
					//log.Info("Proof验证", "value2", value2, "err2", err2)

					//data := new(types.StateAccount)
					//if err := rlp.DecodeBytes(enc, data); err != nil {
					//	log.Error("Failed to decode state object", "addr", addr, "err", err)
					//}
					//log.Info("Account", "blockNumber", blockNumber, "hash", tx.Hash().Hex(), "addr", addr, "stobject", data)

				}
			case "1":
				addr := common.HexToAddress(sArr[1])
				hash := common.HexToHash(sArr[2])
				val := statedb.GetState(addr, hash)
				//txExtra.AddPreStorage(addr, hash, val)
				types.WritePreStorage(blockNumber, addr, hash, val)
				//if proof, err := statedb.GetStorageProof(addr, hash); err == nil {
				//if tx.Hash().String() == "0x10d0e6b3bce2988501490c31a2ba9a2e3b49c27ef8381480cebde87eb8402d96" && addr.String() == "0x0000000000000000000000000000000000001002" {
				//	enc, _ := Ctrie.TryGet(addr.Bytes())
				//	data := new(types.StateAccount)
				//	if err := rlp.DecodeBytes(enc, data); err != nil {
				//		log.Error("Failed to decode state object", "addr", addr, "err", err)
				//	}
				//	log.Info("GetStorageProof", "blockNumber", blockNumber, "addr", addr, "hash", hash, "proof", proof, "account", data.Root)
				//}
				//types.WritePreStorageProof(blockNumber, addr, hash, proof)

				//enc, _ := Ctrie.TryGet(addr.Bytes())
				//data := new(types.StateAccount)
				//if err := rlp.DecodeBytes(enc, data); err != nil {
				//	log.Error("Failed to decode state object", "addr", addr, "err", err)
				//}
				//tempProof := types.NewExtraProof(proof)
				//log.Info("Proof信息", "root", data.Root.Hex(), "key", hash, "val", val, "tempProof", tempProof)

				//value1, _ := trie.VerifyProof(data.Root, crypto.Keccak256Hash(hash.Bytes()).Bytes(), tempProof)

				//log.Info("txExtra", "hash", tx.Hash().Hex(), "addr", addr, "key", hash.Hex(), "c", val, "value", value1)
				//log.Info("Proof验证", "value1", value1, "err1", err1)
				//if !bytes.Equal(common.BytesToHash(value1).Bytes(), val.Bytes()) {
				//	log.Info("测试数据", "key", hash.String(), "hash", tx.Hash().Hex(), "value", value1, "val", val.Bytes())
				//}

				//value2, err2 := trie.VerifyProof(roothash, crypto.Keccak256Hash(hash.Bytes()).Bytes(), tempProof)
				//log.Info("Proof验证", "value2", value2, "err2", err2)
				//txExtra.AddPreStorageProof(addr, hash, proof)
				//}
			case "2":
				addr := common.HexToAddress(sArr[1])
				code := statedb.GetCode(addr)
				types.WritePreCode(blockNumber, addr, code)
			default:
			}
		}
		receipt, err := applyTransaction(msg, p.config, gp, statedb, blockNumber, blockHash, tx, usedGas, vmenv)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("could not apply tx %d [%v]: %w", i, tx.Hash().Hex(), err)
		}
		//minerExtra
		types.WritePostRoot(blockNumber, statedb.IntermediateRoot(p.config.IsEIP158(blockNumber)))

		//root2 := statedb.IntermediateRoot_new(p.bc.chainConfig.IsEIP158(block.Number()))
		//log.Info("查看Trie更改", "block", blockNumber, "header.root", block.Header().Root.Hex(), "root", root2.Hex())

		//for _, s := range l {
		//	sArr := strings.Split(s, "\t")
		//	t := sArr[0]
		//	switch t {
		//	case "0":
		//		addr := common.HexToAddress(sArr[1])
		//		if enc, error := statedb.GetAccountEnc(addr); error == nil {
		//			//txExtra.AddPostState(addr, enc)
		//			proof, err := statedb.GetProof(addr)
		//			if err != nil {
		//				return nil, nil, 0, err
		//			}
		//			//txExtra.AddPostStateProof(addr, proof)
		//			types.WritePostState(blockNumber, addr, enc)
		//			types.WritePostStateProof(blockNumber, addr, proof)
		//		}
		//
		//	case "1":
		//		addr := common.HexToAddress(sArr[1])
		//		hash := common.HexToHash(sArr[2])
		//		val := statedb.GetState(addr, hash)
		//		//txExtra.AddPostStorage(addr, hash, val)
		//		types.WritePostStorage(blockNumber, addr, hash, val)
		//		if proof, err := statedb.GetStorageProof(addr, hash); err == nil {
		//			types.WritePostStorageProof(blockNumber, addr, hash, proof)
		//		}
		//	case "2":
		//		addr := common.HexToAddress(sArr[1])
		//		code := statedb.GetCode(addr)
		//		types.WritePostCode(blockNumber, addr, code)
		//	default:
		//
		//	}
		//}
		types.WriteTxHash(blockNumber, tx.Hash())

		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
	}
	// Fail if Shanghai not enabled and len(withdrawals) is non-zero.
	withdrawals := block.Withdrawals()
	if len(withdrawals) > 0 && !p.config.IsShanghai(block.Number(), block.Time()) {
		return nil, nil, 0, errors.New("withdrawals before shanghai")
	}
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	p.engine.Finalize(p.bc, header, statedb, block.Transactions(), block.Uncles(), withdrawals)

	return receipts, allLogs, *usedGas, nil
}

func applyTransaction(msg *Message, config *params.ChainConfig, gp *GasPool, statedb *state.StateDB, blockNumber *big.Int, blockHash common.Hash, tx *types.Transaction, usedGas *uint64, evm *vm.EVM) (*types.Receipt, error) {
	// Create a new context to be used in the EVM environment.
	txContext := NewEVMTxContext(msg)
	evm.Reset(txContext, statedb)

	// Apply the transaction to the current state (included in the env).
	result, err := ApplyMessage(evm, msg, gp)
	if err != nil {
		return nil, err
	}

	// Update the state with pending changes.
	var root []byte
	if config.IsByzantium(blockNumber) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(blockNumber)).Bytes()
	}
	*usedGas += result.UsedGas

	// Create a new receipt for the transaction, storing the intermediate root and gas used
	// by the tx.
	receipt := &types.Receipt{Type: tx.Type(), PostState: root, CumulativeGasUsed: *usedGas}
	if result.Failed() {
		receipt.Status = types.ReceiptStatusFailed
	} else {
		receipt.Status = types.ReceiptStatusSuccessful
	}
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = result.UsedGas

	// If the transaction created a contract, store the creation address in the receipt.
	if msg.To == nil {
		receipt.ContractAddress = crypto.CreateAddress(evm.TxContext.Origin, tx.Nonce())
	}

	// Set the receipt logs and create the bloom filter.
	receipt.Logs = statedb.GetLogs(tx.Hash(), blockNumber.Uint64(), blockHash)
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	receipt.BlockHash = blockHash
	receipt.BlockNumber = blockNumber
	receipt.TransactionIndex = uint(statedb.TxIndex())
	return receipt, err
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *params.ChainConfig, bc ChainContext, author *common.Address, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg vm.Config) (*types.Receipt, error) {
	msg, err := TransactionToMessage(tx, types.MakeSigner(config, header.Number, header.Time), header.BaseFee)
	if err != nil {
		return nil, err
	}
	// Create a new context to be used in the EVM environment
	blockContext := NewEVMBlockContext(header, bc, author)
	vmenv := vm.NewEVM(blockContext, vm.TxContext{}, statedb, config, cfg)
	return applyTransaction(msg, config, gp, statedb, header.Number, header.Hash(), tx, usedGas, vmenv)
}
