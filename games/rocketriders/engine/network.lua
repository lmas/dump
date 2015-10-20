
--[[
this module performs all the low level networking, outside of the game.
this includes handling new clients and logins. parsing game commands
and such is done inside the protocol module instead.
]]--

socket = require('socket')
local sock

function network_client(host, port)
    sock = socket.tcp()
    sock:settimeout(30)
    local ret, err = sock:connect(host, port)
    if not ret then
        print('ERROR while connecting to the server: '..err)
        return false, err
    end
    sock:settimeout(0)
    socket.sleep(1)
    return true, err
end

function network_client_disconnect()
    sock:close()
end

function network_client_read()
    local data, err = sock:receive()
    return data, err
end

function network_client_write(data)
    sock:send(data..'\1\n')
end

----------------------------------------------------------------------

local clients = {}
local bad_clients = {}

function network_server(host, port)
    sock = socket.tcp()
    sock:settimeout(0)
    sock:setoption('keepalive', true)
    sock:setoption('reuseaddr', true)
    sock:setoption('tcp-nodelay', true)
    sock:bind(host, port)
    sock:listen(5)
end

local function network_id(ip, port)
    return tostring(ip)..':'..tostring(port)
end

local function network_new(id, ip, port)
    print('Connect: '..id)
    return {ip=ip, port=port, bi='', bo={}}
end

function network_del(id)
    if not network_valid(id) then return end
    print('Disconnect: '..id)
    clients[id].sock:close()
    clients[id] = nil
    table.insert(bad_clients, id)
end

function network_valid(id)
    if clients[id] == nil then return false end
    return true
end

function network_read_all()
    for id, tmp in pairs(clients) do
        local data, err = tmp.sock:receive()
        if data then
            clients[id].bi = clients[id].bi .. data
        else
            if err == 'closed' then network_del(id) end
        end
    end
end

function network_send_all()
    for id, tmp in pairs(clients) do
        if tmp.bo[1] then
            local data = tmp.bo
            tmp.bo = {}
            for _, line in pairs(data) do
                local total, cur, hits = line:len(), 0, 0
                while true do
                    local bytes, err = tmp.sock:send(line, cur)
                    if err == 'closed' then
                        break
                    elseif err == 'timeout' then
                        socket.sleep(0.01)
                        hits = hits + 1
                        if hits > 20 then
                            -- TODO: make sure not disconnection
                            print('Warning: '..id..' timed out during a send.')
                            table.insert(clients[id].bo, line)
                            return
                        end
                    end
                    if bytes ~= nil then
                        cur = cur + bytes
                        if cur > total then break end
                    end
                end
            end
        end
    end
end

function network_read(id)
    if not network_valid(id) then return end
    local tmp = clients[id].bi
    clients[id].bi = ''
    return tmp
end

function network_send(id, data)
    if not network_valid(id) then return end
    table.insert(clients[id].bo, data .. '\1\n')
end

function network_new_clients()
    local new = sock:accept()
    if not new then return end
    new:settimeout(0)
    ip, port = new:getpeername()
    local id = network_id(ip, port)
    if not network_valid(id) then
        clients[id] = network_new(id, ip, port)
        clients[id].sock = new
        return id
    end
end

function network_bad_clients()
    return function()
        for i, id in pairs(bad_clients) do
            table.remove(bad_clients, i)
            return id
        end
    end
end
