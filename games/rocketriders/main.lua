
require('libs.loveframes')

mainmenu = require('engine.mainmenu')
game = require('engine.game')
mapeditor = require('engine.mapeditor')
gameoptions = require('engine.gameoptions')

gamestate = {}

function state_push(title, state)
    table.insert(gamestate, {title, state})
    loveframes.SetState(title)
    state.load()
end

function state_pop()
    local old = table.remove(gamestate, #gamestate)
    collectgarbage('collect')
    love.graphics.clear()
    state_current()
end

function state_current()
    if #gamestate < 1 then
        love.event.quit()
    else
        local state = gamestate[#gamestate]
        loveframes.SetState(state[1])
        return state[2]
    end
end

function love.load()
    state_push('mainmenu', mainmenu)
end

function love.update(dt)
    if not state_current().update(dt) then
        state_pop()
    end
end

function love.draw()
    local c = state_current()
    if c then c.draw() end
end

function love.keypressed(key, unicode)
    local c = state_current()
    if c then c.keypressed(key, unicode) end
end

function love.keyreleased(key)
    local c = state_current()
    if c then c.keyreleased(key) end
end

function love.mousepressed(x, y, button)
    local c = state_current()
    if c then c.mousepressed(x, y, button) end
end

function love.mousereleased(x, y, button)
    local c = state_current()
    if c then c.mousereleased(x, y, button) end
end

function love.quit()
end

--[[ to create a new game state, copy this block of code into a new lua file:
local function load()
end

local function update(dt)
end

local function draw()
end

local function keypressed(key, unicode)
end

local function keyreleased(key)
end

local function mousepressed(x, y, button)
end

local function mousereleased(x, y, button)
end



return {
    load = load,
    update = update,
    draw = draw,
    keypressed = keypressed,
    keyreleased = keyreleased,
    mousepressed = mousepressed,
    mousereleased = mousereleased,
    --quit = quit,
}
]]
