async function init() {
  const go = new Go();
  const results = await WebAssembly.instantiateStreaming(
    fetch('../wasm/main.wasm'),
    go.importObject,
  );

  go.run(results.instance);
}

console.log('Hello, JS!');
init();
