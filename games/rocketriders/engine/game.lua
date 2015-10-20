
vector = require('libs.vector')
require('libs.anal')

require('engine.map')
require('engine.tileset')
require('engine.animations')
require('engine.network')
require('engine.clientprotocol')

globals = {}

local function translatemouse(x, y)
    x = math.min(math.max(1, x), 671)
    y = math.min(math.max(1, y), 671)
    local x = math.floor(x / tileset_size) + globals.player.pos.x - globals.offset.x
    local y = math.floor(y / tileset_size) + globals.player.pos.y - globals.offset.y
    local tmp = vector(x,y)
    return tmp
end

function chat_msg(msg)
    local text = loveframes.Create('text')
    text:SetMaxWidth(msgframe:GetWidth())
    text:SetFont(globals.font)
    text:SetText(msg)
    msglist:AddItem(text)
    if #msglist:GetChildren() > 499 then
        msglist:RemoveItem(1)
    end
end

local function parse_msg(data)
    -- TODO: move to the server instead
    msginput:Clear()
    if data:len() < 1 then return end
    local cmd = string.match(data, '!(%w+)')
    local arg = string.match(data, ' (%w*)')
    if cmd == 'name' then
        if arg == nil then arg = globals.player.id end
        network_client_write('name:'..arg)
    else
        network_client_write('msg:'..data)
    end
end

function do_jump(obj)
    local x = tonumber(navjumpx:GetText())
    local y = tonumber(navjumpy:GetText())
    if (x == nil) or (y == nil) then return end
    network_client_write(string.format('jump:%d,%d', x, y))
end

local function load()
    globals.running = true
    globals.focus = nil
    globals.offset = vector(10,10)
    globals.player = {
        pos=vector(10,10),
        sight=10,
        direction=vector(0,0),
        frac=0,
    }
    globals.entities = {}
    globals.host = '127.0.0.1'
    globals.port = 8123
    globals.name = 'Player'
    globals.fullscreen = false
    globals.font_name = 'assets/profont.ttf'
    globals.font_size = 11
    globals.font = love.graphics.newFont(globals.font_name, globals.font_size)
    love.graphics.setFont(globals.font)
    love.keyboard.setKeyRepeat(0.5, 0.05)
    local ret, data = pcall(dofile, 'settings.lua')
    if ret then
        globals.host = data.host
        globals.port = data.port
        globals.name = data.name
        globals.fullscreen = data.fullscreen
    end
    data = nil
    if globals.fullscreen then love.graphics.toggleFullscreen() end

    local ret, err = network_client(globals.host, globals.port)
    if not ret then
        local frame = loveframes.Create('frame')
        frame:SetState('mainmenu')
        frame:SetName('Connection Error')
        frame:SetSize(400,75)
        frame:Center()
        local txt = loveframes.Create('text', frame)
        txt:SetFont(globals.font)
        txt:SetPos(10,30)
        txt:SetMaxWidth(400)
        txt:SetText('Could not connect to the server: '..err)
        globals.running = false
        return
    end
    network_client_write('handshake:')

    local sh = love.graphics.getHeight()
    msgframe = loveframes.Create('frame')
    msgframe:SetState('game')
    msgframe:SetName('Messages')
    msgframe:SetDraggable(false)
    msgframe:ShowCloseButton(false)
    msgframe:SetPos(675,sh/2)
    msgframe:SetSize(349, sh/2)
    local mw, mh = msgframe:GetWidth(), msgframe:GetHeight()
    msglist = loveframes.Create('list', msgframe)
    msglist:SetPos(0, 25)
    msglist:SetSize(mw, mh-44)
    msglist:SetAutoScroll(true)
    msglist:SetPadding(2)
    msginput = loveframes.Create('textinput', msgframe)
    msginput:SetLimit(500)
    msginput:SetFont(globals.font)
    msginput:SetPos(0, mh-20)
    msginput:SetSize(mw, 20)
    msginput.OnFocusGained = function(obj) globals.focus = msginput end
    msginput.OnFocusLost = function(obj) globals.focus = nil end
    msginput.OnEnter = function(obj, txt) parse_msg(txt) end

    -- TODO: move to sep. file. be able to tab between input boxes
    navframe = loveframes.Create('frame')
    navframe:SetState('game')
    navframe:SetName('NavCon')
    navframe:ShowCloseButton(false)
    navframe:SetSize(300, 300)
    navframe:Center()
    navframe:SetVisible(false)
    local nw, nh = navframe:GetWidth(), navframe:GetHeight()
    local navclose = loveframes.Create('button', navframe)
    navclose:SetSize(75, 30)
    navclose:SetPos(nw-85, nh-40)
    navclose:SetText('Close')
    navclose.OnClick = function(obj) navframe:SetVisible(false) end
    navpower = loveframes.Create('text', navframe)
    navpower:SetFont(globals.font)
    navpower:SetPos(10, (nh/10)*1)
    navpower:SetText('Navigation Power: 100kW')
    navdrive = loveframes.Create('text', navframe)
    navdrive:SetFont(globals.font)
    navdrive:SetPos(10, (nh/10)*2)
    navdrive:SetText('Jump drive: 0%')
    navpos = loveframes.Create('text', navframe)
    navpos:SetFont(globals.font)
    navpos:SetPos(10, (nh/10)*3)
    navpos:SetText('Ship position: 0, 0')
    local navjump = loveframes.Create('text', navframe)
    navjump:SetFont(globals.font)
    navjump:SetPos(10, (nh/10)*4.5)
    navjump:SetText('Jump coordinates:')
    navjumpx = loveframes.Create('textinput', navframe)
    navjumpx:SetFont(globals.font)
    navjumpx:SetPos(10, (nh/10)*5)
    navjumpx:SetSize(75, 20)
    navjumpx.OnFocusGained = function(obj) globals.focus = navjumpx end
    navjumpx.OnFocusLost = function(obj) globals.focus = nil end
    navjumpy = loveframes.Create('textinput', navframe)
    navjumpy:SetFont(globals.font)
    navjumpy:SetPos(95, (nh/10)*5)
    navjumpy:SetSize(75, 20)
    navjumpy.OnFocusGained = function(obj) globals.focus = navjumpy end
    navjumpy.OnFocusLost = function(obj) globals.focus = nil end
    navdojump = loveframes.Create('button', navframe)
    navdojump:SetSize(75, 30)
    navdojump:SetPos(10, nh-40)
    navdojump:SetText('Jump')
    navdojump.OnClick = do_jump

    load_animations()
    load_tiles()
    map_clear()
