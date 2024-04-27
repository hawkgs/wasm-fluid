const DEFAULT_FPS = 60;

export function initAnimationControls({ onPlay, onPause }) {
  const manUpdateBtn = document.getElementById('manual-update');
  const playBtn = document.getElementById('play-btn');
  const statsBtn = document.getElementById('stats-btn');
  const resetBtn = document.getElementById('reset-btn');
  const fpsInput = document.getElementById('fps-input');

  fpsInput.value = DEFAULT_FPS;

  let isPlaying = false;

  const play = () => {
    playBtn.innerHTML = '⏸️ PAUSE';
    manUpdateBtn.disabled = true;
    fpsInput.disabled = true;

    const fps = parseInt(fpsInput.value, 10);

    onPlay(fps);
  };

  const pause = () => {
    playBtn.innerHTML = '▶️ PLAY';
    manUpdateBtn.disabled = false;
    fpsInput.disabled = false;

    onPause();
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

export function initParametersControls(onParamsUpdate) {
  const fragment = document.createDocumentFragment();
  [
    createSliderCtrl({
      name: 'Gravity',
      range: { min: 0, max: 5000 },
      defaultValue: 1000,
      step: 10,
      onUpdate: (v) => onParamsUpdate('gravityForce', v),
    }),
    createSliderCtrl({
      name: 'Smoothing radius',
      range: { min: 0, max: 10 },
      defaultValue: 0.5,
      step: 0.01,
      onUpdate: (v) => onParamsUpdate('smoothingRadiusH', v),
    }),
    createSliderCtrl({
      name: 'Gas const',
      range: { min: 0, max: 2000 },
      defaultValue: 800,
      step: 1,
      onUpdate: (v) => onParamsUpdate('gasConstK', v),
    }),
    createSliderCtrl({
      name: 'Rest density',
      range: { min: 0, max: 100 },
      defaultValue: 5,
      step: 0.1,
      onUpdate: (v) => onParamsUpdate('restDensity', v),
    }),
    createSliderCtrl({
      name: 'Viscosity const',
      range: { min: 0, max: 50 },
      defaultValue: 0,
      step: 0.1,
      onUpdate: (v) => onParamsUpdate('viscosityConst', v),
    }),
  ].forEach((ctrl) => fragment.appendChild(ctrl));

  document.getElementById('params').appendChild(fragment);
}

function createSliderCtrl({ name, range, defaultValue, step, onUpdate }) {
  const ctrl = document.createElement('div');
  ctrl.className = 'slider-ctrl';

  const nameLabel = document.createElement('span');
  nameLabel.className = 'name';
  nameLabel.innerHTML = name + ':';

  ctrl.appendChild(nameLabel);

  const minLabel = document.createElement('span');
  minLabel.className = 'range-val range-min';
  minLabel.innerHTML = range.min;

  const maxLabel = document.createElement('span');
  maxLabel.className = 'range-val range-max';
  maxLabel.innerHTML = range.max;

  ctrl.appendChild(minLabel);

  const slider = document.createElement('input');
  slider.type = 'range';
  slider.min = range.min;
  slider.max = range.max;
  slider.value = defaultValue;
  slider.step = step || 0.1;

  ctrl.appendChild(slider);
  ctrl.appendChild(maxLabel);

  const currVal = document.createElement('span');
  currVal.className = 'curr-val';
  currVal.innerHTML = defaultValue;
  ctrl.appendChild(currVal);

  slider.addEventListener('input', () => {
    currVal.innerHTML = slider.value;
    onUpdate(slider.value);
  });

  return ctrl;
}
