////////////////////////////////////////////////////////////////////////////////
// mobs.js

function MobManager() {
	Manager.call(this);
	this.sheet = Game.getImage("human");
	this.list = new Map();
}

MobManager.prototype = new Manager();

MobManager.prototype.set = function(mob) {
	this.list.set(mob.id, mob);
	this.dirty();
};

MobManager.prototype.get = function(id) {
	return this.list.get(id);
};

MobManager.prototype.del = function(id) {
	this.list.delete(id);
	this.dirty();
};

MobManager.prototype.render = function(ctx, camera) {
	ctx.textAlign = "center";
	ctx.font = "10px serif";
	ctx.fillStyle = "#fffa";
	for (var mob of this.list.values()) {
		var gender = 0;
		if (mob.m !== true) {
			gender = 1;
		}
		var pos = camera.translate(mob.x, mob.y);
		var x = pos[0];
		var y = pos[1];
		if (camera.visible(x, y) === false) {
			continue;
		}
		ctx.drawImage(
			this.sheet,
			// source x = direction, source y = gender
			mob.d * 2 * tileSize, gender * tileSize, tileSize, tileSize,
			x, y, screenSize, screenSize
		);
		nameTag(ctx, camera, mob.name, x + tileSize, y + (tileSize * 2.5));
	}
};

function nameTag(ctx, camera, name, x, y) {
	if (camera.visible(x, y) === false) {
		return;
	}
	ctx.fillText(name, x, y);
}
