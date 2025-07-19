export function connectSocket() {
	let socket = new WebSocket("wss://javascript.info/article/websocket/demo/hello")
	socket.onopen = function(_e) {
		alert("[open] Connection established");
		socket.send("My name is John");
	}

	socket.onmessage = function(event) {
		alert(`[message] Data received from server: ${event.data}`);
	}

	socket.onclose = function(event) {
		if (event.wasClean) {
			alert(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
		} else {
			alert('[close] Connection died');
		}
	};
}
