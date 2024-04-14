package system

type SystemConfig struct {
	Width            float64
	Height           float64
	Particles        uint
	ParticleUiRadius float64
}

func NewSystemConfig(width int, height int, particles int, particleUiRadius int) *SystemConfig {
	return &SystemConfig{
		Width:            float64(width) / float64(SystemScale),
		Height:           float64(height) / float64(SystemScale),
		Particles:        uint(particles),
		ParticleUiRadius: float64(particleUiRadius),
	}
}
