package main

import "fmt"

// Предполагаемые к использованию в коде типы данных
type Block = [16]byte
type Key256 = [32]byte
type RoundKey = [16]byte
type RoundKeys = [10]RoundKey

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

func XorKey(a, b RoundKey) (res RoundKey) {
	return XorBlock(a, b)
}

func main() {
	fmt.Println("Учебный проект по реализации Кузнечика на go")
	zero := Block{}     // все 0
	ones := Block{0xFF} // все 0xFF (Go заполнит все байты)

	var result Block
	var rkResult RoundKey

	result = XorBlock(zero, ones)
	// result должен быть все 0xFF

	rk1 := RoundKey{0xAA} // все 0xAA
	rk2 := RoundKey{0x55} // все 0x55
	rkResult = XorKey(rk1, rk2)

	fmt.Println(result)
	fmt.Println(rkResult)
}
