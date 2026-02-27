package httpmux

import "testing"

func TestFastRandIntRange(t *testing.T) {
	bounds := []int{1, 2, 3, 7, 10, 63, 64, 127, 255, 1024, 1 << 20}
	for _, n := range bounds {
		for i := 0; i < 10000; i++ {
			v := fastRandInt(n)
			if v < 0 || v >= n {
				t.Fatalf("out of range: n=%d v=%d", n, v)
			}
		}
	}
}

func TestFastRandIntEdgeCases(t *testing.T) {
	if v := fastRandInt(0); v != 0 {
		t.Fatalf("expected 0 for n=0, got %d", v)
	}
	if v := fastRandInt(-1); v != 0 {
		t.Fatalf("expected 0 for n=-1, got %d", v)
	}
}

func TestFastRandIntDistributionSanity(t *testing.T) {
	const (
		n       = 8
		samples = 200000
	)
	var counts [n]int
	for i := 0; i < samples; i++ {
		counts[fastRandInt(n)]++
	}
	expected := samples / n
	tolerance := expected / 10 // ±10%
	for i, c := range counts {
		delta := c - expected
		if delta < 0 {
			delta = -delta
		}
		if delta > tolerance {
			t.Fatalf("distribution skew too high for bucket %d: got=%d expected≈%d tolerance=%d", i, c, expected, tolerance)
		}
	}
}

func BenchmarkFastRandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fastRandInt(8192)
	}
}
