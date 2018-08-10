// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Package dct provides discrete cosine transform support.
//
// dct implements a naive O(n*n) algorithm for reference and
// testing and also an O(n log n) algorithm by Byeong Gi Lee (1984).
//
// see http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.118.3056&rep=rep1&type=pdf#page=34
// and https://www.nayuki.io/page/fast-discrete-cosine-transform-algorithms
package dct /* import "zikichombo.org/dsp/dct" */
