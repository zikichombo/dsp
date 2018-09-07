// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package convol

import (
	"fmt"
	"math/rand"
	"testing"
)

func genLong() []float64 {
	n := rand.Intn(128) + 128
	res := make([]float64, n)
	for i := range res {
		a := float64(rand.Intn(10))
		res[i] = a
	}
	return res
}

func TestOla(t *testing.T) {
	// we just compare ola to k on random inputs.
	for i := 0; i < 64; i++ {
		krn, _ := gen()
		seq := genLong()
		blk := len(seq) / len(krn)
		if blk < 2*len(krn) {
			fmt.Printf("skipping %d, seq not long enough", i)
		}
		ola := NewOla(krn, blk)
		k := NewK(krn, len(seq))
		kRes, e := k.ConvTo(nil, seq)
		if e != nil {
			t.Error(e)
			continue
		}
		kRes = kRes[:len(seq)]
		blkWin := ola.WinSrc(nil)
		dstWin := ola.WinDst(nil)
		n := 0
		for i := 0; i < len(krn)-1; i++ {
			seq = append(seq, 0i)
		}
		//fmt.Printf("ola on seq %d with kernel %d, kRes len is %d\n", len(seq), len(krn), len(kRes))
		for n < len(kRes) {
			end := n + ola.N()
			copy(blkWin, seq[n:end])
			if e := ola.Block(blkWin, dstWin); e != nil {
				t.Fatalf("%d ola error: %s\n", i, e)
			}
			resEnd := end
			if resEnd > len(kRes) {
				resEnd = len(kRes)
			}
			if k := approxEq(kRes[n:resEnd], dstWin[:resEnd-n], 0.001); k != -1 {
				t.Errorf("ola run %d: data error at %d. %.2f v %.2f", i, n+k, kRes[n+k], dstWin[k])
			}
			n += resEnd - n
			//fmt.Printf("\tprocessed %d\n", n)
		}
	}
}

func TestOlaSource(t *testing.T) {
}
