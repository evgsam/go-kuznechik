// Кузнечик — функции шифрования

package main

// Encrypt — функция шифрования с раундовой схемой
// Выполняет 10 раундов с использованием KeySchedule
// 9*(X->S->L)->X
func Encrypt(masterKey Key256, block RoundKey) RoundKey {
	key := KeySchedule(masterKey)
	state := block
	for i := 0; i < 9; i++ {
		state = XorKey(state, key[i]) // X
		state = S(state)              // S
		state = L(state)              // L
	}
	state = XorKey(state, key[9])
	return state
}
