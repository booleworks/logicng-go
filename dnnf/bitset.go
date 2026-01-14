package dnnf

import (
	"math/bits"
)

type bitset struct {
	words []uint64
}

func newBitset(size ...int32) *bitset {
	capacity := int32(0)
	if len(size) > 0 {
		capacity = size[0]
	}
	numWords := (capacity + 63) / 64
	return &bitset{make([]uint64, numWords)}
}

func (b *bitset) set(index int32) {
	wordIndex := index / 64
	bitIndex := index % 64
	if int(wordIndex) >= len(b.words) {
		b.ensureSize(int(index) + 1)
	}
	b.words[wordIndex] |= (1 << bitIndex)
}

func (b *bitset) get(index int) bool {
	wordIndex := index / 64
	bitIndex := index % 64
	return wordIndex < len(b.words) && (b.words[wordIndex]&(1<<bitIndex)) != 0
}

func (b *bitset) or(other *bitset) {
	b.ensureWords(len(other.words))
	for i := 0; i < len(other.words); i++ {
		b.words[i] |= other.words[i]
	}
}

func (b *bitset) and(other *bitset) {
	minLen := min(len(b.words), len(other.words))
	for i := range minLen {
		b.words[i] &= other.words[i]
	}
	for i := minLen; i < len(b.words); i++ {
		b.words[i] = 0
	}
}

func (b *bitset) ensureSize(bitSize int) {
	requiredWords := (bitSize + 63) / 64
	if len(b.words) >= requiredWords {
		return
	}
	newWords := make([]uint64, requiredWords)
	copy(newWords, b.words)
	b.words = newWords
}

func (b *bitset) ensureWords(numWords int) {
	if len(b.words) >= numWords {
		return
	}
	newWords := make([]uint64, numWords)
	copy(newWords, b.words)
	b.words = newWords
}

func (b *bitset) cardinality() int {
	count := 0
	for _, word := range b.words {
		count += bits.OnesCount64(word)
	}
	return count
}

func (b *bitset) nextSetBit(fromIndex int32) int32 {
	wordIndex := int32(fromIndex / 64)
	if int(wordIndex) >= len(b.words) {
		return -1
	}

	word := b.words[wordIndex]
	bitIndex := fromIndex % 64
	word &= (^uint64(0)) << bitIndex

	if word != 0 {
		trailingZeros := bits.TrailingZeros64(word)
		return wordIndex*64 + int32(trailingZeros)
	}

	for i := wordIndex + 1; i < int32(len(b.words)); i++ {
		if b.words[i] != 0 {
			trailingZeros := bits.TrailingZeros64(b.words[i])
			return i*64 + int32(trailingZeros)
		}
	}

	return -1
}

func (b *bitset) toBoolSlice() []bool {
	if len(b.words) == 0 {
		return []bool{}
	}
	maxBit := 0
	for i := len(b.words) - 1; i >= 0; i-- {
		if b.words[i] != 0 {
			maxBit = i*64 + 63 - bits.LeadingZeros64(b.words[i])
			break
		}
	}
	result := make([]bool, maxBit+1)
	for i := 0; i <= maxBit; i++ {
		result[i] = b.get(i)
	}
	return result
}

func (b *bitset) clear() {
	for i := range b.words {
		b.words[i] = 0
	}
}

func (b *bitset) clone() *bitset {
	cln := make([]uint64, len(b.words))
	copy(cln, b.words)
	return &bitset{cln}
}
