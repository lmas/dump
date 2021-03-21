"""
"""

class TimeStruct:
	def __init__(self, starttime = 0):
		''' '''
		self.__current_time = starttime
	
	def increaseTime(self, tick = 1):
		''' '''
		self.__current_time += tick
	
	def currentTime(self):
		''' '''
		days = (self.__current_time // 3600) // 24
		hours = (self.__current_time // 3600) % 24
		minutes = (self.__current_time % 3600) // 60
		seconds = (self.__current_time % 3600) % 60
		return (days, hours, minutes, seconds)
