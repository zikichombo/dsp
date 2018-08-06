// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package fft

import (
	"fmt"
	"testing"

	"zikichombo.org/sound/freq"
)

type tstCfg struct {
	n    int
	s, f freq.T
}

var inDat = []tstCfg{
	{100, 100 * freq.Hertz, freq.Hertz},
	{100, 100 * freq.Hertz, freq.Hertz},
	{100, 100 * freq.Hertz, 2 * freq.Hertz},
	{100, 100 * freq.Hertz, 10 * freq.Hertz},
	{100, 100 * freq.Hertz, 25 * freq.Hertz}}

var alignDat = []tstCfg{
	{1000, freq.Hertz, 2 * freq.Hertz / 3},
	{100, 100 * freq.Hertz, 26*freq.Hertz - freq.MilliHertz},
	{1000, 44100 * freq.Hertz, 440 * freq.Hertz}}

func TestFreqBin(t *testing.T) {
	for i := 0; i < len(inDat); i += 2 {
		cfg := &inDat[i]
		b := FreqBin(cfg.s, cfg.f, cfg.n)
		l, u := BinRange(cfg.s, cfg.n, b)
		fmt.Printf("s=%s n=%d f=%s b=%d [%s, %s)\n", cfg.s, cfg.n, cfg.f, b, l, u)
	}
	for i := 0; i < len(alignDat); i++ {
		cfg := &alignDat[i]
		b := FreqBin(cfg.s, cfg.f, cfg.n)
		l, u := BinRange(cfg.s, cfg.n, b)
		fmt.Printf("s=%s n=%d f=%s b=%d [%s, %s)\n", cfg.s, cfg.n, cfg.f, b, l, u)
	}
}

func TestWinSize(t *testing.T) {
	for _, err := range []freq.T{1 * freq.MicroHertz, 10 * freq.MicroHertz, 100 * freq.MicroHertz, freq.MilliHertz, 10 * freq.MilliHertz} {
		for i := 0; i < len(inDat); i += 2 {
			cfg := &inDat[i]
			_, c := WinSize(cfg.s, cfg.f, err)
			if c != 1 {
				t.Errorf("wrong number of cyles: %d\n", c)
			}
			//fmt.Printf("%s: find signal of %s in %d samples (%d cylcles) at %s\n", err, cfg.f, sz, c, cfg.s)
		}
		for i := 0; i < len(alignDat); i++ {
			cfg := &alignDat[i]
			_, c := WinSize(cfg.s, cfg.f, err)
			if (freq.T(c)*cfg.s)%cfg.f > err {
				t.Errorf("bad error")
			}
			//fmt.Printf("%s: find signal of %s in %d samples (%d cycles) at %s\n", err, cfg.f, sz, c, cfg.s)
		}
	}
}
