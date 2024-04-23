package main

import (
	"github.com/hawkgs/wasm-fluid/fluid/js"
)

func main() {
	js.InitJsApi()

	// Prevent our program from exiting during the browser session
	select {}
}
