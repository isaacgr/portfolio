const net = require('net')

const server = new net.Server() 
server.listen({host: '127.0.0.1', port: 9999})
server.on('connection', (client) => {
    client.write('{"name": "isaac",')
    setTimeout(() => {
        client.write('"age": "28"}\n{"name": "steve",')
    }, 100)
})

const client = new net.Socket()
client.connect(9999, '127.0.0.1')

let received = ""
client.on('data', (data) => {
  const time = new Date()
  received += data
  const messages = received.split('\n') // remember this
  if (messages.length > 1 ){
    for (let message of messages){
      if (message !== ""){
        console.log(JSON.parse(message))
        received = ""
      }
    }
  }
  const end = new Date() - time
  console.log(end)

})

client.on('close', () => {
  console.log('connection closed')
})

