////////////////////////////////////////////////////////////////////////////////
// input.js

Game.onKeyDown = function(event) {
	var key = event.keyCode;
	var el = document.querySelector("#chat input");
	if (document.activeElement === el) {
		return;
	}

	switch (key) {
		case 83: // down
		case 40:
			if (this.moving === false) {
				this.moving = true;
				this.send({
					type: eventMove,
					data: 0
				});
			}
			break;
		case 87: // up
		case 38:
			if (this.moving === false) {
				this.moving = true;
				this.send({
					type: eventMove,
					data: 1
				});
			}
			break;
		case 68: // right
		case 39:
			if (this.moving === false) {
				this.moving = true;
				this.send({
					type: eventMove,
					data: 2
				});
			}
			break;
		case 65: // left
		case 37:
			if (this.moving === false) {
				this.moving = true;
				this.send({
					type: eventMove,
					data: 3
				});
			}
			break;
	}
};

Game.onKeyUp = function(event) {
	var key = event.keyCode;
	var el = document.querySelector("#chat input");
	if (document.activeElement === el) {
		if (key === 27) {
			el.blur();
			return;
		}
		if (key === 13) {
			var msg = el.value.trim();
			el.value = "";
			el.blur();
			if (msg.length > 0) {
				this.send({
					type: eventMessage,
					data: msg
				});
			}
		}
		return;
	}

	switch (key) {
		case 13: // enter
			el.focus();
			break;
		case 83: // down
		case 40:
		case 87: // up
		case 38:
		case 68: // right
		case 39:
		case 65: // left
		case 37:
			// Stop moving
			this.moving = false;
			this.send({
				type: eventStopMove
			});
			break;
		default:
			console.debug("keyup: " + key);
			break;
	}
};
