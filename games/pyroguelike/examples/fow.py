import sys
import curses

class creature_struct:
	def __init__(self, x = 1, y = 1, char = 'a', radius = 10):
		self.x = x
		self.y = y
		self.char = char
		self.radius = radius
	
	def move(self, x, y):
		self.x = x
		self.y = y

###############################################################################

class tile_struct:
	'''A generic tile'''
	def __init__(self, type = 0, blocked = True, sight = False, visited = False):
		self.type = type
		self.blocked = blocked
		self.sight = sight
		self.visited = visited

###############################################################################

class map_struct:
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
				row.append(tile_struct(self.FLOOR[0], False, True))
			self.map.append(row)
		
		#then add a border around it
		for x in range(xsize):
			for y in range(ysize):
				self.setCell(0, y, tile_struct(self.WALL[0], True))
				self.setCell(self.xsize-1, y, tile_struct(self.WALL[0], True))
			self.setCell(x, 0, tile_struct(self.WALL[0], True))
			self.setCell(x, self.ysize-1, tile_struct(self.WALL[0], True))
	
	def setCell(self, x, y, tile = None):
		'''Set a specific tile at x,y in the map to a new tile'''
		if tile == None:
			temp = tile_struct(self.EMPTY[0])
		else:
			temp = tile
		if (-1 < x < self.xsize) and (-1 < y < self.ysize):
			self.map[x][y] = temp
	
	def getCell(self, x, y):
		'''Return the tile at x,y in the map'''
		if (-1 < x < self.xsize) and (-1 < y < self.ysize):
			return self.map[x][y]
		else:
			return tile_struct(self.EMPTY[0])
	
	def showMap(self, stdscr, px, py, radius, player, fow):
		'''Shows the map on the screen'''
		dark, shadow, light = curses.color_pair(8), curses.color_pair(7) | curses.A_NORMAL, curses.color_pair(7) | curses.A_STANDOUT
		j = 0
		for x in range(px-40, px+40):
			i = 0
			for y in range(py-10, py+10):
				if self.getCell(x,y).visited:
					if fow.lit(x, y):
						attr = light
					else:
						attr = shadow
				else:
					attr = dark
				
				if self.getCell(x,y).type == self.FLOOR[0]:
					stdscr.addch(i, j, self.FLOOR[1], attr)
				elif self.getCell(x,y).type == self.WALL[0]:
					stdscr.addch(i, j, self.WALL[1], attr)
				else:
					stdscr.addch(i, j, self.EMPTY[1], attr)
				i += 1
			j += 1
		stdscr.addch(10, 40, player)
		stdscr.refresh()

###############################################################################

