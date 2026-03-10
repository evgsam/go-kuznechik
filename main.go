package main

import "fmt"

// Предполагаемые к использованию в коде типы данных
type Block = [16]byte
type Key256 = [32]byte
type RoundKey = [16]byte
type RoundKeys = [10]RoundKey

// числовое выражение -константа полинома x^8 + x^7 + x^6 + x + 1
const gf8 = 0xc3

// Предполагаемые к использованию функции
// инициалищация контекста шифра
/*func NewKuznyechik(key Key256) RoundKeys

// шифрование одного блока
func EncryptBlock(plaintext Block, rk RoundKeys) (ciphertext Block)

// расшифрование одного блока
func DecryptBlock(ciphertext Block, rk RoundKeys) (plaintext Block)
*/
// X-преобразования
// xor двух блоков c= a XOR b
func XorBlock(a, b Block) (res Block) {
	var i int
	for i = 0; i < 16; i++ {
		res[i] = a[i] ^ b[i]
	}
	return res
}

func XorKey(a, b RoundKey) RoundKey {
	return XorBlock(a, b)
}

// функция умножения чисел в поле Галуа над неприводимым полиномом x^8 + x^7 + x^6 + x + 1 (0xc3)
func GF8Mul(a, b uint8) uint8 {
	var c uint8
	c = 0
	for b != 0 { // проверяем, остались ли биты в b
		if b&1 != 0 {
			c = c ^ a
		}
		if a&0x80 != 0 {
			a = (a << 1) ^ gf8
		} else {
			a = a << 1
		}
		b >>= 1 // переходим к следующему биту
	}
	return c
}

func main() {
	fmt.Println("Учебный проект по реализации Кузнечика на go")
	fmt.Println(GF8Mul(0x01, 0x01)) // Должно быть 0x01
	fmt.Println(GF8Mul(0x02, 0x01)) // Должно быть 0x02
	fmt.Println(GF8Mul(0xFF, 0x01)) // Должно быть 0xFF
}
