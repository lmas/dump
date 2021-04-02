////////////////////////////////////////////////////////////////////////////////
// websocket.js

var eventUnknown = 0;
var eventID = 1;
var eventMap = 2;
var eventConnect = 3;
var eventDisconnect = 4;
var eventMessage = 5;
var eventMove = 6;
var eventStopMove = 7;
var eventSeeMob = 8;
var eventHideMob = 9;

Game.connect = function(addr) {
	this.conn = new WebSocket(addr);
	this.conn.addEventListener("open", this.onOpen.bind(this));
	this.conn.addEventListener("close", this.onClose.bind(this));
	this.conn.addEventListener("error", this.onError.bind(this));
	this.conn.addEventListener("message", this.onEvent.bind(this));
};

Game.send = function(e) {
	var data = JSON.stringify(e);
	//console.debug("send: " + data);
	this.conn.send(data);
};

Game.onOpen = function(e) {
	//console.debug("connected");
};

Game.onClose = function(e) {
	//console.debug("disconnected");
	//showError("Disconnected from server");
	showConnectionError();
};

Game.onError = function(e) {
	//console.debug("ws error: " + e + "(" + e.data + ")");
	//showError("websocket: " + JSON.stringify(e));
};

Game.onEvent = function(e) {
	//console.debug("event: " + e.data);
	var ev = JSON.parse(e.data);
	switch (ev.type) {
		case eventID:
			this.camera.follow(ev.data);
			break;

		case eventMap:
			this.mapManager.set(ev.data);
			break;

			//case eventConnect:
			//addMessage("system", ev.data.name + " connected");
			//break;

			//case eventDisconnect:
			//addMessage("system", ev.data.name + " disconnected");
			//this.mobManager.del(ev.data.id);
			//break;

		case eventMessage:
			addMessage(ev.data.name, ev.data.text);
			var mob = this.mobManager.get(ev.data.id);
			if (mob !== undefined) {
				this.bubbleManager.add(mob.name + ": " + ev.data.text, mob.x, mob.y);
			}
			break;

			//case eventMove:
			//break;

			//case eventStopMove:
			//break;

		case eventSeeMob:
			this.mobManager.set(ev.data);
			if (this.camera.following(ev.data.id)) {
				this.camera.move(ev.data.x, ev.data.y);
				this.mapManager.dirty();
			}
			break;

		case eventHideMob:
			this.mobManager.del(ev.data);
			break;

		default:
			console.debug("Unhandled event: " + e.data);
			break;
	}
};
