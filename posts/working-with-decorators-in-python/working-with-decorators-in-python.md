---
title: "Working with decorators in Python"
created_at: 2019-11-20
description: "My attempt to explain decorators, and show why you would use them into your application"
tags:
  - archive
  - python
skip: true
---

I’ve been using the [klein](https://github.com/twisted/klein) web framework at my job to put together a changelog generation tool. I wanted to have a simple frontend and api to fetch all of the changes between two repo versions on stash and display the tagged jira tickets (as well as a verbose mode to display the commits).

Part of this includes setting up basic routes. For example:

```python
...
  @app.route('/')
  def handle_fetch_data(self, request):
    return self.data
...
```

Here we see klein uses this `route` decorator to pass a request object to my function, and then handle sending the return value of the function as a body message to the client.

One additional thing that I needed to do was to set the `Content-type` header since I’m sending json back. Typically this is just done within the function (something like `request.setHeader('Content-type', 'application/json')`). However I have several routes setup, and needed to set some additional headers, so I figured I could do this a little more efficiently with something like a decorator.

## What is a decorator?

A decorator is a callable function, that takes another function as an argument (this is the *decorated* function).

The decorator typically performs some processing with the decorated function, and then returns it or replaces it with another function, or callable object. In other words **decorators replace the decorated function with a different one**.

An important concept when trying to understand decorators, is that Python functions are actually objects. As an object, you can assign the function to a variable.

```python
def talk():
  print 'Hello'

speak = talk
```

Then calling that funciton will produce the same output

```bash
>>> speak()
Hello
```

Additionally, we are able to define functions within another function.

```python
def talk(word):
  def captialize(word):
    return word.upper()
  return captialize(word)

talk('hello') # outputs 'HELLO'
```

Since we defined `captialize` inside of `talk` its called every time we call `talk`.

Its important to note that the `capitalize` function is locally scoped to `talk` and does not exist outside of that function definition.

### Writing your own decorator

Now we want to take these concepts and write our own decorator. We want to take in our target function, run some commands, then return our original funciton to be called.

A decorator is nothing more than a function that takes another function as an argument. We can pass this target funciton to our decorator using the `@` syntax as shown below.

```python
def decorator(func):
  def wrapper():
    print('This runs before.')
    return func() # then our original function will be called
  return wrapper

@decorator
def target():
  print('Running target function')
```

Which has the effect to print out

```bash
>>> target()
This runs before.
Running target function
```

The use of the `@` symbol is just syntactic sugar. The above code has the same effect as writing:

```python
def decorator(func):
  def wrapper():
    print('This runs before.')
    return func() # then our original function will be called
  return wrapper

def target():
  print('Running target function')

target = decorator(target)
```

The end result is the same. At the end of both snippets, the target name doesnt necessarily refer to the original funciton, but to whatever function was returned by our decorator (in this case the wrapper function). Decorators are just a variant of the [decorator design pattern](https://en.wikipedia.org/wiki/Decorator_pattern).

## Passing arguments to the decorated function

A large part of using decorators is to pass an argument to the function being called. We can pass arguments to the target function using the wrapper function. Since calling the decorated function means calling the wrapper as well, passing arguments to the wrapper will let it pass them to the target function.

```python
def decorator(func):
    def wrapper(arg1, arg2):
        print 'Got args:', arg1, arg2
        return func(arg1, arg2)
    return wrapper

@decorator
def print_name(first_name, last_name):
    print 'My name is', first_name, last_name
```

Which has the output

```bash
>>> print_name('Isaac', 'Rowell')
Got args: Isaac Rowell
My name is Isaac Rowell
```

You could see as well how we would be able to modify these arguments within our wrapper function, before sending them back to the target function

##### My solution for setting the klein request headers

Just because I know you’re dying to find out, here was how I implemented the decorator to set the request headers for my tool:

```python
def set_request_headers(func):
  def wrapper(self, request):
    request.setHeader('Content-Type', 'application/json')
    return func(self, request)
  return wrapper
```

And then using the decorator with our routes

```python
...
  @app.route('/')
  @set_request_headers
  def handle_fetch_data(self, request):
    return self.data
...
```

This gives us the flexibility to add more headers, without changing any of our main code.

Now I’m sure you have a couple of questions about the above implementation.

1. Why didnt I just use a function?
2. Does the order matter when using two decorators?

To answer number 1, you’re right. I could have just used a separate method in each function to set the requests. But this is an exploration on how we can use decorators, so that wouldnt have been very useful. To clarify, decorators allow you to define reusable building blocks that can change or extend the behavior of other functions. They let you do that without permanently modifying the wrapped function itself. The function’s behavior changes only when it’s decorated.

To answer number 2, yes the order of operations does matter when defining multiple decorators. Decorators wrap the function they are decorating. So the top most decorator will wrap the result of the second decorator and so on.

Theres a good example of this on [StackOverflow](https://stackoverflow.com/questions/27342149/decorator-execution-order).

```python
def make_bold(fn):
    return lambda : "<b>" + fn() + "</b>"

def make_italic(fn):
    return lambda : "<i>" + fn() + "</i>"

@make_bold
@make_italic
def hello():
  return "hello world"

helloHTML = hello()
```

The `make_bold` decorator wrapped the result of `make_italic`, whic wrapped our original function. So the output is:

```bash
<b><i>hello world</i></b>
```

## Best practices (The importance of functools)

Some things to keep in mind.

1. Decorators slow down the function call.
2. You cannot un-decorate a function. So once a function is decorated, it’s decorated for all the code.
3. Decorators wrap functions, which can make them hard to debug.
4. The `functools` module was introduced in Python 2.5. It includes the function `functools.wraps()`, which copies the name, module, and docstring of the decorated function to its wrapper.

An important note from above is the use of `functools`. Since a decorator essentially replaces a function with another function, functools allows us to keep information about a function. Lets take a look at an example.

```python
def decorator(func):
  def wrapper():
    print('This runs before.')
    print(func.__name__ + ' was called')
    return func() # then our original function will be called
  return wrapper

@decorator
def target():
  """This is the target function"""
  print('Running target function')

target()
print(target.__name__)
```

This will have the output

```bash
This runs before.
target was called
Running target function
wrapper
```

When we print the target name, it prints `wrapper`. And if we were to print `target.__doc__` it would be blank, since wrapper doesnt have a docstring. If using a decorator always meant losing this information about a function, it would be a serious problem. That’s why we have `functools.wraps`. This takes a function used in a decorator and adds the functionality of copying over the function name, docstring, arguments list, etc.

The following code does the correct thing:

```python
import functools

def decorator(func):
  @functools.wraps(func)
  def wrapper():
    print('This runs before.')
    print(func.__name__ + ' was called')
    return func() # then our original function will be called
  return wrapper

@decorator
def target():
  """This is the target function"""
  print('Running target function')

target()
print(target.__name__) # prints 'target'
```

## What can decorators be used for?

While decorators seem cool, I’m sure you’re still wondering what practical application they might have.

Well, there are literally thousands of possibilities. Things like extending function behaviour from an external library, or for debugging. You can also use them to extend several functions at once in a `DRY` way.

For example, lets take a look at a decorator that prints out the activity of the funciton being called.

```python
def logging(func):
    """
    A decorator that logs the details of a function
    """
    def wrapper(*args, **kwargs):
        res = func(*args, **kwargs)
        print func.__name__, args, kwargs
        return res
    return wrapper

@logging
def say_hello(name):
    print 'Hello ' + name
```

Will output

```bash
>>> say_hello('Isaac')
Hello Isaac
say_hello ('Isaac',) {}
```

Python itself provides several decorators: `property`, `staticmethod`, etc. So checkout the docs and explore the extensive world of decorators.

Thanks for reading!
