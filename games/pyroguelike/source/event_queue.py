"""
"""

class EventQueue:
	def __init__(self):
		'''Initiates the queue.'''
		self.__queue = {}
	
	def addObj(self, obj):
		'''Adds an object to the queue, along with a new counter.'''
		self.__queue[obj] = 0
	
	def delObj(self, obj):
		'''Delete an object from the queue.'''
		if self.__queue.has_key(obj):
			del self.__queue[obj]
	
	def loopQueue(self):
		'''Loops through the queue, increasing each object's turn counter with 
			it's speed. Then checks if the object has reached the turn limit, 
			letting it do an action if it has.'''
		for obj in self.__queue:
			#increase the counter
			self.__queue[obj] -= obj.speed
			#and check if the counter has reached the limit
			if self.__queue[obj] <= 0:
				#ding! time to start over and letting the object do something
				self.__queue[obj] = obj.doAction()
