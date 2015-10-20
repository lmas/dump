
require('engine.utils')
local TILESET = require('engine.tileset')

local tilemap = {}


function tilemap:init()
    self:clear()
    self:load_tileset()
end

function tilemap:clear(tid)
    self.tid = tid or ''
    self.tiles = {}
    self.tileset = TILESET
end

function tilemap:load_tileset(sheet_file, tile_size)
    local tilesheet = love.graphics.newImage(sheet_file or TILESET.sheet)
    self.tile_size = tile_size or TILESET.size
    self.tilebatch = love.graphics.newSpriteBatch(tilesheet, 21 * 21)

    for _, tile in pairs(TILESET.set) do
        local x, y = tile.sheetx or 0, tile.sheety or 0
        tile.quad = love.graphics.newQuad(
            x*self.tile_size, y*self.tile_size, -- The top-left corner of the tile
            self.tile_size, self.tile_size, -- The width and height of the new quad
            tilesheet:getWidth(), tilesheet:getHeight()
        )
    end
end

function tilemap:load_map(filename)
    self:clear()
    local f = io.open(filename, 'rb')
    for line in f:lines() do
        self:parse_tile(line)
    end
    f:close()
end

function tilemap:save_map(filename)
    local data = ''
    for pos, tile in self:iter() do
        if tile.tile > 0 then
            data = data .. string.format('%s: %d\n', pos, tile.tile)
        end
    end
    local f = io.open(filename, 'wb')
    f:write(data)
    f:close()
end

function tilemap:set(x, y, tile)
    if type(tile) == 'table' then
        if tile.tile > 0 then
            self.tiles[x..','..y] = tile
        else
            self.tiles[x..','..y] = nil
        end
    else
        self.tiles[x..','..y] = nil
    end
end

function tilemap:get(x, y)
    local tile = self.tiles[x..','..y]
    if tile == nil then
        tile = {tile=0}
    end
    return tile
end

function tilemap:iter()
    return pairs(self.tiles)
end

function tilemap:parse_tile(data)
    local tmp = split(data, ':')

    local pos = split(tmp[1], ',')
    local x, y = pos[1], pos[2]

    local tile_data = split(tmp[2], ',')
    local tile = {
        tile = tonumber(tile_data[1]),
    }

    self:set(x, y, tile)
end

function tilemap:update()
    self.tilebatch:bind()
    self.tilebatch:clear()

    for pos, tile in self:iter() do
        local tmp = split(pos, ',')
        local x, y = tmp[1] * self.tile_size, tmp[2] * self.tile_size
        local quad = TILESET.set[tile.tile].quad
        self.tilebatch:add(quad, x, y)
    end
    self.tilebatch:unbind()
end

function tilemap:draw()
    love.graphics.draw(self.tilebatch)
end


return tilemap

