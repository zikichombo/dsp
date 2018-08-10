// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"math"
	"math/big"
)

// nb this is unused, but can be useful if/when we add factor based
// general radix cooley tuckey implementations
var primes = []int{
	2, 3, 5, 7, 11, 13, 17, 23, 29, 31, 37, 41, 43, 47, 51, 53, 57, 59,
	61, 67, 71, 73, 79}

func init() {
	factor(101483)
}

func factor(n int) (int, int, int) {
	t := 1
	for n%2 == 0 {
		n /= 2
		t *= 2
	}
	if n == 1 || n == 3 || n == 5 || n == 7 {
		return 1, n, t
	}
	var p int
	L := int(math.Sqrt(float64(n))) + 1
	for _, p = range primes {
		if p > L {
			return 1, n, t
		}
		if n%p == 0 {
			return p, n / p, t
		}
	}
	p += 2
	x := new(big.Int)
	x.SetInt64(int64(p))
	const nb = 2
	for x.Int64() < int64(L) {
		p = int(x.Int64())
		// nb according to docs I think this should be 100% accurate in our range.
		if x.ProbablyPrime(nb) {
			//fmt.Printf("new prime %d\n", p)
			primes = append(primes, p)
			if n%p == 0 {
				return p, n / p, t
			}
		}
		if p > L {
			return 1, n, t
		}
		x.SetInt64(int64(p) + 2)
	}
	return 1, n, t
}

func log2(i int) uint {
	r := uint(0)
	for i > 1<<r {
		r++
	}
	return r
}

func is2pow(i int) bool {
	return i&(i-1) == 0
}
