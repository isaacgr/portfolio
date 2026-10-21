---
title: "The Factory design pattern"
created_at: 2019-08-20
description: "How we can use the Factory creational pattern and protocol instances to handle client connections to a NodeJS server"
tags:
  - archive
  - design patterns
skip: true
---

## The Problem

Imagine you’re creating a server to host requests. Clients connect to the server, send some data and then the server sends back a response.

To track the data, you may keep the messages in a buffer, parse them and send out the response.

But what if you have multiple clients connecting to the server? If your buffer is a class attribute, then all those clients are writing to the same buffer and so you could end up with issues parsing the requests.

Sure, you could keep a buffer for every client, but then you would have to track that and reference it, which makes the code that much less maintainable and readable.

### The Factory Pattern

The Factory pattern deals with the problem of creating objects without the need to specify the exact class of the object being created. A factory can provide a generic interface for creating objects, where we can speicfy the type of factory object we want to be created.

A solution to the problem above would be to use `Factory` classes and `Protocol` instnaces.

When a client connects, the Factory would spin up a protocol instance (another class) of the desired type (http or tcp for example) which handles all the logic for message buffering and parsing. Then when the connection is done that class is destroyed, but the factory remains to track things like state and connection status.

You can find factory methods in frameworks where the code may need to create objects of particular types, which may be subclassed by scripts using the frameworks.

The Factory pattern suggests that you replace direct object construction calls (using the `new` operator) with calls to a special *factory* method.

There are several cases where you would want to consider using the Factory pattern

1. Our object or component setup involves a high level of complexity
2. We dont know ahead of time what object we need
3. To centralize the class selection code, in the case where we’re working with many small objects that share the same properties
4. When the user doesnt need to know every subclass
5. To encapsulate object creation

Lets take a look at an example of this problem.

## Setting up a server

