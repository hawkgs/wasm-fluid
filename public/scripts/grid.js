// Creates a grid that corresponds to the particle grid where the cell radius = smoothing radius
// Note: Has to be improved for smoothing radiuses that don't match the full width and height of the field
export function createGrid({ width, height, showCellKey, parameters }) {
  const scaledH = parameters.systemScale * parameters.smoothingRadiusH;
  const grid = document.getElementById('grid');
  grid.style.width = width + 'px';
  grid.style.height = height + 'px';

  const gridWidth = width / scaledH;
  const gridHeight = height / scaledH;
  const gridSize = gridWidth * gridHeight;
  const cells = [];

  for (let i = 0, x = 0; i < gridSize; i++) {
    const y = i % gridWidth;

    const gCell = document.createElement('div');
    gCell.className = 'grid-cell';
    gCell.style.width = scaledH + 'px';
    gCell.style.height = scaledH + 'px';
    cells.push(gCell);

    if (showCellKey) {
      gCell.innerText = `${x},${y}`;
    }

    grid.appendChild(gCell);

    if (y === gridWidth - 1) {
      x++;
    }
  }

  // Returns a method that can be used for further adjusting
  // of the cell radius based on a provided smoothing radius
  return (smoothinRadius) => {
    const scaledH = parameters.systemScale * smoothinRadius;
    cells.forEach((c) => {
      c.style.width = scaledH + 'px';
      c.style.height = scaledH + 'px';
    });
  };
}
