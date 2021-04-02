////////////////////////////////////////////////////////////////////////////////
// speech_bubbles.js

const maxWidth = 400;
const r = 5;
const pi2 = Math.PI * 2;
const fontSize = 12;

function BubbleManager() {
	Manager.call(this);
	this.list = [];
}

BubbleManager.prototype = new Manager();

BubbleManager.prototype.add = function(text, x, y) {
	var bubble = {
		id: Date.now() + ":" + x + ":" + y,
		text: text.trim(),
		x: x,
		y: y,
	};
	this.list.push(bubble);
	this.dirty();

	// TODO: move up all other bubbles if it's cramped

	setTimeout(function() {
		var i = this.list.indexOf(bubble);
		this.list.splice(i, 1);
		this.dirty();
	}.bind(this), 3000 + (text.length / 2 * 100));
};

BubbleManager.prototype.render = function(ctx, camera) {
	ctx.textAlign = "center";
	ctx.font = fontSize + "px serif";
	for (var i = 0; i < this.list.length; i++) {
		var bubble = this.list[i];
		var pos = camera.translate(bubble.x, bubble.y);
		var x = pos[0] + tileSize;
		var y = pos[1];
		if (camera.visible(x, y) === false) {
			continue;
		}
		drawBubble(ctx, bubble.text, x, y);
	}
};

function drawBubble(ctx, text, x, y) {
	var lines = wrapText(ctx, text, maxWidth);
	var height = lines.length * fontSize;
	var width = 0;
	for (var l = 0; l < lines.length; l++) {
		var w = ctx.measureText(lines[l]).width;
		if (w > width) {
			width = w;
		}
	}

	x = x - (width / 2);
	y = y - height - r;
	ctx.beginPath();
	ctx.stroke();
	ctx.arc(x, y, r, pi2 * 0.5, pi2 * 0.75);
	ctx.arc(x + width, y, r, pi2 * 0.75, pi2);
	ctx.arc(x + width, y + height, r, 0, pi2 * 0.25);
	ctx.arc(x, y + height, r, pi2 * 0.25, pi2 * 0.5);
	ctx.fillStyle = "#444a";
	ctx.fill();

	ctx.fillStyle = "#fffc";
	for (var ll = 0; ll < lines.length; ll++) {
		ctx.fillText(lines[ll], x + (width / 2), y + (ll * fontSize) + (r * 2));
	}
}

function wrapText(ctx, text, maxWidth) {
	var words = text.trim().split(/\s+/);
	var line = "";
	var lines = [];
	for (var w = 0; w < words.length; w++) {
		var width = ctx.measureText(line + words[w] + " ").width;
		if (width > maxWidth) {
			lines.push(line.trim());
			line = "";
		}
		line = line + words[w] + " ";
	}
	lines.push(line.trim());
	return lines;
}
