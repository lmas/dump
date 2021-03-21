class event_queue:
	def __init__(self):
		'''Initiates the queue and defaults the turn limit to 100.'''
		self.__queue = {}
		self.__turn_limit = 100
	
	def addObj(self, obj, speed):
		'''Adds an object to the queue, along with it's current speed.'''
		self.__queue[obj] = [0, speed]
	
	def delObj(self, obj):
		'''Delete an object from the queue.'''
		if self.__queue.has_key(obj):
			del self.__queue[obj]
			
	def modObj(self, obj, speed):
		'''Modifies the speed of an object.'''
		if self.__queue.has_key(obj):
			self.__queue[obj][1] = speed
	
	def loopQueue(self):
		'''Loops through the queue, increasing each object's turn counter with 
			it's speed. Then checks if the object has reached the turn limit, 
			letting it do an action if it has.'''
		temp = []
		for obj in self.__queue:
			#increase the counter
			self.__queue[obj][0] += self.__queue[obj][1]
			#and check if the counter has reached the limit
			if self.__queue[obj][0] >= self.__turn_limit:
				#ding! time to start over!
				self.__queue[obj][0] -= self.__turn_limit
				#and letting the object do something
				obj.doAction()
				temp.append(obj)
		return temp

class mob:
	'''Our generic creature class for monsters and the player'''
	def __init__(self, name, speed = 10):
		self.name = name
		self.speed = speed

	def doAction(self):
		print '\t\t' + self.name + ': IS DOING AN ACTION!'

def show_time(time):
	temp = time % 1440
	hours = temp // 60
	mins = temp % 60
	days =  (time // 60) // 24

	#display modifiers, adds a leading zero when hours or mins < 10
	hour_mod = ''
	if hours < 10:
		hour_mod = '0'
	min_mod = ''
	if mins < 10:
		min_mod = '0'
	
	return 'day ' + str(days) + ', ' + hour_mod + str(hours) + ':' + min_mod + str(mins)

if __name__ == '__main__':
	max_turns = 300 #just for testing purposes, we need some kind of limited loop
	turns_taken = {} #keep count of the creatures turns taken
	eq = event_queue() #our turn queue
	monsters = [mob('aaaaaaaaaaa', 20), mob('bbb', 7), mob('ccc', 33)] #simple list of monsters
	time = 1337 #time testing, each turn roughly equals 1min
	
	for i in monsters:
		#setup the turn counter and add each monster to the turn queue
		turns_taken[i] = 0
		eq.addObj(i, i.speed)
	
	print 'test started at simulated game time: ' + show_time(time)
	##################################################################
	
	for turn in range(max_turns):
		#our "simulated gameplay", a for-loop limited to 100 turns
		time += 1
		print '#' + str(turn+1) + '\ttime: ' + show_time(time)
		for obj in eq.loopQueue(): #loop through a list over mobs who can do an action
			#let the mob do an action here!
			#print '\t' + event.name
			turns_taken[obj] += 1
			#TEST! modifying the current mob's speed on a specific turn,
			#uncomment the following 2lines to turn it on
			#if turn == 40:
			#	eq.modEvent(event, 10) #pulls down the creature's speed to 10
		
	##################################################################
	
	print '\ntest ended at simulated game time: '+ show_time(time)
	print 'turns of each mob:'
	for i in monsters:
		#simply shows us how many turns each mob took
		print i.name + ': ' + str(turns_taken[i]) + '/' + str(max_turns) + ' (' + str((turns_taken[i] / float(max_turns)) * 100) + '%)'
