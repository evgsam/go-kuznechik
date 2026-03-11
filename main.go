package main

import (
	"flag"
	"fmt"
	"os"
)

// Переменные из стандарта
const gf8 = 0xc3

// Таблицы для оптимизации расшифрования
var SL_dec_lookup [16][256]Block

// InitTables — инициализация таблицы SL⁻¹
func InitTables() {
	for i := 0; i < 16; i++ {
		for j := 0; j < 256; j++ {
			var y Block
			y[i] = Pi_inverse_table[j]
			y = L_invers(y)
			SL_dec_lookup[i][j] = y
		}
	}
}

func parseMasterKeyFromHex(keyStr string) Key256 {
	var key Key256
	if len(keyStr) != 64 {
		fmt.Fprintf(os.Stderr, "Неверная длина ключа: %d (ожида 64 символа)\n", len(keyStr))
		os.Exit(1)
	}
	for i := 0; i < 32; i++ {
		var byteVal uint8
		_, err := fmt.Sscanf(keyStr[i*2:i*2+2], "%2x", &byteVal)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка парсинга ключа на позиции %d: %v\n", i, err)
			os.Exit(1)
		}
		key[i] = byteVal
	}
	return key
}

func blockEqual(a, b Block) bool {
	for i := 0; i < 16; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// main — точка входа
func main() {
	// Парсинг аргументов командной строки
	masterKeyFlag := flag.String("k", "", "Master key в hex формате (32 байта = 64 символа)")
	flag.Parse()

	var masterkey Key256
	if *masterKeyFlag != "" {
		masterkey = parseMasterKeyFromHex(*masterKeyFlag)
	} else {
		fmt.Println("Внимание: Master key не указан, используется дефолтный")
		masterkey = Key256{
			0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
			0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
			0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10,
			0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		}
	}

	InitTables()

	plaintext := RoundKey{
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x00,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
	}

	ciphertext := Encrypt(masterkey, plaintext)
	plaintext2 := Decrypt(masterkey, ciphertext)

	expected := RoundKey{
		0x7f, 0x67, 0x9d, 0x90, 0xbe, 0xbc, 0x24, 0x30,
		0x5a, 0x46, 0x8d, 0x42, 0xb9, 0xd4, 0xed, 0xcd,
	}

	fmt.Println("=== Тест шифрования ===")
	fmt.Printf("Открытый текст:  % x\n", plaintext)
	fmt.Printf("Шифртекст:       % x\n", ciphertext)
	fmt.Printf("Расшифровка:     % x\n", plaintext2)

	if blockEqual(ciphertext, expected) {
		fmt.Println("Шифрование работает")
	} else {
		fmt.Println("Ошибка в шифровании")
	}

	if blockEqual(plaintext2, plaintext) {
		fmt.Println("Расшифрование работает")
	} else {
		fmt.Println("Ошибка в расшифровке")
	}
}
