
local TILESHEET = 'assets/tilesheet.png'
local TILE_SIZE = 32

local TILESET = {
}

TILESET[-1] = {
    name = 'ERROR',
    sheetx = 0,
    sheety = 0,
}

TILESET[0] = {
    name = 'Empty',
    sheetx = 1,
    sheety = 0,
}

TILESET[1] = {
    name = 'Floor',
    sheetx = 3,
    sheety = 0,
}

TILESET[2] = {
    name = 'Wall',
    sheetx = 2,
    sheety = 0,
}

return {
    sheet = TILESHEET,
    size = TILE_SIZE,
    set = TILESET,
}
