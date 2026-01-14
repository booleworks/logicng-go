package dnnf

import (
	"testing"
)

func TestNewBitset(t *testing.T) {
	t.Run("empty bitset", func(t *testing.T) {
		bs := newBitset()
		if len(bs.words) != 0 {
			t.Errorf("expected 0 words, got %d", len(bs.words))
		}
	})

	t.Run("with size", func(t *testing.T) {
		bs := newBitset(100)
		expectedWords := (100 + 63) / 64
		if len(bs.words) != expectedWords {
			t.Errorf("expected %d words, got %d", expectedWords, len(bs.words))
		}
	})

	t.Run("with size exactly 64", func(t *testing.T) {
		bs := newBitset(64)
		if len(bs.words) != 1 {
			t.Errorf("expected 1 word, got %d", len(bs.words))
		}
	})

	t.Run("with size 65", func(t *testing.T) {
		bs := newBitset(65)
		if len(bs.words) != 2 {
			t.Errorf("expected 2 words, got %d", len(bs.words))
		}
	})
}

func TestBitsetSet(t *testing.T) {
	t.Run("set single bit", func(t *testing.T) {
		bs := newBitset()
		bs.set(5)
		if !bs.get(5) {
			t.Error("bit 5 should be set")
		}
	})

	t.Run("set multiple bits in same word", func(t *testing.T) {
		bs := newBitset()
		bs.set(0)
		bs.set(10)
		bs.set(63)
		if !bs.get(0) || !bs.get(10) || !bs.get(63) {
			t.Error("bits 0, 10, 63 should be set")
		}
	})

	t.Run("set bits across multiple words", func(t *testing.T) {
		bs := newBitset()
		bs.set(0)
		bs.set(64)
		bs.set(128)
		if !bs.get(0) || !bs.get(64) || !bs.get(128) {
			t.Error("bits 0, 64, 128 should be set")
		}
		if len(bs.words) < 3 {
			t.Errorf("expected at least 3 words, got %d", len(bs.words))
		}
	})

	t.Run("set same bit twice", func(t *testing.T) {
		bs := newBitset()
		bs.set(42)
		bs.set(42)
		if !bs.get(42) {
			t.Error("bit 42 should still be set")
		}
		if bs.cardinality() != 1 {
			t.Errorf("expected cardinality 1, got %d", bs.cardinality())
		}
	})

	t.Run("set large index", func(t *testing.T) {
		bs := newBitset()
		bs.set(1000)
		if !bs.get(1000) {
			t.Error("bit 1000 should be set")
		}
	})
}

func TestBitsetGet(t *testing.T) {
	t.Run("get unset bit", func(t *testing.T) {
		bs := newBitset(100)
		if bs.get(42) {
			t.Error("bit 42 should not be set")
		}
	})

	t.Run("get bit beyond size", func(t *testing.T) {
		bs := newBitset(10)
		if bs.get(100) {
			t.Error("bit 100 should not be set")
		}
	})

	t.Run("get from empty bitset", func(t *testing.T) {
		bs := newBitset()
		if bs.get(0) {
			t.Error("bit 0 should not be set in empty bitset")
		}
	})
}

func TestBitsetOr(t *testing.T) {
	t.Run("or with empty bitset", func(t *testing.T) {
		bs1 := newBitset()
		bs1.set(5)
		bs2 := newBitset()
		bs1.or(bs2)
		if !bs1.get(5) {
			t.Error("bit 5 should still be set")
		}
	})

	t.Run("or with non-overlapping bits", func(t *testing.T) {
		bs1 := newBitset()
		bs1.set(5)
		bs2 := newBitset()
		bs2.set(10)
		bs1.or(bs2)
		if !bs1.get(5) || !bs1.get(10) {
			t.Error("bits 5 and 10 should be set")
		}
		if bs1.cardinality() != 2 {
			t.Errorf("expected cardinality 2, got %d", bs1.cardinality())
		}
	})

	t.Run("or with overlapping bits", func(t *testing.T) {
		bs1 := newBitset()
		bs1.set(5)
		bs1.set(10)
		bs2 := newBitset()
		bs2.set(5)
		bs2.set(15)
		bs1.or(bs2)
		if !bs1.get(5) || !bs1.get(10) || !bs1.get(15) {
			t.Error("bits 5, 10, 15 should be set")
		}
		if bs1.cardinality() != 3 {
			t.Errorf("expected cardinality 3, got %d", bs1.cardinality())
		}
	})

	t.Run("or across word boundaries", func(t *testing.T) {
		bs1 := newBitset()
		bs1.set(63)
		bs2 := newBitset()
		bs2.set(64)
		bs1.or(bs2)
		if !bs1.get(63) || !bs1.get(64) {
			t.Error("bits 63 and 64 should be set")
		}
	})
}

