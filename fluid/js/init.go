package js

import (
	"syscall/js"

	"github.com/hawkgs/wasm-fluid/fluid/system"
)

const GoApi = "GoApi"

var fluidSystem *system.System

func initCreateFluidSystem() {
	createFluidSystem := js.FuncOf(func(this js.Value, args []js.Value) any {
		jsCfg := args[0]
		width := jsCfg.Get("width").Int()
		height := jsCfg.Get("height").Int()

		cfg := &system.SystemConfig{Width: uint16(width), Height: uint16(height)}
		fluidSystem = system.NewSystem(cfg)

		return nil
	})

	js.Global().Get(GoApi).Set("goCreateFluidSystem", createFluidSystem)
}

func initRequestUpdate() {
	requestUpdate := js.FuncOf(func(this js.Value, args []js.Value) any {
		particles := fluidSystem.Update()
		js.Global().Get(GoApi).Call("goUpdateHandler", particles[0].ToMap())
		return nil
	})

	js.Global().Get(GoApi).Set("goRequestUpdate", requestUpdate)
}

func InitJsApi() {
	initCreateFluidSystem()
	initRequestUpdate()
}
