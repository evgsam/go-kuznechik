package main

import (
	"bytes"
	"fmt"
)

// Предполагаемые к использованию в коде типы данных
type Block = [16]byte
type Key256 = [32]byte
type RoundKey = [16]byte
type RoundKeys = [10]RoundKey

// числовое выражение -константа полинома x^8 + x^7 + x^6 + x + 1
const gf8 = 0xc3

// коэфициенты для линейного преобразования
var L_coeffs = [16]byte{
	0x94, 0x20, 0x85, 0x10, 0xC2, 0xC0, 0x01, 0xFB,
	0x01, 0xC0, 0xC2, 0x10, 0x85, 0x20, 0x94, 0x01,
}

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

// обертка для умножения констант
func MulByConst(b, c uint8) uint8 {
	return GF8Mul(b, c)
}

// L-функция
func L(block Block) Block {
	var i, j int
	var x uint8
	for j = 0; j < 16; j++ { // 16 R-итераций
		x = block[15]             // x=a[15]
		for i = 14; i >= 0; i-- { //сдвигаю вправо
			block[i+1] = block[i] // a_i -> a_{i+1}
			x = x ^ GF8Mul(block[i], L_coeffs[i])
		}
		block[0] = x // новый a0 = l(...)
	}
	return block
}

// L-функция инверсная
func L_invers(block Block) Block {
	var i, j int
	var x uint8
	for j = 0; j < 16; j++ { // 16 R-итераций
		x = block[0]             // x=a[0]
		for i = 0; i < 15; i++ { //сдвигаю влево
			block[i] = block[i+1] //  a_{i+1} -> a_i
			x = x ^ GF8Mul(block[i], L_coeffs[i])
		}
		block[15] = x // новый a15 = l(...)
	}
	return block
}

func main() {
	fmt.Println("Учебный проект по реализации Кузнечика на go")

	input := Block{
		0xd4, 0x56, 0x58, 0x4d, 0xd0, 0xe3, 0xe8, 0x4c,
		0xc3, 0x16, 0x6e, 0x4b, 0x7f, 0xa2, 0x89, 0x0d,
	}

	answer := Block{
		0x79, 0xd2, 0x62, 0x21, 0xb8, 0x7b, 0x58, 0x4c,
		0xd4, 0x2f, 0xbc, 0x4f, 0xfe, 0xa5, 0xde, 0x9a,
	}

	result := L(input)
	fmt.Printf("L результат: %s\n", BlockToHex(result)) // ← ДОБАВЬ ЭТО
	fmt.Printf("Ожидаемый:   79d26221b87b584cd42fbc4ffea5de9a\n")

	if bytes.Equal(result[:], answer[:]) {
		fmt.Println("Тест L ПРОВЕДЁН УСПЕШНО!")
	} else {
		fmt.Println("Тест L ПРОВАЛЕН!")
	}

}
