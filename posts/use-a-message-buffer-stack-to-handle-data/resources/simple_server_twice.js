const net = require("net")

const server = new net.Server()
server.listen({ host: "127.0.0.1", port: 9999 })
server.on("connection", client => {
  client.write("Hello")
  setTimeout(() => {
    client.write(" Isaac\n")
  }, 100)
})

const client = new net.Socket()
client.connect(9999, "127.0.0.1")

let received = ""
client.on("data", data => {
  received += data
  console.log(received)
})

client.on("close", () => {
  console.log("connection closed")
})
