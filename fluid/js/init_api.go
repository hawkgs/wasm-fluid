package js

import (
	"syscall/js"

	"github.com/hawkgs/wasm-fluid/fluid/forces"
	"github.com/hawkgs/wasm-fluid/fluid/system"
)

const FluidApi = "FluidApi"

var fluidSystem *system.System

func initCreateFluidSystem() {
	createFluidSystem := js.FuncOf(func(this js.Value, args []js.Value) any {
		jsCfg := args[0]
		width := jsCfg.Get("width").Int()
		height := jsCfg.Get("height").Int()
		particles := jsCfg.Get("particles").Int()
		particleUiRadius := jsCfg.Get("particleUiRadius").Int()

		// Create config
		cfg := &system.SystemConfig{
			Width:            uint(width),
			Height:           uint(height),
			Particles:        uint(particles),
			ParticleUiRadius: uint(particleUiRadius),
		}

		// Create initial forces
		forces := []forces.Force{
			forces.NewGravity(),
		}

		fluidSystem = system.NewSystem(cfg, forces)

		return nil
	})

	js.Global().Get(FluidApi).Set("createFluidSystem", createFluidSystem)
}

func initRequestUpdate() {
	requestUpdate := js.FuncOf(func(this js.Value, args []js.Value) any {
		particles := fluidSystem.Update()

		js.Global().Get(FluidApi).Call("updateHandler", convertParticlesToPositionsArray(particles))

		return nil
	})

	js.Global().Get(FluidApi).Set("requestUpdate", requestUpdate)
}

func InitJsApi() {
	initCreateFluidSystem()
	initRequestUpdate()
}
