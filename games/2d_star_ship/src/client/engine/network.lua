
local socket = require('socket')

local client = {}


function client:connect(host, port)
    self.socket = socket.tcp()
    self.socket:settimeout(30)
    local ret, err = self.socket:connect(host, port)
    if not ret then
        return nil, err
    end

    self.socket:settimeout(0)
    return true, ''
end

function client:disconnect()
    self.socket:close()
end

function client:write(data)
    self.socket:send(data..'\n')
end

function client:read(chunk_size)
    local chunk_size = chunk_size or 5120
    local packet = ''

    local data, err, partial = self.socket:receive(chunk_size)
    while data do
        packet = packet .. data
        data, err, partial = self.socket:receive(chunk_size)
    end

    if not data and partial ~= '' then
        packet = packet .. partial
    elseif not data and err ~= 'timeout' then
        return nil, err
    end

    if packet then
        return packet, ''
    end
    return nil, 'no data'
end


return client
