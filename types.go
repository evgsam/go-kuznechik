// Кузнечик — типы данных

package main

// Block — тип блока данных (16 байт)
type Block = [16]byte

// Key256 — тип 256-битного ключа
type Key256 = [32]byte

// RoundKey — тип раундового ключа (128 бит)
type RoundKey = [16]byte

// RoundKeys — массив раундовых ключей
type RoundKeys = [10]RoundKey
