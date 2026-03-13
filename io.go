// ========== I/O операции ==========

package main

import (
	"fmt"
	"os"
)

// pkcs7Pad — добавление padding PKCS#7
// Блок 16 байт дополняется до длины 16 с помощью PKCS#7 padding
func pkcs7Pad(block []byte, blockSize int) []byte {
	padLen := blockSize - len(block)
	padded := make([]byte, blockSize)
	copy(padded, block)
	for i := len(block); i < blockSize; i++ {
		padded[i] = byte(padLen)
	}
	return padded
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
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("открыть %s: %w", inputPath, err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("создать %s: %w", outputPath, err)
	}
	defer outFile.Close()

	buffer := make([]byte, 16)

	for {
		n, _ := inFile.Read(buffer)
		if n == 0 {
			break
		}

		if n < 16 {
			// Последний блок — добавляем padding
			padded := pkcs7Pad(buffer[:n], 16)
			ciphertext := Encrypt(masterKey, RoundKey(padded))
			_, _ = outFile.Write(ciphertext[:])
			fmt.Printf("Padding: %d байт → %d байт\n", n, 16)
			break
		}

		// Полный блок 16 байт
		block := RoundKey(buffer)
		ciphertext := Encrypt(masterKey, block)
		_, _ = outFile.Write(ciphertext[:])
	}
	return nil
}

// DecryptFileStream — потоковое расшифрование файла
// Читает зашифрованный файл по блокам 16 байт и удаляет padding
func DecryptFileStream(inputPath, outputPath string, masterKey Key256) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("открыть %s: %w", inputPath, err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("создать %s: %w", outputPath, err)
	}
	defer outFile.Close()

	buffer := make([]byte, 16)
	allDecrypted := make([]byte, 0)

	for {
		n, _ := inFile.Read(buffer)
		if n == 0 {
			break
		}
		if n != 16 {
			return fmt.Errorf("неполный блок")
		}

		block := Block(buffer)
		plaintext := Decrypt(masterKey, block)
		allDecrypted = append(allDecrypted, plaintext[:]...)
	}

	// Удаляем padding
	cleanData, padErr := pkcs7Unpad(allDecrypted)
	if padErr != nil {
		return fmt.Errorf("padding: %w", padErr)
	}

	if _, writeErr := outFile.Write(cleanData); writeErr != nil {
		return fmt.Errorf("запись: %w", writeErr)
	}

	return nil
}
