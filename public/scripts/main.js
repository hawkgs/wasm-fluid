/**
 * @typedef {Object} Go
 */

const CANVAS_WIDTH = 600;
const CANVAS_HEIGHT = 400;

const canvas = document.getElementById('canvas');
canvas.width = CANVAS_WIDTH;
canvas.height = CANVAS_HEIGHT;
window.GoApi = {};

const ctx = canvas.getContext('2d');

GoApi.goUpdateHandler = (vector) => {
  console.log('render', vector);
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  ctx.beginPath();
  ctx.arc(vector.x, vector.y, 5, 0, 2 * Math.PI, false);
  ctx.fillStyle = 'red';
  ctx.fill();
};

async function init() {
  const go = new Go();
  const results = await WebAssembly.instantiateStreaming(
    fetch('../wasm/fluid.wasm'),
    go.importObject,
  );

  go.run(results.instance);
  GoApi.goCreateFluidSystem({
    width: CANVAS_WIDTH,
    height: CANVAS_HEIGHT,
  });

  setInterval(() => {
    GoApi.goRequestUpdate();
  }, 3000);
}

init();
