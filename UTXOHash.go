package utils

import (
	"fmt"

	"github.com/libsv/go-bk/crypto"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-p2p/chaincfg/chainhash"
)

func UTXOHash(previousTxid *chainhash.Hash, index uint32, lockingScript []byte, satoshis uint64) (*chainhash.Hash, error) {
	if len(lockingScript) == 0 {
		return nil, fmt.Errorf("locking script is nil")
	}

	if satoshis == 0 {
		return nil, fmt.Errorf("satoshis is 0")
	}

	utxoHash := make([]byte, 0, 200)
	utxoHash = append(utxoHash, previousTxid.CloneBytes()...)
	utxoHash = append(utxoHash, bt.VarInt(index).Bytes()...)
	utxoHash = append(utxoHash, lockingScript...)
	utxoHash = append(utxoHash, bt.VarInt(satoshis).Bytes()...)

	hash := crypto.Sha256(utxoHash)
	chHash, err := chainhash.NewHash(hash)
	if err != nil {
		return nil, err
	}

	return chHash, nil
}
