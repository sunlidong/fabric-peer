/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blockledger

import (
	cb "github.com/hyperledger/fabric-protos-go/common"
	ab "github.com/hyperledger/fabric-protos-go/orderer"
)

// Factory retrieves or creates new ledgers by channelID
// //工厂根据通道ID检索或创建新的分类账
type Factory interface {
	// GetOrCreate gets an existing ledger (if it exists)
	// or creates it if it does not
	// 获取或创建
	GetOrCreate(channelID string) (ReadWriter, error)

	// ChannelIDs returns the channel IDs the Factory is aware of
	// 通道ID
	ChannelIDs() []string

	// Close releases all resources acquired by the factory
	// 关闭
	Close()
}

// Iterator is useful for a chain Reader to stream blocks as they are created

// 迭代
type Iterator interface {
	// Next blocks until there is a new block available, or returns an error if
	// the next block is no longer retrievable
	Next() (*cb.Block, cb.Status)
	// Close releases resources acquired by the Iterator
	Close()
}

// Reader allows the caller to inspect the ledger
// 读
type Reader interface {
	// Iterator returns an Iterator, as specified by an ab.SeekInfo message, and
	// its starting block number
	Iterator(startType *ab.SeekPosition) (Iterator, uint64)
	// Height returns the number of blocks on the ledger
	Height() uint64
}

// Writer allows the caller to modify the ledger

// 写
type Writer interface {
	// Append a new block to the ledger
	Append(block *cb.Block) error
}

//go:generate mockery -dir . -name ReadWriter -case underscore  -output mocks/

// ReadWriter encapsulates the read/write functions of the ledger

// 读写 接口
type ReadWriter interface {
	Reader
	Writer
}
