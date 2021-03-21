import random

TILE_UNUSED = 0
TILE_WALL = 1
TILE_FLOOR = 2

class CaveGenerator:
	def __init__(self, defx, defy, wall = 45):
		self.__map = []
		self.__maxx = defx
		self.__maxy = defy
		for y in range(defy):
			for x in range(defx):
				self.__map.append(TILE_WALL)
		
		for y in range(defy):
			for x in range(defx):
				if (self.__getRand(0, 100) < wall):
					self.__setCell(x, y, TILE_WALL)
				else:
					self.__setCell(x, y, TILE_FLOOR)
				
				self.__setCell(x, 0, TILE_UNUSED)
				self.__setCell(x, defy-1, TILE_UNUSED)
			self.__setCell(0, y, TILE_UNUSED)
			self.__setCell(defx-1, y, TILE_UNUSED)
		
	def __getRand(self, min, max):
		temp = random.randrange(min, max)
		#random.seed(temp)
		return temp
	
	def __setCell(self, x, y, tile):
		self.__map[x + self.__maxx * y] = tile
	
	def __getCell(self, x, y):
		return self.__map[x + self.__maxx * y]
	
	def showMap(self):
		for y in range(self.__maxy):
			line = ''
			for x in range(self.__maxx):
				if (self.__getCell(x, y) == TILE_UNUSED):
					line += '+'
				elif (self.__getCell(x, y) == TILE_WALL):
					line += '#'
				elif (self.__getCell(x, y) == TILE_FLOOR):
					line += '.'
			print line
	
#######################################################################

	def makeMap(self):
		for y in range(1, self.__maxy-1):
			for x in range(1, self.__maxx-1):
				count1 = count2 = 0
				
				for i in range(y-1, y+2):
					for j in range(x-1, x+2):
						if (self.__getCell(j, i) != TILE_FLOOR):
							count1 += 1
				
				for i in range(y-2, y+3):
					for j in range(x-2, x+3):
						if (abs(i-y) == 2 and abs(j-x) == 2):
							continue
						elif (i < 0 or j < 0 or i >= self.__maxy or j >= self.__maxx):
							continue
						elif (self.__getCell(j, i) != TILE_FLOOR):
							count2 += 1
						
				
				#if (self.__getCell(x, y) == TILE_FLOOR):
				#	if (count1 > 5): 
				#		self.__setCell(x, y, TILE_WALL)
				#elif (count1 < 4):
				#	self.__setCell(x, y, TILE_FLOOR)
				if debug:
					if (count1 >= 5 or count2 <= 2):
						self.__setCell(x, y, TILE_WALL)
					else:
						self.__setCell(x, y, TILE_FLOOR)
				else:
					if (count1 >= 5):
						self.__setCell(x, y, TILE_WALL)
					else:
						self.__setCell(x, y, TILE_FLOOR)
						
		return True
#######################################################################

if __name__ == '__main__':
	test = CaveGenerator(80, 20)
	generations = 4
	debug = True
	for i in range(generations):
		test.makeMap()
	
	generations = 3
	debug = False
	for i in range(generations):
		test.makeMap()

	test.showMap()

