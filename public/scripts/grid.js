export function createGrid({ width, height, showCellKey, parameters }) {
  const scaledH = parameters.systemScale * parameters.smoothingRadiusH;
  const grid = document.getElementById('grid');
  grid.style.width = width + 'px';
  grid.style.height = height + 'px';

  const gridWidth = width / scaledH;
  const gridHeight = height / scaledH;
  const gridSize = gridWidth * gridHeight;

  for (let i = 0, x = 0; i < gridSize; i++) {
    const y = i % gridWidth;

    const gCell = document.createElement('div');
    gCell.className = 'grid-cell';
    gCell.style.width = scaledH + 'px';
    gCell.style.height = scaledH + 'px';

    if (showCellKey) {
      gCell.innerText = `${x},${y}`;
    }

    grid.appendChild(gCell);

    if (y === gridWidth - 1) {
      x++;
    }
  }
}
