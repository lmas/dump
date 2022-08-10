
local entity = {}


function entity:init()
    self.entities = {}
end

function entity:add(gid, pos)
    self.entities[gid] = {
        pos = pos,
    }
end

function entity:del(gid)
    self.entities[gid] = nil
end

function entity:update(gid, key, value)
    self.entities[gid][key] = value
end

function entity:iter()
    return pairs(self.entities)
end

function entity:get(gid)
    return self.entities[gid]
end


return entity
