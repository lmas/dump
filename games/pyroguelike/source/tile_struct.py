import decor_struct
import constants
import item_struct

#A list over terrain types, each type contains: id(unique!), symbol and color
TerrainData = {
	constants.STONE_WALL	:['#', 8, 1],
	constants.WALL		:['#', 0, 1],
	constants.PASSAGE		:['#', 0, 1],
	constants.FLOOR		:['.', 8, 1]
}
################################################################################
#List over misc. tile features, same format as the terrain data
FeatureData = {
	constants.DOOR_CLOSED	:['-', 8, 'Closed door'],
	constants.DOOR_OPEN	:['+', 8, 'Open door'],
	constants.STAIR_UP	:['<', 8, 'Staircase'],
	constants.STAIR_DOWN	:['>', 8, 'Staircase'],
	constants.FOUNTAIN	:['{', 2, 'Fountain']
}
################################################################################

class TileStruct(decor_struct.Decoration):
	def initVars(self, args):
		''' '''
		self.terrain = 0 #The terrain type of the tile
		self.feature = 0 #a special feature, like a door
		self.creature = None #holds one creature
		self.items = item_struct.ContainerStruct() #has a container for items, "drops on the floor"
		self.addDecor(constants.IS_BLOCKED)
		self.addDecor(constants.BLOCK_LIGHT)
		
		if args.has_key('terrain'): self.terrain = args['terrain']
		if args.has_key('feature'): self.feature = args['feature']

	def getTerrain(self):
		''' '''
		if self.isInt(self.terrain):	
			return TerrainData[self.terrain]
		return None
	
	def getFeature(self):
		''' '''
		if self.isInt(self.feature):
			if self.feature > 0:
				return FeatureData[self.feature]
		return None
