---
title: "Understanding function scope in JavaScript"
created_at: 2019-06-13
description: "The use of arrow functions and the 'this' keyword"
tags:
    - archive
    - javascript
skip: true
---

If you’re familiar with object oriented programming, then you’re probably familiar with the various keywords used to reference the object in which they’re working.

For example, theres the `self` keyword in Python, which is used to reference the current class instance.

Many other languages such as C and Java use a `this` keyword to reference the current object instance. JavaScript also uses `this`, but it acts a little differently than in those other languages.

Lets first try to understand how the `this` keyword works in ES5.

```js
var Person = {
  name: "Daniel",
  greeting: function() {
    alert("Hi " + this.name)
  },
}

Person.greeting() // Hi Daniel
```

In the above code, we use `this` inside of an objects method. That means that `this` refers (binds) to that object.

The function is scoped to the object, so when we call the `this` keyword it refers to the `Person` object and our function can access the `name` property of that object using `this`.

Outside of a function, the `this` keyword refers to the global object (this is the `window` object in browsers, or the `global` object in Node).

```js
a = 1
console.log(this.a)
// 1
console.log(window.a)
// 1
```

**Note**: When in ‘strict’ mode the value of `this’ refers to whatever was set when entering the execution context

```js
function test() {
  "use strict"
  return this
}
test() === undefined //true
```

What we can take away from all of the above is that the `this` keyword in JavaScript **has different values depending on the context in which its used**.

As we can see, in a function *this* refers to the object its in.

By itself *this* refers to the window or global object.

When in ‘strict’ mode *this* is undefined.

### Going deeper

What happens to our scope if we went one layer deeper…

*Insert inappropriate meme*

```js
var Person = {
  name: "Daniel",
  greetings: ["Hi", "Hey", "Hello"],
  greeting: function() {
    this.greetings.forEach(function(greeting) {
      alert(greeting + " " + this.name)
    })
  },
}

Person.greeting()
```

We are declaring a function which is calling `this` within the class method. What do you think the expected output will be?

```js
Person.greeting()

// [object Window] Daniel
// [object Window] Daniel
// [object Window] Daniel
```

Now when we call `this` it refers to the window object. So what exactly is happening here?

When a function is called as a method of an object, its `this` is set to the object the method is called on.

However, our `this` keyword within the function in the objects method has now gone out of scope, and is referencing the `global` scope.

That means that our `this` keyword binds to the `window` or `global` object.

### Why does this happen?

Its a quirk of JavaScript, and is something that only happens with this language.

The Javascript developers recognize this as a poor design choice, and are trying to fix the implementation through the use of arrow functions in ES6.

Before we discuss those however, how did people work around this quirk?

There is a clever way we can modify the scope of `this` to fit our requirements.

### The “bind()” method

```js
var Person = {
  name: "Daniel",
  greetings: ["Hi", "Hey", "Hello"],
  greeting: function() {
    this.greetings.forEach(
      function(greeting) {
        alert(greeting + " " + this.name)
      }.bind(this)
    )
  },
}

Person.greeting()
// Hi Daniel
// Hey Daniel
// Hello Daniel
```

When called, the `bind()` method creates a new function with the same body and scope as the original, but has its `this` keyword set to the provided value.

In the case above, our function essentially became a new function where `this` was bound to the `Person` objects `this` keyword.

This is probably a lot to take in, but its always better to understand the basics.
We need to understand where we came from before trying to understand where we’re going and why we need to get there.

## So what’s an arrow function and how does it relate to all this?

ES6 introduced a new syntax for writing functions in JavaScript called **arrow functions**.

```js
const add = (a, b) => {
  return a + b
}
```

This isnt just a fancy new way of writing a function, this syntax comes with some improvements to how we handle scope and makes using the `this` keyword much more manageable.

### How?

Arrow functions use whats called *lexical* scoping. The `this` keyword now only refers to the surrounding scope. So our inner function knows to bind `this` to the inner function, not the object or its method.

```js
var Person = {
  name: "Daniel",
  greetings: ["Hi", "Hey", "Hello"],
  greeting: function() {
    this.greetings.forEach(greeting => {
      alert(greeting + " " + this.name)
    })
  },
}

Person.greeting()
// Hi Daniel
// Hey Daniel
// Hello Daniel
```

Nice, our function works again.

But what would happen if we tried to use an arrow function in place of our function declaration for *greeting*?

```js
var Person = {
  name: "Daniel",
  greetings: ["Hi", "Hey", "Hello"],
  greeting: () => {
    this.greetings.forEach(greeting => {
      alert(greeting + " " + this.name)
    })
  },
}

Person.greeting()
/**
VM355:5 Uncaught TypeError: Cannot read property 'forEach' of undefined
    at Object.greeting (<anonymous>:5:20)
    at <anonymous>:11:8
*/
```

When invoking the method greeting() on the Person object, the context still remains `window`. It happens because the arrow function binds the context *lexically* with the window object.

### How can I migrate existing code to use arrow functions?

Easily! The change in syntax is relatively painless as we can see below

```js
// ES5
var add = function(a, b) {
  return a + b
}

// ES6
var add = (a, b) => a + b
```

If we are not evaluating any expressions, we can *implicitly* return a value from the arrow function.

Otherwise, we can still just use the return statement

```js
var add = (a, b) => {
  return a + b
}
```

I hope this has given you some better insight on function scope and why and how we use `this` in JavaScript (and why it can be so confusing).

Thanks for reading!
