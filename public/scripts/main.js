/**
 * @typedef {Object} Go
 */

import { initAnimationControls, initParametersControls } from './controls.js';
import { createGrid } from './grid.js';

const CANVAS_WIDTH = 600;
const CANVAS_HEIGHT = 400;
const PARTICLES = 1000;
const PARTICLE_UI_RADIUS = 3;

const parameters = {
  gravityForce: 0,
  gasConstK: 380,
  restDensity: 1.7,
  viscosityConst: 0,
  smoothingRadiusH: 0.45,
  timestep: 0.005,
};

// Create grid (for debugging purposes)
createGrid({ showCellKey: false, width: CANVAS_WIDTH, height: CANVAS_HEIGHT });

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
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  FluidApi.createFluidSystem(
    {
      width: CANVAS_WIDTH,
      height: CANVAS_HEIGHT,
      particles: PARTICLES,
      particleUiRadius: PARTICLE_UI_RADIUS,
    },
    parameters,
  );

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
    onStats: () => FluidApi.devPrintSystemStats(),
    onReset: createSystem,
  });

  initParametersControls((paramName, value) => {
    if (['smoothingRadiusH', 'timestep'].includes(paramName)) {
      createSystem();
    }
    parameters[paramName] = value;
    FluidApi.devParamsUpdate(parameters);
  }, parameters);
}

init();
