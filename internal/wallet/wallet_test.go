package wallet

import (
	"math/big"
	"os"
	"testing"

	"crypto-wallet/internal/crypto"
)

func TestNewWallet(t *testing.T) {
	// Тест создания кошелька с невалидным URL
	_, err := NewWallet("invalid-url", "test-wallet.json")
	if err == nil {
		t.Fatal("Должна быть ошибка для невалидного URL")
	}

	// Тест создания кошелька с пустым URL
	_, err = NewWallet("", "test-wallet.json")
	if err == nil {
		t.Fatal("Должна быть ошибка для пустого URL")
	}
}

func TestGenerateNewWallet(t *testing.T) {
	// Создаем временный кошелек для тестирования
	w, err := NewWallet("https://sepolia.infura.io/v3/test", "test-wallet.json")
	if err != nil {
		t.Fatalf("Ошибка создания кошелька: %v", err)
	}
	defer w.Close()

	// Генерируем новый кошелек
	err = w.GenerateNewWallet()
	if err != nil {
		t.Fatalf("Ошибка генерации кошелька: %v", err)
	}

	// Проверяем, что кошелек создан
	if w.KeyPair == nil {
		t.Fatal("KeyPair не должен быть nil после генерации")
	}

	// Проверяем, что адрес не пустой
	address, err := w.GetAddress()
	if err != nil {
		t.Fatalf("Ошибка получения адреса: %v", err)
	}

	if address == "" {
		t.Fatal("Адрес не должен быть пустым")
	}

	// Проверяем, что адрес начинается с 0x
	if len(address) < 2 || address[:2] != "0x" {
		t.Fatal("Адрес должен начинаться с 0x")
	}

	// Очищаем тестовый файл
	os.Remove("test-wallet.json")
}

func TestSaveAndLoadWallet(t *testing.T) {
	// Создаем временный кошелек
	w, err := NewWallet("https://sepolia.infura.io/v3/test", "test-wallet.json")
	if err != nil {
		t.Fatalf("Ошибка создания кошелька: %v", err)
	}
	defer w.Close()

	// Генерируем новый кошелек
	err = w.GenerateNewWallet()
	if err != nil {
		t.Fatalf("Ошибка генерации кошелька: %v", err)
	}

	// Получаем адрес до сохранения
	originalAddress, err := w.GetAddress()
	if err != nil {
		t.Fatalf("Ошибка получения адреса: %v", err)
	}

	// Создаем новый кошелек для загрузки
	w2, err := NewWallet("https://sepolia.infura.io/v3/test", "test-wallet.json")
	if err != nil {
		t.Fatalf("Ошибка создания второго кошелька: %v", err)
	}
	defer w2.Close()

	// Загружаем кошелек
	err = w2.LoadWallet()
	if err != nil {
		t.Fatalf("Ошибка загрузки кошелька: %v", err)
	}

	// Проверяем, что адрес совпадает
	loadedAddress, err := w2.GetAddress()
	if err != nil {
		t.Fatalf("Ошибка получения загруженного адреса: %v", err)
	}

	if originalAddress != loadedAddress {
		t.Fatalf("Адреса не совпадают: оригинал %s, загруженный %s", originalAddress, loadedAddress)
	}

	// Очищаем тестовый файл
	os.Remove("test-wallet.json")
}

func TestGetAddress(t *testing.T) {
	// Создаем временный кошелек
	w, err := NewWallet("https://sepolia.infura.io/v3/test", "test-wallet.json")
	if err != nil {
		t.Fatalf("Ошибка создания кошелька: %v", err)
	}
	defer w.Close()

	// Пытаемся получить адрес без инициализации кошелька
	_, err = w.GetAddress()
	if err == nil {
		t.Fatal("Должна быть ошибка при получении адреса неинициализированного кошелька")
	}

	// Генерируем кошелек
	err = w.GenerateNewWallet()
	if err != nil {
		t.Fatalf("Ошибка генерации кошелька: %v", err)
	}

	// Получаем адрес
	address, err := w.GetAddress()
	if err != nil {
		t.Fatalf("Ошибка получения адреса: %v", err)
	}

	if address == "" {
		t.Fatal("Адрес не должен быть пустым")
	}

	// Очищаем тестовый файл
	os.Remove("test-wallet.json")
}

