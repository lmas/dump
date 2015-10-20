
import logging
from uuid import uuid4 as uuid

from engine.utils import Vector

TILESET = {
    -1: dict(
        name = 'ERROR',
    ),

    0: dict(
        name = 'Empty',
    ),

    1: dict(
        name = 'Floor',
        walkable = True,
    ),

    2: dict(
        name = 'Wall',
    ),
}


class Tilemap(object):
    def __init__(self):
        self.tiles = {}

    def new_instance(self, tiles={}):
        if not isinstance(tiles, dict):
            raise TypeError('Tileset must be a dict with position keys and tile type values!')
        tid = uuid().hex
        self.tiles[tid] = tiles
        return tid

    def load_instance(self, tid, mapfile):
        logging.info('Loading tilemap instance {}: {}'.format(tid, mapfile))
        tiles = {}
        with open(mapfile, 'r') as f:
            for line in f.readlines():
                x, y, tile = line.split(',')
                pos = Vector(int(x.strip()), int(y.strip()))
                tile = int(tile.strip())
                tiles[pos] = dict(
                    tile = tile,
                )
        self.set_instance(tid, tiles)

    def get_instance(self, tid):
        return self.tiles[tid]

    def set_instance(self, tid, tiles):
        self.tiles[tid] = tiles

    def get_tile(self, tid, pos):
        try:
            tile = self.tiles[tid][pos]
            tile_type = tile['tile']
            return TILESET[tile_type]
        except KeyError:
            return None

    def set_tile(self, tid, pos, tile):
        self.tiles[tid][pos] = tile

    def check_collision(self, tid, pos):
        # HACK: This might be a bad idea, might create unnecessary load...
        pos = Vector(int(pos.x), int(pos.y))
        tile = self.get_tile(tid, pos)
        if tile and tile.get('walkable', False):
            return False
        return True

