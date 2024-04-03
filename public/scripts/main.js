/**
 * @typedef {Object} Go
 */

const CANVAS_WIDTH = 600;
const CANVAS_HEIGHT = 400;

const canvas = document.getElementById('canvas');
canvas.width = CANVAS_WIDTH;
canvas.height = CANVAS_HEIGHT;

const ctx = canvas.getContext('2d');

window.InitCanvas = (resp) => {
  console.log('InitCanvas called');

  ctx.beginPath();
  ctx.arc(resp.x, resp.y, 5, 0, 2 * Math.PI, false);
  ctx.fillStyle = 'red';
  ctx.fill();
};

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
