
local commands = {}

function parse_cmds(ent)
    data = network_read(ent.client)
    if not data then return end
    for line in string.gmatch(data, '(%C+)\1') do
        local cmd = string.match(line, '(%w+):')
        local arg = string.match(line, '%w+:(.+)')
        local tmp = commands[cmd]
        if tmp and tmp.state == ent.state then
            tmp.func(ent, arg)
        else
            if ent.state < 2 then
                network_del(ent.client)
            end
        end
    end
end

-----------------------------------------------------------------------

function send_game_world(ent)
    local pkt = ''

    -- map
    network_send(ent.client, 'map:'..map_save())

    -- doors
    for pos in map_loop_all() do
        local cell = map_get(pos)
        local tile = tileset[cell.tile]
        if tile.door and cell.open then
            network_send(ent.client, 'door:'..tostring(pos.x)..','..tostring(pos.y))
        end
    end

    -- entities
    local parts = {'id', 'title', 'sight', 'pos', 'up', 'down', 'left', 'right'}
    for e in entity_loop() do
        pkt = 'ent:'
        for _, k in pairs(parts) do
            v = e[k]
            if v ~= nil then
                local t = type(v)
                if t == 'string' then t = 's'
                elseif t == 'number' then t = 'n'
                --elseif t == 'boolean' then t = 'b'
                elseif isvector(v) then
                    v = tostring(v.x)..':'..tostring(v.y)
                    t = 'v'
                end
                pkt = pkt..tostring(k)..'='..t..tostring(v)..','
            end
        end
        if e.id == ent.id then
            send_all(pkt)
        else
            network_send(ent.client, pkt)
        end
    end
    pkt = ''

    -- player
    pkt = 'player:'..ent.id
    network_send(ent.client, pkt)
    send_msg_all('New client from '..ent.client)
end

function send_msg(ent, msg)
    network_send(ent.client, 'msg:'..msg)
end

function send_msg_all(msg)
    send_all('msg:'..msg)
end

function send_all(data)
    for ent in pairs(system_each_entity('parse_clients')) do
        network_send(ent.client, data)
    end
end

function send_navcon(ent)
    -- TODO: do this automagically
    local x, y = globals.pos.x, globals.pos.y
    local tmp = string.format('navcon:%d,%d,%d', x, y, globals.jumppower)
    network_send(ent.client, tmp)
end
-----------------------------------------------------------------------

function cmd_handshake(ent, arg)
    ent.state = 1
    network_send(ent.client, 'handshakeok:')
end

function cmd_login(ent, arg)
    ent.state = 2
    network_send(ent.client, 'loginok:')
    send_game_world(ent)
end

function cmd_quit(ent, arg)
    network_del(ent.client)
end

function cmd_msg(ent, arg)
    if not arg then return end
    send_msg_all(ent.title..': '..arg)
end

function cmd_noplayer(ent, arg)
    network_send(ent.client, 'player:'..ent.id)
end

function cmd_move(ent, arg)
    local x,y = arg:match('(%C+),(%C+)')
    x, y = tonumber(x), tonumber(y)
    if (x == nil) or (y == nil) then return end
    ent.velocity = vector(x,y)
end

function cmd_name(ent, arg)
    if arg == nil then return end
    ent.title = arg
    send_all('name:'..ent.id..','..arg)
end

function cmd_dir(ent, arg)
    send_all('dir:'..ent.id..','..arg)
end

local function toggle_door(pos)
    local cell = map_get(pos)
    cell.open = not cell.open
    send_all('door:'..tostring(pos.x)..','..tostring(pos.y))

    local check
    function check(pos)
        local group = map_group(pos)
        if group.air ~= true then return end
        group.air = nil
    end

    local suckair = false
    local function check_air(pos)
        local cell = map_get(pos)
        local tile = tileset[cell.tile]
        if tile.space == true then suckair = true end
    end
    check_air(pos + vector(1,0))
    check_air(pos + vector(-1,0))
    check_air(pos + vector(0,1))
    check_air(pos + vector(0,-1))
    if suckair == false then return end
    check(pos + vector(1,0))
    check(pos + vector(-1,0))
    check(pos + vector(0,1))
    check(pos + vector(0,-1))
end

function cmd_use(ent, arg)
    local x, y = arg:match('([%d-]+),([%d-]+)')
    local pos = ent.pos + vector(x, y)
    local cell = map_get(pos)
    local tile = tileset[cell.tile]
    if tile.door then
        toggle_door(pos)
    elseif tile.navcon then
        send_navcon(ent)
    end
end

function cmd_jump(ent, arg)
    if globals.jumppower < 100 then
        send_msg(ent, 'Error: jump drive not fully charged!')
        return
    end
    local x,y = arg:match('([%d-]+),([%d-]+)')
    if (x == nil) or (y == nil) then return end
    local new = vector(x, y)
    local dist = math.ceil(new:dist(globals.pos))
    if dist <= globals.jumpdist then
        globals.pos = new
        globals.jumppower = 0
        send_navcon(ent)
        send_msg_all('Jump completed!')
    else
        send_msg(ent, string.format('Error: Jump distance (%f sectors) is too great!', dist))
    end
end

commands = {
    handshake = {state = 0, func=cmd_handshake},
    login = {state = 1, func=cmd_login},
    quit = {state = 2, func=cmd_quit},
    msg = {state = 2, func=cmd_msg},
    noplayer = {state = 2, func=cmd_noplayer},
    move = {state = 2, func=cmd_move},
    name = {state = 2, func=cmd_name},
    dir = {state = 2, func=cmd_dir},
    use = {state = 2, func=cmd_use},
    jump = {state = 2, func=cmd_jump},
}
