// Copyright 2018 The ZikiChombo Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

// Package lpc provides a linear predictive coding interface.
//
// The LPC modelling algorithm is based on the autocorrelation method with the
// addition of a numerical tweak to enforce stability of the resulting model.
//
// Package lpc supports modelling, predicting and generation/synthesis.
//
// Package lpc does not yet support line spectral frequencies or
// conversion to other coefficient representations.
//
// Package lpc is part of http://zikichombo.org
package lpc /* import "zikichombo.org/dsp/lpc" */
