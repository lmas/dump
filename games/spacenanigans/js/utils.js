////////////////////////////////////////////////////////////////////////////////
// utils.js

function padInt(i) {
	return i.toString().padStart(2, "0");
}

function formatShortTime(unixTime) {
	//var t = new Date(unixTime * 1000);
	var t = new Date(unixTime);
	var time = padInt(t.getHours()) + ":" + padInt(t.getMinutes());
	return time;
}

function formatFullTime(unixTime) {
	//var t = new Date(unixTime * 1000);
	var t = new Date(unixTime);
	var date = t.getFullYear() + "-" + padInt(t.getMonth()) + "-" + padInt(t.getDate());
	var time = padInt(t.getHours()) + ":" + padInt(t.getMinutes()) + ":" + padInt(t.getSeconds());
	return date + " " + time;
}

function shouldScroll(el) {
	return el.scrollHeight - el.scrollTop === el.clientHeight;
}

function scrollToBottom(el) {
	el.scrollTop = el.scrollHeight;
}

function addMessage(user, msg) {
	var now = new Date();
	var li = document.createElement("li");
	li.setAttribute("title", formatFullTime(now));
	var t = document.createElement("span");
	var u = document.createElement("span");
	var m = document.createElement("span");
	t.className = "msgtime";
	u.className = "msguser";
	m.className = "msgtext";
	t.innerText = formatShortTime(now);
	u.innerText = user;
	m.innerText = msg;
	li.appendChild(t);
	li.appendChild(u);
	li.appendChild(m);

	var el = document.querySelector("#chat ul");
	var scroll = shouldScroll(el);
	el.appendChild(li);
	if (scroll) {
		scrollToBottom(el);
	}
	while (el.childElementCount > maxHistory) {
		el.removeChild(el.firstChild);
	}
}

function hideWindow(win) {
	document.getElementById(win).className = "hidden";
}

function showError(msg) {
	document.querySelector("#error p").textContent = msg;
	document.querySelector("#error").className = "window";
}

function showConnectionError(msg) {
	document.querySelector("#connection-error").className = "window";
}

function getWebsocketAddr() {
	var proto = "ws";
	if (window.location.protocol === "https:") {
		proto = "wss";
	}
	return proto + "://" + window.location.host + "/game/ws";
}
