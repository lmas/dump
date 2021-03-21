import sys
import curses
import random

class creature_type:
	def __init__(self, x, y, char = 'a'):
		self.x = x
		self.y = y
		self.char = char
	
	def move(self, x, y):
		self.x = x
		self.y = y

class tile_type:
	'''A generic tile'''
	def __init__(self, type = 0, blocked = True):
		self.type = type
		self.blocked = blocked

class map_type:
	'''A generic map'''
	def __init__(self, xsize = 10, ysize = 10):
		'''Initialize and setup the properties of the map'''
		self.map = [] #the map array
		self.xsize = xsize #max sizes of the map
		self.ysize = ysize
		#default tile types and their characters to show
		self.EMPTY = [0, ' ']
		self.FLOOR = [1, '.']
		self.WALL = [2, '#']
		
		#fill the map array with generic, empty tiles
		for x in range(xsize):
			row = []
			for y in range(ysize):
				row.append(tile_type(self.FLOOR[0]))
			self.map.append(row)
		
		#then add a border around it
		for x in range(xsize):
			for y in range(ysize):
				self.setCell(0, y, tile_type(self.WALL[0]))
				self.setCell(self.xsize-1, y, tile_type(self.WALL[0]))
			self.setCell(x, 0, tile_type(self.WALL[0]))
			self.setCell(x, self.ysize-1, tile_type(self.WALL[0]))
	
	def setCell(self, x, y, tile = None):
		'''Set a specific tile at x,y in the map to a new tile'''
		if tile == None:
			temp = tile_type(self.EMPTY[0])
		else:
			temp = tile
		if (-1 < x < self.xsize) and (-1 < y < self.ysize):
			self.map[x][y] = temp
	
	def getCell(self, x, y):
		'''Return the tile at x,y in the map'''
		if (-1 < x < self.xsize) and (-1 < y < self.ysize):
			return self.map[x][y]
		else:
			return tile_type(self.EMPTY[0])
	
	def showMap(self, stdscr, px, py, player):
		'''Shows the map on the screen'''
		j = 0
		for x in range(px-40, px+40):
			i = 0
			for y in range(py-10, py+10):
				if self.getCell(x,y).type == self.FLOOR[0]:
					stdscr.addch(i, j, self.FLOOR[1])
				elif self.getCell(x,y).type == self.WALL[0]:
					stdscr.addch(i, j, self.WALL[1])
				else:
					stdscr.addch(i, j, self.EMPTY[1])
				i += 1
			j += 1
		stdscr.addch(10, 40, player)
		stdscr.refresh()


def main(stdscr):
	PLAYER = creature_type(1, 1, '@')
	MAP = map_type(100, 100)
	MAP.setCell(10, 10, tile_type(MAP.WALL[0]))
	curses.curs_set(0) #turn of the blinking cursor
	
	game_state = 1
	input = 0
	while game_state == 1:
		x, y = PLAYER.x, PLAYER.y
		stdscr.erase()
		stdscr.addstr(20,0, 'key pressed:' + curses.keyname(input) + ' (' + str(input) + ')')
		stdscr.addstr(20, 30, 'x:' + str(x) + ' y:' + str(y))
		stdscr.addstr(22,0, 'q = quit, w = new map, numbers OR numpad = move')
		MAP.showMap(stdscr, x, y, PLAYER.char)
		input = stdscr.getch()
		if input == ord('q') or input == 27:
			game_state = 2
			break
		elif input == ord('8'): #up
			if MAP.getCell(x, y-1).type == MAP.FLOOR[0]:
				PLAYER.move(x, y-1)
		elif input == ord('6'): #right
			if MAP.getCell(x+1, y).type == MAP.FLOOR[0]:
				PLAYER.move(x+1, y)
		elif input == ord('2'): #down
			if MAP.getCell(x, y+1).type == MAP.FLOOR[0]:
				PLAYER.move(x, y+1)
		elif input == ord('4'): #left
			if MAP.getCell(x-1, y).type == MAP.FLOOR[0]:
				PLAYER.move(x-1, y)
		
		#"diagonal" movement
		#elif input == ord('9'):
		#	if getc(posx+1, posy-1) == FLOOR:
		#		posx += 1
		#		posy -= 1
		#elif input == ord('3'):
		#	if getc(posx+1, posy+1) == FLOOR:
		#		posx += 1
		#		posy += 1
		#elif input == ord('1'):
		#	if getc(posx-1, posy+1) == FLOOR:
		#		posx -= 1
		#		posy += 1
		#elif input == ord('7'):
		#	if getc(posx-1, posy-1) == FLOOR:
		#		posx -= 1
		#		posy -= 1
	sys.exit(0)

if __name__ == '__main__':
	curses.wrapper(main)
