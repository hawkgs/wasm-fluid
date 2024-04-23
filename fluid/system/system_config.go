package system

type SystemConfig struct {
	Width            float64
	Height           float64
	Particles        uint
	ParticleUiRadius float64
}

func NewSystemConfig(width int, height int, particles int, particleUiRadius int) *SystemConfig {
	return &SystemConfig{
		Width:            float64(width) / SystemScale,
		Height:           float64(height) / SystemScale,
		Particles:        uint(particles),
		ParticleUiRadius: float64(particleUiRadius) / SystemScale,
	}
}
