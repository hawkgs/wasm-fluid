# wasm-fluid

_**Under development**_

🌊 Browser-based fluid simulation (Müller et al SPH method).

Calculations: WebAssembly (Go); Visualization: HTML Canvas

![Demo](./assets/demo.gif)

Müller et al SPH method: https://matthias-research.github.io/pages/publications/sca03.pdf

### Issues

As stated above, the simulation is based on smoothed-particle hydrodynamics method by M. Müller and others. It incorporates density, pressure and viscosity calculations, plus external forces – gravity in our case. Surface tension is not implemented. Since the paper focuses on 3D SPH, the used smoothing kernels for this implementation – or more precisely, their normalization scaling factors – are adapted for 2D. Anyway, there are still things that require attention/fixes, some of which are probably a matter of proper parameter tuning:

- The fluid acts as a highly viscous liquid.
- Particles that go outside of the particle stack and then get reintroduced at the edge of the smoothing radus of another particle, produce a critically low density. When the forces are divided by the density, the product acceleration has a very large magnitude which ejects the particle from the stack. There is a temporary fix for this.
- There is an instability at the edge of the fluid sometimes. This is, again, most likely a result of low densities producing high acceleration (but not enough to eject the particle out of the fluid).
