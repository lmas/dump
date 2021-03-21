
MARS MAN
--------------------------------------------------------------------------------
You have crash landed on Mars, survive it's harsh surface and escape it.

GAMEPLAY
--------------------------------------------------------------------------------
ACT 1:
- Start at crashed escape pod, on Mars, and with a starting tool.
- Gather resources and make better tools.
- Build first antenna and computer to contact Mission Control.
- Receive the first blueprint for a producer.

ACT 2:
- Construct producers to harvest resources for you.
- Store resources in containers.
- Manually feed fuel to the producers until you can build power plants.
- You now generate resources automatically.

ACT 3:
- Receive enough data for a rocket blueprint.
- Gather enough resources to build it.
- Escape Mars in built rocket.

Research:
- Blueprint data is sent from Mission Control on Earth.
- Ticks over time, to simulate data transfer.
- More data = more/better blueprints.
- Build better antennaes and computers to speed up transfer rate.

BUILDINGS
--------------------------------------------------------------------------------
Construction:
- Create a "building command block" (aka BCB)
- Place the BCB somewhere
- When placed, make two lists: requirements and slots
- Requirements = list of blocks that must be built near the BCB
- Slots = list of built and active slots used with the BCB
- When placing a new, required block nearby check if it's in the req. list
- If so, remove block from req. list and add it to the slots list
- Building done when all blocks from the req. list is gone and the slot list is full

Usage:
- Press Use key on the BCB for the building.
- BCB check all blocks in the slots list and make sure they're in good condition
- If so, run the buildings actuall function





OLD SHIT DOWN BELOW



OVERALL
--------------------------------------------------------------------------------
- Generate map
- Player can move on map
- Can gather resources from map tiles, using tools
- Construct buildings from resources
- Buildings will produce more resources
- Research will enable new buildings/tools and upgrade old ones

Game will not have:
- Life support (no point as the player will be the only human around)

MAP
--------------------------------------------------------------------------------
- Generate height map based on opensimplex noise
- Use height map to make terrain and terrain resources

Tile types:
- Rock
- Soil/dirt
- Biomass
- Liquid?
- Heat source

RESOURCES
--------------------------------------------------------------------------------
Mineral Matter:
- Used for building stuff
- Mined from tiles (ground = low yield, rocks/mountains = high yield)

Organic Matter:
- Used for research

Energy:
- Powers buildings and tools
- Place building near terrain source to produce energy

TOOLS
--------------------------------------------------------------------------------
- Used by the player to do things manually

BUILDINGS
--------------------------------------------------------------------------------
- Does things automatically

Building types:
- Living (main player hub)
- Producers (placed near tile source, produces resource)
- Storage (stores resources)
- Research

Energy producers:
- Solar cell (low cost, low yield, placed everywhere)
- Thermal plant (placed near heat source)
- Nuclear plant (high cost)

RESEARCH
--------------------------------------------------------------------------------
