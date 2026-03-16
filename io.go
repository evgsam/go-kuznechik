// ========== I/O операции ==========

package main

import (
	"fmt"
	"os"
)

// pkcs7Pad — добавление padding PKCS#7
// Блок 16 байт дополняется до длины 16 с помощью PKCS#7 padding
func pkcs7Pad(data []byte, blockSize int) []byte {
	if len(data) == 0 {
		// Пустой файл → полный блок padding'а
		result := make([]byte, blockSize)
		for i := 0; i < blockSize; i++ {
			result[i] = byte(blockSize)
		}
		return result
	}

	padLen := blockSize - (len(data) % blockSize)
	result := make([]byte, len(data)+padLen)
	copy(result, data)

	// Заполняем padding
	for i := len(data); i < len(result); i++ {
		result[i] = byte(padLen)
	}
	return result
}

// pkcs7Unpad — удаление padding PKCS#7
// Проверяет и удаляет PKCS#7 padding из расшифрованных данных
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 || len(data)%16 != 0 {
		return nil, fmt.Errorf("invalid data length")
	}

	padLen := int(data[len(data)-1])
	if padLen == 0 || padLen > 16 {
		return nil, fmt.Errorf("invalid pad length: %d", padLen)
	}

	for i := len(data) - padLen; i < len(data); i++ {
		if data[i] != byte(padLen) {
			return nil, fmt.Errorf("invalid padding bytes")
		}
	}

	return data[:len(data)-padLen], nil
}

// EncryptFileStream — потоковое шифрование файла
// Читает входной файл по блокам 16 байт, добавляет padding и шифрует
func EncryptFileStream(inputPath, outputPath string, masterKey Key256) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("открыть %s: %w", inputPath, err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("создать %s: %w", outputPath, err)
	}
	defer outFile.Close()
	//добавляю padding
	padded := pkcs7Pad(data, 16)
	for i := 0; i < len(padded); i += 16 {
		block := RoundKey(padded[i : i+16])
		ciphertext := Encrypt(masterKey, block)
		_, _ = outFile.Write(ciphertext[:])
	}
	return nil
}

// DecryptFileStream — потоковое расшифрование файла
// Читает зашифрованный файл по блокам 16 байт и удаляет padding
func DecryptFileStream(inputPath, outputPath string, masterKey Key256) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("открыть %s: %w", inputPath, err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("создать %s: %w", outputPath, err)
	}
	defer outFile.Close()

	decryptedData := make([]byte, 0)
	for i := 0; i < len(data); i += 16 {
		block := Block(data[i : i+16])
		plaintext := Decrypt(masterKey, block)
		decryptedData = append(decryptedData, plaintext[:]...)
	}

	// Удаляю padding
	cleanData, err := pkcs7Unpad(decryptedData)
	if err != nil {
		return fmt.Errorf("padding: %w", err)
	}
	outFile.Write(cleanData)
	return nil
}
