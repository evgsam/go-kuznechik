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

// L-преобразование и вспомогательные таблицы

// L — линейное преобразование блока
// Использует 16 итераций линейного сдвига и GF(2^8) умножения
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

// L_invers — обратное L-преобразование
// Использует LFSR-подобную схему для восстановления исходного блока
func L_invers(block Block) Block {
	var x uint8
	for j := 0; j < 16; j++ { // 16 итераций LFSR
		x = block[0]
		for i := 0; i < 15; i++ { // Сдвиг влево
			block[i] = block[i+1]
			x ^= GF8Mul(block[i], L_coeffs[i])
		}
		block[15] = x
	}
	return block
}

// S — подстановка (S-блок)
// Применяет таблицу подстановки Pi к каждому байту блока
func S(block Block) Block {
	result := block
	for i := 0; i < 16; i++ {
		result[i] = Pi_table[block[i]]
	}
	return result
}

// S_invers — обратная S-функция
// Применяет обратную таблицу подстановки Pi⁻¹ к каждому байту
func S_invers(block Block) Block {
	result := block
	for i := 0; i < 16; i++ {
		result[i] = Pi_inverse_table[block[i]]
	}
	return result
}

// XorBlock — побитовое сложение двух блоков
func X(a, b Block) (res Block) {
	var i int
	for i = 0; i < 16; i++ {
		res[i] = a[i] ^ b[i]
	}
	return res
}