I’m going to use Node’s `net` module to spin up a server. If you arent familiar with it, checkout the documentation [here](https://nodejs.org/api/net.html). Its a very simple library to create server and client instances that’s built into node, so no external libraries are required. This server is going to listen for data from incoming connections and print it to the console.

```js
const net = require("net")

class Server {
  constructor() {
    this.server = new net.Server()
  }
  listen() {
    this.server.listen({ port: 8100 }, () => {
      console.log(`server listening on ${server.address().port}`)
    })
  }
  start() {
    this.server.on("connection", client => {
      console.log("client connected")
      client.on("data", data => {
        console.log(data.toString())
      })
      client.on("end", () => {
        console.log("client disconnected")
      })
    })
  }
}
const server = new Server()
server.listen()
server.start()
```

If we run this and use [telnet](https://en.wikipedia.org/wiki/Telnet) to connect to the interface, then whatever we type to the screen will be displayed by our server.

```bash
telnet localhost 8100
```

And we can do that from multiple terminals, so we’ll have multiple clients connected to this server instance.

Now lets say those clients are sending data to our server. If we wanted to do something with that data, then we would likely want to store it (lets say in a class attribute)to be used in other methods. The problem with this is that its a global value for this class, so if data from client 1 is fragmented and data from client 2 is fragmented, then its possible that we could recieve a mix of data from each client in our buffer. Depending on what we are doing with this data, it could lead to errors when trying to parse it.

## Setting up a Protocol

I’m using the concept of a protocol instance from Pythons [Twisted](https://twisted.readthedocs.io/en/latest/core/howto/servers.html) library. This library includes built in methods for handling network data, client connection status, and classes you can inherit from to build factories to instantiate protocols.

We can adapt these concepts to make a similar structure for our server application.

The whole purpose of a protocol instance is to represent a single client connection. The protocol instance is created when a client connects, responds to data events as they come in and makes calls to methods in our protocol class to handle the data. The persistent state is kept in the factory (since this will be created once).

Lets see how we can modify the above code to use a protocol instead. We’ll also store our data in a `buffer` class attribute in that protocol, so we can reference it in other methods later on if we wanted.

```js
const net = require("net")

class ServerProtocol {
  constructor(client) {
    this.client = client
    this.buffer = ""
  }
  handleData() {
    this.client.on("data", data => {
      this.buffer += data
      console.log(this.buffer)
    })
    this.client.on("end", () => {
      console.log("client disconnected")
    })
  }
}

class Server {
  constructor() {
    this.server = new net.Server()
  }
  listen() {
    this.server.listen({ port: 8100 }, () => {
      console.log(`server listening on ${this.server.address().port}`)
    })
  }
  start() {
    this.server.on("connection", client => {
      console.log("client connected")
      const serverProtocol = new ServerProtocol(client)
      serverProtocol.handleData()
    })
  }
}
const server = new Server()
server.listen()
server.start()
```

Now a new `clientProtocol` class is created for every client that gets connected, and we can see how we can use that to manage data received by each client. Each protocol class is keeping track of its own message buffer and connection state.

## Setting up the Factory

So how does the factory come into play here? Seems like we have everything we need.

Well sort of.

The design pattern calls for a single class that instantiates our protocol for every connection. This class would be responsible for keeping connection state, such as which clients are still connected and which have disconnected. We could also keep other state here if the application called for it (if we were making requests to an external api for example).

How would this look? Well we would want to create an instance of our factory when the server starts up, so we can move the connection logic from the `Server` class to the factory class.

```js
class ServerFactory {
  constructor() {
    this.server = null
  }
  start() {
    this.server.on("connection", client => {
      console.log("client connected")
      const serverProtocol = new ServerProtocol(client)
      serverProtocol.handleData()
    })
  }
}
```

We can then use this class to manage state for our clients. We can do this by passing a reference to the factory to the protocol. This way the protocol instance can call methods on the factory class, which allows our factory to store those in some persistent state. Expanding on the factory above:

```js
class ServerFactory {
  constructor(server) {
    this.server = server
    this.connectedClients = []
  }
  start() {
    this.server.on("connection", client => {
      console.log("client connected")
      this.connectedClients.push(client)
      const serverProtocol = new ServerProtocol(client)
      serverProtocol.factory = this
      serverProtocol.handleData()
    })
  }
  connectionClosed(client) {
    this.connectedClients.splice(this.connectedClients.indexOf(client)) // remove the client from the connected clients list
  }
}
```

## Putting it all together

So if we bring all of this code together, we can actually replace our inital server class with the factory (they effectively do the same thing at this point, by creating the protocol instance and managing connection state) and we have something like this:

```js
const net = require("net")

class ServerProtocol {
  constructor(client) {
    this.client = client
    this.buffer = ""
    this.factory = null
  }
  handleData() {
    this.client.on("data", data => {
      this.buffer += data
      console.log(this.buffer)
    })
    this.client.on("end", () => {
      this.factory.connectionClosed(this.client)
      console.log("client disconnected")
    })
  }
}

class ServerFactory {
  constructor() {
    this.server = new net.Server()
    this.connectedClients = []
  }
  listen() {
    this.server.listen({ port: 8100 }, () => {
      console.log(`server listening on ${this.server.address().port}`)
    })
  }
  start() {
    this.server.on("connection", client => {
      console.log("client connected")
      this.connectedClients.push(client)
      const serverProtocol = new ServerProtocol(client)
      serverProtocol.factory = this
      serverProtocol.handleData()
    })
  }
  connectionClosed(client) {
    this.connectedClients.splice(this.connectedClients.indexOf(client)) // remove the client from the connected clients list
  }
}
const serverFactory = new ServerFactory()
serverFactory.listen()
serverFactory.start()
```

We can replace our initial `Server` class with a Factory which creates `Protocol` instances whenever a client connects.

These instances exist for every client that connects to our application, and are destroyed when the client disconnects. We use these to handle receiving data, and can send that data back to our factory instance using `this.factory`. That way, all of our state is kept somewhere persistent.

The benefit here is that we can use a single factory to spin up different protocol instances conditionally. Depending on the client connecting, we can handle the connection and the data independent of our base class. This helps us manage class instances with a common interface.

Thanks for reading!
