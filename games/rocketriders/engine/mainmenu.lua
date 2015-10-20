

local running = true

local function load()
    local z = love.graphics.getHeight()/10

    local logo = loveframes.Create('image')
    logo:SetState('mainmenu')
    logo:SetImage('assets/images/logo.png')
    logo:SetSize(love.graphics.getWidth() - 100,z*5)
    logo:CenterX()
    logo:SetY(z*0)

    local play = loveframes.Create('imagebutton')
    play:SetState('mainmenu')
    play:SetImage('assets/images/play.png')
    play:SetText('')
    play:SizeToImage()
    play:CenterX()
    play:SetY(z*5)
    play.OnClick = function(obj)
        state_push('game', game)
    end
    local options = loveframes.Create('imagebutton')
    options:SetState('mainmenu')
    options:SetImage('assets/images/options.png')
    options:SetText('')
    options:SizeToImage()
    options:CenterX()
    options:SetY(z*6)
    options.OnClick = function(obj)
        state_push('gameoptions', gameoptions)
    end
    local editor = loveframes.Create('imagebutton')
    editor:SetState('mainmenu')
    editor:SetImage('assets/images/editor.png')
    editor:SetText('')
    editor:SizeToImage()
    editor:CenterX()
    editor:SetY(z*7)
    editor.OnClick = function(obj)
        state_push('mapeditor', mapeditor)
    end
    local exit = loveframes.Create('imagebutton')
    exit:SetState('mainmenu')
    exit:SetImage('assets/images/exit.png')
    exit:SetText('')
    exit:SizeToImage()
    exit:CenterX()
    exit:SetY(z*8)
    exit.OnClick = function(obj)
        running = false
    end
end

local function update(dt)
    loveframes.update(dt)
    return running
end

local function draw()
    loveframes.draw()
end

local function keypressed(key, unicode)
    loveframes.keypressed(key, unicode)
    if key == 'escape' then running = false end
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
