package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	client *ethclient.Client
	url    string
}

func NewClient(url string) (*Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("error connecting to blockchain: %w", err)
	}

	return &Client{
		client: client,
		url:    url,
	}, nil
}

func (c *Client) GetBalance(address common.Address) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	balance, err := c.client.BalanceAt(ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting balance: %w", err)
	}

	return balance, nil
}

func (c *Client) GetBalanceInEther(address common.Address) (*big.Float, error) {
	balance, err := c.GetBalance(address)
	if err != nil {
		return nil, err
	}

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(1e18))

	return ethValue, nil
}

func (c *Client) GetGasPrice() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting gas price: %w", err)
	}

	return gasPrice, nil
}

func (c *Client) GetNonce(address common.Address) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	nonce, err := c.client.PendingNonceAt(ctx, address)
	if err != nil {
		return 0, fmt.Errorf("error getting nonce: %w", err)
	}

	return nonce, nil
}

func (c *Client) GetNetworkID() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	networkID, err := c.client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting network ID: %w", err)
	}

	return networkID, nil
}

func (c *Client) SendTransaction(tx *types.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := c.client.SendTransaction(ctx, tx)
	if err != nil {
		return fmt.Errorf("error sending transaction: %w", err)
	}

	return nil
}

func (c *Client) GetTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	receipt, err := c.client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("error getting transaction receipt: %w", err)
	}

	return receipt, nil
}

func (c *Client) WaitForTransaction(txHash common.Hash, maxAttempts int) (*types.Receipt, error) {
	for i := 0; i < maxAttempts; i++ {
		receipt, err := c.GetTransactionReceipt(txHash)
		if err == nil && receipt != nil {
			return receipt, nil
		}

		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("transaction not confirmed after %d attempts", maxAttempts)
}

func (c *Client) CreateTransaction(
	from common.Address,
	to common.Address,
	value *big.Int,
	gasLimit uint64,
	gasPrice *big.Int,
	nonce uint64,
	data []byte,
) *types.Transaction {
	return types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)
}

func (c *Client) SignTransaction(tx *types.Transaction, privateKey *big.Int) (*types.Transaction, error) {
	networkID, err := c.GetNetworkID()
	if err != nil {
		return nil, fmt.Errorf("error getting network ID for signing: %w", err)
	}

	privateKeyBytes := privateKey.Bytes()
	if len(privateKeyBytes) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(privateKeyBytes):], privateKeyBytes)
		privateKeyBytes = padded
	}

	ecdsaPrivateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error converting private key: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(networkID), ecdsaPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error signing transaction: %w", err)
	}

	return signedTx, nil
}

func (c *Client) EstimateGas(from common.Address, to *common.Address, value *big.Int, data []byte) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msg := ethereum.CallMsg{
		From:  from,
		To:    to,
		Value: value,
		Data:  data,
	}

	gasLimit, err := c.client.EstimateGas(ctx, msg)
	if err != nil {
		return 0, fmt.Errorf("error estimating gas: %w", err)
	}

	return gasLimit, nil
}

func (c *Client) Close() {
	if c.client != nil {
		c.client.Close()
	}
}
