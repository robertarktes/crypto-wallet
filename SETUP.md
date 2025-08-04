# Simple Crypto Wallet Setup

## Prerequisites

1. **Go 1.21 or higher**
   ```bash
   go version
   ```

2. **Internet access** for Ethereum blockchain interaction

## Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd crypto-wallet
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Build the project:**
   ```bash
   go build -o crypto-wallet
   ```

## Blockchain Configuration

### Option 1: Using Infura (Recommended)

1. Register at [Infura](https://infura.io/)
2. Create a new project
3. Get your Project ID
4. Update the URL in `cmd/main.go`:

```go
const (
    defaultBlockchainURL = "https://sepolia.infura.io/v3/YOUR-PROJECT-ID"
    defaultWalletFile    = "wallet.json"
)
```

### Option 2: Using Alchemy

1. Register at [Alchemy](https://www.alchemy.com/)
2. Create a new project
3. Get your API Key
4. Update the URL in `cmd/main.go`:

```go
const (
    defaultBlockchainURL = "https://eth-sepolia.g.alchemy.com/v2/YOUR-API-KEY"
    defaultWalletFile    = "wallet.json"
)
```

### Option 3: Using local node

If you have a local Ethereum node:

```go
const (
    defaultBlockchainURL = "http://localhost:8545"
    defaultWalletFile    = "wallet.json"
)
```

## Testing

1. **Run tests:**
   ```bash
   go test ./...
   ```

2. **Check build:**
   ```bash
   ./crypto-wallet help
   ```

3. **Create test wallet:**
   ```bash
   ./crypto-wallet generate
   ```

## Getting Test ETH

For testing in Sepolia network:

1. **Sepolia Faucet:**
   - Go to [sepoliafaucet.com](https://sepoliafaucet.com/)
   - Enter your wallet address
   - Get test ETH

2. **Alchemy Faucet:**
   - Go to [sepoliafaucet.com](https://sepoliafaucet.com/)
   - Connect wallet or enter address
   - Get test ETH

## Security

**IMPORTANT WARNINGS:**

1. **This wallet is intended for testing only**
2. **Do not use it for storing real funds**
3. **Private keys are stored unencrypted**
4. **Use only in test networks (Sepolia, Goerli)**

## Troubleshooting

### Blockchain connection error

```
Error creating wallet: blockchain connection error
```

**Solution:**
- Check blockchain URL correctness
- Ensure internet access
- Verify Project ID or API Key

### Wallet loading error

```
Error loading wallet: wallet file not found
```

**Solution:**
- Ensure `wallet.json` file exists
- Check file permissions
- Create new wallet with `generate` command

### Transaction sending error

```
Error sending transaction: insufficient funds
```

**Solution:**
- Ensure sufficient ETH on wallet
- Get test ETH through faucet
- Check transaction amount correctness

## Project Structure

```
crypto-wallet/
├── cmd/
│   └── main.go              # Main application file
├── internal/
│   ├── wallet/              # Wallet logic
│   │   ├── wallet.go
│   │   └── wallet_test.go
│   ├── blockchain/          # Blockchain interaction
│   │   ├── client.go
│   │   └── client_test.go
│   └── crypto/              # Cryptographic functions
│       ├── keys.go
│       └── keys_test.go
├── go.mod                   # Go dependencies
├── go.sum                   # Dependency hashes
├── README.md                # Main documentation
├── SETUP.md                 # This file
└── .gitignore              # Git exclusions
```

## Support

If you encounter issues:

1. Check the "Troubleshooting" section
2. Ensure all dependencies are installed
3. Check Go version (should be 1.21+)
4. Create an issue in the project repository 