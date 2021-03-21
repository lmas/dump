"""
The most basic form of objects, it can have different decorations (but only 
one of each type)."""

class Decoration:
	def __init__(self, args = {}):
		#private vars, not to be used outside of this class
		self.__decorations = set() #the list over the object's current decorations
		
		self.initVars(args)
	
	def isInt(self, num):
		try:
			var = int(num)
			return True
		except ValueError:
			return False
	
	def initVars(self, args):
		"""In order to initiate stuff in a subclass during init, this funct
			should be overrided"""
		#raise NotImplementedError, "must be implemented in subclass"
		pass
	
	def addDecor(self, decor):
		"""Add a decoration to the object"""
		if self.isInt(decor):
			self.__decorations.add(decor)
			return True
		return False
	
	def delDecor(self, decor):
		"""Remove a decoration"""
		if self.isInt(decor):
			self.__decorations.discard(decor)
	
	def hasDecor(self, decor):
		"""Check if the object has a decoration"""
		if self.isInt(decor):
			for item in self.__decorations:
				if item == decor:
					return True
		return False
	
	def allDecors(self):
		"""Returns all aviable decorations for the object"""
		return self.__decorations
