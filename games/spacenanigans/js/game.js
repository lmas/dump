////////////////////////////////////////////////////////////////////////////////
// game.js

const maxHistory = 500;
const ratio = window.devicePixelRatio;
const tileSize = 32;
const screenSize = 64;

////////////////////////////////////////////////////////////////////////////////

var screenWidth = 0;
var screenHeight = 0;
var Game = {};

window.onload = function() {
	var canvas = document.getElementById("game");
	Game.run(canvas);
};

Game.run = function(canvas) {
	this.canvas = canvas;
	this.ctx = this.canvas.getContext("2d");
	this.images = {};
	var load = [
		this.loadImage("tiles", "/static/assets/tiles.png"),
		this.loadImage("human", "/static/assets/human.png")
	];
	Promise.all(load).then(function(loaded) {
		this.init();
		window.requestAnimationFrame(this.tick.bind(this));
	}.bind(this));
};

Game.init = function() {
	this.camera = new Camera();
	this.mapManager = new MapManager();
	this.mobManager = new MobManager();
	this.bubbleManager = new BubbleManager();
	this.moving = false;
	this.connect(getWebsocketAddr());
	window.addEventListener("keydown", this.onKeyDown.bind(this));
	window.addEventListener("keyup", this.onKeyUp.bind(this));
};

Game.tick = function(elapsed) {
	window.requestAnimationFrame(this.tick.bind(this)); // keep us runnign
	var dirty = resizeCanvas(this.canvas);
	if (dirty === true) {
		this.camera.update(screenWidth, screenHeight);
	}
	dirty = this.mapManager.update(this.camera, dirty);
	dirty = this.mobManager.update(this.camera, dirty);
	dirty = this.bubbleManager.update(this.camera, dirty);
	if (dirty === false) {
		return;
	}
	this.ctx.clearRect(0, 0, screenWidth, screenHeight); // clear previous frame
	this.mapManager.draw(this.ctx);
	this.mobManager.draw(this.ctx);
	this.bubbleManager.draw(this.ctx);
};

function resizeCanvas(canvas) {
	screenWidth = canvas.clientWidth * ratio;
	screenHeight = canvas.clientHeight * ratio;
	if (canvas.width !== screenWidth || canvas.height !== screenHeight) {
		canvas.width = screenWidth;
		canvas.height = screenHeight;
		return true;
	}
	return false;
}
////////////////////////////////////////////////////////////////////////////////

Game.loadImage = function(title, src) {
	var img = new Image();
	var d = new Promise(function(resolve, reject) {
		img.onload = function() {
			this.images[title] = img;
			resolve(img);
		}.bind(this);

		img.onerror = function() {
			reject("Could not load image: " + src);
		};
	}.bind(this));
	img.src = src;
	return d;
};

Game.getImage = function(title) {
	if (title in this.images) {
		return this.images[title];
	}
	return undefined;
};
