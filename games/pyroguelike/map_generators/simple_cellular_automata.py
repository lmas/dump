#http://pixelenvy.ca/wa/ca_cave.html
import random

TILE_UNUSED = 0
TILE_WALL = 1
TILE_FLOOR = 2

class CaveGenerator:
	def __init__(self, defx, defy, floor_count = 0.40):
		self.__map = []
		self.__maxx = defx
		self.__maxy = defy
		for y in range(defy):
			for x in range(defx):
				self.__map.append(TILE_WALL)
		
		for y in range(defy):
			for x in range(defx):
				self.__setCell(x, 0, TILE_UNUSED)
				self.__setCell(x, defy-1, TILE_UNUSED)
			self.__setCell(0, y, TILE_UNUSED)
			self.__setCell(defx-1, y, TILE_UNUSED)
		
		amount = int((defx * defy) * floor_count)
		while amount > 0:
			x = self.__getRand(1, (defx-1))
			y = self.__getRand(1, (defy-1))
						
			if self.__getCell(x, y) == TILE_WALL:
				self.__setCell(x, y, TILE_FLOOR)
				amount -= 1
				#print 'amount:', amount, '\tx:', x, '\ty:', y

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
				count = 0
				for i in (-1, 0, 1):
					for j in (-1, 0, 1):
						if (self.__getCell(x+j, y+i) != TILE_FLOOR and not(i == 0 and j == 0)):
							count += 1
				
				if (self.__getCell(x, y) == TILE_FLOOR):
					if (count > 5): self.__setCell(x, y, TILE_WALL)
				elif (count < 4):
					self.__setCell(x, y, TILE_FLOOR)
				
		return True

#######################################################################

if __name__ == '__main__':
	test = CaveGenerator(80, 20)
	for i in range(5):
		test.makeMap()

	test.showMap()
