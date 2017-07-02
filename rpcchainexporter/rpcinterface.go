package rpcchainexporter

import (
	"github.com/decred/dcrd/dcrjson"
	"fmt"
	"os"
	"strings"
	"bytes"
	"encoding/json"
)

func executeRpcCommand(method string, params []interface{}) string {

	// Attempt to create the appropriate command using the arguments
	// provided by the user.
	cmd, err := dcrjson.NewCmd(method, params...)
	if err != nil {
		// Show the error along with its error code when it's a
		// dcrjson.Error as it reallistcally will always be since the
		// NewCmd function is only supposed to return errors of that
		// type.
		if jerr, ok := err.(dcrjson.Error); ok {
			fmt.Fprintf(os.Stderr, "%s command: %v (code: %s)\n",
				method, err, jerr.Code)
			os.Exit(1)
		}

		// The error is not a dcrjson.Error and this really should not
		// happen.  Nevertheless, fallback to just showing the error
		// if it should happen due to a bug in the package.
		fmt.Fprintf(os.Stderr, "%s command: %v\n", method, err)
		os.Exit(1)
	}

	// Marshal the command into a JSON-RPC byte slice in preparation for
	// sending it to the RPC server.
	marshalledJSON, err := dcrjson.MarshalCmd(1, cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Send the JSON-RPC request to the server using the user-specified
	// connection configuration.
	result, err := sendPostRequest(marshalledJSON, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}


	// Choose how to display the result based on its type.
	strResult := string(result)
	if strings.HasPrefix(strResult, "{") || strings.HasPrefix(strResult, "[") {
		var dst bytes.Buffer
		if err := json.Indent(&dst, result, "", "  "); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to format result: %v",
				err)
			os.Exit(1)
		}
		return dst.String()

	} else if strings.HasPrefix(strResult, `"`) {
		var str string
		if err := json.Unmarshal(result, &str); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to unmarshal result: %v",
				err)
			os.Exit(1)
		}
		return str

	} else if strResult != "null" {
		fmt.Println(strResult)
	}
	return ""
}


func getbestblock() *BestBlock {
	response := BestBlock{}
	jsonResult := executeRpcCommand("getbestblock", nil)
	json.Unmarshal([]byte(jsonResult), &response)
	return &response
}


func getblockhash(height int) string {
	result := executeRpcCommand("getblockhash", []interface{}{height})
	return result
}


func getblock(blockhash string) *Block {
	response := Block{}
	result := executeRpcCommand("getblock", []interface{}{blockhash, true})
	json.Unmarshal([]byte(result), &response)
	return &response
}


func getrawtransaction(txHash string) *RawTransaction {
	response := RawTransaction{}
	result := executeRpcCommand("getrawtransaction", []interface{}{txHash, 1})
	json.Unmarshal([]byte(result), &response)
	return &response
}

