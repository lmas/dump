
from engine.utils import Vector

class CNetwork(object):
    def __init__(self, client=None):
        self.client = client
        self.data = ''

class CPosition(object):
    def __init__(self, x=0, y=0):
        self.vector = Vector(x,y)

class CVelocity(object):
    def __init__(self):
        self.current = Vector(0,0)
        self.target = Vector(0,0)

