
function log(data)
    if settings['DEBUG'] == true then
        print('DEBUG: '..data)
    end
end

function trim(str)
    -- Source: http://lua-users.org/wiki/CommonFunctions
    -- from PiL2 20.4
    return (str:gsub("^%s*(.-)%s*$", "%1"))
end

function get_lines(str)
    return str:gmatch("[^\n]+") 
end

function split(str, pat)
    -- Source: http://lua-users.org/wiki/SplitJoin
    local t = {}
    local fpat = "(.-)" .. pat
    local last_end = 1
    local s, e, cap = str:find(fpat, 1)
    while s do
        if s ~= 1 or cap ~= "" then
            table.insert(t,cap)
        end
        last_end = e+1
        s, e, cap = str:find(fpat, last_end)
    end
    if last_end <= #str then
        cap = str:sub(last_end)
        table.insert(t, cap)
    end
    return t
end

