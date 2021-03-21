"use strict";

const maxHistory = 1000;

var elHistory = document.querySelector("#history");
var elInput = document.querySelector("#message input");

window.onload = function() {
	elInput.addEventListener("keyup", keyEvent);
	wsConnect("ws://127.0.0.1:8080/ws");
};

function keyEvent(event) {
	var key = event.keyCode;
	if (key !== 13) {
		return;
	}
        wsSend(elInput.value.trim());
	elInput.value = "";
	scrollToBottom(elHistory, true);
}

function addMessage(msg) {
	var li = document.createElement("li");
	//li.innerText = msg;
	li.innerHTML = msg; // Enables us to use html in msgs, which is unsafe

	var scroll = shouldScroll(elHistory);
	elHistory.appendChild(li);
	scrollToBottom(elHistory, scroll);
	while (elHistory.childElementCount > maxHistory) {
		elHistory.removeChild(elHistory.firstChild);
	}
}

function shouldScroll(el) {
	return el.scrollHeight - el.scrollTop === el.clientHeight;
}

function scrollToBottom(el, scroll) {
	if (scroll) {
		el.scrollTop = el.scrollHeight;
	}
}

////////////////////////////////////////////////////////////////////////////////
// websock stuff

// socket ready states:
// 0    CONNECTING      Socket has been created. The connection is not yet open.
// 1    OPEN            The connection is open and ready to communicate.
// 2    CLOSING         The connection is in the process of closing.
// 3    CLOSED          The connection is closed or couldn't be opened.

var sock;

function wsConnect(addr) {
        sock = new WebSocket(addr);
        sock.addEventListener("open", wsOnOpen);
        sock.addEventListener("close", wsOnClose);
        sock.addEventListener("error", wsOnError);
        sock.addEventListener("message", wsOnMsg);
}

function wsSend(msg) {
        if (!sock || msg.length < 1) {
                return;
        }
        //console.debug("send: " + msg);
        sock.send(msg);
}

function wsClose() {
        sock.close();
}

function wsOnOpen(event) {
        //console.debug("open: " + sock.readyState);
        //addMessage("Connected");
}

function wsOnClose(event) {
        //console.debug("close: " + sock.readyState);
        addMessage("Disconnected");
}

function wsOnError(event) {
        //console.debug("err: " + sock.readyState + " " + event + "(" + JSON.stringify(event.data) + ")");
        addMessage("Error: " + event + "(" + JSON.stringify(event.data) + ")");
}

function wsOnMsg(event) {
        //console.debug("msg: " + msg);
        var lines = event.data.split("\n");
        for (var i = 0; i < lines.length; i++) {
                addMessage(lines[i]);
        }
}
