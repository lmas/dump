
require('libs.anal')

animation_size = 32
animations = {}
animations['computeron'] = 'assets/animations/computeron.png'
animations['playerleft'] = 'assets/animations/playerleft.png'
animations['playerright'] = 'assets/animations/playerright.png'
animations['playerdown'] = 'assets/animations/playerdown.png'
animations['playerup'] = 'assets/animations/playerup.png'
animations['bugup'] = 'assets/animations/bugup.png'
animations['bugdown'] = 'assets/animations/bugdown.png'
animations['bugleft'] = 'assets/animations/bugleft.png'
animations['bugright'] = 'assets/animations/bugright.png'
--animations['dooropen'] = 'assets/animations/dooropen.png'

function load_animations()
    if loaded then return end
    for title, anim in pairs(animations) do
        animations[title] = newAnimation(love.graphics.newImage(anim),
        animation_size, animation_size, 0.5, 0)
    end
    loaded = true
end

function update_animations(dt)
    for _, anim in pairs(animations) do
        anim:update(dt)
    end
end
