
vector = require('libs.vector')
require('libs.anal')

require('engine.map')
require('engine.tileset')

globals = {}

local function translatemouse(x, y)
    x = math.min(math.max(1, x), 671)
    y = math.min(math.max(1, y), 671)
    local x = math.floor(x / tileset_size) + globals.pos.x
    local y = math.floor(y / tileset_size) + globals.pos.y
    local tmp = vector(x,y)
    return tmp
end

local function move(delta)
    globals.pos = globals.pos + delta
end

local function loadmap()
    if not map_file('ship.dat') then
        map_clear()
    end
end

local function savemap()
    local map = map_save()
    local f = io.open('ship.dat', 'wb')
    f:write(map)
    --f:flush()
    f:close()
end

local function update_view()
    tilebatch:clear()
    local start = vector(0,0) + globals.pos
    local stop = vector(20,20) + globals.pos
    for pos in map_loop(start, stop) do
            local tile = tileset[map_get(pos).tile]
            if tile == nil then tile = tileset[-1] end
            local tmp = (pos - globals.pos) * tileset_size
            tilebatch:addq(tile.quad, tmp.x, tmp.y)
    end
end

local function load()
    load_tiles()
    loadmap()

    globals.running = true
    globals.pos = vector(0,0)
    globals.tile = 0
    globals.mapchange = false
    globals.font_name = 'assets/profont.ttf'
    globals.font_size = 11
    globals.font = love.graphics.newFont(globals.font_name, globals.font_size)
    love.graphics.setFont(globals.font)
    love.keyboard.setKeyRepeat(0.5, 0.05)
end

local function update(dt)
    loveframes.update(dt)
    update_view()
    return globals.running
end

local function draw()
    love.graphics.draw(tilebatch)

    -- show selected tile
    local tm = translatemouse(love.mouse.getX(), love.mouse.getY())
    tm = (tm - globals.pos) * tileset_size
    quad = tileset[globals.tile].quad
    love.graphics.drawq(tilesheet, quad, tm.x, tm.y)

    -- map borders
    love.graphics.line(0, 673, 673, 673)
    love.graphics.line(673, 0, 673,673)

    -- show some debug status
    local tm = translatemouse(love.mouse.getX(), love.mouse.getY())
    love.graphics.print('fps: '..love.timer.getFPS(), 0, 0)
    debugtext = 'Coord: '..tostring(globals.pos)..'\n'
    debugtext = debugtext .. 'Mouse: '..tostring(tm) .. '\n'
    debugtext = debugtext .. '\n\n'
    local cell = map_get(tm)
    debugtext = debugtext ..'Tile: '..tileset[cell.tile].title..'\n'
    local tile = tileset[globals.tile].title
    debugtext = debugtext .. 'Selected: '..tile..'\n\n'
    debugtext = debugtext .. 'Map state: '.. (globals.mapchanges and 'UNSAVED' or 'Saved')..'\n'
    debugtext = debugtext .. [[

=================================================

CONTROLS:
esc         - quit
w, a, s, d  - move around the map
q and e     - change tile selection
f           - toggle fullscreen
ctrl + s    - save map
ctrl + r    - reload map

MOUSE:
left        - place tile
right       - select tile

=================================================

Map file: C:\Users\alpha\AppData\Roaming\LOVE\..]]
    --debugtext = debugtext .. '\n'
    love.graphics.print(debugtext, 675, 0)

    -- render the GUI
    loveframes.draw()
end

local function keypressed(key)
    loveframes.keypressed(key, unicode)
    if key == 'escape' then globals.running = false
    elseif key == 'f' then love.graphics.toggleFullscreen()
    elseif key == 's' and love.keyboard.isDown('lctrl') then
        if globals.mapchanges then
            savemap()
            globals.mapchanges = false
        end
    elseif key == 'r' and love.keyboard.isDown('lctrl') then
        loadmap()
        globals.mapchanges = false
    elseif key == 'q' then
        globals.tile = math.min(math.max(globals.tile - 1, 0), #tileset)
    elseif key == 'e' then
        globals.tile = math.min(math.max(globals.tile + 1, 0), #tileset)
    elseif key == 'w' then move(vector(0,-1))
    elseif key == 's' then move(vector(0,1))
    elseif key == 'a' then move(vector(-1,0))
    elseif key == 'd' then move(vector(1,0))
    end
end

local function keyreleased(key)
    loveframes.keyreleased(key)
end

local function mousepressed(x, y, button)
    loveframes.mousepressed(x, y, button)
    local pos = translatemouse(x, y)
    if button == 'l' then
        local tile = map_get(pos)
        tile.tile = globals.tile
        map_set(pos, tile)
        globals.mapchanges = true
    elseif button == 'r' then
        globals.tile = map_get(pos).tile
    end
end

local function mousereleased(x, y, button)
    loveframes.mousereleased(x, y, button)
end

return {
    load = load,
    update = update,
    draw = draw,
    keypressed = keypressed,
    keyreleased = keyreleased,
    mousepressed = mousepressed,
    mousereleased = mousereleased,
}
