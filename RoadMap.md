# DSP RoadMap
The code here represents a start to a broadly featured audio DSP library.

However, there is lots of work to be done.  Below is our RoadMap for much 
of this work.

## STFT
### Analysis
#### Spectral features for analysis
1. Spectral Flux, Centroids, etc 
1. Mel Spectrogram
1. Cepstrum
### Synthesis
1. Deal with windowing correctly.
1. PV
1. Spectral whitening
1. median filter
### Gabor Frames

# Filtering
## FIR (feedforward) filtering toolbox
### Basic brick wall, band pass/stop FIR filtering design
### Implementation, using convolution for larger order filters
### Tools supporting window method design
### Tools supporting automatic design (eg Remez exchange)

## IIR (feedback) Filtering
### Biquad
Simple filter design interface and implementation (need
to decide which type(s) (II/IV?) to support.
1. Notch, 
1. Peak, 
1. Brick wall, 
1. Shelf
### Parallel and Series composition
### Butterworth, etc  (requires some general filtering tools below)

## General filtering tools (combined feedforward/feedback)
## like matlab freqz
## Pole Zero filter specification
## Bilinear transform
## Z transform and polynomial solver

## Perceptual filtering
### A weighting, etc






