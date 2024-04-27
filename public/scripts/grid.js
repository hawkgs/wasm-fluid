// Correlates to the vals in parameters.go
const SYS_SCALE = 40;
const SMOOTHING_RADIUS_H = 0.5;
const SCALED_H = SYS_SCALE * SMOOTHING_RADIUS_H;

export function createGrid({ width, height, showCellKey }) {
  const grid = document.getElementById('grid');
  grid.style.width = width + 'px';
  grid.style.height = height + 'px';

  const gridWidth = width / SCALED_H;
  const gridHeight = height / SCALED_H;
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