end

local function update_view()
    tilebatch:clear()
    local start = vector(0,0) + globals.player.pos - globals.offset
    local stop = vector(20,20) + globals.player.pos - globals.offset
    for pos in map_loop(start, stop) do
            local tile = tileset[map_get(pos).tile]
            if tile == nil then tile = tileset[-1] end
            local tmp = (pos - globals.player.pos + globals.offset) * tileset_size
            tilebatch:addq(tile.quad, tmp.x, tmp.y)
    end
end

local function update(dt)
    -- TODO: show a disconnect msg
    if parse_cmds() then globals.running = false end
    update_animations(dt)
    update_view() -- TODO: only call this when needed
    loveframes.update(dt)
    return globals.running
end

local function draw_entity(ent)
    local delta = vector(0,0) + ent.direction * ent.frac
    local pos = (ent.pos * tileset_size) - delta
    if ent.direction == vector(-1,0) then
        img = ent.left
    elseif ent.direction == vector(1,0) then
        img = ent.right
    elseif ent.direction == vector(0,-1) then
        img = ent.up
    else
        img = ent.down
    end
    if img ~= nil then
        animations[img]:draw(pos.x, pos.y)
        love.graphics.print(ent.title, pos.x, pos.y-10)
    end
    if ent.frac > 0 then ent.frac = ent.frac - 1 end
end

local function draw_entities()
    for ent in entity_loop() do
        if ent.id ~= globals.player.id then
            draw_entity(ent)
        end
    end
    draw_entity(globals.player)
end

