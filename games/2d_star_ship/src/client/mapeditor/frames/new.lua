
loveframes = require('libs.loveframes')


STATE = 'newmap'

local frame = loveframes.Create('frame')
frame:SetState(STATE)
frame:SetName('Create new, empty map?')
frame:SetSize(200,60)
frame:Center()

local yesbutton = loveframes.Create('button', frame)
yesbutton:SetWidth(90)
yesbutton:SetPos(5, 30)
yesbutton:SetText('OK')
yesbutton.OnClick = function(obj)
    tilemap:clear()
    loveframes.SetState('main')
end

local nobutton = loveframes.Create('button', frame)
nobutton:SetWidth(90)
nobutton:SetPos(105, 30)
nobutton:SetText('Cancel')
nobutton.OnClick = function(obj)
    loveframes.SetState('main')
end


return {
    state = STATE,
    frame = frame,
}
