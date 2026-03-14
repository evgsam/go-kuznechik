// ========== Генерация ключей ==========

package main

// Vec128 — конструктор RoundKey из uint8
// Все байты равны 0, кроме последнего (индекс 0 = младший)
func Vec128(i uint8) RoundKey {
	var vec RoundKey
	// Все байты 0, кроме последнего (индекс 0 = младший)
	vec[15] = i
	return vec
}

// GenConstants — генерация 32 раундовых констант
// Ci = L(Vec128(i)) для i = 1..32
func GenConstants() [32]RoundKey {
	var constants [32]RoundKey
	for i := 1; i <= 32; i++ {
		vec := Vec128(uint8(i))
		constants[i-1] = L(vec) // Ci = L(Vec128(i))
	}
	return constants
}

// F — раундовая функция для расширения ключа
// Используется в KeySchedule для генерации раундовых ключей
func F(a, b, c RoundKey) (RoundKey, RoundKey) {
	temp := XorKey(a, c)    // a ⊕ c
	temp = S(temp)          // S(a ⊕ c)
	temp = L(temp)          // L(S(a ⊕ c))
	newA := XorKey(b, temp) // b ⊕ L(S(a ⊕ c))
	return newA, a
}

// KeySchedule — генерация 10 раундовых ключей из 256-битного ключа
// K1, K2 берутся из masterKey
// K3..K10 генерируются через 4 группы по 8 F-функций
func KeySchedule(masterKey Key256) RoundKeys {
	roundConstants := GenConstants() // 32 константы C1..C32
	var k0, k1 RoundKey              // текущая пара (a,b)
	copy(k0[:], masterKey[:16])      // K1
	copy(k1[:], masterKey[16:32])    // K2
	var rkeys RoundKeys
	rkeys[0] = k0 // K1
	rkeys[1] = k1 // K2
	// 4 группы по 8 F-функций
	for group := 0; group < 4; group++ {
		startC := group * 8 // C1-8, C9-16, C17-24, C25-32

		for step := 0; step < 8; step++ {
			cIdx := startC + step
			k0, k1 = F(k0, k1, roundConstants[cIdx])
		}
		// После 8 F сохраняем новую пару
		rkeys[2+2*group] = k0 // K3,K5,K7,K9
		rkeys[3+2*group] = k1 // K4,K6,K8,K10
	}

	return rkeys
}

// Decrypt — функция расшифрования блока
func Decrypt(masterKey Key256, ciphertext Block) Block {
	key := KeySchedule(masterKey)
	pt := ciphertext
	for i := 9; i >= 1; i-- {
		pt = XorKey(pt, key[i]) // L⁻¹(Ki)
		pt = L_invers(pt)
		pt = S_invers(pt)
	}
	pt = XorKey(pt, key[0])
	return pt
}
