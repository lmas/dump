
local commands = {}
buffer = ''

function parse_cmds()
    while true do
        local data, err = network_client_read(1024)
        if data == nil then
            if err == 'closed' then
                network_client_disconnect()
                return true
            elseif err == 'timeout' then
                break
            else
                break
            end
        end
        buffer = buffer .. data
    end
    for line in string.gmatch(buffer, '(%C+)\1') do
        local cmd, arg = string.match(line, '(%w+):(.*)')
        local tmp = commands[cmd]
        if tmp then
            tmp(arg)
        else
            print('Bad pkt: '..line)
        end
    end
    buffer = ''
end

-----------------------------------------------------------------------

function send_msg(msg)
    network_client_write('msg:'..msg)
end

-----------------------------------------------------------------------

function cmd_handshakeok(arg)
    network_client_write('login:')
    socket.sleep(0.1)
end

function cmd_loginok(arg)
end

function cmd_msg(arg)
    chat_msg(arg)
end

function cmd_map(arg)
    map_load(arg)
    globals.entities = {}
end

function cmd_ent(arg)
    local ent = {}
    for k, t, v in string.gmatch(arg, '(%w+)=(%a)([%w+-:<>]+),') do
        if t == 'n' then
            v = tonumber(v)
        elseif t == 'v' then
            local x, y = string.match(v, '([%d-]+):([%d-]+)')
            v = vector(tonumber(x), tonumber(y))
        else
            v = tostring(v)
        end
        ent[k] = v
    end
    ent.frac = 0
    ent.direction = vector(0,1)
    globals.entities[ent] = true
    map_insert(ent)
end

function cmd_entremove(arg)
    for ent in entity_loop() do
        if ent.id == arg then
            globals.entities[ent] = nil
            return
        end
    end
end
function entity_loop()
    return pairs(globals.entities)
end

function cmd_player(arg)
    for ent in entity_loop() do
        if ent.id == arg then
            globals.player = ent
            if globals.name then
                network_client_write('name:'..tostring(globals.name))
            end
            return
        end
    end
end

function cmd_move(arg)
    local id, x, y = arg:match('(%w+),([%d-]+),([%d-]+)')
    local pos = vector(tonumber(x), tonumber(y))
    for ent in entity_loop() do
        if ent.id == id then
            ent.direction = pos
            map_remove(ent)
            ent.pos = ent.pos + pos
            map_insert(ent)
            ent.frac = tileset_size --/ 2
            return
        end
    end
end

function cmd_name(arg)
    local id, name = arg:match('(%w+),(%C+)')
    for ent in entity_loop() do
        if ent.id == id then
            ent.title = name
            return
        end
    end
end

function cmd_dir(arg)
    local id, x, y = arg:match('(%w+),([%d-]+),([%d-]+)')
    local dir = vector(tonumber(x), tonumber(y))
    for ent in entity_loop() do
        if ent.id == id then
            ent.direction = dir
            return
        end
    end
end

function cmd_door(arg)
    local x, y = arg:match('([%d-]+),([%d-]+)')
    local pos = vector(tonumber(x), tonumber(y))
    local cell = map_get(pos)
    cell.open = not cell.open
    if cell.open then
        cell.tile = tileset[cell.tile].open
    else
        cell.tile = tileset[cell.tile].closed
    end
end

function cmd_navcon(arg)
    navframe:SetVisible(true)
    local x, y, drive = arg:match('([%d-]+),([%d-]+),(%d*)')
    navpos:SetText(string.format('Ship Position: %d, %d', x, y))
    navdrive:SetText(string.format('Jump drive: %d%%', drive))
end

commands = {
    handshakeok = cmd_handshakeok,
    loginok = cmd_loginok,
    msg = cmd_msg,
    map = cmd_map,
    clear = cmd_clear,
    ent = cmd_ent,
    entremove = cmd_entremove,
    player = cmd_player,
    move = cmd_move,
    name = cmd_name,
    dir = cmd_dir,
    door = cmd_door,
    navcon = cmd_navcon,
}
