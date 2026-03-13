// ========== L-преобразование ==========

package main

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

// S_inverse — S⁻¹ для финального шага расшифрования
// Alias для S_invers
func S_inverse(block Block) Block {
	return S_invers(block)
}

// S_inv_L_inv — оптимизированное вычисление (S⁻¹∘L⁻¹) через lookup-таблицы
// Использует предвычисленную таблицу SL_dec_lookup для ускорения
func S_inv_L_inv(block Block) Block {
	var result Block
	copy(result[:], SL_dec_lookup[0][block[0]][:])
	for j := 1; j < 16; j++ {
		result = XorBlock(result, SL_dec_lookup[j][block[j]])
	}
	return result
}
