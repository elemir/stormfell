package algo

type Noise interface {
	Noise(x, y float64) float64
}

type FractalNoise struct {
	noise   Noise
	octaves uint8
	persist float64
}

func NewFractalNoise(noise Noise, octaves uint8, persist float64) *FractalNoise {
	return &FractalNoise{
		noise:   noise,
		octaves: octaves,
		persist: persist,
	}
}

func (f *FractalNoise) Noise(x, y float64) float64 {
	maxAmplitude := 0.0
	amplitude := 1.0
	frequency := 1.0
	noise := 0.0

	for range f.octaves {
		noise += amplitude * f.noise.Noise(x*frequency, y*frequency)
		maxAmplitude += amplitude
		amplitude *= f.persist
		frequency *= 2
	}

	return noise / maxAmplitude
}
