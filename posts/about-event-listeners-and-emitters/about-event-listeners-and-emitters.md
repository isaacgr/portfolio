---
title: "Event Emitters"
created_at: 2019-07-02
description: "What event emitters are, and how can you use them in your application"
tags:
    - javascript
    - archive
    - nodejs
---

Not too long ago I started working on an NPM library I hope to publish soon, and I found that I was using a lot of functions that had the following syntax:

```js
socket.on("data", data => {
  // ...do something
})
```

It looks like this `on` method is waiting for an event, and then firing a `callback` when that event is emitted.

The function that triggers the `on` method is called an *event emitter* and the function that listens for this method is called the *event listener*.

Whats actually happening here is that this `socket` instance is *listening* for an event called `data`, and when that event name is *emitted* by some other source, it executes the callback with the paramters it passed.

## What are events?

Events are objects called `emitters` that emit an event by some name. Functions that are waiting for these events to be emitted are called `listeners`.

All objects that emit events are instances of the `EventEmitter` class. This class exposes an `eventEmitter.on()` function which allows other functions to listen for events emitted by that object.

Its important to note that events are called synchronously, so if you have multiple functions listening for that event they will be called one after the other.

An example of an an `EventEmitter` instance with a single listener is shown below. The `eventEmitter.on()` method is used to register listeners, and the `eventEmitter.emit()` method triggers the event.

```js
// Create a class with custom emitters
const EventEmitter = require("events")

class Emitter extends EventEmitter {
  constructor() {
    super()
    this.count = 0
    // initialize method to listen for events
    // after instatiation
    this.listenForEvent()
  }

  listenForEvent() {
    this.on("customEvent", () => {
      this.count++
      console.log(this.count)
    })
  }
}

const emitter = new Emitter()
// returns true if the event had listeners
emitter.emit("customEvent")
// 1
// true
```

## Why would I use these?

Good question!

If it makes sense for your code to subscribe to something rather than receive a callback, you should consider event emitters. Typically if multiple blocks of code in your application need to do something when an event happens.

An example from my library would be in the case of listening for a client connection being made to the server. In this case, when the `clientConnected` event is emitted, it triggers the callback for this function. It will only be called if at some point the user has instantiated the `clientConnected()` method and will do nothing otherwise.

```js
class Server extends EventEmitter {
  ...
  this.on('connection', (client) => {
    ...
    this.emit('clientConnected', client)
    ...
  })
  ...
  clientConnected(cb) {
    this.on("clientConnected", client => cb({
      host: client.remoteAddress,
      port: client.remotePort
    }));
  }
}

// So if we have called that function elsewhere in our code
server.clientConnected(() => {
  console.log('client connected')
})
```

So the methods that get instantiated are `subscribed` to the events, and when the named event is emitted, it triggers an action on those methods. If we wanted to trigger other actions when a client connected, we would simply have to add the listener to that function. This allows us to attach multiple listeners to a particular event, without having seperate callback functions for each one.

I’d like to note that this also saved me from importing a bunch of functions/methods into files of my library. I was simply able to emit the event and since the same class instance was listening for those events, they triggered different actions.

This was a short and sweet post but I hope it gave some basic insight on event emitters and listeners, and how they can simplify different workflow requirements in your application.

Thanks for reading!
