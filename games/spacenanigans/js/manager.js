////////////////////////////////////////////////////////////////////////////////
// manager.js

function Manager() {
	this._dirty = true;
	this._canvas = document.createElement("canvas");
	this._ctx = this._canvas.getContext("2d");
}

Manager.prototype.dirty = function() {
	this._dirty = true;
};

Manager.prototype.draw = function(ctx) {
	ctx.drawImage(this._canvas, 0, 0);
};

Manager.prototype._render = function(camera) {
	// MUST BE IMPLEMENTED
};

Manager.prototype.update = function(camera, dirty) {
	if (dirty === false && this._dirty === false) {
		return false;
	}
	// Resizing causes the canvas to be cleared I think
	this._canvas.width = camera.width();
	this._canvas.height = camera.height();
	//this._ctx.clearRect(0, 0, screenWidth, screenHeight);
	this.render(this._ctx, camera);
	this._dirty = false;
	return true;
};
