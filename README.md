# Astral
A procedurally generated pixel-art galaxy


## How it works
For the planet, the screen is divided into a grid of squares. Then, any squares that are within a certain radius
from the center are considered part of the planet. To generate landforms, we use multiple octaves of simplex noise. The number of
octaves is randomly selected from a predefined range. The parameters of the noise (frequency, amplitude, persistence) are also randomly 
selected from a predefined range. The colors are also randomly generated- there are certain "base" colors (based on
common Earth colors; blue for oceans, green for land, white for snow, etc.) which are then interpolated with randomly chosen colors. This
is to ensure that the planet does not become too ridiculously colored, although I'm not too sure that it works. Finally, multiple planets
are generated into 2D color arrays and placed in randomly sized "orbits" around the center.
