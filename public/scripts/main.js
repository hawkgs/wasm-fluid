/**
 * @typedef {Object} Go
 */

const CANVAS_WIDTH = 600;
const CANVAS_HEIGHT = 400;
const PARTICLES = 600;
const PARTICLE_UI_RADIUS = 3;
const DEFAULT_FPS = 60;

// Create grid (for debugging purposes)

// Correlates to the vals in parameters.go
const SYS_SCALE = 40;
const SMOOTHING_RADIUS_H = 0.5;
const SCALED_H = SYS_SCALE * SMOOTHING_RADIUS_H;

function createGrid(showCellKey) {
  const grid = document.getElementById('grid');
  grid.style.width = CANVAS_WIDTH + 'px';
  grid.style.height = CANVAS_HEIGHT + 'px';

  const gridWidth = CANVAS_WIDTH / SCALED_H;
  const gridHeight = CANVAS_HEIGHT / SCALED_H;
  const gridSize = gridWidth * gridHeight;

  for (let i = 0, x = 0; i < gridSize; i++) {
    const y = i % gridWidth;

    const gCell = document.createElement('div');
    gCell.className = 'grid-cell';
    gCell.style.width = SCALED_H + 'px';
    gCell.style.height = SCALED_H + 'px';

    if (showCellKey) {
      gCell.innerText = `${x},${y}`;
    }

    grid.appendChild(gCell);

    if (y === gridWidth - 1) {
      x++;
    }
  }
}

createGrid(false);

// Animation

const canvas = document.getElementById('canvas');
canvas.width = CANVAS_WIDTH;
canvas.height = CANVAS_HEIGHT;
window.FluidApi = {};

const ctx = canvas.getContext('2d');

FluidApi.updateHandler = (particles) => {
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  particles.forEach((particle) => {
    ctx.beginPath();
    ctx.arc(particle.x, particle.y, PARTICLE_UI_RADIUS, 0, 2 * Math.PI, false);
    ctx.fillStyle = particle.vColor;
    ctx.fill();
  });
};

function createSystem() {
  FluidApi.createFluidSystem({
    width: CANVAS_WIDTH,
    height: CANVAS_HEIGHT,
    particles: PARTICLES,
    particleUiRadius: PARTICLE_UI_RADIUS,
  });

  console.log('%cFluid system initialized!', 'color: lightgreen');
}

function initControls() {
  const manUpdateBtn = document.getElementById('manual-update');
  const playBtn = document.getElementById('play-btn');
  const statsBtn = document.getElementById('stats-btn');
  const resetBtn = document.getElementById('reset-btn');
  const fpsInput = document.getElementById('fps-input');

  fpsInput.value = DEFAULT_FPS;

  let isPlaying = false,
    interval;

  const play = () => {
    playBtn.innerHTML = '⏸️ PAUSE';
    manUpdateBtn.disabled = true;
    fpsInput.disabled = true;

    const fps = parseInt(fpsInput.value, 10);

    interval = setInterval(() => {
      requestAnimationFrame(() => FluidApi.requestUpdate());
    }, 1000 / fps);
  };

  const pause = () => {
    playBtn.innerHTML = '▶️ PLAY';
    manUpdateBtn.disabled = false;
    fpsInput.disabled = false;

    clearInterval(interval);
  };

  playBtn.addEventListener('click', () => {
    if (isPlaying) {
      pause();
    } else {
      play();
    }

    isPlaying = !isPlaying;
  });

  manUpdateBtn.addEventListener('click', () => {
    requestAnimationFrame(() => FluidApi.requestUpdate());
  });

  resetBtn.addEventListener('click', () => {
    if (isPlaying) {
      pause();
      isPlaying = false;
    }
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    createSystem();
  });

  // For debugging
  statsBtn.addEventListener('click', () => FluidApi.devPrintSystemStats());
}

async function init() {
  const go = new Go();
  const results = await WebAssembly.instantiateStreaming(
    fetch('../wasm/fluid.wasm'),
    go.importObject,
  );
  go.run(results.instance);

  createSystem();
  initControls();
}

init();
