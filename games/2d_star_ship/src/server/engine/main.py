
import logging
from settings import settings
logging.basicConfig(
    format = settings.get('LOGGING_FORMAT'),
    level = settings.get('LOGGING_LEVEL'),
)

from tornado.ioloop import IOLoop, PeriodicCallback
from engine.entity.components import CNetwork, CPosition, CVelocity


class Game(object):
    _instance = None

    @staticmethod
    def instance():
        # The singleton pattern, mainly to be used inside the engine itself.
        # I kinda stole this idea from tornado's IOLoop...
        if Game._instance == None:
            Game._instance = Game()
        return Game._instance

    def __init__(self):
        Game._instance = self

        logging.info('Loading server...')
        self.settings = settings
        self.running = False
        self.ioloop = IOLoop.instance()
        PeriodicCallback(self.update, 1000.0 / self.settings.get('FPS', 30), self.ioloop).start()

        # Hold off loading modules until now, too solve the problems with
        # circular imports
        from engine.entity.main import EntityManager
        from engine.entity.systems import Calc_movement, Parse_commands
        from engine.tilemaps import Tilemap
        from engine.network import NetworkServer

        self.entities = EntityManager()
        self.entities.add_system(Calc_movement().update, CPosition)
        self.entities.add_system(Parse_commands().update, CPosition)

        self.tilemap = Tilemap()
        for tid, mapfile in self.settings.get('TILEMAPS', {}).items():
            self.tilemap.load_instance(tid, mapfile)

        self.server = NetworkServer(
            self.client_connected,
            self.client_disconnected,
            self.client_data,
            self.settings.get('SERVER_PORT'),
            self.settings.get('SERVER_ADDRESS'),
        )

        logging.info('Server loaded.')

    def start(self):
        logging.info('Running server...')
        self.running = True
        try:
            self.ioloop.start()
        except KeyboardInterrupt:
            logging.info('Keyboard interrupt.')
            self.stop()

    def stop(self):
        logging.info('Stopping server...')
        self.running = False
        self.ioloop.stop()

    def update(self):
        if not self.running: return
        self.entities.update_systems(1)

    def client_connected(self, client):
        if not self.running: return
        entity = self.entities.add_entity()
        tmp = CNetwork()
        # map the client with the entity and store refs on both of 'em
        tmp.client = client
        self.entities.add_component(entity, tmp)
        self.entities.add_component(entity, CPosition())
        self.entities.add_component(entity, CVelocity())
        client.entity = entity
        self.send_tilemap(client, self.settings.get('STARTMAP'))
        for c in self.server.clients:
            entity = c.entity
            pos = self.entities.get_component(entity, CPosition)
            client.write('new {} {}'.format(entity, str(pos.vector)))
        client.write('pla {}'.format(client.entity))
        pos = self.entities.get_component(entity, CPosition).vector
        self.server.broadcast_all_but('new {} {}'.format(entity, str(pos)), client)
        self.server.broadcast('msg Client connected: {}'.format(client.ip))

    def client_disconnected(self, client):
        if not self.running: return
        self.entities.del_entity(client.entity)
        self.server.broadcast('del {}'.format(client.entity))
        self.server.broadcast('msg Client disconnected: {}'.format(client.ip))

    def client_data(self, client, data):
        if not self.running: return
        self.entities.get_component(client.entity, CNetwork).data = data

    def send_tilemap(self, client, tid):
        tiles = self.tilemap.get_instance(tid)
        client.write('map {}'.format(tid))
        for pos, tile in tiles.items():
            client.write('til {}, {}: {}'.format(pos.x, pos.y, tile['tile']))

