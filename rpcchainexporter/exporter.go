package rpcchainexporter

import (
	"os"
	"encoding/json"
	"fmt"
)

var (
	fileBlocks        *os.File
	fileTx            *os.File
	fileVin           *os.File
	fileVout          *os.File
)

type BestBlock struct {
	Hash   string     `json:"hash"`
	Height int        `json:"height"`
}

type Block struct {
	Time              uint64          `json:"time"`
	Freshstake        uint32          `json:"freshstake"`
	Nonce             uint64          `json:"nonce"`
	Sbits             float64         `json:"sbits"`
	Extradata         string          `json:"extradata"`
	Votebits          uint32          `json:"votebits"`
	Confirmations     uint32          `json:"confirmations"`
	Height            uint32          `json:"height"`
	Size              uint32          `json:"size"`
	Merkleroot        string          `json:"merkleroot"`
	Bits              string          `json:"bits"`
	Poolsize          uint32          `json:"poolsize"`
	Finalstate        string          `json:"finalstate"`
	Revocations       uint32          `json:"revocations"`
	Stakeroot         string          `json:"stakeroot"`
	Difficulty        float64         `json:"difficulty"`
	Stx               []string        `json:"stx"`
	Tx                []string        `json:"tx"`
	Hash              string          `json:"hash"`
	Stakeversion      uint32          `json:"stakeversion"`
	Previousblockhash string          `json:"previousblockhash"`
	Voters            uint32          `json:"voters"`
	Version           uint32          `json:"version"`
}

type VinScriptSig struct {
	Hex string        `json:"hex"`
	Asm string        `json:"asm"`
}

type VoutScriptPubKey struct {
	Asm       string        `json:"asm"`
	Hex       string        `json:"hex"`
	Type      string        `json:"type"`
	ReqSigs   int32         `json:"reqSigs"`
	Addresses []string      `json:"addresses"`
}

type Vin struct {
	Txid        string        `json:"txid"`
	Vout        uint64        `json:"vout"`
	Tree        uint64        `json:"tree"`
	Sequence    uint64        `json:"sequence"`
	Amountin    float64       `json:"amountin"`
	Blockheight uint64        `json:"blockheight"`
	Blockindex  uint64        `json:"blockindex"`
	ScriptSig   VinScriptSig  `json:"scriptSig"`
}

type Vout struct {
	Value        float64          `json:"value"`
	N            int32            `json:"n"`
	Version      int32            `json:"version"`
	ScriptPubKey VoutScriptPubKey `json:"scriptPubKey"`
}

type RawTransaction struct {
	Hex           string  `json:"hex"`
	Txid          string  `json:"txid"`
	Version       uint64  `json:"version"`
	Locktime      uint64  `json:"locktime"`
	Expiry        uint64  `json:"expiry"`
	Blockhash     string  `json:"blockhash"`
	Blockheight   uint64  `json:"blockheight"`
	Confirmations uint64  `json:"confirmations"`
	Time          uint64  `json:"time"`
	Blocktime     uint64  `json:"blocktime"`
	Vin           []Vin   `json:"vin"`
	Vout          []Vout  `json:"vout"`
}

func createOutputFiles() {
	fileTx, _ = os.Create("tx.json")
	fileVin, _ = os.Create("vin.json")
	fileVout, _ = os.Create("vout.json")
	fileBlocks, _ = os.Create("blocks.json")
}

func closeOutputFiles() {
	fileVin.Close()
	fileVout.Close()
	fileBlocks.Close()
	fileTx.Close()
}

func exportBlock(block *Block){
	jsonValue, _ := json.Marshal(block)
	fileBlocks.WriteString(string(jsonValue))
	fileBlocks.WriteString("\n")

	for _, txHash := range block.Tx {

		fmt.Printf("  Tx %s\n", txHash)
		rawTx := getrawtransaction(txHash)

		exportRawTransaction(rawTx)

		for _, vin := range rawTx.Vin {
			//fmt.Printf("    vin.txid:      %s\n", vin.Txid)
			//fmt.Printf("    vin.Amountin:  %f\n", vin.Amountin)
			exportVin(&vin)
		}
		for _, vout := range rawTx.Vout {
			//fmt.Printf("    vout.Amountin: %f\n", vout.Value)
			exportVout(&vout)
		}
		//fmt.Printf("\n")
	}
}

func exportRawTransaction(rawTx *RawTransaction){
	jsonValue, _ := json.Marshal(rawTx)
	fileTx.WriteString(string(jsonValue))
	fileTx.WriteString("\n")
}

func exportVin(vin *Vin){
	jsonValue, _ := json.Marshal(vin)
	fileVin.WriteString(string(jsonValue))
	fileVin.WriteString("\n")
}

func exportVout(vout *Vout){
	jsonValue, _ := json.Marshal(vout)
	fileVout.WriteString(string(jsonValue))
	fileVout.WriteString("\n")
}
