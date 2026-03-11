package main

// Типы данных для блочного шифра Кузнечик
type Block = [16]byte
type Key256 = [32]byte
type RoundKey = [16]byte
type RoundKeys = [10]RoundKey
