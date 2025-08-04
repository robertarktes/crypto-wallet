package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type KeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    common.Address
}

func GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("error generating private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKeyECDSA,
		Address:    address,
	}, nil
}

func (kp *KeyPair) GetPrivateKeyHex() string {
	return hex.EncodeToString(crypto.FromECDSA(kp.PrivateKey))
}

func (kp *KeyPair) GetPublicKeyHex() string {
	return hex.EncodeToString(crypto.FromECDSAPub(kp.PublicKey))
}

func (kp *KeyPair) GetAddressHex() string {
	return kp.Address.Hex()
}

func (kp *KeyPair) SignMessage(message []byte) ([]byte, error) {
	hash := crypto.Keccak256Hash(message)
	signature, err := crypto.Sign(hash.Bytes(), kp.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error signing message: %w", err)
	}
	return signature, nil
}

func VerifySignature(message []byte, signature []byte, address common.Address) bool {
	hash := crypto.Keccak256Hash(message)
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		return false
	}

	var signerAddr common.Address
	copy(signerAddr[:], crypto.Keccak256(sigPublicKey[1:])[12:])

	return signerAddr == address
}

func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("error generating random bytes: %w", err)
	}
	return bytes, nil
}

func IsValidAddress(address string) bool {
	return common.IsHexAddress(address)
}

func HexToAddress(hex string) (common.Address, error) {
	if !common.IsHexAddress(hex) {
		return common.Address{}, fmt.Errorf("invalid Ethereum address: %s", hex)
	}
	return common.HexToAddress(hex), nil
}

func WeiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236)
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetInt(wei)
	return new(big.Float).Quo(fWei, big.NewFloat(1e18))
}

func EtherToWei(ether *big.Float) *big.Int {
	truncInt, _ := ether.Int(nil)
	truncInt.Mul(truncInt, big.NewInt(1e18))
	fracStr := ether.Text('f', 18)
	if len(fracStr) > 2 && fracStr[0:2] == "0." {
		fracStr = fracStr[2:]
		if len(fracStr) > 18 {
			fracStr = fracStr[:18]
		}
		if len(fracStr) < 18 {
			fracStr = fracStr + "000000000000000000"[:18-len(fracStr)]
		}
		fracInt, _ := new(big.Int).SetString(fracStr, 10)
		truncInt.Add(truncInt, fracInt)
	}
	return truncInt
}
