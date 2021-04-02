////////////////////////////////////////////////////////////////////////////////
// map.js

// Might want to split up the map into prerendered canvas tiles and use
// AABB to find out which tiles to draw for the player, see:
// https://stackoverflow.com/questions/25342237/fast-min-max-aabb-collision-detection

function MapManager() {
	Manager.call(this);
	this.sheet = Game.getImage("tiles");
	this.width = 0;
	this.height = 0;
	this.tileset = {};
	this.raw = [];
}

MapManager.prototype = new Manager();

MapManager.prototype.set = function(data) {
	this.width = data.width;
	this.height = data.height;
	this.tileset = data.tileset;
	this.raw = data.raw;
	this.dirty();
};

MapManager.prototype.getTile = function(x, y) {
	if (x < 0 || x >= this.width || y < 0 || y >= this.height) {
		return undefined;
	}

	var t = this.raw[Math.floor(y) * this.width + Math.floor(x)];
	return t;
};

MapManager.prototype.setTile = function(x, y, tile) {
	if (x < 0 || x >= this.width || y < 0 || y >= this.height) {
		return undefined;
	}

	this.raw[Math.floor(y) * this.width + Math.floor(x)] = tile;
};

MapManager.prototype.bitmaskTile = function(x, y, tile) {
	var bitmask = 0;
	if (this.getTile(x, y - 1) === tile) bitmask += 1;
	if (this.getTile(x, y + 1) === tile) bitmask += 4;
	if (this.getTile(x - 1, y) === tile) bitmask += 8;
	if (this.getTile(x + 1, y) === tile) bitmask += 2;
	return bitmask;
};

MapManager.prototype.render = function(ctx, camera) {
	var camWidth = camera.halfWidth();
	var camHeight = camera.halfHeight();
	var startX = 0;
	var startY = -1; // make this -1 to showing missing tile on top of screen
	var stopX = (camWidth * 2) + 1;
	var stopY = (camHeight * 2) + 1;
	var offsetX = Math.floor((camera.x - Math.floor(camera.x)) * screenSize);
	var offsetY = Math.floor((camera.y - Math.floor(camera.y)) * screenSize) - screenSize / 2;
	// Added small pixel adjustment upwards to y, so we don't clip tiles
	// underneath us, from above

	for (var x = startX; x <= stopX; x++) {
		for (var y = startY; y <= stopY; y++) {
			var tx = x + camera.x - camWidth;
			var ty = y + camera.y - camHeight;
			var tile = this.getTile(tx, ty);
			if (tile === undefined) {
				continue;
			}

			// if it's a wall, let's bitmask it according to it's neighbors
			if (tile === 8) {
				tile = this.bitmaskTile(tx, ty, tile) + 100;
			}
			if (tile === 9) {
				tile = this.bitmaskTile(tx, ty, tile) + 120;
			}

			var tileset = this.tileset[tile];
			var sx = tileset.sheetx * tileSize;
			var sy = tileset.sheety * tileSize;
			ctx.drawImage(
				this.sheet,
				sx, sy, tileSize, tileSize,
				x * screenSize - offsetX, y * screenSize - offsetY,
				screenSize, screenSize
			);
		}
	}
};
