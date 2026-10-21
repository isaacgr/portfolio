---
title: "Classes vs. Prototypes"
created_at: 2019-05-19
description: "A brief explanation of the differences between protoypes and classes in JavaScript"
image: ./resources/class-proto.jpg
tags:
  - archive
  - javascript
skip: true
---

[Details of the object model - JavaScript | MDN](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Details_of_the_Object_Model)

---

## Classes vs Prototypes

You have probably read that Javascript is an object oriented language, like Java or C++. This is not entirely true, however. Javascript is a *prototype* based OOP language whereas Java or C++ are *class* based languages.

This can cause confusion on how to create objects, heirarchies, inheritence and methods with JavaScript.

Class based languages are founded on the concept of two entities: **classes** and **instances**.

A `class` defines all of the properties that characterize a set of objects.

An `instance` is the instantiation of a class. Ex. `Isaac` could be an instance of the `Person` class. This instance has the same properties as its parent.

A prototype based language like JavaScript does not have this distinciton. It just has objects.

Any object can specify its own properties, when created or during run time. As well, any object can be associated as the `prototype` for another object, allowing a second object to share the first’s properties (inheritence).

An example of how to implement a class using the prototype method would look something like this:

```js
function Person(name) {
  this.name = name
}

Person.prototype.greeting = function greeting() {
  console.log("Hello " + this.name)
}

let isaac = new Person("Isaac")

isaac.greeting()
// Hello Isaac
```

In this case, our `Person` object has `greeting` as a prototype method. We *instatiate* the `Person` obejct as `Isaac`, and then call the `greeting` method.

If we wanted to create an object which inherited properties from another, it would looks something like the below:

```js
function Developer(language) {
  this.language = language
}

let dev = new Developer("JavaScript")

dev.__proto__ = isaac
```

### Sorry, so whats the prototype?

Right. Its probably not entirely obvious from the examples above still.

All JavaScript objects inherit properties and methods from a prototype. Almost everything in JavaScript is an object. Only 6 things are not objects:

- null
- undefined
- strings
- numbers
- boolean
- symbols

These are called primitive values or primitive types.

Anything else is an Object. That includes arrays, functions, constructors and objects.

Thats right, functions are objects too. Specifically they are objects with an internal `[[Call]]` method.

Objects dont just have key-value pairs, they can store a link (a reference) to another object. This parent object is consulted when a property is read. This parent can also have a parent and so on until we reach an object that has no parent.

This chain is the `prototype chain`. By default, all objects have their parents as `Object.prototype`. There is no parent to this and typically where the chain ends.

So prototypes are just references to other objects. Its an internal property, so it cant be accessed directly. Note that We can use `Object.getPrototypeOf` to get the parent.

![prototype-of.png](resources/prototype-of.png)

### But where does the prototype live?

If we checked the console for the `isaac` object, we could see the protoypal chain:

![isaac-proto.png](resources/isaac-proto.png)

**Person** is an object with a name, it has an object `__proto__` which contains a function `greeting` and the `Person` constructor. The object also contains a `__proto__` object `Object`, which is the base prototype of all JavaScript objects.

Calling `dev` from the example above would return

![dev-inherit.png](resources/dev-inherit.png)

**Developer** is an object with a *language* and *name* property. It has a `prototype` of Person, which has a `Developer` constructor, which has a `__proto__` object.

We can see why this is called a protoype chain!

A big difference between prototypes and classical methods, is that the prototypes are *references*. This means if you do something like below

```js
delete isaac.__proto__.greeting
```

You will have removed the ability to call a greeting for isaac (and any other *Person* instances you had) because they all share the same reference to the **Person** prototype.

### Class it up

![giphy.gif](resources/giphy.gif)

Now since ES6, we have the `class` keyword.

According to MDN:

> JavaScript classes, introduced in ECMAScript 2015, are primarily syntactical sugar over JavaScript’s existing prototype-based inheritance. The class syntax does not introduce a new object-oriented inheritance model to JavaScript.

Lets revisit the Person example, using this new syntax

![person-class.png](resources/person-class.png)

Now we have something that looks like any other object oriented language.

Heres what it would look like if we extended the Person class to a *Developer* class

![dev-class.png](resources/dev-class.png)

We have a **class**, a **constructor**, theres also **super**, **extends** (inheritence from another class) and **static** methods using the new syntax!

Using `extends` means that we added depth to the protoype chain.

![chain.png](resources/chain.png)

Look familiar?

We still have an object with a name and a language, it has a `__proto__` object with the Developer contructor which has a `__proto__` object.

So the `class` keyword is essentially the same as the prototype approach as before, but it makes our code cleaner and easier to read as you scale up your application.

It should be noted, that the `class` keyword does not provide any actual benefits over the prototype method. There is likely something that you will be able to do extra in ES2019 or ES2020, such as [private fields](https://github.com/tc39/proposal-class-fields) or [private methods](https://github.com/tc39/proposal-private-methods).

Thanks for reading!
