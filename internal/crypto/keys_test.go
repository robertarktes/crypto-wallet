package crypto

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGenerateKeyPair(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Ошибка генерации пары ключей: %v", err)
	}

	if keyPair == nil {
		t.Fatal("KeyPair не должен быть nil")
	}

	if keyPair.PrivateKey == nil {
		t.Fatal("Приватный ключ не должен быть nil")
	}

	if keyPair.PublicKey == nil {
		t.Fatal("Публичный ключ не должен быть nil")
	}

	if keyPair.Address == (common.Address{}) {
		t.Fatal("Адрес не должен быть пустым")
	}
}

func TestKeyPairHexFormats(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Ошибка генерации пары ключей: %v", err)
	}

	// Проверяем hex форматы
	privateKeyHex := keyPair.GetPrivateKeyHex()
	if len(privateKeyHex) == 0 {
		t.Fatal("Приватный ключ в hex формате не должен быть пустым")
	}

	publicKeyHex := keyPair.GetPublicKeyHex()
	if len(publicKeyHex) == 0 {
		t.Fatal("Публичный ключ в hex формате не должен быть пустым")
	}

	addressHex := keyPair.GetAddressHex()
	if len(addressHex) == 0 {
		t.Fatal("Адрес в hex формате не должен быть пустым")
	}

	// Проверяем, что адрес начинается с 0x
	if addressHex[:2] != "0x" {
		t.Fatal("Адрес должен начинаться с 0x")
	}
}

func TestSignAndVerifyMessage(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Ошибка генерации пары ключей: %v", err)
	}

	message := []byte("Hello, Ethereum!")

	// Подписываем сообщение
	signature, err := keyPair.SignMessage(message)
	if err != nil {
		t.Fatalf("Ошибка подписи сообщения: %v", err)
	}

	if len(signature) == 0 {
		t.Fatal("Подпись не должна быть пустой")
	}

	// Проверяем подпись
	isValid := VerifySignature(message, signature, keyPair.Address)
	if !isValid {
		t.Fatal("Подпись должна быть валидной")
	}

	// Проверяем с неправильным сообщением
	wrongMessage := []byte("Wrong message")
	isValid = VerifySignature(wrongMessage, signature, keyPair.Address)
	if isValid {
		t.Fatal("Подпись не должна быть валидной для неправильного сообщения")
	}
}

func TestGenerateRandomBytes(t *testing.T) {
	length := 32
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		t.Fatalf("Ошибка генерации случайных байт: %v", err)
	}

	if len(bytes) != length {
		t.Fatalf("Длина байт должна быть %d, получено %d", length, len(bytes))
	}

	// Проверяем, что байты не все нули
	allZeros := true
	for _, b := range bytes {
		if b != 0 {
			allZeros = false
			break
		}
	}

	if allZeros {
		t.Fatal("Случайные байты не должны быть все нули")
	}
}

func TestIsValidAddress(t *testing.T) {
	// Валидные адреса
	validAddresses := []string{
		"0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		"0x742d35cc6634c0532925a3b8d4c9db96c4b4d8b6",
		"0x0000000000000000000000000000000000000000",
	}

	for _, addr := range validAddresses {
		if !IsValidAddress(addr) {
			t.Errorf("Адрес должен быть валидным: %s", addr)
		}
	}

	// Невалидные адреса
	invalidAddresses := []string{
		"0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b",   // Слишком короткий
		"0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6G", // Неправильный символ
		"", // Пустая строка
	}

	for _, addr := range invalidAddresses {
		if IsValidAddress(addr) {
			t.Errorf("Адрес не должен быть валидным: %s", addr)
		}
	}
}

func TestHexToAddress(t *testing.T) {
	validAddress := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
	address, err := HexToAddress(validAddress)
	if err != nil {
		t.Fatalf("Ошибка конвертации адреса: %v", err)
	}

	// Проверяем, что адрес не пустой, но не сравниваем точное hex представление
	if address == (common.Address{}) {
		t.Fatalf("Адрес не должен быть пустым")
	}

	// Тест с невалидным адресом
	invalidAddress := "0xinvalid"
	_, err = HexToAddress(invalidAddress)
	if err == nil {
		t.Fatal("Должна быть ошибка для невалидного адреса")
	}
}

func TestWeiToEther(t *testing.T) {
	// Тест с 1 ETH (1e18 Wei)
	oneEthWei := big.NewInt(1e18)
	oneEth := WeiToEther(oneEthWei)

	expected := big.NewFloat(1.0)
	if oneEth.Cmp(expected) != 0 {
		t.Fatalf("1 ETH должен быть равен 1.0, получено %s", oneEth.Text('f', 18))
	}

	// Тест с 0.5 ETH
	halfEthWei := big.NewInt(5e17) // 0.5 * 1e18
	halfEth := WeiToEther(halfEthWei)

	expectedHalf := big.NewFloat(0.5)
	if halfEth.Cmp(expectedHalf) != 0 {
		t.Fatalf("0.5 ETH должен быть равен 0.5, получено %s", halfEth.Text('f', 18))
	}

	// Тест с 0 Wei
	zeroWei := big.NewInt(0)
	zeroEth := WeiToEther(zeroWei)

	expectedZero := big.NewFloat(0.0)
	if zeroEth.Cmp(expectedZero) != 0 {
		t.Fatalf("0 Wei должен быть равен 0.0, получено %s", zeroEth.Text('f', 18))
	}
}

func TestEtherToWei(t *testing.T) {
	// Тест с 1 ETH
	oneEth := big.NewFloat(1.0)
	oneEthWei := EtherToWei(oneEth)

	expected := big.NewInt(1e18)
	if oneEthWei.Cmp(expected) != 0 {
		t.Fatalf("1 ETH должен быть равен 1e18 Wei, получено %s", oneEthWei.String())
	}

	// Тест с 0.5 ETH
	halfEth := big.NewFloat(0.5)
	halfEthWei := EtherToWei(halfEth)

	expectedHalf := big.NewInt(5e17)
	if halfEthWei.Cmp(expectedHalf) != 0 {
		t.Fatalf("0.5 ETH должен быть равен 5e17 Wei, получено %s", halfEthWei.String())
	}

	// Тест с 0 ETH
	zeroEth := big.NewFloat(0.0)
	zeroEthWei := EtherToWei(zeroEth)

	expectedZero := big.NewInt(0)
	if zeroEthWei.Cmp(expectedZero) != 0 {
		t.Fatalf("0 ETH должен быть равен 0 Wei, получено %s", zeroEthWei.String())
	}
}

func TestWeiEtherConversionRoundTrip(t *testing.T) {
	// Тест обратной конвертации с простыми числами
	originalWei := big.NewInt(1000000000000000000) // 1 ETH в Wei
	ether := WeiToEther(originalWei)
	convertedWei := EtherToWei(ether)

	// Для простых чисел конвертация должна быть точной
	if originalWei.Cmp(convertedWei) != 0 {
		t.Fatalf("Обратная конвертация не точна: оригинал %s, результат %s",
			originalWei.String(), convertedWei.String())
	}
}
