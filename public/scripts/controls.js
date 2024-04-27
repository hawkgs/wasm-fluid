const DEFAULT_FPS = 60;

const animEvents = new EventTarget();

export function initAnimationControls({ onPlay, onPause, onReset, onStats }) {
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

    animEvents.dispatchEvent(new CustomEvent('anim', { detail: 'play' }));
    onPlay(fps);
  };

  const pause = () => {
    playBtn.innerHTML = '▶️ PLAY';
    manUpdateBtn.disabled = false;
    fpsInput.disabled = false;

    animEvents.dispatchEvent(new CustomEvent('anim', { detail: 'pause' }));
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
    onReset();
  });

  statsBtn.addEventListener('click', onStats);
}

export function initParametersControls(onParamsUpdate, defaults) {
  const fragment = document.createDocumentFragment();
  const controls = [
    createSliderCtrl({
      name: 'Smoothing radius (h)',
      range: { min: 0, max: 2 },
      defaultValue: defaults.smoothingRadiusH,
      step: 0.001,
      onUpdate: (v) => onParamsUpdate('smoothingRadiusH', v),
    }),
    createSliderCtrl({
      name: 'Timestep (𝚫t)',
      range: { min: 0, max: 0.1 },
      defaultValue: defaults.timestep,
      step: 0.0001,
      onUpdate: (v) => onParamsUpdate('timestep', v),
    }),
    document.createElement('hr'),
    createSliderCtrl({
      name: 'Gravity (G)',
      range: { min: 0, max: 5000 },
      defaultValue: defaults.gravityForce,
      step: 10,
      onUpdate: (v) => onParamsUpdate('gravityForce', v),
    }),
    createSliderCtrl({
      name: 'Gas const (k)',
      range: { min: 0, max: 2000 },
      defaultValue: defaults.gasConstK,
      step: 1,
      onUpdate: (v) => onParamsUpdate('gasConstK', v),
    }),
    createSliderCtrl({
      name: 'Rest density (⍴₀)',
      range: { min: 0, max: 100 },
      defaultValue: defaults.restDensity,
      step: 0.1,
      onUpdate: (v) => onParamsUpdate('restDensity', v),
    }),
    createSliderCtrl({
      name: 'Viscosity const (μ)',
      range: { min: 0, max: 200 },
      defaultValue: defaults.viscosityConst,
      step: 0.1,
      onUpdate: (v) => onParamsUpdate('viscosityConst', v),
    }),
  ];

  controls.forEach((ctrlObj) => {
    const el = ctrlObj.ctrl ? ctrlObj.ctrl : ctrlObj;
    fragment.appendChild(el);
  });

  const [smRadiusCtrlObj, timestepCtrlObj] = controls;

  animEvents.addEventListener('anim', ({ detail }) => {
    switch (detail) {
      case 'play':
        smRadiusCtrlObj.slider.disabled = true;
        timestepCtrlObj.slider.disabled = true;
        break;
      case 'pause':
        smRadiusCtrlObj.slider.disabled = false;
        timestepCtrlObj.slider.disabled = false;
        break;
    }
  });

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
  slider.step = step.toString();
  slider.value = defaultValue.toString();

  ctrl.appendChild(slider);
  ctrl.appendChild(maxLabel);

  const currVal = document.createElement('span');
  currVal.className = 'curr-val';
  currVal.innerHTML = defaultValue;
  ctrl.appendChild(currVal);

  slider.addEventListener('input', () => {
    currVal.innerHTML = slider.value;
    onUpdate(parseInt(slider.value, 10));
  });

  return { ctrl, slider };
}
