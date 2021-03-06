// api/index.js
var socket = new WebSocket("ws://"
			+process.env.REACT_APP_BASE_BACKEND_URI+":"
			+process.env.REACT_APP_BASE_BACKEND_PORT+"/v1/ws");
const IN_SYNC = "database_in_sync";
const OUT_OF_SYNC = "database_out_of_sync";

let connect = cb => {
  console.log("Attempting Connection...");

  socket.onopen = () => {
    console.log("Successfully Connected");
    cb(IN_SYNC);
  };

  socket.onmessage = msg => {
    cb(msg);
  };

  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
    cb(OUT_OF_SYNC)
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  };
};

let sendMsg = msg => {
  console.log("sending msg: ", msg);
  socket.send(msg);

};

export { connect, sendMsg };
