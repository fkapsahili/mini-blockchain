package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fkapsahili/mini-blockchain/internal/blockchain"
)

var (
	dataDir string
	port    uint
)

func main() {
	flag.StringVar(&dataDir, "datadir", "./data", "Data directory for blockchain")
	flag.UintVar(&port, "port", 8333, "Port for P2P communication")

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	createBlockCmd := flag.NewFlagSet("createBlock", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	blockCmd := flag.NewFlagSet("block", flag.ExitOnError)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])
		handleStart()
	case "createblock":
		createBlockCmd.Parse(os.Args[2:])
		handleCreateBlock()
	case "status":
		statusCmd.Parse(os.Args[2:])
		handleStatus()
	case "block":
		blockCmd.Parse(os.Args[2:])
		handleBlock(blockCmd)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  node [--datadir <dir>] [--port <port>] <command>")
	fmt.Println("\nCommands:")
	fmt.Println("  start        Start the blockchain node")
	fmt.Println("  createblock  Create a new block")
	fmt.Println("  status       Show blockchain status")
	fmt.Println("  block        Show block information")
}

func handleStart() {
	chain, err := blockchain.NewChain(dataDir)
	if err != nil {
		fmt.Printf("Failed to initialize blockchain: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Blockchain initialized at height %d\n", chain.GetHeight())
}

func handleCreateBlock() {
	fmt.Println("Creating new block...")
}

func handleStatus() {
	fmt.Println("Blockchain status:")
}

func handleBlock(cmd *flag.FlagSet) {
	if cmd.NArg() < 1 {
		fmt.Println("Please provide block height")
		os.Exit(1)
	}
	height := cmd.Arg(0)
	fmt.Printf("Showing block info for height: %s\n", height)
}
