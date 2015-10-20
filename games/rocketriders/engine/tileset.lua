
local loaded = false
tileset_size = 32
tileset = {}

tileset[-1] = {title='ERROR', x=0, y=0,
}
tileset[0] = { title='Empty', x=1, y=0,
    space=true,
}

tileset[1] = { title='Floor', x=3, y=0,
    group=true,
}
tileset[2] = {title='Wall', x=2, y=0,
    block=true,
    wall=true,
}
tileset[3] = {title='Wall-n', x=1, y=3,
    block=true,
    wall=true,
}
tileset[4] = {title='Wall-s', x=0, y=3,
    block=true,
    wall=true,
}
tileset[5] = {title='Wall-e', x=3, y=3,
    block=true,
    wall=true,
}
tileset[6] = {title='Wall-w', x=2, y=3,
    block=true,
    wall=true,
}

tileset[7] = { title='Teleportpad', x=0, y=4,
    spawn=true,
    group=true,
}
tileset[8] = { title='NavCon', x=1, y=4,
    block=true,
    navcon=true,
    group=true,
}
tileset[9] = { title='Doorclosed', x=2, y=4,
    block=true,
    door=true,
    open=10,
}
tileset[10] = {title='Dooropen', x=3, y=4,
    door=true,
    closed=9,
    group=true,
}

function load_tiles()
    if loaded then return end
    tilesheet = love.graphics.newImage('assets/tiles.png')
    tilebatch = love.graphics.newSpriteBatch(tilesheet, 21*21)

    for _, tile in pairs(tileset) do
        local x, y = tile.x, tile.y
        if x == nil then x=0 end
        if y == nil then y=0 end
        tile.quad = love.graphics.newQuad(x*tileset_size, y*tileset_size,
        tileset_size, tileset_size,
        tilesheet:getWidth(), tilesheet:getHeight())
    end
    loaded = true
end
