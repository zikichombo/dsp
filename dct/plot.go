// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package dct

import (
	"image"
	"image/color"
	"math"
)

func minMax(d []float64) (min float64, max float64) {
	min = math.Inf(1)
	max = math.Inf(-1)
	for _, v := range d {
		if v < min {
			min = v
		}
		if v >= max {
			max = v
		}
	}
	return
}

func Plot(b image.Rectangle, d []float64) *image.RGBA {
	mAcc := 0.0
	wAcc := 0.0
	rAcc := 0.0
	im := image.NewRGBA(b)
	x := 0
	black := color.RGBA{A: 255}
	for x := b.Min.X; x < b.Max.X; x++ {
		im.Set(x, b.Min.Y, black)
		im.Set(x, b.Max.Y-1, black)
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		im.Set(b.Min.X, y, black)
		im.Set(b.Max.X-1, y, black)
	}
	bb := image.Rect(b.Min.X+1, b.Min.Y+1, b.Max.X-1, b.Max.Y-1)
	subIm := im.SubImage(bb).(*image.RGBA)
	r := float64(bb.Dx()) / float64(len(d))

	color := color.RGBA{
		A: 180,
		R: 0,
		G: 200,
		B: 255}

	min, max := minMax(d)
	if min == 0 {
		min = 1e-10
	}

	//minDb := sample.ToDb(min)
	//maxDb := sample.ToDb(max)
	minDb, maxDb := min, max
	for j := 0; j < len(d); j++ {
		//mdb := sample.ToDb(d[j])
		mdb := d[j]
		rAcc += r
		if rAcc < 1.0 {
			mAcc += mdb * r
			wAcc += r
			continue
		}
		_, ar := math.Modf(rAcc)
		mAcc += mdb * (r - ar)
		wAcc += r - ar
		v := mAcc / wAcc
		yr := (v - minDb) / (maxDb - minDb)
		Y := int(math.Floor(yr*float64(bb.Dy()) + 0.5))
		for rAcc >= 1.0 {
			for y := 1; y <= Y; y++ {
				subIm.Set(b.Min.X+x, b.Max.Y-y, color)
			}
			rAcc -= 1.0
			x++
		}
		mAcc = mdb * ar
		rAcc = ar
		wAcc = ar
	}
	return im
}
