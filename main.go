// Кузнечик — блочный шифр ГОСТ Р 34.12-2015
// Реализация на Go в рамках учебного проекта по
// криптографическим методам защиты информации
// Самарин Евгений

package main

import (
	"flag"
	"fmt"
	"os"
)

// ========== Инициализация ==========

// SL_dec_lookup — таблица для оптимизации SL⁻¹ (S⁻¹∘L⁻¹)
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

// ========== Конвертация ключей ==========

// parseMasterKeyFromHex — парсинг ключа из hex-строки
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

// ========== Утилиты ==========

// blockEqual — проверка равенства двух блоков
func blockEqual(a, b Block) bool {
	for i := 0; i < 16; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ========== CLI ==========

// printHelp — вывод справки по использованию
func printHelp() {
	fmt.Println("Кузнечик — блочный шифр ГОСТ Р 34.12-2015")
	fmt.Println()
	fmt.Println("Использование:")
	fmt.Println("  kuznechik -e -i <input> -o <output> -k <key>")
	fmt.Println("  kuznechik -d -i <input> -o <output> -k <key>")
	fmt.Println()
	fmt.Println("Параметры:")
	fmt.Println("  -e            шифрование")
	fmt.Println("  -d            расшифровка")
	fmt.Println("  -i <file>     входной файл")
	fmt.Println("  -o <file>     выходной файл")
	fmt.Println("  -k <key>      ключ (64 hex символа)")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  kuznechik -e -i data.txt -o data.enc -k 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	fmt.Println("  kuznechik -d -i data.enc -o data.txt -k 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
}

// main — точка входа с поддержкой CLI-интерфейса
func main() {
	flag.Usage = func() { printHelp() }

	keyFlag := flag.String("k", "", "ключ (64 hex)")
	inputFlag := flag.String("i", "", "входной файл")
	outputFlag := flag.String("o", "", "выходной файл")
	encryptFlag := flag.Bool("e", false, "шифрование")
	decryptFlag := flag.Bool("d", false, "расшифровка")

	flag.Parse()

	// CLI режим файлов
	if *inputFlag != "" && *outputFlag != "" {
		var masterKey Key256
		if *keyFlag == "" {
			fmt.Fprintf(os.Stderr, "Ключ не указан. Используйте -k <key>\n")
			os.Exit(1)
		}
		masterKey = parseMasterKeyFromHex(*keyFlag)

		InitTables()

		if *encryptFlag {
			if err := EncryptFileStream(*inputFlag, *outputFlag, masterKey); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		} else if *decryptFlag {
			if err := DecryptFileStream(*inputFlag, *outputFlag, masterKey); err != nil {
				fmt.Fprintf(os.Stderr, " %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Укажите -e (шифрование) или -d (расшифровка)")
			os.Exit(1)
		}
		fmt.Println("Операция выполнена!")
		return
	}
	fmt.Println("Укажите -h для справки")
	os.Exit(1)
}
