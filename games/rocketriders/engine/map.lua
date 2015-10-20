
local map = {}
local groups = {}

function map_file(file)
    local ret, err = io.open(file, 'rb')
    if not ret then
        print('ERROR while loading map: '..err)
        return false, err
    else
        local data = ret:read('*a')
        ret:close()
        map_load(data)
        data = nil
        return true
    end
end

function map_load(data)
    if data == nil then
        return
    end
    map_clear()
    for x, y, tile in string.gmatch(data, '(%d+),(%d+),(%d+)|') do
        x, y, tile = tonumber(x), tonumber(y), tonumber(tile)
        local cell = {tile=tile, group=0, entities={}}
        map_set(vector(x,y), cell)
    end
    map_areas()
end

function map_save()
    local pkt = ''
    for pos in map_loop_all() do
        local cell = map_get(pos)
        pkt = pkt .. string.format('%d,%d,%d|', pos.x, pos.y, cell.tile)
    end
    return pkt
end

function map_clear()
    map = {}
end

function map_get(pos)
    local tmp = pos.x ..','..pos.y
    local cell = map[tmp]
    if cell ~= nil then
        return map[tmp]
    else
        local cell = {tile=0, group=0, entities={}}
        map_set(pos, cell)
        return cell
    end
end

function map_set(pos, cell)
    local tmp = pos.x ..','..pos.y
    map[tmp] = cell
    local tile = tileset[cell.tile]
    if tile.door then
        map[tmp].open = false
    end
end

function map_del(pos)
    local tmp = pos.x ..','..pos.y
    map[tmp] = nil
end

function map_loop(start, stop)
    if start == nil then start = vector(1,1) end
    if stop == nil then stop = vector(100,100) end
    return coroutine.wrap(function()
    for x=start.x, stop.x do
        for y=start.y, stop.y do
            local pos = vector(x,y)
            local cell = map_get(pos)
            if cell.tile ~= 0 then coroutine.yield(pos) end
        end
    end
    end)
end

function map_loop_all()
    return coroutine.wrap(function()
        for tmp, cell in pairs(map) do
            if cell.tile ~= 0 then
                local x, y = tmp:match('([%d-+]+),([%d-+]+)')
                local pos = vector(x, y)
                coroutine.yield(pos)
            end
        end
    end)
end

function map_group(pos)
    local cell = map_get(pos)
    if cell.group == 0 then return {} end
    local group = groups[cell.group]
    return group
end

function map_areas()
    groups = {}
    local group = 1
    local visited = {}

    local check
    function check(pos)
        local tmp = tostring(pos.x)..','..tostring(pos.y)
        if visited[tmp] ~= nil then return end
        visited[tmp] = true
        local cell = map_get(pos)
        local tile = tileset[cell.tile]
        if tile.group ~= true then return end
        cell.group = group
        check(pos + vector(1,0))
        check(pos + vector(-1,0))
        check(pos + vector(0,1))
        check(pos + vector(0,-1))
        return true
    end

    groups[group] = {air=true, power=true}
    for pos in map_loop_all() do
        if check(pos) then
            group = group + 1
            groups[group] = {air=true, power=true}
        end
    end
end

function map_spawn()
    local list = {}
    for pos in map_loop_all() do
        local cell = map_get(pos)
        local tile = tileset[cell.tile]
        if tile.spawn then
            table.insert(list, pos)
        end
    end
    if #list > 0 then
        return list[math.random(#list)]
    else
        return vector(0,0)
    end
end

function map_insert(ent)
    local cell = map_get(ent.pos)
    cell.entities[ent] = true
end

function map_remove(ent)
    local cell = map_get(ent.pos)
    cell.entities[ent] = nil
    if cell.tile == 0 and #cell.entities < 1 then
        map_del(ent.pos)
    end
end

function map_move(ent, delta, ignoreblock)
    local pos = ent.pos + delta
    local cell = map_get(pos)
    local tile = tileset[cell.tile]
    local block = false
    if ignoreblock ~= true then
        for e in pairs(cell.entities) do
            if e.blocks ~= nil and e.blocks == true then
                block = true
                break
            end
        end
    end
    if not block then
        if (tile.door and cell.open) or (not tile.block)then
            map_remove(ent)
            ent.pos = pos
            map_insert(ent)
            return pos
        end
    end
    return nil
end
