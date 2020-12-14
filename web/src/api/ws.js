const socket = new WebSocket("ws://localhost:5000/ws");

export const connect = callback => {
  console.log("Attempting Connection...");

  socket.onopen = () => {
    console.log("Successfully Connected");
  };

  socket.onmessage = msg => {
    callback(msg);
  };

  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  };

  console.log(socket)
  return socket;
};

export const sendMsg = msg => {
  socket.send(msg);
};
