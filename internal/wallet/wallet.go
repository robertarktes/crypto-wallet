package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"crypto-wallet/internal/blockchain"
	"crypto-wallet/internal/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethereumCrypto "github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	KeyPair    *crypto.KeyPair
	Blockchain *blockchain.Client
	WalletFile string
}

type WalletData struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Address    string `json:"address"`
}

func NewWallet(blockchainURL string, walletFile string) (*Wallet, error) {
	client, err := blockchain.NewClient(blockchainURL)
	if err != nil {
		return nil, fmt.Errorf("error creating blockchain client: %w", err)
	}

	return &Wallet{
		Blockchain: client,
		WalletFile: walletFile,
	}, nil
}

func (w *Wallet) GenerateNewWallet() error {
	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("error generating keys: %w", err)
	}

	w.KeyPair = keyPair

	err = w.SaveWallet()
	if err != nil {
		return fmt.Errorf("error saving wallet: %w", err)
	}

	return nil
}

func (w *Wallet) LoadWallet() error {
	if _, err := os.Stat(w.WalletFile); os.IsNotExist(err) {
		return fmt.Errorf("wallet file not found: %s", w.WalletFile)
	}

	data, err := os.ReadFile(w.WalletFile)
	if err != nil {
		return fmt.Errorf("error reading wallet file: %w", err)
	}

	var walletData WalletData
	err = json.Unmarshal(data, &walletData)
	if err != nil {
		return fmt.Errorf("error parsing wallet data: %w", err)
	}

	keyPair, err := w.restoreKeyPair(walletData)
	if err != nil {
		return fmt.Errorf("error restoring keys: %w", err)
	}

	w.KeyPair = keyPair
	return nil
}

func (w *Wallet) SaveWallet() error {
	if w.KeyPair == nil {
		return fmt.Errorf("wallet not initialized")
	}

	walletData := WalletData{
		PrivateKey: w.KeyPair.GetPrivateKeyHex(),
		PublicKey:  w.KeyPair.GetPublicKeyHex(),
		Address:    w.KeyPair.GetAddressHex(),
	}

	data, err := json.MarshalIndent(walletData, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing wallet data: %w", err)
	}

	dir := filepath.Dir(w.WalletFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	err = os.WriteFile(w.WalletFile, data, 0600)
	if err != nil {
		return fmt.Errorf("error writing wallet file: %w", err)
	}

	return nil
}

func (w *Wallet) GetAddress() (string, error) {
	if w.KeyPair == nil {
		return "", fmt.Errorf("wallet not initialized")
	}

	return w.KeyPair.GetAddressHex(), nil
}

func (w *Wallet) GetBalance() (*big.Float, error) {
	if w.KeyPair == nil {
		return nil, fmt.Errorf("wallet not initialized")
	}

	balance, err := w.Blockchain.GetBalanceInEther(w.KeyPair.Address)
	if err != nil {
		return nil, fmt.Errorf("error getting balance: %w", err)
	}

	return balance, nil
}

func (w *Wallet) SendTransaction(toAddress string, amount *big.Float) (string, error) {
	if w.KeyPair == nil {
		return "", fmt.Errorf("wallet not initialized")
	}

	if !crypto.IsValidAddress(toAddress) {
		return "", fmt.Errorf("invalid recipient address: %s", toAddress)
	}

	toAddr := common.HexToAddress(toAddress)

	amountWei := crypto.EtherToWei(amount)

	nonce, err := w.Blockchain.GetNonce(w.KeyPair.Address)
	if err != nil {
		return "", fmt.Errorf("error getting nonce: %w", err)
	}

	gasPrice, err := w.Blockchain.GetGasPrice()
	if err != nil {
		return "", fmt.Errorf("error getting gas price: %w", err)
	}

	gasLimit, err := w.Blockchain.EstimateGas(w.KeyPair.Address, &toAddr, amountWei, nil)
	if err != nil {
		gasLimit = 21000
	}

	tx := w.Blockchain.CreateTransaction(
		w.KeyPair.Address,
		toAddr,
		amountWei,
		gasLimit,
		gasPrice,
		nonce,
		nil,
	)

	privateKeyInt := new(big.Int)
	privateKeyInt.SetString(w.KeyPair.GetPrivateKeyHex(), 16)

	signedTx, err := w.Blockchain.SignTransaction(tx, privateKeyInt)
	if err != nil {
		return "", fmt.Errorf("error signing transaction: %w", err)
	}

	err = w.Blockchain.SendTransaction(signedTx)
	if err != nil {
		return "", fmt.Errorf("error sending transaction: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}

func (w *Wallet) GetTransactionStatus(txHash string) (*types.Receipt, error) {
	hash := common.HexToHash(txHash)
	receipt, err := w.Blockchain.GetTransactionReceipt(hash)
	if err != nil {
		return nil, fmt.Errorf("error getting transaction status: %w", err)
	}

	return receipt, nil
}

func (w *Wallet) WaitForTransaction(txHash string, maxAttempts int) (*types.Receipt, error) {
	hash := common.HexToHash(txHash)
	receipt, err := w.Blockchain.WaitForTransaction(hash, maxAttempts)
	if err != nil {
		return nil, fmt.Errorf("error waiting for transaction confirmation: %w", err)
	}

	return receipt, nil
}

func (w *Wallet) SignMessage(message []byte) ([]byte, error) {
	if w.KeyPair == nil {
		return nil, fmt.Errorf("wallet not initialized")
	}

	signature, err := w.KeyPair.SignMessage(message)
	if err != nil {
		return nil, fmt.Errorf("error signing message: %w", err)
	}

	return signature, nil
}

func (w *Wallet) VerifyMessage(message []byte, signature []byte, address string) bool {
	if !crypto.IsValidAddress(address) {
		return false
	}

	addr := common.HexToAddress(address)
	return crypto.VerifySignature(message, signature, addr)
}

func (w *Wallet) Close() {
	if w.Blockchain != nil {
		w.Blockchain.Close()
	}
}

func (w *Wallet) restoreKeyPair(walletData WalletData) (*crypto.KeyPair, error) {
	privateKeyBytes, err := hex.DecodeString(walletData.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error decoding private key: %w", err)
	}

	privateKey, err := ethereumCrypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error restoring private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	address := ethereumCrypto.PubkeyToAddress(*publicKeyECDSA)

	if strings.ToLower(address.Hex()) != strings.ToLower(walletData.Address) {
		return nil, fmt.Errorf("address mismatch: expected %s, got %s", walletData.Address, address.Hex())
	}

	return &crypto.KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKeyECDSA,
		Address:    address,
	}, nil
}
