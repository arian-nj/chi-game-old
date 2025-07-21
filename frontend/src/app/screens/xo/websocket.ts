export function connectSocket(url: string, protocol?: string | string[]) {
  const socket = new WebSocket(url, protocol);
  console.log(socket);

  socket.onopen = function (_e: Event) {
    alert("[open] Connection established");
    console.log(_e);
    socket.send("My name is John");
  };

  socket.onmessage = function (event) {
    alert(`[message] Data received from server: ${event.data}`);
  };

  socket.onclose = function (event) {
    if (event.wasClean) {
      alert(
        `[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`,
      );
    } else {
      alert("[close] Connection died");
    }
  };
}