func TestBitsetAnd(t *testing.T) {
	t.Run("and with empty bitset", func(t *testing.T) {
		bs1 := newBitset()
		bs1.set(5)
		bs2 := newBitset()
		bs1.and(bs2)
		if bs1.get(5) {
			t.Error("bit 5 should be cleared after AND with empty")
		}
		if bs1.cardinality() != 0 {
			t.Errorf("expected cardinality 0, got %d", bs1.cardinality())
		}
	})

	t.Run("and with non-overlapping bits", func(t *testing.T) {
		bs1 := newBitset()
		bs1.set(5)
		bs2 := newBitset()
		bs2.set(10)
		bs1.and(bs2)
		if bs1.cardinality() != 0 {
			t.Errorf("expected cardinality 0, got %d", bs1.cardinality())
		}
	})

	t.Run("and with overlapping bits", func(t *testing.T) {
		bs1 := newBitset()
		bs1.set(5)
		bs1.set(10)
		bs1.set(15)
		bs2 := newBitset()
		bs2.set(5)
		bs2.set(15)
		bs2.set(20)
		bs1.and(bs2)
		if !bs1.get(5) || !bs1.get(15) {
			t.Error("bits 5 and 15 should be set")
		}
		if bs1.get(10) || bs1.get(20) {
			t.Error("bits 10 and 20 should not be set")
		}
		if bs1.cardinality() != 2 {
			t.Errorf("expected cardinality 2, got %d", bs1.cardinality())
		}
	})

	t.Run("and with identical bitsets", func(t *testing.T) {
		bs1 := newBitset()
		bs1.set(5)
		bs1.set(10)
		bs2 := newBitset()
		bs2.set(5)
		bs2.set(10)
		bs1.and(bs2)
		if !bs1.get(5) || !bs1.get(10) {
			t.Error("bits 5 and 10 should still be set")
		}
		if bs1.cardinality() != 2 {
			t.Errorf("expected cardinality 2, got %d", bs1.cardinality())
		}
	})
}

func TestBitsetCardinality(t *testing.T) {
	t.Run("empty bitset", func(t *testing.T) {
		bs := newBitset()
		if bs.cardinality() != 0 {
			t.Errorf("expected cardinality 0, got %d", bs.cardinality())
		}
	})

	t.Run("single bit", func(t *testing.T) {
		bs := newBitset()
		bs.set(42)
		if bs.cardinality() != 1 {
			t.Errorf("expected cardinality 1, got %d", bs.cardinality())
		}
	})

	t.Run("multiple bits in same word", func(t *testing.T) {
		bs := newBitset()
		for i := range int32(64) {
			bs.set(i)
		}
		if bs.cardinality() != 64 {
			t.Errorf("expected cardinality 64, got %d", bs.cardinality())
		}
	})

	t.Run("bits across multiple words", func(t *testing.T) {
		bs := newBitset()
		bs.set(0)
		bs.set(63)
		bs.set(64)
		bs.set(127)
		bs.set(128)
		if bs.cardinality() != 5 {
			t.Errorf("expected cardinality 5, got %d", bs.cardinality())
		}
	})
}

func TestBitsetNextSetBit(t *testing.T) {
	t.Run("empty bitset", func(t *testing.T) {
		bs := newBitset()
		if bs.nextSetBit(0) != -1 {
			t.Error("expected -1 for empty bitset")
		}
	})

	t.Run("single bit at start", func(t *testing.T) {
		bs := newBitset()
		bs.set(0)
		if bs.nextSetBit(0) != 0 {
			t.Errorf("expected 0, got %d", bs.nextSetBit(0))
		}
		if bs.nextSetBit(1) != -1 {
			t.Error("expected -1 when searching after the only set bit")
		}
	})

	t.Run("iterate through all bits", func(t *testing.T) {
		bs := newBitset()
		bs.set(5)
		bs.set(10)
		bs.set(63)
		bs.set(64)
		bs.set(100)

		expected := []int32{5, 10, 63, 64, 100}
		index := 0
		for i := bs.nextSetBit(0); i != -1; i = bs.nextSetBit(i + 1) {
			if index >= len(expected) {
				t.Fatal("found more bits than expected")
			}
			if i != expected[index] {
				t.Errorf("expected bit %d, got %d", expected[index], i)
			}
			index++
		}
		if index != len(expected) {
			t.Errorf("expected %d bits, found %d", len(expected), index)
		}
	})

	t.Run("start search from middle", func(t *testing.T) {
		bs := newBitset()
		bs.set(5)
		bs.set(10)
		bs.set(15)
		if bs.nextSetBit(7) != 10 {
			t.Errorf("expected 10, got %d", bs.nextSetBit(7))
		}
	})

	t.Run("start search from set bit", func(t *testing.T) {
		bs := newBitset()
		bs.set(10)
		if bs.nextSetBit(10) != 10 {
			t.Errorf("expected 10, got %d", bs.nextSetBit(10))
		}
	})

	t.Run("word boundary", func(t *testing.T) {
		bs := newBitset()
		bs.set(63)
		bs.set(64)
		if bs.nextSetBit(0) != 63 {
			t.Errorf("expected 63, got %d", bs.nextSetBit(0))
		}
		if bs.nextSetBit(64) != 64 {
			t.Errorf("expected 64, got %d", bs.nextSetBit(64))
		}
	})

	t.Run("start beyond all bits", func(t *testing.T) {
		bs := newBitset()
		bs.set(10)
		if bs.nextSetBit(100) != -1 {
			t.Error("expected -1 when starting beyond all set bits")
		}
	})
}

