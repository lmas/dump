
local entities = {}
local systems = {}

-- create a new entity
function entity_new(pos)
    local ent = {pos=pos}
    ent.id = entity_id(ent)
    entities[ent] = true
    map_insert(ent)
    return ent
end

-- remove an entity
function entity_del(ent)
    map_remove(ent)
    entities[ent] = nil
    for sys in pairs(systems) do
        system_del(sys, ent)
    end
end

function entity_id(ent)
    local tmp = string.gmatch(tostring(ent), 'table: (%w+)')()
    return tmp
end

function entity_count()
    return #entities
end

function entity_loop()
    return pairs(entities)
end

-- register a new system
function system_new(title, func)
    systems[title] = {func = func, entities = {}}
end

-- add an entity to a system
function system_add(sys, ent)
    systems[sys].entities[ent] = true
end

-- remove an entity from a system
function system_del(sys, ent)
    systems[sys].entities[ent] = nil
end

-- update specific system
function system_update(sys)
    local tmp = systems[sys]
    if tmp ~= nil then tmp.func(tmp.entities) end
end

-- return all entities in a system
function system_each_entity(sys)
    local tmp = systems[sys]
    if tmp ~= nil then return tmp.entities end
    return nil
end

----------------------------------------------------------------------

function entity_move(ent, pos, ignoreblock)
    pos.x = math.max(math.min(1, pos.x), -1)
    pos.y = math.max(math.min(1, pos.y), -1)
    if map_move(ent, pos, ignoreblock) then
        send_all('move:'..ent.id..','..tostring(pos.x)..','..tostring(pos.y))
        return true
    else
        send_all('dir:'..ent.id..','..tostring(pos.x)..','..tostring(pos.y))
    end
    return false
end

function load_entities()
    entities = {}
    systems = {}
    system_new('parse_clients', function(entities)
        for ent in pairs(entities) do
            parse_cmds(ent)
        end
    end)
    system_new('move_clients', function(entities)
        for ent in pairs(entities) do
            if ent.velocity ~= vector(0,0) then
                entity_move(ent, ent.velocity)
                ent.velocity = vector(0,0)
            end
        end
    end)
    system_new('air', function(entities)
        for ent in pairs(entities) do
            local group = map_group(ent.pos)
            if group.air == nil then
                ent.air = math.max(ent.air - 5, 0)
            elseif ent.air < 100 then
                ent.air = math.min(ent.air + 5, 100)
            end
            print(group.air, ent.air)
        end
    end)
    system_new('example', function(entities)
        for ent in pairs(entities) do
            if math.random() < 0.25 then
                if math.random(0,1) == 0 then
                    delta= vector(math.random(-1,1), 0)
                else
                    delta = vector(0, math.random(-1,1))
                end
                entity_move(ent, delta)
            end
        end
    end)

    for i = 1, 50 do -- add some example entities
        local pos = vector(math.random(1, 20), math.random(1, 20))
        local ent = entity_new(pos)
        ent.title = 'Spaceman'
        ent.up = 'playerup'
        ent.down = 'playerdown'
        ent.left = 'playerleft'
        ent.right = 'playerright'
        ent.blocks = true
        system_add('example', ent)
    end
    for i = 1, 50 do -- add some example entities
        local pos = vector(math.random(1, 20), math.random(1, 20))
        local ent = entity_new(pos)
        ent.title = 'Bug'
        ent.up = 'bugup'
        ent.down = 'bugdown'
        ent.left = 'bugleft'
        ent.right = 'bugright'
        ent.blocks = false
        system_add('example', ent)
    end
end

function connect_player(id)
    local new = entity_new(map_spawn())
    new.title = 'player'--new.id
    new.client = id
    new.state = 0
    new.up = 'playerup'
    new.down = 'playerdown'
    new.left = 'playerleft'
    new.right = 'playerright'
    new.sight = 10
    new.blocks = true
    new.velocity = vector(0,0)
    new.air = 100
    system_add('parse_clients', new)
    system_add('move_clients', new)
    system_add('air', new)
end

function disconnect_player(id)
    for ent in entity_loop() do
        if ent.client == id then
            entity_del(ent)
            send_msg_all('Client disconnect: '..id)
            send_all('entremove:'..ent.id)
            return
        end
    end
end
