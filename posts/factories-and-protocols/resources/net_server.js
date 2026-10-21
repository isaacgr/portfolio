const net = require('net')
const server = new net.Server()
server.listen({port: 8100},() => {
  console.log(`server listening on ${server.address().port}`)
})
server.on('connection', (client) => {
  console.log('client connected')
  client.on('data', (data) => {
    console.log(data.toString())
  })
  client.on('end', () => {
    console.log('client disconnected')
  })
})