func TestBitsetClear(t *testing.T) {
	t.Run("clear empty bitset", func(t *testing.T) {
		bs := newBitset()
		bs.clear()
		if bs.cardinality() != 0 {
			t.Errorf("expected cardinality 0, got %d", bs.cardinality())
		}
	})

	t.Run("clear bitset with bits", func(t *testing.T) {
		bs := newBitset()
		bs.set(5)
		bs.set(64)
		bs.set(128)
		bs.clear()
		if bs.cardinality() != 0 {
			t.Errorf("expected cardinality 0, got %d", bs.cardinality())
		}
		if bs.get(5) || bs.get(64) || bs.get(128) {
			t.Error("all bits should be cleared")
		}
	})

	t.Run("set after clear", func(t *testing.T) {
		bs := newBitset()
		bs.set(10)
		bs.clear()
		bs.set(20)
		if bs.get(10) {
			t.Error("bit 10 should still be cleared")
		}
		if !bs.get(20) {
			t.Error("bit 20 should be set")
		}
		if bs.cardinality() != 1 {
			t.Errorf("expected cardinality 1, got %d", bs.cardinality())
		}
	})
}

func TestBitsetClone(t *testing.T) {
	t.Run("clone empty bitset", func(t *testing.T) {
		bs := newBitset()
		clone := bs.clone()
		if clone.cardinality() != 0 {
			t.Errorf("expected cardinality 0, got %d", clone.cardinality())
		}
	})

	t.Run("clone with bits", func(t *testing.T) {
		bs := newBitset()
		bs.set(5)
		bs.set(64)
		bs.set(128)
		clone := bs.clone()

		if !clone.get(5) || !clone.get(64) || !clone.get(128) {
			t.Error("clone should have same bits set")
		}
		if clone.cardinality() != 3 {
			t.Errorf("expected cardinality 3, got %d", clone.cardinality())
		}
	})

	t.Run("clone is independent", func(t *testing.T) {
		bs := newBitset()
		bs.set(5)
		clone := bs.clone()

		// Modify original
		bs.set(10)
		if clone.get(10) {
			t.Error("clone should not be affected by changes to original")
		}

		// Modify clone
		clone.set(15)
		if bs.get(15) {
			t.Error("original should not be affected by changes to clone")
		}
	})

	t.Run("clone after clear", func(t *testing.T) {
		bs := newBitset()
		bs.set(5)
		bs.clear()
		clone := bs.clone()
		if clone.cardinality() != 0 {
			t.Errorf("expected cardinality 0, got %d", clone.cardinality())
		}
	})
}

func TestBitsetEnsureSize(t *testing.T) {
	t.Run("grow from empty", func(t *testing.T) {
		bs := newBitset()
		bs.set(100) // Should trigger ensureSize
		if !bs.get(100) {
			t.Error("bit 100 should be set")
		}
	})

	t.Run("grow multiple times", func(t *testing.T) {
		bs := newBitset()
		bs.set(10)
		bs.set(100)
		bs.set(1000)
		if !bs.get(10) || !bs.get(100) || !bs.get(1000) {
			t.Error("all bits should be set")
		}
		if bs.cardinality() != 3 {
			t.Errorf("expected cardinality 3, got %d", bs.cardinality())
		}
	})
}

func TestBitsetComplex(t *testing.T) {
	t.Run("multiple operations", func(t *testing.T) {
		bs1 := newBitset()
		bs2 := newBitset()

		// Set some bits
		for i := int32(0); i < 100; i += 10 {
			bs1.set(i)
		}
		for i := int32(5); i < 100; i += 10 {
			bs2.set(i)
		}

		// OR them
		bs1.or(bs2)
		if bs1.cardinality() != 20 {
			t.Errorf("expected cardinality 20, got %d", bs1.cardinality())
		}

		// Clear and AND
		bs3 := newBitset()
		for i := int32(0); i < 50; i += 10 {
			bs3.set(i)
		}
		bs1.and(bs3)
		if bs1.cardinality() != 5 {
			t.Errorf("expected cardinality 5, got %d", bs1.cardinality())
		}
	})

	t.Run("clone and modify", func(t *testing.T) {
		original := newBitset()
		for i := int32(0); i < 64; i++ {
			original.set(i)
		}

		clone1 := original.clone()
		clone2 := original.clone()

		clone1.clear()
		clone1.set(100)

		clone2.set(200)

		if original.cardinality() != 64 {
			t.Errorf("original cardinality should be 64, got %d", original.cardinality())
		}

		if clone1.cardinality() != 1 || !clone1.get(100) {
			t.Error("clone1 should only have bit 100 set")
		}
		if clone2.cardinality() != 65 || !clone2.get(200) {
			t.Error("clone2 should have original bits plus bit 200")
		}
	})
}
