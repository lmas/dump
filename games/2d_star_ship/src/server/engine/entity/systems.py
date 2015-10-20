
import logging
import math

from engine.utils import sign, Vector
from engine.entity.components import *

from engine.main import Game
GAME = Game.instance()


class Calc_movement(object):
    component_types = (CPosition, CVelocity)
    # Used to correct the calculations so everything runs at the same speed,
    # no matter the FPS. Or that's the theory, anyway...
    seconds_per_frame = 1.0 / GAME.settings.get('FPS', 30)
    max_speed = seconds_per_frame * GAME.settings.get('SPEED_MAX', 4)
    acceleration = seconds_per_frame / GAME.settings.get('SPEED_ACCELERATION', 0.1)
    # Prevents the player from "wobbling" when his current acceleration is below 0.1%
    velocity_threshold = acceleration / 100

    def update(self, manager, entities, dt=0):
        for entity in entities:
            self.update_entity(manager, entity)

    def update_entity(self, manager, entity):
        position = manager.get_component(entity, CPosition)
        velocity = manager.get_component(entity, CVelocity)
        current = velocity.current
        target = velocity.target

        if target == (0,0) and current == (0,0):
            # Do nothing if we're not moving
            return

        target_speed = target * self.max_speed
        x = self.acceleration * target_speed.x + (1 - self.acceleration) * current.x
        y = self.acceleration * target_speed.y + (1 - self.acceleration) * current.y

        # Prevent the wobbling
        if math.fabs(x) < self.velocity_threshold:
            x = 0
        if math.fabs(y) < self.velocity_threshold:
            y = 0

        # So we know to what tile we're travelling to
        direction = Vector(sign(x), sign(y))

        # Do collision check separately, mainly to avoid bugs when moving on
        # both axis at the same time.
        xaxis = position.vector + (x,0) + direction
        if GAME.tilemap.check_collision(GAME.settings.get('STARTMAP'), xaxis):
            x = 0

        yaxis = position.vector + (0,y) + direction
        if GAME.tilemap.check_collision(GAME.settings.get('STARTMAP'), yaxis):
            y = 0

        velocity.current = Vector(x, y)

        # If we haven't changed position at all, do nothing else
        if x == 0 and y == 0:
            return

        # Everything good, do the updates!
        position.vector += (x, y)
        net = manager.get_component(entity, CNetwork)
        net.client.server.broadcast(
            'pos {} {}'.format(entity, str(position.vector))
        )

class Parse_commands(object):
    component_types = (CNetwork)

    def update(self, manager, entities, dt=0):
        for entity in entities:
            self.update_entity(manager, entity)

    def update_entity(self, manager, entity):
        tmp = manager.get_component(entity, CNetwork)
        client = tmp.client
        if not tmp.data:
            return
        logging.debug('{} data: {}'.format(entity, tmp.data))
        cmd, arg = tmp.data[:3].lower(), tmp.data[3:].strip()
        tmp.data = ''

        if cmd == 'mov':
            x, y = arg.split(',')
            x, y = int(x.strip()), int(y.strip())
            velocity = manager.get_component(entity, CVelocity)
            velocity.target = Vector(x, y)
        elif cmd == 'qui':
            client.close()

