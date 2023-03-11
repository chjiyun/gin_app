
const socket = io("ws://localhost:8008", {
  transports: ["websocket", "polling"], // use WebSocket first, if available
  protocols: ["chat"],
  reconnectionDelayMax: 10000,
  path: "/",
  reconnection: false,
  // auth: {
  //   token: "20221202"
  // },
  query: {
  }
});

// socket.emit("test", {hello: "world"})

socket.on("connect", () => {
  console.log("connected: ", socket.connected);
  console.log("id: ", socket.id);
});

// const conn = new WebSocket("ws://localhost:8008", ["chat"])
// conn.addEventListener("open", ev => {
//   console.info("websocket connected")
// })

// // This is where we handle messages received.
// conn.addEventListener("message", ev => {
//   console.log(ev.data)
//   if (typeof ev.data !== "string") {
//     console.error("unexpected message type", typeof ev.data)
//     return
//   }
// })

// conn.addEventListener("close", ev => {
//   console.log(`WebSocket Disconnected code: ${ev.code}, reason: ${ev.reason}`)
//   if (ev.code !== 1001) {
//     // setTimeout(dial, 1000)
//   }
// })