func TestSendTransactionValidation(t *testing.T) {
	// Создаем временный кошелек
	w, err := NewWallet("https://sepolia.infura.io/v3/test", "test-wallet.json")
	if err != nil {
		t.Fatalf("Ошибка создания кошелька: %v", err)
	}
	defer w.Close()

	// Пытаемся отправить транзакцию без инициализации кошелька
	amount := big.NewFloat(0.001)
	_, err = w.SendTransaction("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6", amount)
	if err == nil {
		t.Fatal("Должна быть ошибка при отправке транзакции неинициализированным кошельком")
	}

	// Генерируем кошелек
	err = w.GenerateNewWallet()
	if err != nil {
		t.Fatalf("Ошибка генерации кошелька: %v", err)
	}

	// Пытаемся отправить на невалидный адрес
	_, err = w.SendTransaction("invalid-address", amount)
	if err == nil {
		t.Fatal("Должна быть ошибка для невалидного адреса")
	}

	// Пытаемся отправить отрицательную сумму
	negativeAmount := big.NewFloat(-0.001)
	_, err = w.SendTransaction("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6", negativeAmount)
	if err == nil {
		t.Fatal("Должна быть ошибка для отрицательной суммы")
	}

	// Очищаем тестовый файл
	os.Remove("test-wallet.json")
}

func TestSignAndVerifyMessage(t *testing.T) {
	// Создаем временный кошелек
	w, err := NewWallet("https://sepolia.infura.io/v3/test", "test-wallet.json")
	if err != nil {
		t.Fatalf("Ошибка создания кошелька: %v", err)
	}
	defer w.Close()

	// Генерируем кошелек
	err = w.GenerateNewWallet()
	if err != nil {
		t.Fatalf("Ошибка генерации кошелька: %v", err)
	}

	message := []byte("Hello, Ethereum!")

	// Подписываем сообщение
	signature, err := w.SignMessage(message)
	if err != nil {
		t.Fatalf("Ошибка подписи сообщения: %v", err)
	}

	if len(signature) == 0 {
		t.Fatal("Подпись не должна быть пустой")
	}

	// Получаем адрес кошелька
	address, err := w.GetAddress()
	if err != nil {
		t.Fatalf("Ошибка получения адреса: %v", err)
	}

	// Проверяем подпись
	isValid := w.VerifyMessage(message, signature, address)
	if !isValid {
		t.Fatal("Подпись должна быть валидной")
	}

	// Проверяем с неправильным сообщением
	wrongMessage := []byte("Wrong message")
	isValid = w.VerifyMessage(wrongMessage, signature, address)
	if isValid {
		t.Fatal("Подпись не должна быть валидной для неправильного сообщения")
	}

	// Проверяем с неправильным адресом
	wrongAddress := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
	isValid = w.VerifyMessage(message, signature, wrongAddress)
	if isValid {
		t.Fatal("Подпись не должна быть валидной для неправильного адреса")
	}

	// Очищаем тестовый файл
	os.Remove("test-wallet.json")
}

func TestWalletDataStructure(t *testing.T) {
	// Тест структуры WalletData
	walletData := WalletData{
		PrivateKey: "test-private-key",
		PublicKey:  "test-public-key",
		Address:    "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
	}

	if walletData.PrivateKey != "test-private-key" {
		t.Fatal("PrivateKey не совпадает")
	}

	if walletData.PublicKey != "test-public-key" {
		t.Fatal("PublicKey не совпадает")
	}

	if walletData.Address != "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6" {
		t.Fatal("Address не совпадает")
	}
}

func TestCryptoIntegration(t *testing.T) {
	// Тест интеграции с криптографическими функциями
	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Ошибка генерации пары ключей: %v", err)
	}

	// Проверяем, что адрес валидный
	address := keyPair.GetAddressHex()
	if !crypto.IsValidAddress(address) {
		t.Fatal("Сгенерированный адрес должен быть валидным")
	}

	// Проверяем подпись сообщения
	message := []byte("Test message")
	signature, err := keyPair.SignMessage(message)
	if err != nil {
		t.Fatalf("Ошибка подписи сообщения: %v", err)
	}

	// Проверяем подпись
	isValid := crypto.VerifySignature(message, signature, keyPair.Address)
	if !isValid {
		t.Fatal("Подпись должна быть валидной")
	}
}

func TestBigFloatOperations(t *testing.T) {
	// Тест операций с большими числами с плавающей точкой
	amount1 := big.NewFloat(1.5)
	amount2 := big.NewFloat(0.5)

	// Сложение
	sum := new(big.Float).Add(amount1, amount2)
	expectedSum := big.NewFloat(2.0)
	if sum.Cmp(expectedSum) != 0 {
		t.Fatalf("Сложение неверно: ожидалось %s, получено %s", expectedSum.Text('f', 18), sum.Text('f', 18))
	}

	// Вычитание
	diff := new(big.Float).Sub(amount1, amount2)
	expectedDiff := big.NewFloat(1.0)
	if diff.Cmp(expectedDiff) != 0 {
		t.Fatalf("Вычитание неверно: ожидалось %s, получено %s", expectedDiff.Text('f', 18), diff.Text('f', 18))
	}

	// Умножение
	product := new(big.Float).Mul(amount1, amount2)
	expectedProduct := big.NewFloat(0.75)
	if product.Cmp(expectedProduct) != 0 {
		t.Fatalf("Умножение неверно: ожидалось %s, получено %s", expectedProduct.Text('f', 18), product.Text('f', 18))
	}
}
