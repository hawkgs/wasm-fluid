/**
 * @typedef {Object} Go
 */

const CANVAS_WIDTH = 600;
const CANVAS_HEIGHT = 400;
const PARTICLES = 200;
const PARTICLE_UI_RADIUS = 3;
const DEFAULT_FPS = 60;

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
  const stopAnimBtn = document.getElementById('stop-anim');
  const fpsInput = document.getElementById('fps-input');

  fpsInput.value = DEFAULT_FPS;
  stopAnimBtn.disabled = true;
  let interval;

  manUpdateBtn.addEventListener('click', () => {
    requestAnimationFrame(() => FluidApi.requestUpdate());
  });

  document.addEventListener('keypress', (e) => {
    if (e.key === 'Space') {
      requestAnimationFrame(() => FluidApi.requestUpdate());
    }
  });

  startAnimBtn.addEventListener('click', () => {
    startAnimBtn.disabled = true;
    manUpdateBtn.disabled = true;
    fpsInput.disabled = true;
    stopAnimBtn.disabled = false;
    const fps = parseInt(fpsInput.value, 10);

    interval = setInterval(() => {
      requestAnimationFrame(() => FluidApi.requestUpdate());
    }, 1000 / fps);
  });

  stopAnimBtn.addEventListener('click', () => {
    startAnimBtn.disabled = false;
    manUpdateBtn.disabled = false;
    fpsInput.disabled = false;
    stopAnimBtn.disabled = true;
    clearInterval(interval);
  });
}

init();
