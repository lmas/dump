
local TILESHEET = 'tilesheet.png'
local TILE_SIZE = 32

local TILESET = {}

TILESET[-1] = {
    name = 'ERROR',
    sheetx = 0,
    sheety = 0
}

TILESET[0] = {
    name = 'Empty',
    sheetx = 1,
    sheety = 0
}

TILESET[1] = {
    name = 'Floor',
    sheetx = 3,
    sheety = 0
}

TILESET[2] = {
    name = 'Wall',
    sheetx = 2,
    sheety = 0
}

map_sizex, map_sizey = 10, 10
map = {}

function love.load()
    -- Add random tiles to the map
    for x = 0, map_sizex do
        map[x] = {}
        for y = 0, map_sizey do
            map[x][y] = math.random(0, 2)
        end
    end

    -- Load the tileset
    tilesheet = love.graphics.newImage(TILESHEET)
    tilebatch = love.graphics.newSpriteBatch(tilesheet, 21 * 21)
    for _, tile in pairs(TILESET) do
        local x, y = tile.sheetx or 0, tile.sheety or 0
        print (x, y, tile.name)
        tile.quad = love.graphics.newQuad(
            x * TILE_SIZE, y * TILE_SIZE,
            TILE_SIZE, TILE_SIZE,
            tilesheet.getWidth(), tilesheet.getHeight()
        )
    end
end

function love.update()
    tilebatch:bind()
    tilebatch:clear()

    for x = 0, map_sizex do
        for y = 0, map_sizey do
            local tile = map[x][y]
            tilebatch:addq(tile, x * TILE_SIZE, y * TILE_SIZE)
        end
    end

    tilebatch:unbind()
end

function love.draw()
    love.graphics.draw(tilebatch)
end
