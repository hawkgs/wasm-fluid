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
		cfg := createConfigObject(args[0])
		fluidSystem = system.NewSystem(cfg)

		return nil
	})

	js.Global().Get(FluidApi).Set("createFluidSystem", createFluidSystem)
}

// initRequestUpdate exposes fluid system updating to our JS context
func initRequestUpdate() {
	requestUpdate := js.FuncOf(func(this js.Value, args []js.Value) any {
		particles := fluidSystem.Update()

		js.Global().Get(FluidApi).Call(
			"updateHandler",
			convertParticlesToJsArray(particles, fluidSystem.GetConfig().SystemScale),
		)

		return nil
	})

	js.Global().Get(FluidApi).Set("requestUpdate", requestUpdate)
}

// initUpdateDynamicParams exposes dynamic parameters updating to our JS context
func initUpdateDynamicParams() {
	updateDynamicParams := js.FuncOf(func(this js.Value, args []js.Value) any {
		obj := args[0]
		particleMass := obj.Get("particleMass").Float()
		gravityForce := obj.Get("gravityForce").Float()
		gasConstK := obj.Get("gasConstK").Float()
		restDensity := obj.Get("restDensity").Float()
		viscosityConst := obj.Get("viscosityConst").Float()
		velocityLimit := obj.Get("velocityLimit").Float()
		collisionDamping := obj.Get("collisionDamping").Float()

		fluidSystem.UpdateDynamicParams(&system.DynamicParams{
			ParticleMass:     particleMass,
			GravityForce:     gravityForce,
			GasConstK:        gasConstK,
			RestDensity:      restDensity,
			ViscosityConst:   viscosityConst,
			VelocityLimit:    velocityLimit,
			CollisionDamping: collisionDamping,
		})

		return nil
	})

	js.Global().Get(FluidApi).Set("updateDynamicParams", updateDynamicParams)
}

// initDevPrintSystemStats prints a snapshot of the system's state (used for debugging); Public/JS
func initDevPrintSystemStats() {
	devPrintSystemStats := js.FuncOf(func(this js.Value, args []js.Value) any {
		fluidSystem.DevPrintStats()

		return nil
	})

	js.Global().Get(FluidApi).Set("devPrintSystemStats", devPrintSystemStats)
}

// createConfigObject creates a SystemConfig based on the provided JS object
func createConfigObject(obj js.Value) *system.SystemConfig {
	width := obj.Get("width").Int()
	height := obj.Get("height").Int()
	particles := obj.Get("particles").Int()
	particleUiRadius := obj.Get("particleUiRadius").Float()
	systemScale := obj.Get("systemScale").Float()

	smoothingRadiusH := obj.Get("smoothingRadiusH").Float()
	timestep := obj.Get("timestep").Float()

	particleMass := obj.Get("particleMass").Float()
	gravityForce := obj.Get("gravityForce").Float()
	gasConstK := obj.Get("gasConstK").Float()
	restDensity := obj.Get("restDensity").Float()
	viscosityConst := obj.Get("viscosityConst").Float()
	velocityLimit := obj.Get("velocityLimit").Float()
	collisionDamping := obj.Get("collisionDamping").Float()

	return system.NewSystemConfig(
		width,
		height,
		particles,
		particleUiRadius,

		systemScale,
		smoothingRadiusH,
		timestep,

		particleMass,
		gravityForce,
		gasConstK,
		restDensity,
		viscosityConst,
		velocityLimit,
		collisionDamping,
	)
}

// InitJsApi initializes the whole JS API
func InitJsApi() {
	initCreateFluidSystem()
	initRequestUpdate()
	initUpdateDynamicParams()
	initDevPrintSystemStats()
}
