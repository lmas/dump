import creature_struct
import constants

class PlayerStruct(creature_struct.CreatureStruct):
	def setup(self, args):
		''' '''
		self.addDecor(constants.IS_PLAYER)
		self.symbol = '@'
		self.name = 'Lolipop'
		self.title = 'Gametester'
	
	def setMethods(self, io, time, map, state):
		self.io = io
		self.time = time
		self.map = map
		self.setState = state

################################################################################

	def doAction(self):
		self.cool_down_action_timer = 0
		while self.cool_down_action_timer == 0:
			x, y = self.xpos, self.ypos
			moremsgs = self.io.showMsg()
			self.updateStatus()
			self.map.getMap(x, y, self.radius, self.io.showCh)
			self.lastkey = key = self.io.getKey()
			self.io.clearScreen()
		
			if key == 27:
				self.setState(2)
				break
			elif key == ord('8'): #up
				self.movePlayer(x, y-1)
			elif key == ord('6'): #right
				self.movePlayer(x+1, y)
			elif key == ord('2'): #down
				self.movePlayer(x, y+1)
			elif key == ord('4'): #left
				self.movePlayer(x-1, y)
			else:
				#no (valid) action was taken by the player, so don't update 
				#the world, instead go back to the player and ask for something
				#else to do
				pass
		return self.cool_down_action_timer

################################################################################

	def movePlayer(self, x, y):
		if self.map.moveCreature(x, y, self):
			desc = self.map.getDescription(x, y)
			self.io.clearMsg()
			self.io.addMsg(desc)
			self.cool_down_action_timer = 1000
		else:
			self.io.clearMsg()
			#check what's blocking the player and take action acording to that
			self.io.addMsg("You can't walk that way, it's blocked!")

	def updateStatus(self):
		x, y = self.xpos, self.ypos
		self.io.showStr(1, 1, self.name + ' the ' + self.title)
		self.io.showStr(1, 3, 'Health: ' + str(self.health) + '/' + str(self.max_health), 5)
		self.io.showStr(1, 4, 'Mana:\t ' + str(self.mana) + '/' + str(self.max_mana), 2)
		
		self.io.showStr(1, 6, 'Muscles:   ' + str(self.muscles.getBase()) + ' (+' + str(self.muscles.getBonus()) + ')', 7)
		self.io.showStr(1, 7, 'Coordin.:  ' + str(self.coordination.getBase()) + ' (+' + str(self.coordination.getBonus()) + ')', 7)
		self.io.showStr(1, 8, 'Tolerance: ' + str(self.tolerance.getBase()) + ' (+' + str(self.tolerance.getBonus()) + ')', 7)
		self.io.showStr(1, 9, 'Knowledge: ' + str(self.knowledge.getBase()) + ' (+' + str(self.knowledge.getBonus()) + ')', 7)
		self.io.showStr(1, 10, 'Mind:\t    ' + str(self.mind.getBase()) + ' (+' + str(self.mind.getBonus()) + ')', 7)
		
		msg = 'Day %d, %02d:%02d:%02d' % self.time.currentTime()
		self.io.showStr(1, 17, msg)
		self.io.showStr(1, 18, 'Coordinates: ' + str(x) + ', ' + str(y))
		#self.io.showStr(1, 19, 'key pressed: ' + str(self.lastkey) + ' (' + self.io.keyName(self.lastkey) + ')')
