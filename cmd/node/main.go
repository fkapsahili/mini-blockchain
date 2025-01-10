package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fkapsahili/mini-blockchain/internal/blockchain"
	"github.com/fkapsahili/mini-blockchain/internal/types"
)

var (
	dataDir string
	port    uint
	chain   *blockchain.Chain
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
	var err error
	chain, err := blockchain.NewChain(dataDir)
	if err != nil {
		fmt.Printf("Failed to initialize blockchain: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Blockchain initialized at height %d\n", chain.GetHeight())

	fmt.Println("Node is running. Press Ctrl+C to stop.")
	select {} // Keep running for now
}

func handleCreateBlock() {
	if chain == nil {
		fmt.Println("Blockchain not initialized. Please start the node first.")
		return
	}

	latestBlock, err := chain.GetLatestBlock()
	if err != nil {
		fmt.Printf("Failed to get latest block: %v\n", err)
		return
	}

	newBlock := &types.Block{
		Header: types.BlockHeader{
			Version:       1,
			PrevBlockHash: latestBlock.Hash,
			Timestamp:     time.Now(),
			Difficulty:    1, // Simple difficulty for now
			Nonce:         0, // Should be calculated with proper PoW
		},
		Transactions: []types.Transaction{}, // Empty transactions for now
		Height:       latestBlock.Height + 1,
	}

	// TODO: Calculate proper hash
	// For now, just use a dummy hash
	newBlock.Hash = [32]byte{} // This should be properly calculated

	if err := chain.AddBlock(newBlock); err != nil {
		fmt.Printf("Failed to add block: %v\n", err)
		return
	}

	fmt.Printf("Created new block at height %d\n", newBlock.Height)
}

func handleStatus() {
	if chain == nil {
		fmt.Println("Blockchain not initialized. Please start the node first.")
		return
	}

	height := chain.GetHeight()
	latestBlock, err := chain.GetLatestBlock()
	if err != nil {
		fmt.Printf("Failed to get latest block: %v\n", err)
		return
	}

	fmt.Printf("Current Height: %d\n", height)
	fmt.Printf("Latest Block Hash: %x\n", latestBlock.Hash)
	fmt.Printf("Latest Block Time: %v\n", latestBlock.Header.Timestamp)
	fmt.Printf("Number of Transactions: %d\n", len(latestBlock.Transactions))
}

func handleBlock(cmd *flag.FlagSet) {
	if chain == nil {
		fmt.Println("Blockchain not initialized. Please start the node first.")
		return
	}

	if cmd.NArg() < 1 {
		fmt.Println("Please provide block height")
		os.Exit(1)
	}

	height, err := strconv.ParseUint(cmd.Arg(0), 10, 64)
	if err != nil {
		fmt.Printf("Invalid height: %v\n", err)
		return
	}

	block, err := chain.GetBlock(height)
	if err != nil {
		fmt.Printf("Failed to get block: %v\n", err)
		return
	}

	fmt.Printf("Block Height: %d\n", block.Height)
	fmt.Printf("Block Hash: %x\n", block.Hash)
	fmt.Printf("Previous Block Hash: %x\n", block.Header.PrevBlockHash)
	fmt.Printf("Timestamp: %v\n", block.Header.Timestamp)
	fmt.Printf("Difficulty: %d\n", block.Header.Difficulty)
	fmt.Printf("Nonce: %d\n", block.Header.Nonce)
	fmt.Printf("Number of Transactions: %d\n", len(block.Transactions))
}
