const net = require('net')

class ServerProtocol {
  constructor(client){
    this.client = client
    this.buffer = ""
    this.factory = null
  }
  handleData(){
    this.client.on('data', (data) => {
      this.buffer += data
      console.log(this.buffer)
    })
    this.client.on('end', () => {
      this.factory.connectionClosed(this.client)
      console.log('client disconnected')
    })
  }
}

class ServerFactory {
  constructor(){
    this.server = new net.Server()
    this.connectedClients = []
  }
  listen(){
    this.server.listen({port: 8100},() => {
      console.log(`server listening on ${this.server.address().port}`)
    })
  }
  start(){
    this.server.on('connection', (client) => {
      console.log('client connected')
      this.connectedClients.push(client)
      const serverProtocol = new ServerProtocol(client)
      serverProtocol.factory = this
      serverProtocol.handleData()
    })
  }
  connectionClosed(client){
    this.connectedClients.splice(this.connectedClients.indexOf(client)) // remove the client from the connected clients list
  }
}
const serverFactory = new ServerFactory()
serverFactory.listen()
serverFactory.start()
