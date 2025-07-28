export function ConnectSocket(url: string, protocol?: string | string[]) {
  const socket = new WebSocket(url, protocol);

  socket.onopen = function (_e: Event) {
    console.log("[open] Connection established");
    console.log(_e);
  };

  socket.onmessage = function (event) {
    try {
      const data = JSON.parse(event.data);
      console.log(data);
    } catch (err) {
      console.error("[WebSocket] Failed to parse message:", event.data);
    }
    // socket.close(1000, "Normal Closure");
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
