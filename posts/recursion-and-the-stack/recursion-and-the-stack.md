---
title: "Recursion and the Stack"
created_at: 2019-07-05
description: "An overview of what recursion is and how it works."
tags:
    - archive
skip: true
---

Recursion is a pretty common topic when you start programming, and it can be a little confusing to grasp especially for new developers.

I’m going to cover what recursion is and how it works.

## Recursion

So what is recursion?

Recursion is a method of solving a problem where the solution depends on the accumulated solutions of smaller instances of the same problem. Its important to note that this is different from iteration, which is where a loop is repeatedly executed to build on the solution until some condition is met.

A classic example of a recursive function in computer science is *Factorial*. This is the sequence where the next value depends on the sum of the previous values.

```js
const fact = n => {
  if (n === 1) {
    return 1
  }
  return n * fact(n - 1)
}
```

So we first check if `n` is 1, and just return that if so. Then it gets interesting…

```js
return return n * fact(n-1)
```

We return `n` multiplied by the same function with 1 subtracted from `n`. This will continue until we call `fact` with 1 at which point we just return 1. If we run this function with `n=4` we get `24` as the result. But how do we get a result at all? How does the function actually return the number we need?

## The Stack

The stack is a data structure used to store a collection of function calls and parameters. You can think of function execution like a stack of plates being picked up by a waiter. The plates are grabbed and stacked from the bottom up. The last plate picked up is the first one to be taken out. This is called **LIFO** (last in first out), and its how the stack operates.

![the_stack](resources/the_stack.jpg)

The stack is based on two principle operations: *push*, meaning to add an element to the stack, and *pop*, meaning to remove the most recently added element. Elements are pushed onto the stack when they’re called, and are popped off when they return a value. The stack then moves onto the previous element until it reaches the end of the stack and there are no more values to return.

In our example above, this is what the stack would look like.

![stack_recursion](resources/stack_recursion.jpg)

1. We call `fact(4)`, which returns `4xfact(3)`
2. `fact(3)` gets called, which returns `3xfact(2)`
3. `fact(2)` gets called, which returns `2xfact(1)`
4. Finally, `fact(1)` gets called, which returns 1

After the 4th step returns the value of 1, the stack starts poping off the remaining items until we get the final result of 24.

And thats effectively what happens with every recursive function. It calls itself until the final call returns a value, and then starts a chain reaction down the stack until the final result is given.

A real world example of this would be in `JSON.parse()`. It recursively calls the parse function on the json object that has been passed, until it has gone through the whole message at which point it returns the fully parsed message. It errors out if at any point its unable to parse the string correctly.

## Making your own recursive functions

There are some simple ways that you can make your own recursive functions. As we discussed above the main concept is that the function will call itself until some condition is met.

One that I often find comes up is for retry attempts.

Lets say you have a function that calls an external API to fetch some data. If the server hosting that API is no good, then you might experience connection type errors when trying to query it. I’ll use a real life example that I came across where the server would return a `503` response indicating its too busy to fufill the request.

We would want to retry the call if we receive this error, and we would want to do it for *x* amount of times until we decide to give up.

```js
const fetchData = (limit = 3) =>
  getData()
    .then(data => {
      return data;
    })
    .catch(error => {
      if (
        (error.statusCode === 503) && --limit) {
        console.log(`Retrying...${limit} attempts left.`);
        return fetchData(limit);
      } else {
        return [];
      }
    });
  });
```

So we call `fetchData` initally with a default `limit` value of 3. We call our `getData` function which returns a promise.

What does this syntax do?

```js
error.statusCode === 503 && --limit
```

Essentially what we’re saying is **if the status code is 503 and the limit hasn’t reached 0 *after* we decrement it then do something**. What we do is recall `fetchData` with the decremented `limit` value. Once we reach 0 for limit, we just return an empty list.

So we just keep calling the our function from within itself until the limit runs out.

Hopefully this has taken the veil off of recursion, recursive functions and how they relate to the stack. Try implementing some recursion in your future projects, if done correctly it can really simplify your codes workflow.

Thanks for reading!
