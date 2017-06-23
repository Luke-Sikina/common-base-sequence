package main

type Sequence [16]rune

func (sequence *Sequence) hash() (hash uint32) {
	hash = 0
	for _, char := range sequence {
		switch char {
		case 'A':
			hash = hash << 2
		case 'G':
			hash = hash<<2 + 1
		case 'C':
			hash = hash<<2 + 2
		default:
			hash = hash<<2 + 3
		}
	}
	return
}
