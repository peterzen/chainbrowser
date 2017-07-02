package rpcchainexporter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	showHelpMessage = "Specify -h to show available options"
	listCmdMessage = "Specify -l to list available commands"
)


// usage displays the general usage when the help flag is not displayed and
// and an invalid command was specified.  The commandUsage function is used
// instead when a valid command was specified.
func usage(errorMessage string) {
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	fmt.Fprintln(os.Stderr, errorMessage)
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintf(os.Stderr, "  %s [OPTIONS] <command> <args...>\n\n",
		appName)
	fmt.Fprintln(os.Stderr, showHelpMessage)
	fmt.Fprintln(os.Stderr, listCmdMessage)
}

var (
	cfg               *config
)



func Export() {

	var err error
	var args []string

	cfg, args, err = loadConfig()
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("ARGS", args)

	createOutputFiles()
	defer closeOutputFiles()

	bestblock := getbestblock()

	for height := (bestblock.Height-4); height < bestblock.Height; height++ {

		blockhash := getblockhash(height)
		block := getblock(blockhash)
		fmt.Printf("Block #%d -> %s\n", block.Height, block.Hash)

		exportBlock(block)

		fmt.Printf("\n")
	}
}
