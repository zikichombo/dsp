// Copyright 2018 The ZikiChomgo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package dct

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestZ(t *testing.T) {
	N := 1024
	d := make([]float64, N)
	z := NewZ(N)
	for i := range d {
		d[i] = 0.99*math.Cos(float64(i)*0.0234) + 0.01*(rand.Float64()-0.5)
	}
	ct := New(N)
	ct.Do(d)
	orgPwr := z.Init(d)
	t.Logf("initial power %f\n", orgPwr)
	i, ci := 0, 0
	var lpr, pr, r float64
	for i < N {
		ci, pr, r = z.Top()
		if r < 0.2 {
			break
		}
		t.Logf("%d: %d %f r %f rate %f\n", i, ci, d[ci], pr, r)
		z.Pop()
		i++
		lpr = pr
	}
	z.Zero()
	pwr := 0.0
	for _, v := range d {
		pwr += v * v
	}
	pwr /= float64(len(d))
	pwr = math.Sqrt(pwr)
	if math.Abs(pwr/orgPwr-lpr) > 1e-10 {
		t.Errorf("final power %f ratio %f expected %f\n", pwr, pwr/orgPwr, lpr)
	}
}

func TestProblem(t *testing.T) {
	var d = []float64{2.126846609940003e-08,
		-1.9429238307111518e-07, 4.0083187968775746e-07, -3.18344945071658e-07,
		-4.189423066236486e-07, 1.9653280105558224e-06, -3.877667495544301e-06,
		5.364339813240804e-06, -5.428420081443619e-06, 2.944657808257034e-06,
		2.1607013422908494e-06, -8.945888112066314e-06, 1.6051806596806273e-05,
		-2.084455445583444e-05, 2.093829789373558e-05, -1.5789008102728985e-05,
		5.26436133441166e-06, 8.81815685715992e-06, -2.2116873878985643e-05,
		3.1580013455823064e-05, -3.43189385603182e-05, 2.7483094527269714e-05,
		-1.2780014913005289e-05, -6.7707242124015465e-06, 2.8567645131261088e-05,
		-4.6346620365511626e-05, 5.5720924137858674e-05, -5.617501665255986e-05,
		4.532853563432582e-05, -2.4820908947731368e-05, 7.29386329112458e-07,
		2.471469997544773e-05, -4.663503204938024e-05, 5.795817924081348e-05,
		-5.905329089728184e-05, 4.895852543995716e-05, -2.568614399933722e-05,
		-3.848707819997799e-06, 3.5588047467172146e-05, -6.814824155298993e-05,
		9.227096597896889e-05, -0.00010481775098014623, 0.00010734285751823336,
		-9.362259152112529e-05, 6.522569310618564e-05, -2.8237169317435473e-05,
		-2.0315472283982672e-05, 7.428289973177016e-05, -0.00012250729196239263,
		0.0001633717183722183, -0.00018472722149454057, 0.00017164868768304586,
		-0.0001244281738763675, 3.862930680043064e-05, 8.390571747440845e-05,
		-0.00021764192206319422, 0.00034023079206235707, -0.0004265505413059145,
		0.00043920031748712063, -0.00036495301173999906, 0.0002092648355755955,
		1.6686981325619854e-05, -0.0002667483640834689, 0.0004822083574254066}
	for i := range d {
		d[i] *= 1e6
	}
	ct := New(len(d))
	e := make([]float64, len(d))
	copy(e, d)
	ct.Do(d)
	z := NewZ(len(d))
	z.Init(d)
	last := math.Inf(1)
	for i := range d {
		ci, r, s := z.Top()
		fmt.Printf("%d ratio %f speed %f\n", i, r, s)
		if r >= 0.95 {
			break
		}
		if math.Abs(d[ci]) > last {
			t.Errorf("%d: got %f > %f", i, d[ci], last)
		}
		last = math.Abs(d[ci])
		z.Pop()
	}
}