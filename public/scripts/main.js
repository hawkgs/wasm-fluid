/**
 * @typedef {Object} Go
 */

const CANVAS_WIDTH = 600;
const CANVAS_HEIGHT = 400;
const PARTICLES = 200;
const PARTICLE_UI_RADIUS = 3;
const UPDATE_FREQ = 1000 / 10;

let updateType = 'none'; // none | auto | manual

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
    ctx.arc(particle.x, particle.y, PARTICLE_UI_RADIUS, 0, 2 * Math.PI, false);
    ctx.fillStyle = 'blue';
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
    particleUiRadius: PARTICLE_UI_RADIUS,
  });

  // Testing helpers

  const manUpdateBtn = document.getElementById('manual-update');
  const startAnimBtn = document.getElementById('start-anim');

  manUpdateBtn.addEventListener('click', () => {
    if (updateType !== 'manual') {
      updateType = 'manual';
      startAnimBtn.disabled = true;
    }

    requestAnimationFrame(() => FluidApi.requestUpdate());
  });

  document.addEventListener('keypress', (e) => {
    if (updateType !== 'manual') {
      updateType = 'manual';
      startAnimBtn.disabled = true;
    }

    if (e.key === 'Space') {
      requestAnimationFrame(() => FluidApi.requestUpdate());
    }
  });

  document.getElementById('start-anim').addEventListener('click', () => {
    manUpdateBtn.disabled = true;
    startAnimBtn.disabled = true;
    updateType = 'auto';

    setInterval(() => {
      requestAnimationFrame(() => FluidApi.requestUpdate());
    }, UPDATE_FREQ);
  });
}

init();
