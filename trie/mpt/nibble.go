package mpt

// Nibble represents a 4-bit half-byte used in trie keys.
type Nibble byte

// convertToNibbles converts a byte slice to a nibble slice.
// Each byte is split into two nibbles (high and low).
func convertToNibbles(key []byte) []Nibble {
	nibbles := make([]Nibble, 0, len(key)*2)
	for _, b := range key {
		nibbles = append(nibbles, Nibble(b>>4))   // 高四位
		nibbles = append(nibbles, Nibble(b&0x0F)) // 低四位
	}
	return nibbles
}

// convertToBytes converts a nibble slice back to a byte slice.
// Pairs of adjacent nibbles are combined into bytes.
func convertToBytes(nibbles []Nibble) []byte {
	if len(nibbles)%2 != 0 {
		panic("nibble slice length must be even")
	}

	bytes := make([]byte, len(nibbles)/2)
	for i := 0; i < len(bytes); i++ {
		high := nibbles[i*2] << 4   // Shift high nibble to bits 4-7
		low := nibbles[i*2+1]       // Low nibble stays in bits 0-3
		bytes[i] = byte(high | low) // Combine into one byte
	}
	return bytes
}

func nibbleToBytes(nibbles []Nibble) []byte {
	if len(nibbles)%2 != 0 {
		nibbles = append(nibbles, 0x0)
	}
	out := make([]byte, len(nibbles)/2)
	for i := 0; i < len(out); i++ {
		out[i] = byte(nibbles[2*i]<<4 | nibbles[2*i+1])
	}
	return out
}
