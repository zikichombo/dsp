// Copyright 2017 The ZikiChombo Authors. All rights reserved.  Use of
// this source code is governed by a license that can be found in the License
// file.

// Package resample implements resampling/changes in resolution.
//
// Package resample uses interpolation for resampling which provides easy
// control over the quality/cost tradeoff and can produce very high quality
// resampling.  Other resampling methods may be more appropriate for a given
// calling context, package resample doesn't yet provide other mechanisms.
//
// When resampling audio, any decrease in sample rate from rate S to a rate R
// must be applied to a signal which does not contain frequencies at or above
// R/2, or aliasing will produce strange results.
//
// This is often achieved by first applying a low pass filter and then
// resampling.  As ZikiChombo does not yet have filter design support, we
// recommend in the meantime simply taking a moving average of the signal with
// a window size W = ceil(S/R) before decreasing the sample rate if you do not
// have access to or knowledge about low pass filtering design.
//
// BUG(wsc) the shift size, effecting interpolation order limits and
// latency of implementations is constant (64 frames).
package resample /* import "zikichombo.org/dsp/resample" */
