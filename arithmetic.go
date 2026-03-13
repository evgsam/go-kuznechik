// Кузнечик — арифметика

package main

// gf8 — неприводимый полином поля Галуа GF(2^8)
const gf8 = 0xc3

// GF8Mul — умножение байтов в поле Галуа GF(2^8)
// Использует неприводимый полином 0xC3 (x^8 + x^7 + x^6 + x + 1)
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

// XorBlock — побитовое сложение двух блоков
func XorBlock(a, b Block) (res Block) {
	var i int
	for i = 0; i < 16; i++ {
		res[i] = a[i] ^ b[i]
	}
	return res
}

// XorKey — побитовое сложение двух раундовых ключей
func XorKey(a, b RoundKey) RoundKey {
	return XorBlock(a, b)
}
