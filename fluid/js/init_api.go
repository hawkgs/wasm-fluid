package js

import (
	"syscall/js"

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

		cfg := &system.SystemConfig{
			Width:     uint16(width),
			Height:    uint16(height),
			Particles: uint16(particles),
		}

		fluidSystem = system.NewSystem(cfg)

		return nil
	})

	js.Global().Get(FluidApi).Set("createFluidSystem", createFluidSystem)
}

func initRequestUpdate() {
	requestUpdate := js.FuncOf(func(this js.Value, args []js.Value) any {
		particles := fluidSystem.Update()

		js.Global().Get(FluidApi).Call("updateHandler", convertVectorsToArray(particles))

		return nil
	})

	js.Global().Get(FluidApi).Set("requestUpdate", requestUpdate)
}

func InitJsApi() {
	initCreateFluidSystem()
	initRequestUpdate()
}
