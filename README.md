# Simple Crypto Wallet

A cryptocurrency wallet for Ethereum written in Go using the go-ethereum library.

## Features

- Generate new key pairs (private and public keys)
- Display wallet address
- Sign and send transactions to Ethereum test network (Sepolia)
- Check wallet balance via API
- Tests for key functions

## Requirements

- Go 1.21 or higher
- Internet access for blockchain interaction

## Installation

1. Clone the repository:
```bash
git clone https://github.com/robertarktes/crypto-wallet
cd crypto-wallet
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the project:
```bash
go build -o crypto-wallet
```

## Usage

### Generate new wallet

```bash
./crypto-wallet generate
```

This command will create a new key pair and save it to `wallet.json`.

### Display wallet address

```bash
./crypto-wallet address
```

### Check balance

```bash
./crypto-wallet balance <address>
```

### Send transaction

```bash
./crypto-wallet send <to_address> <amount_in_eth>
```

Example:
```bash
./crypto-wallet send 0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 0.001
```

### Get test ETH

For testing in Sepolia network, you can get test ETH through:
- [Sepolia Faucet](https://sepoliafaucet.com/)
- [Alchemy Sepolia Faucet](https://sepoliafaucet.com/)

## Project Structure

```
crypto-wallet/
├── cmd/
│   └── main.go          # Main application file
├── internal/
│   ├── wallet/          # Wallet logic
│   │   ├── wallet.go
│   │   └── wallet_test.go
│   ├── blockchain/      # Blockchain interaction
│   │   ├── client.go
│   │   └── client_test.go
│   └── crypto/          # Cryptographic functions
│       ├── keys.go
│       └── keys_test.go
├── go.mod
├── go.sum
└── README.md
```

## Security

**Important**: This wallet is intended for educational purposes and testing only. Do not use it for storing real funds.

- Private keys are stored unencrypted
- No additional security measures
- Use only in test networks

## Testing

Run tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## License

MIT License

## Contributing

Pull requests and issues for project improvements are welcome.