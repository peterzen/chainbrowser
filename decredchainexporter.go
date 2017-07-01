package chainbrowser

import (
	"fmt"
	"encoding/json"
	"github.com/decred/dcrd/blockchain"
	"github.com/decred/dcrd/database"
	"github.com/decred/dcrd/wire"
	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/txscript"
	"github.com/decred/dcrutil"
	_ "github.com/decred/dcrd/database/ffldb"
	"time"
	"os"
)

var (
	fileBlocks        *os.File
	fileTx            *os.File
	fileTxin          *os.File
	fileTxout         *os.File
)

func marshalBlock(block *wire.MsgBlock) {

	b := block.Header
	blockHash := b.BlockHash().String()

	jValue, _ := json.Marshal(&struct {
		Hash         string
		Version      int32
		PrevBlock    string
		MerkleRoot   string
		StakeRoot    string
		VoteBits     uint16
		//FinalState string
		Voters       uint16
		FreshStake   uint8
		Revocations  uint8
		PoolSize     uint32
		Bits         uint32
		SBits        int64
		Height       uint32
		Size         uint32
		Timestamp    time.Time
		Nonce        uint32
		//ExtraData string
		StakeVersion uint32
	}{
		Hash:           blockHash,
		Version:        b.Version,
		PrevBlock:      b.PrevBlock.String(),
		MerkleRoot:     b.MerkleRoot.String(),
		StakeRoot:      b.StakeRoot.String(),
		VoteBits:       b.VoteBits,
		//FinalState [6]byte
		Voters:         b.Voters,
		FreshStake:     b.FreshStake,
		Revocations:    b.Revocations,
		PoolSize:       b.PoolSize,
		Bits:           b.Bits,
		SBits:          b.SBits,
		Height:         b.Height,
		Size:           b.Size,
		Timestamp:      b.Timestamp,
		Nonce:          b.Nonce,
		//ExtraData [32]byte
		StakeVersion:   b.StakeVersion,
	})
	fileBlocks.WriteString(string(jValue))
	fileBlocks.WriteString("\n")

	marshalTransactions(block.Transactions, blockHash)
}

func marshalTransactions(txs []*wire.MsgTx, blockHash string) {

	for _, tx := range txs {
		txFormatted := &struct {
			BlockHash string
			Hash      string
			Version   int32
			LockTime  uint32
			Expiry    uint32
		}{
			BlockHash:      blockHash,
			Hash:           tx.TxHash().String(),
			Version:        tx.Version,
			LockTime:       tx.LockTime,
			Expiry:         tx.Expiry,
		}
		jValue, _ := json.Marshal(txFormatted)
		fileTx.WriteString(string(jValue))
		fileTx.WriteString("\n")
		marshalTxIn(tx.TxIn, tx.TxHash().String())
		marshalTxOut(tx.TxOut, tx.TxHash().String())
	}
}

func marshalTxIn(txIn []*wire.TxIn, parentTxHash string) {

	for _, tx := range txIn {
		signatureScriptDisasm, _ := txscript.DisasmString(tx.SignatureScript)
		txFormatted := &struct {
			ParentTxHash string
			// Non-witness
			Hash         string
			Index        uint32
			Tree         int8
			Sequence     uint32
			// Witness
			ValueIn      int64
			BlockHeight  uint32
			BlockIndex   uint32
			SignatureScript	string
		}{
			ParentTxHash:      parentTxHash,
			Hash:              tx.PreviousOutPoint.Hash.String(),
			Index:             tx.PreviousOutPoint.Index,
			Tree:              tx.PreviousOutPoint.Tree,
			Sequence:          tx.Sequence,
			// Witness
			ValueIn:           tx.ValueIn,
			BlockHeight:       tx.BlockHeight,
			BlockIndex:        tx.BlockIndex,
			SignatureScript:   signatureScriptDisasm,
		}
		jValue, _ := json.Marshal(txFormatted)
		fileTxin.WriteString(string(jValue))
		fileTxin.WriteString("\n")
	}
}

func marshalTxOut(txOut []*wire.TxOut, parentTxHash string) {

	for _, vout := range txOut {

		pkScriptDisasm, _ := txscript.DisasmString(vout.PkScript)

		sc, addrs, _, _ := txscript.ExtractPkScriptAddrs(
			vout.Version, vout.PkScript, &chaincfg.MainNetParams)

		fmt.Println("SC", sc)
		fmt.Println("ADDR", addrs)

		txFormatted := &struct {
			ParentTxHash string
			Value        int64
			Version      uint16
			PkScript     string
		}{
			ParentTxHash:      parentTxHash,
			Value:             vout.Value,
			Version:           vout.Version,
			PkScript:	   pkScriptDisasm,
		}
		jValue, _ := json.Marshal(txFormatted)
		fileTxout.WriteString(string(jValue))
		fileTxout.WriteString("\n")
	}
}

func Export() {

	fileTx, _ = os.Create("tx.json")
	defer fileTx.Close()

	fileTxin, _ = os.Create("txin.json")
	defer fileTxin.Close()

	fileTxout, _ = os.Create("txout.json")
	defer fileTxout.Close()

	fileBlocks, _ = os.Create("blocks.json")
	defer fileBlocks.Close()

	db, err := database.Open("ffldb", "/Users/peter/Library/Application Support/Dcrd/data/mainnet/blocks_ffldb/", wire.MainNet)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	var block *dcrutil.Block

	var height int64

	for height = 147100; height < 147133; height++ {

		fmt.Printf("Height #%d\n", height)

		err = db.View(func(tx database.Tx) error {
			block, _ = blockchain.DBFetchBlockByHeight(tx, height)

			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			fmt.Println("ERROR", err)
			return
		}

		block := block.MsgBlock()

		if err != nil {
			fmt.Println("FromBytes", err)
			return
		}

		marshalBlock(block)
	}
}
