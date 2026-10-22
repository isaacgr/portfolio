---
title: "How to access the correct 'this' inside of a function or class method"
created_at: 2020-03-01
description: "I explore how to access the correct 'this' when passing a function to a class method"
tags:
  - archive
  - javascript
skip: true
---

I was recently working through some problems with an app I’ve been writing, a [JSON-RPC Web client](https://github.com/isaacgr/jsonrpc-client).

I had two classes that I was working with: one to get instantiated as new clients connected to the server, and one which handled sending updates to clients. The second class would get instantiated as clients subscribed to updates.

A short example of the setup:

```js
const UpdateHandler = require("updatehandler")
class ClientCreator {
  constructor() {
    this.client = new Client()
  }
  handleClientUpdates() {
    const Handler = new UpdateHandler()
    return Handler.handleUpdates()
  }
}
```

Now the `Handler.handleUpdates()` method looks something like this

```js
handleUpdates(){
  ...
  this.server.doSomething()
}
```

And when I would call `handleClientUpdates()` I would get an error

```bash
TypeError: Cannot read property 'doSomething' of undefined
```

Now at first glance this doesnt appear to make sense right? My `Handler` class has `server` as a class attribute and that in turn has `doSomething` as a method, but the error implies that `server` is undefined.

### Investigation

If we were to log `this` in the `handleUpdates()` method, we see that it returns a reference to the `ClientCreator` class, not the `UpdateHandler` class. So the `this` keyword that we’re accessing is of the parent class, not the child class.

In fact, `this` is a keyword inside each function and its value depends on how the function was *called*, not how or where its *defined*.

A **method** in Javascript is just a term for a function that is a value of an object property. But that function doesn’t have a specific link to its **method** object. Consider the example:

```js
function Foo() {
  this.data = 42
  document.body.onclick = this.method
}

Foo.prototype.method = function() {
  console.log(this.data)
}
```

`this.method` is assigned as a click event handler, however if `document.body` is clicked then the value we log will be undefined since within the event handler `this`refers to`document.body`.

If we were to write the code a slightly different way, it may be more obvious why `this.data` is returning undefined.

```js
function method() {
  console.log(this.data)
}

function Foo() {
  this.data = 42
  document.body.onclick = this.method
}

Foo.prototype.method = method
```

The scope of `this` in the function `method()` is that method itself, and it has no `data` property. So when its called, it returns undefined.

### Solution

So what can we do about this?

There are a few things we can try.

1. You can use the `bind()` function.
2. Store a reference to the context/this inside another variable.
3. Use Arrow functions (introduced in ES6).

##### Using bind()

Every function has the method `.bind`, which returns a new function with `this` bound to a value. The function has exactly the same behaviour as the one that `.bind` was called on, except `this` has been set to whatever is passed to it.

```js
function Foo() {
  this.data = 42
  document.body.onclick = this.method.bind(this)
}

Foo.prototype.method = function() {
  console.log(this.data) // 42
}
```

##### Storing a reference to `this` inside another variable

Since `self` is a normal variable, it obeys lexical scope rules and is accessible inside the `onClick` handler.

```js
function Foo() {
  this.data = 42
  document.body.onclick = this.method
}

Foo.prototype.method = function() {
  var self = this
  console.log(self.data) // 42
}
```

##### Arrow functions

Arrow funcitons dont have their own `this` binding. Instead, `this` is looked up in scope like a normal variable.

```js
function Foo() {
  this.data = 42
  document.body.onclick = () => this.method
}

Foo.prototype.method = function() {
  console.log(this.data) // 42
}
```

### My solution

So for the original problem, I chose to use the `bind` method. Since `this` will use the reference thats been passed to the `bind` function, I bound the `handleUpdates` method to the `UpdateHandler` class. My function now looks like this:

```js
const UpdateHandler = require("updatehandler")
class ClientCreator {
  constructor() {
    this.client = new Client()
  }
  handleClientUpdates() {
    const Handler = new UpdateHandler()
    return Handler.handleUpdates().bind(Handler)
  }
}
```

Which solved my inital problem.
