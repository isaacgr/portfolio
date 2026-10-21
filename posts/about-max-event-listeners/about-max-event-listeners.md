---
title: "Max event listeners in NodeJS"
created_at: 2019-07-19
description: "Node sends an error when you add too many listeners to an event"
tags:
    - javascript
    - archive
    - nodejs
---

I came across an interesting problem when I was working on a project not too long ago.

My tests kept throwing an error, and at the time I had no idea what was causing it or what it meant.

```sh
(node:21613) MaxListenersExceededWarning: Possible EventEmitter memory leak detected.
11 event listeners added. Use emitter.setMaxListeners() to increase limit
```

With some research I found that this was a warning being thrown by Node. Specifically, it was being thrown because I kept adding event handlers to an event emitter and never removed them.

Ill give an example of how to cause this error. If you want to refresh yourself on what event emitters are, checkout my [previous post](https://blog.irowell.io/til/about-event-listeners-and-emitters/) about them.

We’re going to setup a class that inherits from `EventEmitter`. Its going to have a method that will print the current count and then increase it when called.

To trigger this method, we are going to add it as an event listener. So whenever we call the `event` event, it calls this function.

```js
const EventEmitter = require("events")

class Emitter extends EventEmitter {
  constructor() {
    super()
    this.count = 0
  }
  someFunction() {
    console.log(this.count)
    this.count++
  }
}

const e = new Emitter()

for (let i = 0; i < 11; i++) {
  e.on("event", e.someFunction)
}

e.emit("event")
```

And when you run it we get the following output.

```sh
0
1
2
3
4
5
6
7
8
9
10
(node:22863) MaxListenersExceededWarning: Possible EventEmitter memory leak detected. 11 event listeners added. Use emitter.setMaxListeners() to increase limit
```

Node throws this when we add more than 10 listeners for *one* specific event. If we were to have a single listener across multiple events, that would be fine. You can actually see the warning in their [source code](https://github.com/nodejs/node/blob/d022cb1bdd20f10483cbe988a950548df2dfd48c/lib/events.js#L192).

Now the warning contains information on how to handle this problem.

```sh
Use emitter.setMaxListeners() to increase limit
```

Oh great, so we can increase the maximum number of listeners a particular event has using this method…

But this would just be a workaround to the overall problem, and in general this is considered bad practice. Even if those handlers were called with data being successfully read they are no longer relevent and should be cleaned up after use.

To do this, we would want to remove the listeners after each call using the `removeListeners()` method.

> Note that removeListeners() requires a named function to remove an event handler from a listener, so if you are calling anonymous functions with the .on() method then this would not work.

To illustrate, Im going to emit `event` from within the for loop, and then immediately remove the listner. This will let us see each call being made before the listener is removed.

```js
const EventEmitter = require("events")

class Emitter extends EventEmitter {
  constructor() {
    super()
    this.count = 0
  }
  someFunction() {
    console.log(this.count)
    this.count++
  }
}

const e = new Emitter()

for (let i = 0; i < 11; i++) {
  e.on("event", e.someFunction)
  e.emit("event")
  e.removeListener("event", e.someFunction)
}
```

And now our output is simply

```sh
0
1
2
3
4
5
6
7
8
9
10
```

So we can clean up our event listeners after they are used to prevent the error, which is a more workable solution than just increasing the number of listeners. Our program is better for it, and we can expand in the future without people needing to know that they have to increase this arbitrary limit.

I hope this proves useful if you ever encounter this error in the future.

Thanks for reading!
