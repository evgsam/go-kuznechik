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

var SL_dec_lookup [16][256]Block // S⁻¹ ○ L⁻¹

func InitTables() {
	for i := 0; i < 16; i++ {
		for j := 0; j < 256; j++ {
			var y Block
			y[i] = Pi_inverse_table[j]
			y = L_invers(y)
			SL_dec_lookup[i][j] = y
		}
	}
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

// S-функция
func S(block Block) Block {
	result := block
	for i := 0; i < 16; i++ {
		result[i] = Pi_table[block[i]]
	}
	return result
}

// S-функция инверсная
func S_invers(block Block) Block {
	result := block
	for i := 0; i < 16; i++ {
		result[i] = Pi_inverse_table[block[i]]
	}
	return result
}

// S_inv_L_inv через lookup таблицы
func S_inv_L_inv(block Block) Block {
	var result Block
	copy(result[:], SL_dec_lookup[0][block[0]][:])
	for j := 1; j < 16; j++ {
		result = XorBlock(result, SL_dec_lookup[j][block[j]])
	}
	return result
}

// S⁻¹ для финального шага
func S_inverse(block Block) Block {
	return S_invers(block)
}

// функция одного раунда шифрования
func Round(state Block, roundKey RoundKey) Block {
	state = XorKey(state, roundKey) // X
	state = S(state)                // S
	state = L(state)                // L
	return state
}

// функция блочного шифрования
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

// разложение числа на 128 битный вектор
func Vec128(i uint8) RoundKey {
	var vec RoundKey
	// Все байты 0, кроме последнего (индекс 0 = младший)
	vec[15] = i
	return vec
}

// генерация раундовых констант
func GenConstants() [32]RoundKey {
	var constants [32]RoundKey
	for i := 1; i <= 32; i++ {
		vec := Vec128(uint8(i))
		constants[i-1] = L(vec) // Ci = L(Vec128(i))
	}
	return constants
}

// F-функция
func F(a, b, c RoundKey) (RoundKey, RoundKey) {
	temp := XorKey(a, c)    // a ⊕ c
	temp = S(temp)          // S(a ⊕ c)
	temp = L(temp)          // L(S(a ⊕ c))
	newA := XorKey(b, temp) // b ⊕ L(S(a ⊕ c))
	return newA, a
}

// Генерация раундовых ключей
func KeySchedule(masterKey Key256) RoundKeys {
	constants := GenConstants() // 32 константы C1..C32

	var k0, k1 RoundKey           // текущая пара (a,b)
	copy(k0[:], masterKey[:16])   // K1
	copy(k1[:], masterKey[16:32]) // K2

	var rkeys RoundKeys
	rkeys[0] = k0 // K1
	rkeys[1] = k1 // K2

	// 4 группы по 8 F-функций
	for group := 0; group < 4; group++ {
		startC := group * 8 // C1-8, C9-16, C17-24, C25-32

		for step := 0; step < 8; step++ {
			cIdx := startC + step
			k0, k1 = F(k0, k1, constants[cIdx])
		}

		// После 8 F сохраняем новую пару
		rkeys[2+2*group] = k0 // K3,K5,K7,K9
		rkeys[3+2*group] = k1 // K4,K6,K8,K10
	}

	return rkeys
}

// Функция шифрования
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

func GetDecryptRoundKeys(rkeys [10][16]uint8) [10][16]uint8 {
	var rkeys_L [10][16]uint8

	// K1 (индекс 0) — БЕЗ ИЗМЕНЕНИЙ
	rkeys_L[0] = rkeys[0]

	// K2..K10 (индексы 1-9) — L⁻¹(Kᵢ)
	for k := 1; k < 10; k++ {
		rkeys_L[k] = L_invers(rkeys[k])
	}
	return rkeys_L
}

// функция расшифровки
func Decrypt(masterKey Key256, ciphertext Block) Block {
	encKeys := KeySchedule(masterKey)
	decKeys := GetDecryptRoundKeys(encKeys)

	pt := ciphertext

	// ШАГ 1: L⁻¹ на входе
	pt = L_invers(pt)

	// ШАГ 2: 8 раундов K10→K3
	for i := 9; i > 1; i-- {
		pt = XorKey(pt, decKeys[i]) // L⁻¹(Ki)
		pt = S_inv_L_inv(pt)        // SL⁻¹
	}

	// ШАГ 3: Финал
	pt = XorKey(pt, decKeys[1]) // L⁻¹(K2)
	pt = S_inverse(pt)          // S⁻¹
	pt = XorKey(pt, decKeys[0]) // K1

	return pt
}

func main() {
	fmt.Println("Учебный проект по реализации Кузнечика на go")

	InitTables()

	masterkey := Key256{
		0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
		0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
	}

	plaintext := RoundKey{
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x00,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
	}

	ciphertext := Encrypt(masterkey, plaintext)
	plaintext2 := Decrypt(masterkey, ciphertext)

	expected := RoundKey{
		0x7f, 0x67, 0x9d, 0x90, 0xbe, 0xbc, 0x24, 0x30,
		0x5a, 0x46, 0x8d, 0x42, 0xb9, 0xd4, 0xed, 0xcd,
	}

	fmt.Println("=== Тест шифрования ===")
	fmt.Printf("Открытый текст:  % x\n", plaintext)
	fmt.Printf("Шифртекст:       % x\n", ciphertext)
	fmt.Printf("Расшифровка:     % x\n", plaintext2)

	if bytes.Equal(ciphertext[:], expected[:]) {
		fmt.Println("Шифрование работает")
	} else {

		fmt.Println("Ошибка в шифровании")
	}

	if bytes.Equal(plaintext2[:], plaintext[:]) {
		fmt.Println("Расшифрование работает")
	} else {

		fmt.Println("Ошибка в расшифровке")
	}

}