local function draw()
    love.graphics.push()
    local tmp = ((globals.offset*2)+ vector(1,1))*tileset_size
    love.graphics.setScissor(1, 1, tmp.x, tmp.y)

    -- add some smooth map movement with this delta
    local delta = vector(0,0) + globals.player.direction * globals.player.frac
    love.graphics.draw(tilebatch, delta.x, delta.y)

    tmp = (globals.player.pos*tileset_size) - (globals.offset*tileset_size)
    tmp = tmp - (globals.player.direction * globals.player.frac)
    love.graphics.translate(-tmp.x, -tmp.y)
    draw_entities()
    love.graphics.setScissor()
    love.graphics.pop()

    -- player sight mask
    love.graphics.setColor(0,0,0,230)
    for x=0, 20 do
        for y=0, 20 do
            local pos = vector(x,y)
            if globals.offset:dist(pos) > globals.player.sight-1 then
                pos = pos * tileset_size
                love.graphics.rectangle('fill', pos.x, pos.y, tileset_size, tileset_size)
            end
        end
    end
    love.graphics.setColor(255,255,255,255)

    -- map borders
    love.graphics.line(0, 673, 673, 673)
    love.graphics.line(673, 0, 673,673)

    -- show some debug status
    mx, my = love.mouse.getX(), love.mouse.getY()
    tm = translatemouse(mx, my)
    love.graphics.print('fps: '..love.timer.getFPS(), 0, 0)
    debugtext = 'Coord: '..tostring(globals.player.pos)..'\n'
    debugtext = debugtext .. 'Mouse: '..tostring(tm) .. '\n'
    debugtext = debugtext..'sight: '..tostring(globals.player.sight).. '\n'
    debugtext = debugtext..'Dir: '..tostring(globals.player.direction).. '\n'
    debugtext = debugtext .. '\n\n'
    debugtext = debugtext .. 'Selections\n'
    local cell = map_get(tm)
    local title, ents = '', ''
    title = tileset[cell.tile].title
    for ent in pairs(cell.entities) do
        ents = ent.title .. ' ('..ent.id..')'
    end
    debugtext = debugtext .. 'Tile: '..title .. '\n'
    debugtext = debugtext .. 'Entity: '..ents..'\n'
    debugtext = debugtext .. '\n'
    love.graphics.print(debugtext, 675, 0)

    -- render the GUI
    loveframes.draw()
end

local function keypressed(key, unicode)
    loveframes.keypressed(key, unicode)
    if globals.focus ~= nil then
        if key == 'escape' then globals.focus:SetFocus(false) end
        return
    end
    if key == 'escape' then
        -- TODO: add confirmation window
        network_client_write('quit:')
        globals.running = false
    elseif key == 'return' then msginput:SetFocus(true)
    elseif key == 'w' then
        local dir = vector(0,-1)
        if love.keyboard.isDown('lshift') then
            network_client_write('dir:'..dir.x..','..dir.y)
        else
            network_client_write('move:'..dir.x..','..dir.y)
        end
    elseif key == 's' then
        local dir = vector(0,1)
        if love.keyboard.isDown('lshift') then
            network_client_write('dir:'..dir.x..','..dir.y)
        else
            network_client_write('move:'..dir.x..','..dir.y)
        end
    elseif key == 'a' then
        local dir = vector(-1,0)
        if love.keyboard.isDown('lshift') then
            network_client_write('dir:'..dir.x..','..dir.y)
        else
            network_client_write('move:'..dir.x..','..dir.y)
        end
    elseif key == 'd' then
        local dir = vector(1,0)
        if love.keyboard.isDown('lshift') then
            network_client_write('dir:'..dir.x..','..dir.y)
        else
            network_client_write('move:'..dir.x..','..dir.y)
        end
    elseif key == 'e' then
        local dir = globals.player.direction
        network_client_write('use:'..dir.x..','..dir.y)
    elseif key == 'f' then
        love.graphics.toggleFullscreen()
        globals.fullscreen = not globals.fullscreen
    end
end

local function keyreleased(key)
    loveframes.keyreleased(key)
end

local function mousepressed(x, y, button)
    loveframes.mousepressed(x, y, button)
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
