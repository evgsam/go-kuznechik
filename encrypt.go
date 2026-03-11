package main

// Функции шифрования

// Round — один раунд шифрования (X → S → L)
func Round(state Block, roundKey RoundKey) Block {
	state = XorKey(state, roundKey) // X
	state = S(state)                // S
	state = L(state)                // L
	return state
}

// EncryptBlock — блочное шифрование
// 9 полных раундов (1-9) + финал (X + S без L)
func EncryptBlock(plaintext Block, roundKeys RoundKeys) Block {
	state := plaintext

	// 9 полных раундов (1-9)
	for r := 0; r < 9; r++ {
		state = Round(state, roundKeys[r])
	}

	// 10-й раунд: только X + S (БЕЗ L)
	state = XorKey(state, roundKeys[9])
	state = S(state)

	return state
}

// Encrypt — функция шифрования с раундовой схемой
func Encrypt(masterKey Key256, block RoundKey) RoundKey {
	decKeys := KeySchedule(masterKey)
	state := block

	state = XorKey(state, decKeys[0]) // X[K1]

	for i := 0; i < 9; i++ {
		state = S(state)
		state = L(state)
		state = XorKey(state, decKeys[i+1]) // K2..K10
	}

	return state
}
