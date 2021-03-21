class mob:
	'''Our generic creature class for monsters and the player'''
	def __init__(self, name, speed = 10):
		self.name = name
		self.speed = speed

class turn_queue:
	def __init__(self):
		'''Creates an empty creature dict and set freeaction limit to default(100)'''
		self.__creatures = {}
		self.__freeaction = 100
		
	def add(self, obj, speed):
		'''Add a new creature to the dict, along with it's current speed 
		and empty counter'''
		self.__creatures[obj] = [0, speed]
	
	def mod(self, obj, new_speed):
		'''Modifies a creature's speed'''
		self.__creatures[obj][1] = new_speed
	
	def rem(self, obj):
		'''Removes a specific creature from the dict'''
		del self.__creatures[obj]
	
	def get(self):
		'''Check each creature in the list if it's counter is >= freeaction,
		returns all whose counters is and then increases each creature's
		counter by it's speed'''
		temp = []
		for creature in self.__creatures:
			if self.__creatures[creature][0] >= self.__freeaction:
				self.__creatures[creature][0] -= self.__freeaction
				temp.append(creature)
			self.__creatures[creature][0] += self.__creatures[creature][1]
		return temp
				

if __name__ == '__main__':
	max_turns = 100 #just for testing purposes, we need some kind of limited loop
	turns_taken = {} #keep count of the creatures turns taken
	eq = turn_queue() #our turn queue
	monsters = [mob('aaaaaaaaaaa', 25), mob('bbb', 13), mob('ccc', 32)] #simple list of monsters
	
	
	for i in monsters:
		#setup the turn counter and add each monster to the turn queue
		turns_taken[i] = 0
		eq.add(i, i.speed)
		
	##################################################################
	
	for turn in range(max_turns):
		#our "simulated gameplay", a for-loop limited to 100 turns
		print '#' + str(turn+1) + ':\t' #print each turn
		for creature in eq.get(): #loop through a list over mobs who can do an action
			#let the mob do an action here!
			print '\t' + creature.name
			turns_taken[creature] += 1
			#TEST! modifying the current mob's speed on a specific turn,
			#uncomment the following 2lines to turn it on
			#if turn == 41:
			#	eq.mod(creature, 10) #pulls down the creature's speed to 10
		
	##################################################################
	
	print 'turns of each mob:'
	for i in monsters:
		#simply shows us how many turns each mob took
		print i.name, turns_taken[i]
