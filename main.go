package main

import "fmt"

// Предполагаемые к использованию в коде типы данных
type Block = [16]byte
type Key256 = [32]byte
type RoundKey = [16]byte
type RoundKeys = [10]RoundKey

// Предполагаемые к использованию функции
// инициалищация контекста шифра
func NewKuznyechik(key Key256) RoundKeys

// шифрование одного блока
func EncryptBlock(plaintext Block, rk RoundKeys) (ciphertext Block)

// расшифрование одного блока
func DecryptBlock(ciphertext Block, rk RoundKeys) (plaintext Block)

// X-преобразования
func XorBlock(a, b Block) (res Block)
func XorKey(a, b RoundKey) (res RoundKey)

func main() {
	fmt.Println("Учебный проект по реализации Кузнечика на go")
}
