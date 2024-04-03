package main

import (
	"fmt"

	"github.com/hawkgs/wasm-fluid/fluid/js"
)

func main() {
	fmt.Println("Hello, Go!")
	js.InitJsApi()

	// Prevent exit
	select {}
}
