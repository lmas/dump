
loveframes = require('libs.loveframes')
Camera = require('libs.hump.camera')
tilemap = require('engine.tilemap')
save_frame = require('frames.save')
load_frame = require('frames.load')
new_frame = require('frames.new')


function love.load()
    selected_tile = 0
    map_changed = false

    love.keyboard.setKeyRepeat(0.5, 0.05)
    camera = Camera(0,0)
    tilemap:init()
    update_selected_tile()
    loveframes.SetState('main')
end

function love.update()
    tilemap:update()
    loveframes.update()
end

function love.textinput(text)
    loveframes.textinput(text)
end

function love.draw()
    camera:attach()

    love.graphics.setColor(255,255,255,255)
    tilemap:draw()

    -- Mark the current tile under the mouse
    mx, my = mouse_pos()
    tx, ty = mx * tilemap.tile_size, my * tilemap.tile_size
    love.graphics.setColor(100,0,0,100)
    love.graphics.rectangle('fill', tx, ty, tilemap.tile_size, tilemap.tile_size)

    camera:detach()

    -- GUI
    cx, cy = camera_pos()
    mx, my = mouse_pos()
    gui_text = 'FPS: '..love.timer.getFPS()..'\n'
    gui_text = gui_text..'POS: '..cx..', '..cy..'\n'
    gui_text = gui_text..'MOUSE: '..mx..', '..my..'\n'
    gui_text = gui_text..'TILE: #'..selected_tile..' '..tile_data.name..'\n'
    love.graphics.setColor(255,255,255,255)
    love.graphics.print(gui_text, 0, 0)

    loveframes.draw()
end

function love.keypressed(key)
    if loveframes.GetState() == 'main' then
        if key == 'escape' then
            love.event.quit()
        elseif key == 's' and love.keyboard.isDown('lctrl') then
            loveframes.SetState(save_frame.state)
            save_frame.input:SetFocus(true)
        elseif key == 'o' and love.keyboard.isDown('lctrl') then
            loveframes.SetState(load_frame.state)
            load_frame.input:SetFocus(true)
        elseif key == 'n' and love.keyboard.isDown('lctrl') then
            loveframes.SetState(new_frame.state)
        elseif key == 'w' then
            camera:move(0, -tilemap.tile_size)
        elseif key == 's' then
            camera:move(0, tilemap.tile_size)
        elseif key == 'a' then
            camera:move(-tilemap.tile_size,0)
        elseif key == 'd' then
            camera:move(tilemap.tile_size,0)
        elseif key == 'q' then
            selected_tile = math.min(math.max(selected_tile - 1, 0), #tilemap.tileset.set)
            update_selected_tile()
        elseif key == 'e' then
            selected_tile = math.min(math.max(selected_tile + 1, 0), #tilemap.tileset.set)
            update_selected_tile()
        end
    else
        loveframes.keypressed(key)
        if key == 'escape' then
            loveframes.SetState('main')
        end
    end
end

function love.keyreleased(key)
    loveframes.keyreleased(key)
end

function love.mousepressed(x, y, button)
    loveframes.mousepressed(x, y, button)

    mx, my = mouse_pos()
    if button == 'l' then
        local tile = {tile=selected_tile}
        tilemap:set(mx, my, tile)
    elseif button == 'r' then
        selected_tile = tilemap:get(mx, my).tile
        update_selected_tile()
    end
end

function love.mousereleased(x, y, button)
    loveframes.mousereleased(x, y, button)
end

function update_selected_tile()
    tile_data = tilemap.tileset.set[selected_tile]
end

function camera_pos()
    local x, y = camera:pos()
    return x / tilemap.tile_size, y / tilemap.tile_size
end

function mouse_pos()
    local mx, my = camera:mousepos()
    mx, my = math.floor(mx / tilemap.tile_size), math.floor(my / tilemap.tile_size)
    return mx, my
end
