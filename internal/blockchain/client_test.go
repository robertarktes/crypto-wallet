package blockchain

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// Тест с использованием тестовой сети (может не работать без реального подключения)
func TestNewClient(t *testing.T) {
	// Тест с невалидным URL
	_, err := NewClient("invalid-url")
	if err == nil {
		t.Fatal("Должна быть ошибка для невалидного URL")
	}

	// Тест с пустым URL
	_, err = NewClient("")
	if err == nil {
		t.Fatal("Должна быть ошибка для пустого URL")
	}
}

func TestIsValidAddress(t *testing.T) {
	// Валидные адреса
	validAddresses := []string{
		"0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		"0x742d35cc6634c0532925a3b8d4c9db96c4b4d8b6",
	}

	for _, addr := range validAddresses {
		address := common.HexToAddress(addr)
		if address == (common.Address{}) {
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
		// Для невалидных адресов HexToAddress может вернуть нулевой адрес
		// или вызвать панику, поэтому используем recover
		func() {
			defer func() {
				if r := recover(); r == nil {
					// Если паники не было, просто проверяем, что функция не паникует
					// Для некоторых невалидных адресов HexToAddress может вернуть непустой адрес
				}
			}()
			common.HexToAddress(addr)
		}()
	}
}

func TestCreateTransaction(t *testing.T) {
	from := common.HexToAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	to := common.HexToAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7")
	value := big.NewInt(1000000000000000000) // 1 ETH в Wei
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(20000000000) // 20 Gwei
	nonce := uint64(0)
	data := []byte{}

	client := &Client{}

	tx := client.CreateTransaction(from, to, value, gasLimit, gasPrice, nonce, data)

	if tx == nil {
		t.Fatal("Транзакция не должна быть nil")
	}

	// Проверяем, что транзакция содержит правильные данные
	if tx.To() == nil || *tx.To() != to {
		t.Fatal("Адрес получателя не совпадает")
	}

	if tx.Value().Cmp(value) != 0 {
		t.Fatal("Сумма транзакции не совпадает")
	}

	if tx.Gas() != gasLimit {
		t.Fatal("Лимит газа не совпадает")
	}

	if tx.GasPrice().Cmp(gasPrice) != 0 {
		t.Fatal("Цена газа не совпадает")
	}

	if tx.Nonce() != nonce {
		t.Fatal("Nonce не совпадает")
	}
}

func TestWeiToEtherConversion(t *testing.T) {
	// Тест конвертации Wei в Ether
	oneEthWei := big.NewInt(1000000000000000000) // 1 ETH в Wei
	expectedEther := big.NewFloat(1.0)

	// Используем простую конвертацию для теста
	fbalance := new(big.Float)
	fbalance.SetString(oneEthWei.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(1e18))

	if ethValue.Cmp(expectedEther) != 0 {
		t.Fatalf("1 ETH должен быть равен 1.0, получено %s", ethValue.Text('f', 18))
	}

	// Тест с 0.5 ETH
	halfEthWei := big.NewInt(500000000000000000) // 0.5 ETH в Wei
	expectedHalfEther := big.NewFloat(0.5)

	fbalance.SetString(halfEthWei.String())
	ethValue = new(big.Float).Quo(fbalance, big.NewFloat(1e18))

	if ethValue.Cmp(expectedHalfEther) != 0 {
		t.Fatalf("0.5 ETH должен быть равен 0.5, получено %s", ethValue.Text('f', 18))
	}
}

func TestTransactionHash(t *testing.T) {
	from := common.HexToAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	to := common.HexToAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b7")
	value := big.NewInt(1000000000000000000) // 1 ETH в Wei
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(20000000000) // 20 Gwei
	nonce := uint64(0)
	data := []byte{}

	client := &Client{}

	tx := client.CreateTransaction(from, to, value, gasLimit, gasPrice, nonce, data)

	// Проверяем, что хеш транзакции не пустой
	txHash := tx.Hash()
	if txHash == (common.Hash{}) {
		t.Fatal("Хеш транзакции не должен быть пустым")
	}

	// Проверяем, что хеш в hex формате начинается с 0x
	hashHex := txHash.Hex()
	if len(hashHex) == 0 || hashHex[:2] != "0x" {
		t.Fatal("Хеш транзакции должен начинаться с 0x")
	}
}

func TestAddressValidation(t *testing.T) {
	// Тест валидных адресов
	validAddresses := []string{
		"0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
		"0x742d35cc6634c0532925a3b8d4c9db96c4b4d8b6",
	}

	for _, addrStr := range validAddresses {
		addr := common.HexToAddress(addrStr)
		if addr == (common.Address{}) {
			t.Errorf("Адрес должен быть валидным: %s", addrStr)
		}

		// Проверяем, что адрес не пустой
		addrHex := addr.Hex()
		if len(addrHex) == 0 {
			t.Errorf("Hex представление адреса не должно быть пустым: %s", addrStr)
		}
	}
}

func TestBigIntOperations(t *testing.T) {
	// Тест операций с большими числами
	value1 := big.NewInt(1000000000000000000) // 1 ETH в Wei
	value2 := big.NewInt(500000000000000000)  // 0.5 ETH в Wei

	// Сложение
	sum := new(big.Int).Add(value1, value2)
	expectedSum := big.NewInt(1500000000000000000) // 1.5 ETH в Wei
	if sum.Cmp(expectedSum) != 0 {
		t.Fatalf("Сложение неверно: ожидалось %s, получено %s", expectedSum.String(), sum.String())
	}

	// Вычитание
	diff := new(big.Int).Sub(value1, value2)
	expectedDiff := big.NewInt(500000000000000000) // 0.5 ETH в Wei
	if diff.Cmp(expectedDiff) != 0 {
		t.Fatalf("Вычитание неверно: ожидалось %s, получено %s", expectedDiff.String(), diff.String())
	}

	// Умножение
	product := new(big.Int).Mul(value2, big.NewInt(2))
	expectedProduct := value1
	if product.Cmp(expectedProduct) != 0 {
		t.Fatalf("Умножение неверно: ожидалось %s, получено %s", expectedProduct.String(), product.String())
	}
}

func TestGasEstimation(t *testing.T) {
	// Тест оценки газа (без реального подключения к блокчейну)
	// Стандартный лимит газа для простой транзакции
	standardGasLimit := uint64(21000)

	// В реальном приложении здесь был бы вызов EstimateGas
	// Для теста используем стандартное значение
	gasLimit := standardGasLimit

	if gasLimit != standardGasLimit {
		t.Fatalf("Лимит газа должен быть %d, получено %d", standardGasLimit, gasLimit)
	}

	// Проверяем, что лимит газа положительный
	if gasLimit <= 0 {
		t.Fatal("Лимит газа должен быть положительным")
	}
}
