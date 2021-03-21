import decor_struct
import constants
import item_struct

class CreatureList:
	def __init__(self):
		self.__creatures = []
	
	def addCreature(self, creature):
		item = self.hasCreature(creature)
		if item == None:
			self.__creatures.append(creature)
	
	def delCreature(self, creature):
		item = self.hasCreature(creature)
		while item != None:
			self.__creatures.remove(item)
			item = self.hasCreature(creature)
	
	def hasCreature(self, creature):
		#if isinstance(creature, CreatureStruct):
		for item in self.__creatures:
			if item == creature:
				return item
		return None
	
	def allCreatures(self):
		return self.__creatures

class AttributeStruct:
	''' '''
	def __init__(self, base=0, bonus=0):
		self.__base = 0 #base points (the "lvl")
		self.__bonus = 0 #bonus points added
		self.__exp = 0 #current advancement in the training
		self.__limit = self.__base * 100 #exp needed to advance one "lvl"
	
	def getPoints(self):
		''' '''
		return self.__base + self.__bonus
	
	def getBase(self):
		''' '''
		return self.__base
	
	def getBonus(self):
		''' '''
		return self.__bonus

	def addBonus(self, points):
		''' '''
		self.__bonus += points
	
	def delBonus(self, points):
		''' '''
		self.__bonus -= points
		if self.__bonus < 0: self.__bonus = 0
	
	def addExp(self, points):
		''' '''
		self.__exp += points
		lvlup = False
		while self.__exp >= self.__limit:
			self.__exp -= self.__limit
			self.setBase(self.__base+1)
			lvlup = True
		return lvlup
	
	def setBase(self, points):
		''' '''
		self.__base = points
		self.__limit = self.__base * 100
	

class CreatureStruct(decor_struct.Decoration):
	'''The most basic form of creatures, it can have different decorations 
		(but only one of each type) and has some basic properties like name, 
		equipment and health.'''
	
	def initVars(self, args):
		#global vars, should be included in every instance (sub or not!)
		#think carefully about the vars, about what is needed as global...
		self.name = 'dummy'
		self.symbol = 'a'
		self.color = 8
		self.radius = 5
		self.xpos = 0
		self.ypos = 0
		self.speed = 1000
		self.cool_down_action_timer = 0
		
		#stuff related to misc. stats
		self.weight = 1
		self.max_health = 1
		self.health = self.max_health
		self.max_mana = 1
		self.mana = self.max_mana
		
		#stuff related to attributes
		self.muscles = AttributeStruct()
		self.coordination = AttributeStruct()
		self.tolerance = AttributeStruct()
		self.knowledge = AttributeStruct()
		self.mind = AttributeStruct()
				
		#some item related stuff
		self.money = 0
		self.items = item_struct.ContainerStruct() #inventory
		
		#and here's the equipment slots, starting with the gear
		self.head = item_struct.ItemStruct()
		self.arms = item_struct.ItemStruct()
		self.torso = item_struct.ItemStruct()
		self.legs = item_struct.ItemStruct()
		self.feet = item_struct.ItemStruct()
		#and ending with the combat slots
		self.weapon = item_struct.ItemStruct()
		self.shield = item_struct.ItemStruct()
		
		self.setup(args)
	
	def setup(self, args):
		'''In order to initiate stuff in a subclass during init, this funct
			should be overrided'''
		#raise NotImplementedError, "must be implemented in subclass"
		pass
	
	def doAction(self):
		raise NotImplementedError, "must be implemented in subclass"
	
	def doHit(self, fromcreature):
		pass
	
	def doExercise(self, attr, points):
		pass
