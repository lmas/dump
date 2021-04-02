////////////////////////////////////////////////////////////////////////////////
// camera.js

function Camera() {
	this.id = 0;
	this.x = 0;
	this.y = 0;
}

Camera.prototype.update = function(width, height) {
	this.w = width;
	this.h = height;
	this.cw = Math.floor(width / 2 / screenSize);
	this.ch = Math.floor(height / 2 / screenSize);
};

Camera.prototype.width = function() {
	return this.w;
};

Camera.prototype.height = function() {
	return this.h;
};

Camera.prototype.halfWidth = function() {
	return this.cw;
};

Camera.prototype.halfHeight = function() {
	return this.ch;
};

////////////////////////////////////////////////////////////////////////////////

Camera.prototype.translate = function(x, y) {
	var cx = Math.floor(x) - this.x + this.halfWidth();
	var cy = Math.floor(y) - this.y + this.halfHeight();
	var offsetX = Math.floor((x - Math.floor(x)) * screenSize);
	var offsetY = Math.floor((y - Math.floor(y)) * screenSize);
	return [cx * screenSize + offsetX, cy * screenSize + offsetY];
};

Camera.prototype.visible = function(x, y) {
	if (x < -screenSize || x > this.width() || y < -screenSize || y > this.height()) {
		return false;
	}
	return true;
};

////////////////////////////////////////////////////////////////////////////////

Camera.prototype.follow = function(id) {
	this.id = id;
};

Camera.prototype.following = function(id) {
	return this.id === id;
};

Camera.prototype.move = function(x, y) {
	this.x = x;
	this.y = y;
};
