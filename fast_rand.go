package httpmux

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/bits"
	"sync/atomic"
	"time"
)

var fastRandState atomic.Uint64

func init() {
	seed := seedFastRand()
	fastRandState.Store(seed)
}

func seedFastRand() uint64 {
	var b [8]byte
	if _, err := cryptorand.Read(b[:]); err == nil {
		seed := binary.LittleEndian.Uint64(b[:])
		if seed != 0 {
			return seed
		}
	}
	seed := uint64(time.Now().UnixNano())
	if seed == 0 {
		seed = 0x9e3779b97f4a7c15
	}
	return seed
}

func nextFastRand64() uint64 {
	// SplitMix64 sequence: fast and suitable for non-crypto randomization.
	state := fastRandState.Add(0x9e3779b97f4a7c15)
	z := state
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

// fastRandInt returns a value in [0, n). It is safe for concurrent use.
// Uses Lemire's unbiased reduction to avoid modulo bias.
func fastRandInt(n int) int {
	if n <= 0 {
		return 0
	}
	bound := uint64(n)

	for {
		x := nextFastRand64()
		hi, lo := bits.Mul64(x, bound)
		if lo < bound {
			threshold := uint64(0) - bound
			threshold %= bound
			if lo < threshold {
				continue
			}
		}
		return int(hi)
	}
}
