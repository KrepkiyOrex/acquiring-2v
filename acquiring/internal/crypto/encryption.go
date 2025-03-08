package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"os"
)

var encryptionKey = []byte(os.Getenv("ENCRYPTION_KEY"))
var fixedIV = []byte(os.Getenv("FIXED_IV"))

// используется только для локального запуска, в docker compose
// .env загружается автоматом
// func init() {
// 	err := godotenv.Load("../.env") // Загружаем .env перед инициализацией переменных
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	encryptionKey = []byte(os.Getenv("ENCRYPTION_KEY"))
// 	fixedIV = []byte(os.Getenv("FIXED_IV"))

// 	if len(encryptionKey) == 0 || len(fixedIV) == 0 {
// 		log.Fatal("ENCRYPTION_KEY or FIXED_IV is missing in .env")
// 	}
// }

type CardData struct {
	ID                  uint    `json:"id" gorm:"primaryKey"`
	Balance             float64 `json:"balance"`
	EncryptedCardNumber string  `json:"encryptedCardNumber"`
	EncryptedExpiryDate string  `json:"encryptedExpiryDate"`
	EncryptedCVV        string  `json:"encryptedCVV"`
	EncryptedCardName   string  `json:"encryptedCardName"`
}

func ProcessEncrypt(details *CardData) error {
	if err := Encrypt(&details.EncryptedCardNumber); err != nil {
		return err
	}

	if err := Encrypt(&details.EncryptedExpiryDate); err != nil {
		return err
	}

	if err := Encrypt(&details.EncryptedCVV); err != nil {
		return err
	}

	if err := Encrypt(&details.EncryptedCardName); err != nil {
		return err
	}

	return nil
}

func Encrypt(data *string) error {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return err
	}

	ciphertext := make([]byte, aes.BlockSize+len(*data))
	iv := fixedIV
	copy(ciphertext[:aes.BlockSize], iv)

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(*data))

	*data = base64.StdEncoding.EncodeToString(ciphertext)

	return nil
}

func ProcessDecrypt(details *CardData) error {
	if err := Decrypt(&details.EncryptedCardNumber); err != nil {
		return err
	}

	if err := Decrypt(&details.EncryptedExpiryDate); err != nil {
		return err
	}

	if err := Decrypt(&details.EncryptedCVV); err != nil {
		return err
	}

	if err := Decrypt(&details.EncryptedCardName); err != nil {
		return err
	}

	return nil
}

func Decrypt(data *string) error {
	ciphertext, _ := base64.StdEncoding.DecodeString(*data)
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return err
	}

	if len(ciphertext) < aes.BlockSize {
		return fmt.Errorf("ciphertext too short: length %d, expected at least %d", len(ciphertext), aes.BlockSize)
	}

	iv := fixedIV
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	*data = string(ciphertext)

	return nil
}
