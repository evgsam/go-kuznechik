package main

import "fmt"

// Предполагаемые к использованию в коде типы данных
type Block = [16]byte
type Key256 = [32]byte
type RoundKey = [16]byte
type RoundKeys = [10]RoundKey

func main() {
	fmt.Println("Учебный проект по реализации Кузнечика на go")
}
