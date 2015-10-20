
local function savechanges()
    globals.fullscreen = chkfullscreen:GetChecked()
    local fullscreen = globals.fullscreen and 'true' or 'false'
    globals.host = inhost:GetText()
    globals.name = inname:GetText()
    local port = tonumber(inport:GetText())
    if port ~= nil then
        globals.port = math.max(1, math.min(65535, port))
    else
        port = 8123
    end

    local tmp = string.format([[return {
    host = '%s',
    port = %d,
    name = '%s',
    fullscreen = %s,
}
]], globals.host, port, globals.name, fullscreen)
    local f = io.open('settings.lua', 'wb')
    f:write(tmp)
    f:close()
    globals.running = false
end

local function load()
    globals.running = true
    globals.host = 'localhost'
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

    local sw, sh = love.graphics.getWidth(), love.graphics.getHeight()
    local frame = loveframes.Create('frame')
    frame:SetState('gameoptions')
    frame:SetName('Game Options')
    frame:SetDraggable(false)
    frame:ShowCloseButton(false)
    frame:SetSize(400, 400)
    frame:Center()
    local fw, fh = frame:GetWidth(), frame:GetHeight()

    local txthost = loveframes.Create('text', frame)
    txthost:SetPos(10, (fh/10)*1)
    txthost:SetText('Server address:')
    txthost:SetFont(globals.font)
    inhost = loveframes.Create('textinput', frame)
    inhost:SetSize(300, 20)
    inhost:SetPos(10, (fh/10)*1.5)
    inhost:SetLimit(200)
    inhost:SetFont(globals.font)
    inhost:SetText(globals.host)

    local txtport = loveframes.Create('text', frame)
    txtport:SetPos(10, (fh/10)*2.5)
    txtport:SetText('Server port:')
    txtport:SetFont(globals.font)
    inport = loveframes.Create('textinput', frame)
    inport:SetSize(100, 20)
    inport:SetPos(10, (fh/10)*3)
    inport:SetFont(globals.font)
    inport:SetText(globals.port)

    local txtname = loveframes.Create('text', frame)
    txtname:SetPos(10, (fh/10)*4)
    txtname:SetText('Nickname:')
    txtname:SetFont(globals.font)
    inname = loveframes.Create('textinput', frame)
    inname:SetSize(300, 20)
    inname:SetPos(10, (fh/10)*4.5)
    inhost:SetLimit(200)
    inname:SetFont(globals.font)
    inname:SetText(globals.name)

    local txtfullscreen = loveframes.Create('text', frame)
    txtfullscreen:SetPos(10, (fh/10)*5.5)
    txtfullscreen:SetText('Fullscreen:')
    txtfullscreen:SetFont(globals.font)
    chkfullscreen = loveframes.Create('checkbox', frame)
    chkfullscreen:SetSize(20, 20)
    chkfullscreen:SetPos(10, (fh/10)*6)
    chkfullscreen:SetChecked(globals.fullscreen)

    local btnok = loveframes.Create('button', frame)
    btnok:SetSize(100, 30)
    btnok:SetPos(10, (fh/10)*9)
    btnok:SetText('OK')
    btnok.OnClick = function(obj) savechanges() end

    local btncancel = loveframes.Create('button', frame)
    btncancel:SetSize(100, 30)
    btncancel:SetPos(fw - 110, (fh/10)*9)
    btncancel:SetText('Cancel')
    btncancel.OnClick = function(obj) globals.running = false end
end

local function update(dt)
    loveframes.update(dt)
    return globals.running
end

local function draw()
    loveframes.draw()
end

local function keypressed(key, unicode)
    loveframes.keypressed(key, unicode)
    if key == 'escape' then globals.running = false
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