class fow_struct:
	def __init__(self, map):
		# Multipliers for transforming coordinates to other octants:
		self.mult = [[1,  0,  0, -1, -1,  0,  0,  1],
				[0,  1, -1,  0,  0, -1,  1,  0],
				[0,  1,  1,  0,  0, -1, -1,  0],
				[1,  0,  0,  1, -1,  0,  0, -1]]
		
		self.map = map
		self.xsize = map.xsize*2
		self.ysize = map.ysize*2
		self.flag = 1
		self.unflag = 0
		
		self.light_map = []
		for x in range(self.xsize):
			row = []
			for y in range(self.ysize):
				row.append(self.unflag)
			self.light_map.append(row)
	
	def blocked(self, x, y):
		return not ((-1 < x < self.xsize) and 
					(-1 < y < self.ysize) and 
					(self.map.getCell(x,y).sight == True))
	
	def lit(self, x, y):
		return ((-1 < x < self.xsize) and (-1 < y < self.ysize) and (self.light_map[x][y] == self.flag))
	
	def set_lit(self, x, y):
		if (-1 < x < self.xsize) and (-1 < y < self.ysize):
			self.light_map[x][y] = self.flag
			old = self.map.getCell(x, y)
			temp = tile_struct(old.type, old.blocked, old.sight, True)
			self.map.setCell(x, y, temp)
	
	def calculate_fow(self, x, y, radius):
		for tx in range(self.xsize):
			for ty in range(self.ysize):
				self.light_map[tx][ty] = self.unflag

		for oct in range(8):
			self.cast_light(x, y, radius, 1, 1.0, 0.0, 
							self.mult[0][oct], self.mult[1][oct],
							self.mult[2][oct], self.mult[3][oct])

	def cast_light(self, cx, cy, radius, row, start, end, xx, xy, yx, yy):
		if start < end:
			return
		
		radius_squared = (radius+0.5) * (radius+0.5)
		
		for j in xrange(row, radius+1):
			dx, dy = -j-1, -j
			blocked = False
			
			while dx <= 0:
				dx += 1
				# Translate the dx, dy coordinates into map coordinates:
				X, Y = cx + dx * xx + dy * xy, cy + dx * yx + dy * yy
				# l_slope and r_slope store the slopes of the left and right
				# extremities of the square we're considering:
				l_slope, r_slope = (dx-0.5)/(dy+0.5), (dx+0.5)/(dy-0.5)
                
				if start < r_slope:
					continue
				elif end > l_slope:
					break
				else:
					# Our light beam is touching this square; light it:
					if dx*dx + dy*dy < radius_squared:
						self.set_lit(X, Y)
		
					if blocked:
						# we're scanning a row of blocked squares:
						if self.blocked(X, Y):
							new_start = r_slope
							continue
				
						else:
							blocked = False
							start = new_start
		
					else:
						if self.blocked(X, Y) and j < radius:
							# This is a blocking square, start a child scan:
							blocked = True
							self.cast_light(cx, cy, radius, j+1, start, l_slope, xx, xy, yx, yy)
							new_start = r_slope
			# Row is scanned; do next row unless last square was blocked:
			if blocked:
				break
 
###############################################################################

def main(stdscr):
	MAP = map_struct(80, 20)
	MAP.setCell(10, 10, tile_struct(MAP.WALL[0], True, True))
	MAP.setCell(15, 10, tile_struct(MAP.WALL[0], False))
	MAP.setCell(12, 13, tile_struct(MAP.WALL[0], True))
	fow = fow_struct(MAP)
	PLAYER = creature_struct(78, 19, '@', 6)

	curses.curs_set(0) #turn off the blinking cursor
	#sets the color pairs for curses, foreground colors:
	#1=red, 2=green, 3=yellow, 4=blue, 5=magenta, 6=cyan, 7=white and 8=black
	for i in range(1, 8):
		curses.init_pair(i, i, 0) # 0 = black background color
	game_state = 1
	input = 0

	while game_state == 1:
		x, y = PLAYER.x, PLAYER.y
		stdscr.erase()
		stdscr.addstr(20,0, 'key pressed:' + curses.keyname(input) + ' (' + str(input) + ')')
		stdscr.addstr(20, 30, 'x:' + str(x) + ' y:' + str(y))
		stdscr.addstr(22,0, 'ESC = quit, numbers OR numpad = move')
		fow.calculate_fow(x, y, PLAYER.radius)
		MAP.showMap(stdscr, x, y, PLAYER.radius, PLAYER.char, fow)
		
		input = stdscr.getch()
		if input == 27:
			game_state = 2
			break
		elif input == ord('8'): #up
			if MAP.getCell(x, y-1).blocked == False:
				PLAYER.move(x, y-1)
		elif input == ord('6'): #right
			if MAP.getCell(x+1, y).blocked == False:
				PLAYER.move(x+1, y)
		elif input == ord('2'): #down
			if MAP.getCell(x, y+1).blocked == False:
				PLAYER.move(x, y+1)
		elif input == ord('4'): #left
			if MAP.getCell(x-1, y).blocked == False:
				PLAYER.move(x-1, y)
	sys.exit(0)

###############################################################################

if __name__ == '__main__':
	curses.wrapper(main)
