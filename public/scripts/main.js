/**
 * @typedef {Object} Go
 */

import { initAnimationControls, initParametersControls } from './controls.js';

const CANVAS_WIDTH = 600;
const CANVAS_HEIGHT = 400;
const PARTICLES = 1000;
const PARTICLE_UI_RADIUS = 3;

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

async function init() {
  const go = new Go();
  const results = await WebAssembly.instantiateStreaming(
    fetch('../wasm/fluid.wasm'),
    go.importObject,
  );
  go.run(results.instance);

  createSystem();

  let interval;

  initAnimationControls({
    onPlay: (fps) => {
      interval = setInterval(() => {
        requestAnimationFrame(() => FluidApi.requestUpdate());
      }, 1000 / fps);
    },
    onPause: () => clearInterval(interval),
  });

  initParametersControls((paramName, value) => {
    // console.log(paramName, value);
  });
}

init();
