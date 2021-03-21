"""
"""

import constants
import tile_struct
import fow_struct

class MapStruct:
	""" """
	def __init__(self, creaturelist, x=50, y=50):
		""" """
		self.__map = []
		self.__creaturelist = creaturelist
		self.xsize = x
		self.ysize = y

		for x in range(self.xsize):
			row = []
			for y in range(self.ysize):
				row.append(tile_struct.TileStruct())
			self.__map.append(row)

		self.__fow = fow_struct.FOWStruct(self.xsize, self.ysize, self.setVisited, self.lightBlocked)

	def setCell(self, x, y, tile = tile_struct.TileStruct()):
		""" """
		if (-1 < x < self.xsize) and (-1 < y < self.ysize):
			self.__map[x][y] = tile

	def getCell(self, x, y):
		""" """
		if (-1 < x < self.xsize) and (-1 < y < self.ysize):
			temp = self.__map[x][y]
			return temp
		return None

################################################################################

	def setVisited(self, x, y):
		if (-1 < x < self.xsize) and (-1 < y < self.ysize):
			self.getCell(x, y).addDecor(constants.PLAYER_VISITED)

	def isVisited(self, x, y):
		if (-1 < x < self.xsize) and (-1 < y < self.ysize):
			return self.getCell(x, y).hasDecor(constants.PLAYER_VISITED)
		return False

	def lightBlocked(self, x, y):
		if (-1 < x < self.xsize) and (-1 < y < self.ysize):
			return self.getCell(x, y).hasDecor(constants.BLOCK_LIGHT)
		return False

	def isLit(self, x, y):
		if (-1 < x < self.xsize) and (-1 < y < self.ysize):
			return self.__fow.isLit(x, y)
		return False

	def isBlocked(self, x, y):
		temp = self.getCell(x, y)
		if temp != None:
			if (temp.hasDecor(constants.IS_BLOCKED) or
					temp.creature != None):
				return True
		return False

	def moveCreature(self, x, y, creature):
		if not self.isBlocked(x, y):
			self.setCreature(x, y, creature)
			return True
		return False

	def setCreature(self, x, y, creature):
		#remove creature from old tile
		oldx, oldy = creature.xpos, creature.ypos
		self.getCell(oldx, oldy).creature = None

		#update the creature's position
		creature.xpos = x
		creature.ypos = y

		#add creature to new tile
		self.getCell(x, y).creature = creature

################################################################################

	def getDescription(self, x, y):
		tile = self.getCell(x, y)
		msg = ''
		if tile.feature > 0:
			msg += 'There is a ' + tile.getFeature()[2] + ' here. '

		if len(tile.items.allItems()) > 0:
			msg += 'Lying here is '
			for item in tile.items.allItems():
				msg += 'a ' + item.name + ', '

		return msg

	def getMap(self, px, py, pradius, printfunct):
		""" """

		self.__fow.calculateFow(px, py, pradius)
		dark = 8
		shadow = 5
		light = 8

		#symbol order:
		#tile>feature>item>creature>shadow

		j = 30
		for x in range(px-25, px+25):
			i = 0
			for y in range(py-10, py+10):
				symbol = ' '
				shadow_symbol = ' '
				color = 8
				modifier = 1
				tile = self.getCell(x, y)
				if tile != None and self.isVisited(x, y):
					temp = tile.getTerrain()
					if temp != None:
						shadow_symbol = symbol = temp[0]
						color = temp[1]
						modifier = temp[2]

					temp = tile.getFeature()
					if temp != None:
						shadow_symbol = symbol = temp[0]
						color = temp[1]

					temp = tile.items.allItems()
					if len(temp) > 0:
						symbol = temp[0].symbol
						color = temp[0].color

					temp = tile.creature
					if temp != None:
						symbol = temp.symbol
						color = temp.color
						if temp.hasDecor(constants.IS_PLAYER):
							modifier = 2

					temp = self.isLit(x, y)
					if temp == False:
						symbol = shadow_symbol
						color = shadow
						#modifier = dim
				else:
					color = dark


				printfunct(j, i, symbol, color, modifier)
				i += 1
			j += 1

################################################################################

	def newStoneWall(self, args = {}):
		args['terrain'] = constants.STONE_WALL
		return tile_struct.TileStruct(args)

	def newWall(self, args = {}):
		args['terrain'] = constants.WALL
		return tile_struct.TileStruct(args)

	def newFloor(self, args = {}):
		args['terrain'] = constants.FLOOR
		temp = tile_struct.TileStruct(args)
		temp.delDecor(constants.IS_BLOCKED)
		temp.delDecor(constants.BLOCK_LIGHT)
		return temp

	def makeMap(self):
		""" """
		xsize, ysize = self.xsize, self.ysize
		for x in range(xsize):
			for y in range(ysize):
				self.setCell(x, y, self.newFloor())
				self.setCell(0, y, self.newStoneWall())
				self.setCell(xsize-1, y, self.newStoneWall())
			self.setCell(x, 0, self.newStoneWall())
			self.setCell(x, ysize-1, self.newStoneWall())

		#import monster
		#temp = self.newFloor()
		#temp.creature = monster.MonsterStruct()
		#self.setCell(5, 5, temp)

		temp = self.newFloor({'feature':constants.FOUNTAIN})
		self.setCell(10, 3, temp)

		self.setCell(11, 5, self.newWall())

		import item_struct
		temp = self.newFloor({'feature':constants.STAIR_DOWN})
		for i in range(20):
			temp.items.addItem(item_struct.ItemStruct({'symbol':'|', 'color':7, 'name':'sword' + str(i)}))
		self.setCell(15, 3, temp)

