
socket = require('socket')
vector = require('libs.vector')

require('engine.map')
require('engine.tileset')
require('engine.entities')
require('engine.network')
require('engine.serverprotocol')

globals = {}

local function gettime()
    return socket.gettime()
end

local function sleep(seconds)
    socket.sleep(seconds)
end

function load()
    math.randomseed(os.time())
    globals.running = false
    globals.max_fps = 20
    globals.fps = 0
    globals.tick = gettime()
    globals.dt = 0
    globals.starttick = 0
    globals.stoptick = 0

    globals.pos = vector(0,0)
    globals.jumpdist = 10
    globals.jumppower = 0
    globals.jumpcharge = 10

    map_file('ship.dat')
    load_entities()

    network_server('*', 8123)
    return true
end

local function new_clients()
    id = network_new_clients()
    if id then
        connect_player(id)
    end
end

local function bad_clients()
    for id in network_bad_clients() do
        disconnect_player(id)
    end
end

local function update_jump_power()
    if globals.jumppower == 100 then return end
    globals.jumppower = globals.jumppower + globals.jumpcharge
end

local function loop_once()
    network_read_all()
    new_clients()
    bad_clients()
    system_update('parse_clients')

    globals.tick = globals.tick + (1 / globals.max_fps)
    globals.dt = globals.tick - gettime()
    if globals.dt > 0 then sleep(globals.dt) end
    if globals.dt > 1 then
        globals.fps = globals.max_fps / globals.dt
    else
        globals.fps = globals.max_fps - globals.dt
    end

    -- update gameworld once a sec
    globals.stoptick = gettime()
    if (globals.stoptick - globals.starttick) > 0.5 then
        globals.starttick = globals.stoptick
        system_update('move_clients')
        system_update('example')
        system_update('air')
        update_jump_power()
    end
    network_send_all()
end

function main()
    print('Loading...')
    if not load() then return end

    print('Server is now online.')
    globals.running = true
    while globals.running do
        loop_once()
    end
end

main()
