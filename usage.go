package main

const usage = `nozzle-isen: isentropic converging-diverging nozzle audit.

Reads a stagnation state (total temperature T0, total pressure p0, heat
capacity ratio gamma, specific gas constant R), an area ratio Ae/A* and a
throat area A* from a JSON case file, then reports the throat state, the
subsonic and supersonic exit branches, whether the throat is choked, and the
mass flow rate. All isentropic relations, the area-Mach function and the mass
flow use the same gamma and the same stagnation state.

usage:
  nozzle-isen design <input.json>

input.json fields:
  t0           total temperature (K, must be > 0)
  p0           total pressure (Pa, must be > 0)
  gamma        heat capacity ratio (must be > 1, margin from 1 enforced)
  r            specific gas constant (J/(kg.K), default 287.0 for air)
  area_ratio   exit-to-throat area ratio Ae/A* (must be >= 1)
  throat_area  throat area A* (m2, must be > 0)
  branch       optional: "subsonic" or "supersonic"; empty prints both
  back_pressure optional exit back pressure pb (Pa); when above the sonic
               pressure ratio the flow is unchoked and the exit falls
               isentropically to pb with the subsonic solution

boundaries: T0<=0, p0<=0, gamma<=1, R<=0, Ae/A*<1 and throat area<=0 are
reported on stderr with a non-zero exit code; Newton inversion that exhausts
its iteration limit reports non-convergence; gamma too close to 1 reports an
error instead of dividing by zero; a specified back pressure is reported as
overexpanded or underexpanded with the note that a normal shock is not
modelled.

example:
  nozzle-isen design example/gamma14.json
`
