package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"

	"crypto-wallet/internal/wallet"
)

const (
	defaultBlockchainURL = "https://sepolia.infura.io/v3/your-project-id"
	defaultWalletFile    = "wallet.json"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	var blockchainURL, walletFile string
	flag.StringVar(&blockchainURL, "url", defaultBlockchainURL, "Blockchain URL")
	flag.StringVar(&walletFile, "wallet", defaultWalletFile, "Wallet file")
	flag.Parse()

	w, err := wallet.NewWallet(blockchainURL, walletFile)
	if err != nil {
		fmt.Printf("Error creating wallet: %v\n", err)
		os.Exit(1)
	}
	defer w.Close()

	switch command {
	case "generate":
		err = handleGenerate(w)
	case "address":
		err = handleAddress(w)
	case "balance":
		err = handleBalance(w)
	case "send":
		err = handleSend(w)
	case "status":
		err = handleStatus(w)
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func handleGenerate(w *wallet.Wallet) error {
	fmt.Println("Generating new wallet...")

	err := w.GenerateNewWallet()
	if err != nil {
		return fmt.Errorf("error generating wallet: %w", err)
	}

	address, err := w.GetAddress()
	if err != nil {
		return fmt.Errorf("error getting address: %w", err)
	}

	fmt.Printf("New wallet created!\n")
	fmt.Printf("Address: %s\n", address)
	fmt.Printf("Data saved to file: %s\n", w.WalletFile)
	fmt.Println("\nIMPORTANT: Save private key in a secure location!")

	return nil
}

func handleAddress(w *wallet.Wallet) error {
	err := w.LoadWallet()
	if err != nil {
		return fmt.Errorf("error loading wallet: %w", err)
	}

	address, err := w.GetAddress()
	if err != nil {
		return fmt.Errorf("error getting address: %w", err)
	}

	fmt.Printf("Wallet address: %s\n", address)
	return nil
}

func handleBalance(w *wallet.Wallet) error {
	err := w.LoadWallet()
	if err != nil {
		return fmt.Errorf("error loading wallet: %w", err)
	}

	balance, err := w.GetBalance()
	if err != nil {
		return fmt.Errorf("error getting balance: %w", err)
	}

	fmt.Printf("Balance: %s ETH\n", balance.Text('f', 18))
	return nil
}

func handleSend(w *wallet.Wallet) error {
	if len(os.Args) < 4 {
		return fmt.Errorf("usage: send <recipient_address> <amount_in_eth>")
	}

	err := w.LoadWallet()
	if err != nil {
		return fmt.Errorf("error loading wallet: %w", err)
	}

	toAddress := os.Args[2]
	amountStr := os.Args[3]

	amount, ok := new(big.Float).SetString(amountStr)
	if !ok {
		return fmt.Errorf("invalid ETH amount: %s", amountStr)
	}

	if amount.Sign() <= 0 {
		return fmt.Errorf("ETH amount must be positive")
	}

	fmt.Printf("Sending %s ETH to address %s...\n", amountStr, toAddress)

	txHash, err := w.SendTransaction(toAddress, amount)
	if err != nil {
		return fmt.Errorf("error sending transaction: %w", err)
	}

	fmt.Printf("Transaction sent!\n")
	fmt.Printf("Transaction hash: %s\n", txHash)
	fmt.Printf("Check status: ./crypto-wallet status %s\n", txHash)

	return nil
}

func handleStatus(w *wallet.Wallet) error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: status <transaction_hash>")
	}

	txHash := os.Args[2]

	fmt.Printf("Checking transaction status %s...\n", txHash)

	receipt, err := w.GetTransactionStatus(txHash)
	if err != nil {
		return fmt.Errorf("error getting transaction status: %w", err)
	}

	if receipt == nil {
		fmt.Println("Transaction not yet confirmed")
		return nil
	}

	if receipt.Status == 1 {
		fmt.Println("Transaction confirmed!")
		fmt.Printf("Block number: %d\n", receipt.BlockNumber.Uint64())
		fmt.Printf("Gas used: %d\n", receipt.GasUsed)
	} else {
		fmt.Println("Transaction failed")
	}

	return nil
}

func printUsage() {
	fmt.Println("Simple Crypto Wallet - Ethereum cryptocurrency wallet")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  ./crypto-wallet <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  generate                    Generate new wallet")
	fmt.Println("  address                     Show wallet address")
	fmt.Println("  balance                     Show wallet balance")
	fmt.Println("  send <address> <amount>     Send ETH")
	fmt.Println("  status <hash>               Check transaction status")
	fmt.Println("  help                        Show this help")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -url <url>                  Blockchain URL (default: Sepolia)")
	fmt.Println("  -wallet <file>              Wallet file (default: wallet.json)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ./crypto-wallet generate")
	fmt.Println("  ./crypto-wallet address")
	fmt.Println("  ./crypto-wallet balance")
	fmt.Println("  ./crypto-wallet send 0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 0.001")
	fmt.Println("  ./crypto-wallet status 0x123...")
	fmt.Println()
	fmt.Println("IMPORTANT: This wallet is intended for testing only!")
	fmt.Println("  Do not use it for storing real funds.")
}
