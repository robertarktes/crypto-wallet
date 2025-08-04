# Simple Crypto Wallet Usage Examples

## Quick Start

### 1. Generate new wallet

```bash
./crypto-wallet generate
```

**Output:**
```
Generating new wallet...
New wallet created!
Address: 0xdDa43abc53563D138A26fD6a7Be0AD2E4Ef0f907
Data saved to file: wallet.json

IMPORTANT: Save private key in a secure location!
```

### 2. Display wallet address

```bash
./crypto-wallet address
```

**Output:**
```
Wallet address: 0xdDa43abc53563D138A26fD6a7Be0AD2E4Ef0f907
```

### 3. Check balance

```bash
./crypto-wallet balance
```

**Output (if funds available):**
```
Balance: 0.5 ETH
```

**Output (if no funds):**
```
Balance: 0.0 ETH
```

### 4. Send transaction

```bash
./crypto-wallet send 0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 0.001
```

**Output:**
```
Sending 0.001 ETH to address 0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6...
Transaction sent!
Transaction hash: 0x1234567890abcdef...
Check status: ./crypto-wallet status 0x1234567890abcdef...
```

### 5. Check transaction status

```bash
./crypto-wallet status 0x1234567890abcdef...
```

**Output (if transaction confirmed):**
```
Checking transaction status 0x1234567890abcdef...
Transaction confirmed!
Block number: 1234567
Gas used: 21000
```

**Output (if transaction not yet confirmed):**
```
Checking transaction status 0x1234567890abcdef...
Transaction not yet confirmed
```

## Usage with custom settings

### Using different wallet file

```bash
./crypto-wallet -wallet my-wallet.json generate
```

### Using different blockchain URL

```bash
./crypto-wallet -url https://mainnet.infura.io/v3/YOUR-PROJECT-ID balance
```

### Combining flags

```bash
./crypto-wallet -url https://sepolia.infura.io/v3/YOUR-PROJECT-ID -wallet test-wallet.json generate
```

## Complete working example

```bash
# 1. Create new wallet
./crypto-wallet generate

# 2. Check address
./crypto-wallet address

# 3. Get test ETH through faucet
# Go to https://sepoliafaucet.com/ and enter address

# 4. Check balance
./crypto-wallet balance

# 5. Send small amount to another address
./crypto-wallet send 0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 0.001

# 6. Check transaction status (use hash from previous command)
./crypto-wallet status 0x1234567890abcdef...

# 7. Check updated balance
./crypto-wallet balance
```

## Error examples and solutions

### Error: "wallet file not found"

```bash
./crypto-wallet address
# Error: error loading wallet: wallet file not found: wallet.json
```

**Solution:**
```bash
./crypto-wallet generate
```

### Error: "insufficient funds"

```bash
./crypto-wallet send 0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 1.0
# Error: error sending transaction: insufficient funds
```

**Solution:**
1. Get test ETH through faucet
2. Check balance: `./crypto-wallet balance`
3. Send smaller amount

### Error: "invalid recipient address"

```bash
./crypto-wallet send invalid-address 0.001
# Error: error sending transaction: invalid recipient address: invalid-address
```

**Solution:**
Use correct Ethereum address:
```bash
./crypto-wallet send 0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 0.001
```

### Error: "ETH amount must be positive"

```bash
./crypto-wallet send 0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 -0.001
# Error: error sending transaction: ETH amount must be positive
```

**Solution:**
Use positive number:
```bash
./crypto-wallet send 0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 0.001
```

## Wallet file structure

File `wallet.json` contains:

```json
{
  "private_key": "c714e9a509f46d217f6ccbd75d50d6a683b374aa5ada498a332c22f0a6d76877",
  "public_key": "0484415287010b987c29ab065b61cf51bf92f056581be65e945d9e43e8daf916d847a3f6b06417a8bb4eac5798d543410fef1402d0d391fe653bb657aabeadb51a",
  "address": "0x9a512Fd5A2bDB6e8Ba5d23bEd04dFDB7853Ad3B7"
}
```

**IMPORTANT:** This file contains private key unencrypted. Store it securely!

## Integration with other tools

### Using curl to check balance

```bash
# Get wallet address
ADDRESS=$(./crypto-wallet address | grep "Address:" | cut -d' ' -f3)

# Check balance via Etherscan API
curl "https://api-sepolia.etherscan.io/api?module=account&action=balance&address=$ADDRESS&tag=latest"
```

### Automation with scripts

```bash
#!/bin/bash
# Script for automatic balance checking

WALLET_ADDRESS=$(./crypto-wallet address | grep "Address:" | cut -d' ' -f3)
BALANCE=$(./crypto-wallet balance | grep "Balance:" | cut -d' ' -f2)

echo "Wallet: $WALLET_ADDRESS"
echo "Balance: $BALANCE ETH"

if [[ $BALANCE == "0.0" ]]; then
    echo "Balance empty! Get test ETH through faucet."
fi
```

## Security

1. **Never share private key**
2. **Use only in test networks**
3. **Do not store real funds in this wallet**
4. **Regularly backup wallet file**
5. **Use separate wallets for different purposes** 