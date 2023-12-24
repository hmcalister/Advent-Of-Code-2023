import numpy as np
from sympy import Symbol
from sympy import solve_poly_system

handle = open("puzzleInput", "r")

hailstones = []
for line in handle:
    pos, vel = line.strip().split(" @ ")
    x, y, z = pos.split(", ")
    vx, vy, vz = vel.split(", ")
    hailstones.append((int(x), int(y), int(z), int(vx), int(vy), int(vz)))

x = Symbol('x')
y = Symbol('y')
z = Symbol('z')
vx = Symbol('vx')
vy = Symbol('vy')
vz = Symbol('vz')

variables = [x, y, z, vx, vy, vz]
equations = []

# We only need three hailstones to solve this problem.
#
# We have six unknowns (above), and each new hailstone gives up another one (the time to intersections).
# Each hailstone gives us three equations (for x, y, and z intersections). So 3 hailstones is the
# smallest number that gives at least as many equations as unknowns.
for index, hail in enumerate(hailstones[:3]):
    print(hail)
    x0, y0, z0, vx0, vy0, vz0 = hail
    t = Symbol(f"t{index}")

    equations.extend([
        (x + vx*t) - (x0 + vx0*t),
        (y + vy*t) - (y0 + vy0*t),
        (z + vz*t) - (z0 + vz0*t),
    ])

    variables.append(t)

print(variables)

result = solve_poly_system(equations, )
print(result[0][0]+result[0][1]+result[0][2])
