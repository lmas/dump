
-- Hack to load modules outside this directory
package.path = "../?.lua;" .. package.path

Camera = require('libs.hump.camera')

function love.load()
    -- Global constants
    screen_width = love.graphics.getWidth()
    screen_height = love.graphics.getHeight()
    world_x, world_y = 1000, 1000

    acceleration = 0.1
    walking_speed = 3
    running_speed = 8
    velocity_threshold = 0.1

    -- Player constants
    player_collided = false
    player_size = 50
    player_speed = running_speed
    player_x, player_y = 100,100
    velocity_x, velocity_y = 0, 0
    target_velocity_x, target_velocity_y = 0, 0

    love.keyboard.setKeyRepeat(0.5, 0.05)
    camera = Camera(0, 0)
end

function love.update(time_delta)
    update_player_keys(time_delta)
    update_player_position(time_delta)
end

function love.draw()
    -- Set up the drawing stack for moving around the camera and drawing the world
    camera:attach()

    -- Keeping it simple and only drawing a color grid as the world
    draw_color_grid()

    -- Done drawing now, reset the drawing stack so we can draw the GUI more
    -- easily.
    camera:detach()

    draw_gui()
end

--------------------------------------------------------------------------------
-- Update utils

function update_player_keys(time_delta)
    -- Shortcut
    local key = love.keyboard.isDown

    if key('escape') then
        love.event.quit()
    end

    -- Movement keys
    -- Check the velocity first and prevent the keys from interrupting each other
    if target_velocity_x < 1 and key('a') then
        target_velocity_x = -player_speed
    elseif key('d') then
        target_velocity_x = player_speed
    else
        target_velocity_x = 0
    end

    -- Same thing here
    if target_velocity_y < 1 and key('w') then
        target_velocity_y = -player_speed
    elseif key('s') then
        target_velocity_y = player_speed
    else
        target_velocity_y = 0
    end

    -- Running key
    if key('lshift') then
        player_speed = walking_speed
    else
        player_speed = running_speed
    end

end

function update_player_position(time_delta)
    velocity_x = acceleration * target_velocity_x + (1 - acceleration) * velocity_x
    velocity_y = acceleration * target_velocity_y + (1 - acceleration) * velocity_y

    -- Prevent "wobbling" the player, when low on velocity
    if math.abs(velocity_x) < velocity_threshold then
        velocity_x = 0
    end
    if math.abs(velocity_y) < velocity_threshold then
        velocity_y = 0
    end

    local tmp_x = player_x + (velocity_x )
    local tmp_y = player_y + (velocity_y )
    if check_collision(tmp_x, tmp_y) then
        player_x, player_y = tmp_x, tmp_y
    else
        player_collided = true
        velocity_x, velocity_y = 0, 0
    end

    -- We're not moving around the player on screen, only the camera
    camera:lookAt(player_x, player_y)
end

function check_collision(pos_x, pos_y)
    -- Keeping it simple for now, make sure the circle is inside the grid
    if (pos_x - player_size < -world_x or pos_x + player_size > world_x) then
        return false
    elseif (pos_y - player_size < - world_y or pos_y + player_size > world_y) then
        return false
    end
    return true
end

--------------------------------------------------------------------------------
--Drawing utils

function draw_gui()
    -- Draw the player in the middle of the screen
    love.graphics.setColor(255,255,255, 100)
    love.graphics.circle('fill', screen_width/2, screen_height/2, player_size, 100)

    -- Grab what world position the mouse is pointing at
    local mouse_x, mouse_y = camera:mousepos()

    -- Draw some debug info
    local fps = 'FPS: ' .. love.timer.getFPS()
    local pos = 'POSITION: '..player_x..', '..player_y
    local vel = 'VELOCITY: '..velocity_x..', '..velocity_y
    local mouse = 'MOUSE: '..mouse_x..', '..mouse_y
    love.graphics.setColor(255,255,255,255)
    love.graphics.print(fps, 10, 10)
    love.graphics.print(pos, 10, 10 * 3)
    love.graphics.print(vel, 10, 10 * 5)
    love.graphics.print(mouse, 10, 10 * 7)

    -- Draw some mouse thingy, at the actuall mouse position
    local mouse_x, mouse_y = love.mouse.getX(), love.mouse.getY()
    love.graphics.setColor(255,0,0,255)
    love.graphics.circle('fill', mouse_x, mouse_y, 10, 100)

    -- And some collision warning all over the screen
    if player_collided then
        player_collided = False
        love.graphics.setColor(255, 0, 0, 50)
        love.graphics.rectangle('fill', 0, 0, screen_width, screen_height)
    end
end

function draw_color_grid()
    local cell_size = 50
    local r, g, b = 0, 0, 0
    for grid_x = -world_x, world_x, cell_size do
        if grid_x < 0 then
            r, g, b = 255, 0, 0
        else
            r, g, b = 0, 255, 0
        end
        love.graphics.setColor(r, g, b)
        love.graphics.line(grid_x, -world_y, grid_x, world_y)
    end
    for grid_y = -world_y, world_y, cell_size do
        if grid_y < 0 then
            r, g, b = 0, 0, 255
        else
            r, g, b = 255, 255, 255
        end
        love.graphics.setColor(r, g, b)
        love.graphics.line(-world_x, grid_y, world_x, grid_y)
    end
end

