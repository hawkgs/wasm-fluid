/**
 * @typedef {Object} Go
 */

const CANVAS_WIDTH = 600;
const CANVAS_HEIGHT = 400;
const PARTICLES = 100;

const canvas = document.getElementById('canvas');
canvas.width = CANVAS_WIDTH;
canvas.height = CANVAS_HEIGHT;
window.FluidApi = {};

const ctx = canvas.getContext('2d');

FluidApi.updateHandler = (particles) => {
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  ctx.beginPath();
  particles.forEach((particle) => {
    ctx.moveTo(particle.x, particle.y);
    ctx.arc(particle.x, particle.y, 3, 0, 2 * Math.PI, false);
    ctx.fillStyle = 'red';
  });
  ctx.fill();
};

async function init() {
  const go = new Go();
  const results = await WebAssembly.instantiateStreaming(
    fetch('../wasm/fluid.wasm'),
    go.importObject,
  );
  go.run(results.instance);

  FluidApi.createFluidSystem({
    width: CANVAS_WIDTH,
    height: CANVAS_HEIGHT,
    particles: PARTICLES,
  });

  setInterval(() => {
    requestAnimationFrame(() => FluidApi.requestUpdate());
  }, 1000 / 60);
}

init();
