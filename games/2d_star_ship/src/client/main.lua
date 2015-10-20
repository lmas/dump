
Camera = require('libs.hump.camera')

client = require('engine.network')
entity = require('engine.entity')
tilemap = require('engine.tilemap')
require('engine.utils')

settings = {
    DEBUG = true,
}


function love.load()
    screen_width, screen_height = love.graphics.getWidth(), love.graphics.getHeight()
    player = {}
    player.x, player.y = 0, 0
    player.gid = nil
    player.movingx, player.movingy = 0, 0
    player.oldx, player.oldy = 0, 0

    camera = Camera(0,0)
    entity:init()
    tilemap:init()
    love.keyboard.setKeyRepeat(0.5, 0.05)

    ret, err = client:connect('localhost', 12345)
    if not ret then
        print('NETWORK ERROR: '..err..', Couldn\' connect to the server!')
        love.event.quit()
    end
end

function love.update(dt)
    network_read()
    update_player_info()
    update_camera()
    update_keys()
    tilemap:update() -- TODO: only do it once, after the map was sent
end

function love.draw()
    camera:attach()

    tilemap:draw()

    love.graphics.setColor(255,255,255)
    for gid, ent in entity:iter() do
        if gid ~= player.gid then
            x, y = ent['pos'][1], ent['pos'][2]
            x = x * tilemap.tile_size + (tilemap.tile_size/2)
            y = y * tilemap.tile_size + (tilemap.tile_size/2)
            love.graphics.circle('fill', x, y, 10, 100)
        end
    end

    camera:detach()

    love.graphics.setColor(255,0,0, 100)
    love.graphics.rectangle('fill',
        screen_width / 2, screen_height / 2,
        tilemap.tile_size, tilemap.tile_size
    )

    love.graphics.setColor(255,255,255)
    love.graphics.print('POS: '..player.x..', '..player.y, 10, 10)
    love.graphics.print('FPS: ' .. love.timer.getFPS(), 10, 30)
end

function love.quit()
    client:disconnect()
end

function update_keys()
    local key = love.keyboard.isDown

    if key('escape') then
        love.event.quit()
    end

    if player.movingy < 1 and key('w') then
        player.movingy = -1
    elseif key('s') then
        player.movingy = 1
    else
        player.movingy = 0
    end

    if player.movingx < 1 and key('a') then
        player.movingx = -1
    elseif key('d') then
        player.movingx = 1
    else
        player.movingx = 0
    end

    if player.oldx ~= player.movingx or player.oldy ~= player.movingy then
        player.oldx, player.oldy = player.movingx, player.movingy
        client:write('mov '..player.movingx..','..player.movingy)
    end
end

function network_read()
    local data, err = client:read()
    if not data then
        print('NETWORK ERROR: '..err)
        love.event.quit()
        return nil
    end
    for line in get_lines(data) do
        cmd = trim(line:sub(1, 3))
        arg = trim(line:sub(4))

        if cmd == 'msg' then
            print(arg)
        elseif cmd == 'new' then
            ent = arg:match('(%d+) ')
            x, y = arg:match('%d+ (%C+),(%C+)')
            entity:add(ent, {x,y})
        elseif cmd == 'pla' then
            player.gid = arg
        elseif cmd == 'del' then
            entity:del(arg)
        elseif cmd == 'pos' then
            ent = arg:match('(%d+) ')
            x, y = arg:match('%d+ (%C+),(%C+)')
            entity:update(ent, 'pos', {x,y})
        elseif cmd == 'map' then
            log('Loading new tilemap '..arg..'...')
            tilemap:clear(arg)
        elseif cmd == 'til' then
            tilemap:parse_tile(arg)
        else
            log('UNKNOWN: "'..cmd..' '..arg..'"')
        end
    end
end

function update_player_info()
    local ent = entity:get(player.gid)
    if not ent then
        return
    end

    player.x, player.y = ent['pos'][1], ent['pos'][2]
end

function update_camera()
    camera:lookAt(player.x * tilemap.tile_size, player.y * tilemap.tile_size)
end

