"""
"""

import decor_struct
import constants

class ItemStruct(decor_struct.Decoration):
	"""The most basic form of items, it can have different decorations (but only 
		one of each type) and has some basic properties like weight and name"""
	
	def initVars(self, args):
		#global vars, should be included in every instance (sub or not!)
		#think carefully about the vars, about what is needed as global...
		self.name = 'dummy_item' #should be more or less unique for each item
		self.symbol = '?'
		self.color = 0
		self.weight = 1 #pretty selfexplaining..
		self.stack_count = 1 #for stacking items, like 100 arrows in one item object
		
		if args.has_key('name'): self.name = args['name']
		if args.has_key('weight'): self.weight = args['weight']
		if args.has_key('stack_count'): self.stack_count = args['stack_count']
		if args.has_key('symbol'): self.symbol = args['symbol']
		if args.has_key('color'): self.color = args['color']
		
		self.setup(args)
		
	def setup(self, args):
		"""In order to initiate stuff in a subclass during init, this funct
			should be overrided"""
		#raise NotImplementedError, "must be implemented in subclass"
		pass
	

class ContainerStruct(ItemStruct):
	"""Reminds of the item_struct, but can also contain different items.
		This "containing" behavior should be limited however, ie. when trying 
		to place a container into another one."""
	
	def setup(self, args):
		"""Replacement for the baseclass's funct, so we can setup stuff during
			class initialization"""
		self.__items = []
		self.addDecor(constants.IS_CONTAINER)

	def addItem(self, obj):
		"""Add an item to the container, as long as it's not another container
			and it's actually based on item_struct.
			Returns true on success (item added), otherwise false."""
		#check if the item is a container or is based on item_struct
		if isinstance(obj, ItemStruct):
			if obj.hasDecor(constants.IS_CONTAINER): return False
		else: return False
		
		duplicate = self.hasItem(obj)
		if duplicate != None:
			#remove the duplicate from the container, increase it's stack count
			#with the obj's count and add the duplicate (with updated stack)
			#into the container again
			self.eraseItem(duplicate)
			duplicate.stack_count += obj.stack_count
			self.__items.append(duplicate)
			self.weight += duplicate.weight * duplicate.stack_count
		else:
			#simply add the non-duplicate item to the container
			self.__items.append(obj)
			self.weight += obj.weight
		return True

	def delItem(self, obj):
		"""Remove one item from the container or a stack"""
		item = self.hasItem(obj)
		if (item != None and item.stack_count > 1):
			#remove one item from it's stack
			self.eraseItem(item)
			item.stack_count -= 1
			self.__items.append(item)
			self.weight += item.weight * item.stack_count
		else:
			#not stacked, remove the item completely
			self.eraseItem(item)

	def eraseItem(self, obj):
		"""Completely removes an item, aswell as the whole stack(if stacked)"""
		item = self.hasItem(obj)
		while item != None:
			self.__items.remove(item)
			self.weight -= item.weight * item.stack_count
			if self.weight < 0: self.weight = 0
			item = self.hasItem(obj)
	
	def hasItem(self, obj):
		"""Check if there's an item in the container"""
		if isinstance(obj, ItemStruct):		
			for item in self.__items:
				if (item.allDecors() == obj.allDecors() and
					item.name == obj.name):
					#we have it, return it
					return item
		return None
	
	def allItems(self):
		"""Returns a list over the items being contained"""
		return self.__items
