// Кузнечик — функции расшифрования

package main

// Decrypt — функция расшифрования с раундовой схемой
// Выполняет 10 раундов с использованием KeySchedule
// 9*(L⁻¹->S⁻¹->X)->X
func Decrypt(masterKey Key256, ciphertext Block) Block {
	key := KeySchedule(masterKey)
	pt := ciphertext
	for i := 9; i >= 1; i-- {
		pt = X(pt, key[i]) // L⁻¹(Ki)
		pt = LInvers(pt)
		pt = SInvers(pt)
	}
	pt = X(pt, key[0])
	return pt
}
