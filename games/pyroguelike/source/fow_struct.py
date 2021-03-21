class FOWStruct:
	def __init__(self, x, y, visitfunc, lightfunc):
		# Multipliers for transforming coordinates to other octants:
		self.__mult = [[1,  0,  0, -1, -1,  0,  0,  1],
						[0,  1, -1,  0,  0, -1,  1,  0],
						[0,  1,  1,  0,  0, -1, -1,  0],
						[1,  0,  0,  1, -1,  0,  0, -1]]
		
		#self.map = map
		self.__visitfunc = visitfunc
		self.__lightfunc = lightfunc
		self.__xsize = x
		self.__ysize = y
		self.__flag = 1
		self.__unflag = 0
		
		self.__light_map = []
		for x in range(self.__xsize):
			row = []
			for y in range(self.__ysize):
				row.append(self.__unflag)
			self.__light_map.append(row)
	
	def isBlocked(self, x, y):
		return ((-1 < x < self.__xsize) and 
					(-1 < y < self.__ysize) and 
					(self.__lightfunc(x, y)))
	
	def isLit(self, x, y):
		return ((-1 < x < self.__xsize) and (-1 < y < self.__ysize) and (self.__light_map[x][y] == self.__flag))
	
	def setLit(self, x, y):
		if (-1 < x < self.__xsize) and (-1 < y < self.__ysize):
			self.__light_map[x][y] = self.__flag
			self.__visitfunc(x, y)
	
	def calculateFow(self, x, y, radius):
		for tx in range(self.__xsize):
			for ty in range(self.__ysize):
				self.__light_map[tx][ty] = self.__unflag

		for oct in range(8):
			self.__castLight(x, y, radius, 1, 1.0, 0.0, 
							self.__mult[0][oct], self.__mult[1][oct],
							self.__mult[2][oct], self.__mult[3][oct])

################################################################################

	def __castLight(self, cx, cy, radius, row, start, end, xx, xy, yx, yy):
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
						self.setLit(X, Y)
		
					if blocked:
						# we're scanning a row of blocked squares:
						if self.isBlocked(X, Y):
							new_start = r_slope
							continue
				
						else:
							blocked = False
							start = new_start
		
					else:
						if self.isBlocked(X, Y) and j < radius:
							# This is a blocking square, start a child scan:
							blocked = True
							self.__castLight(cx, cy, radius, j+1, start, l_slope, xx, xy, yx, yy)
							new_start = r_slope
			# Row is scanned; do next row unless last square was blocked:
			if blocked:
				break
		self.setLit(cx, cy) #cheap bugfix! for some reason the tile the player
							 #is standing on is in shadow
