"""
"""

#import sys

import constants
import io_wrapper
import map_struct
import player
import creature_struct
import event_queue
import time_struct

class GameLoop:
	def __init__(self):
		io_wrapper.IOWrapper(self.setup)
	
	def setup(self, wrapper):
		self.lastkey = 0
		self.game_state = 0	
		
		self.io = wrapper
		self.event_queue = event_queue.EventQueue()
		self.time = time_struct.TimeStruct()
		self.creatures = creature_struct.CreatureList()
		self.map = map_struct.MapStruct(self.creatures, 50, 20)
		self.map.makeMap()

		self.player = player.PlayerStruct()
		self.player.setMethods(self.io, self.time, self.map, self.setState)
		self.map.setCreature(25,10, self.player)
		self.creatures.addCreature(self.player)	

		self.startLoop()

################################################################################
			
	def startLoop(self):
		for creature in self.creatures.allCreatures():
			self.event_queue.addObj(creature)
		
		self.game_state = 1
		while self.game_state == 1:
			self.event_queue.loopQueue()
			self.time.increaseTime()

################################################################################

	def setState(self, state):
		self.game_state = state
