package js

import (
	"syscall/js"

	"github.com/hawkgs/wasm-fluid/fluid/system"
)

// Name of the global JS object that is going to be attached to `globalThis` (i.e. `window` in our case)
const FluidApi = "FluidApi"

// Keeps the current runtime's fluid system
var fluidSystem *system.System

// initCreateFluidSystem exposes fluid system creation to our JS context
func initCreateFluidSystem() {
	// Creates a fluid system by the provided config params
	createFluidSystem := js.FuncOf(func(this js.Value, args []js.Value) any {
		jsCfg := args[0]
		width := jsCfg.Get("width").Int()
		height := jsCfg.Get("height").Int()
		particles := jsCfg.Get("particles").Int()
		particleUiRadius := jsCfg.Get("particleUiRadius").Int()

		// Build the config and initialize the system
		cfg := system.NewSystemConfig(width, height, particles, particleUiRadius)
		fluidSystem = system.NewSystem(cfg)

		return nil
	})

	js.Global().Get(FluidApi).Set("createFluidSystem", createFluidSystem)
}

// initRequestUpdate exposes fluid system updating to our JS context
func initRequestUpdate() {
	requestUpdate := js.FuncOf(func(this js.Value, args []js.Value) any {
		particles := fluidSystem.Update()

		js.Global().Get(FluidApi).Call("updateHandler", convertParticlesToJsArray(particles))

		return nil
	})

	js.Global().Get(FluidApi).Set("requestUpdate", requestUpdate)
}

// initDevPrintSystemStats is used for debugging purposes
func initDevPrintSystemStats() {
	devPrintSystemStats := js.FuncOf(func(this js.Value, args []js.Value) any {
		fluidSystem.DevPrintStats()

		return nil
	})

	js.Global().Get(FluidApi).Set("devPrintSystemStats", devPrintSystemStats)
}

// InitJsApi initializes the whole JS API
func InitJsApi() {
	initCreateFluidSystem()
	initRequestUpdate()
	initDevPrintSystemStats()
}
