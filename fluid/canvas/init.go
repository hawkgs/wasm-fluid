package canvas

import (
	"syscall/js"

	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

func InitCanvas() {
	v := vectors.NewVector(10, 20)
	js.Global().Call("InitCanvas", v.ToMap())
}
