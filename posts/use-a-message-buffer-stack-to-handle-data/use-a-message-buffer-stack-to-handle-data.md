---
title: "A better way to handle incoming data from a NodeJS socket"
created_at: 2019-07-31
description: "How we can use a stack to buffer messages and deal with TCP fragmentation"
tags:
  - javascript
  - node
  - tcp
  - archive
---

I want to explore what I believe is a very efficient and scalable way to buffer messages coming in on a socket in node. This can be extended to your client or server, and is a method I didnt find was often discussed when searching for how a user can work with the incoming data.

If you are familiar with using NodeJS for server/client socket connections, then you are likely familiar with the below basic example for setting up a client and server TCP connection (or something similar) using the `net` module.

```js
const net = require("net")

const server = new net.Server()
server.listen({ host: "127.0.0.1", port: 9999 })
server.on("connection", client => {
  client.write("Hello\n")
})

const client = new net.Socket()
client.connect(9999, "127.0.0.1")

let received = ""
client.on("data", data => {
  received += data
  console.log(received) // Hello
})
client.on("close", () => {
  console.log("connection closed")
})
```

The basics of this example are the following

1. Setup a tcp server listening on port 9999 of our local machine
2. Setup a client to connect to that host and port
3. Have the client listen for incoming data. When data is recieved, it comes in as either a `Buffer` or a `String`. We concat this to an existing string to do the type conversion.
4. When the client is connected to the server, send the data. Our client will log it as it comes in.

> Note that I’m adding a `\n` as a delimiter in our message. Adding something like this to your TCP client/server is standard practice, since the machine needs a way to distinguish between messages. I’ll discuss this later on.

## Analyzing the code

I want to focus on what happens in step 3. The chunk of code we want to focus on is below

```js
let received = ""
client.on("data", data => {
  received += data
  console.log(received) // Hello
})
```

When data is received by the client (or server), we concat the data to a string as it comes in and then do something with that data. We could parse it, pass it to a function, or just print it as we’re doing here.

If we were to do an [ngrep](http://ngrep.sourceforge.net/usage.html) or a [tcpdump](https://www.hugeserver.com/kb/install-use-tcpdump-capture-packets/) we could see exactly how the packets are coming in to the server.

For this exmple, I’ll use `tcpdump`.

```bash
tcpdump -i any port 9999 -X -s0
```

Run the above command, then run our line of code. You should see something like the folowing in your terminal.

```bash
07:22:03.113371 IP localhost.9999 > localhost.49904: Flags [P.], seq 1:7, ack 1, win 342, options [nop,nop,TS val 1197270971 ecr 1197270971], length 6
	0x0000:  4500 003a da87 4000 4006 6234 7f00 0001  E..:..@.@.b4....
	0x0010:  7f00 0001 270f c2f0 2d14 847a 355e 3122  ....'...-..z5^1"
	0x0020:  8018 0156 fe2e 0000 0101 080a 475c e7bb  ...V........G\..
	0x0030:  475c e7bb 4865 6c6c 6f0a                 G\..Hello.
```

Without diving into specifics, we see the client request a connection to the server. The server acknowledges. The client acknowledges. The server then sends the data to client (you can see the `Hello` message printed there.)

So we saw our data come in on a single IP packet here

```bash
07:22:03.113371 IP localhost.9999 > localhost.49904: Flags [P.], seq 1:7, ack 1, win 342, options [nop,nop,TS val 1197270971 ecr 1197270971], length 6
	0x0000:  4500 003a da87 4000 4006 6234 7f00 0001  E..:..@.@.b4....
	0x0010:  7f00 0001 270f c2f0 2d14 847a 355e 3122  ....'...-..z5^1"
	0x0020:  8018 0156 fe2e 0000 0101 080a 475c e7bb  ...V........G\..
	0x0030:  475c e7bb 4865 6c6c 6f0a                 G\..Hello.
```

And we can infer that the `.on('data')` listener is triggered everytime one of these packets is sent to our client. This poses the issue that is the topic of this blog post

## The problem

Lets say that you’re dealing with larger amounts of data, sent over a wider range connection (not our local machine like in this example). Depending on the configuration of your router or network stack, you might experience what is known as *TCP fragmentation*.

TCP fragmentation is where our network automatically splits the data packets into smaller sizes. This is done so the packets can pass through a link with a smaller maximum transmission unit (MTU) than the original packet size.

[IP fragmentation - Wikipedia](https://en.wikipedia.org/wiki/IP_fragmentation)

Lets say that happens with our example above. To simulate, I’ll send two messages. One shortly after the other is sent.

```js
const net = require("net")

const server = new net.Server()
server.listen({ host: "127.0.0.1", port: 9999 })
server.on("connection", client => {
  client.write("Hello")
  setTimeout(() => {
    client.write(" Isaac\n")
  }, 100)
})

const client = new net.Socket()
client.connect(9999, "127.0.0.1")

let received = ""
client.on("data", data => {
  received += data
  console.log(received)
})

client.on("close", () => {
  console.log("connection closed")
})
```

Running this example produces the following output.

```bash
Hello
Hello Isaac
```

Our data handler is called twice. Once because we sent `Hello`, and again because we sent `Isaac\n`. Since we are concatenating the string, the second log statement prints the full message`Hello Isaac`, however the first just prints `Hello`.

This might seem fine, but consider the case where we are receiving some JSON data and trying to parse it but it comes in separately like the example above did.

```js
...
client.write('{"name": "isaac",')
setTimeout(() =>{
  client.write('"age": "28"}\n')
}, 100)
...

let received = ""
client.on('data', (data) => {
  received += data
  console.log(JSON.parse(received))
})
```

When we run the code with these changes, it throws an error.

```bash
undefined:1
{"name": "isaac",

SyntaxError: Unexpected end of JSON input
    at JSON.parse (<anonymous>)
    at Socket.client.on (/mnt/c/Users/irowell/Dropbox/blog/gists/simple_server.js:18:20)
    at Socket.emit (events.js:189:13)
    at addChunk (_stream_readable.js:284:12)
    at readableAddChunk (_stream_readable.js:265:11)
    at Socket.Readable.push (_stream_readable.js:220:10)
    at TCP.onStreamRead [as onread] (internal/stream_base_commons.js:94:17)
```

We attempt to parse the first packet before the message is complete. So our program throws an error. This is because our `.on('data')` handler is called twice. Once for the first half of the messge, and again for the second half.

Just for completeness, this is how the messages are coming in on the wire

```bash
10:59:13.230333 IP localhost.9999 > localhost.53152: Flags [P.], seq 1:18, ack 1, win 342, options [nop,nop,TS val 4120608138 ecr 4120608137], length 17
	0x0000:  4500 0045 c9e8 4000 4006 72c8 7f00 0001  E..E..@.@.r.....
	0x0010:  7f00 0001 270f cfa0 108c 06b8 ce94 a913  ....'...........
	0x0020:  8018 0156 fe39 0000 0101 080a f59b 7d8a  ...V.9........}.
	0x0030:  f59b 7d89 7b22 6e61 6d65 223a 2022 6973  ..}.{"name":."is
	0x0040:  6161 6322 2c                             aac",
10:59:13.331214 IP localhost.9999 > localhost.53152: Flags [P.], seq 18:31, ack 1, win 342, options [nop,nop,TS val 4120608239 ecr 4120608138], length 13
	0x0000:  4500 0041 c9e9 4000 4006 72cb 7f00 0001  E..A..@.@.r.....
	0x0010:  7f00 0001 270f cfa0 108c 06c9 ce94 a913  ....'...........
	0x0020:  8018 0156 fe35 0000 0101 080a f59b 7def  ...V.5........}.
	0x0030:  f59b 7d8a 2261 6765 223a 2022 3238 227d  ..}."age":."28"}
	0x0040:  0a                                       .
10:59:13.331229 IP localhost.53152 > localhost.9999: Flags [.], ack 31, win 342, options [nop,nop,TS val 4120608239 ecr 4120608239], length 0
	0x0000:  4500 0034 6645 4000 4006 d67c 7f00 0001  E..4fE@.@..|....
	0x0010:  7f00 0001 cfa0 270f ce94 a913 108c 06d6  ......'.........
	0x0020:  8010 0156 fe28 0000 0101 080a f59b 7def  ...V.(........}.
	0x0030:  f59b 7def                                ..}.
```

## The solution

There are a couple of solutions that you may find online for this issue, but I think only one is able to scale well. Lets explore them.

### Split the incoming data

A common solution to this problem is to accumulate our `received` data string, and try to split it when new data comes in. If the split is successful, we know that the message is complete. If not, then we know to expect more data.

> Note that a split occurred if the length of the resulting list is > 1, and that the last character of the string was our delimiter if the last index of the array is an empty string.

For example lets update our data handler. We will also need clear out the received buffer after we have successfully parsed a message.

```js
let received = ""
client.on("data", data => {
  received += data
  const messages = received.split("\n") // remember this
  if (messages.length > 1) {
    for (let message of messages) {
      if (message !== "") {
        console.log(JSON.parse(message))
        received = ""
      }
    }
  }
})
```

So we check if the array length is greater than one, if it is then loop the array and parse the messages skipping the empty string if it exists. We then empty our buffer for the next messages.

And now we get

```bash
{ name: 'isaac', age: '28' }
```

Perfect. However if we send another incomplete message, this time along side our completed messsage

```js
client.write('{"name": "isaac",')
setTimeout(() => {
  client.write('"age": "28"}\n{"name": "steve",')
}, 100)
```

We are back to the same issue

```bash
{ name: 'isaac', age: '28' }
undefined:1
{"name": "steve",

SyntaxError: Unexpected end of JSON input
    at JSON.parse (<anonymous>)
    at Socket.client.on (/home/isaac/Dropbox/blog/gists/simple_server.js:23:26)
    at Socket.emit (events.js:189:13)
    at addChunk (_stream_readable.js:284:12)
    at readableAddChunk (_stream_readable.js:265:11)
    at Socket.Readable.push (_stream_readable.js:220:10)
    at TCP.onStreamRead [as onread] (internal/stream_base_commons.js:94:17)
```

This occurred because technically our `split()` function did work, but one of the strings in our array is invalid JSON.

While we could get this to work by adding additional logic to check if the last string in the array is empty (so we know that all messages are complete), there are definitely tradeoffs with efficiency. We have to split the string (which may result in nothing happening), loop the array and check if the string is empty or not before attempting to parse it, and also check if the empty string is in the last index of the array. It also has the potential to put us in a never ending cycle if we get a bunch of invalid messages, since at some point we would want to gracefully error out to the client indicating bad messages are coming from the server.

We can do better.

### The message buffer stack

What we want is to accumulate our buffer as the data comes in and check for our delimiter. If the delimiter isnt there, dont do anything. If it is, pop off the string and send it out for parsing.

If you want a good tutorial on how a stack works and how to write one, checkout [this awesome tutorial series by LucidProgramming.](https://www.youtube.com/playlist?list=PL5tcWHG-UPH32wtNecX9WPLIbHJ4-UnUs)

Lets define our `MessageBuffer` class. Its instantiated with the delimiter that we are using to distinguish between messages.

```js
class MessageBuffer {
  constructor(delimiter) {
    this.delimiter = delimiter
    this.buffer = ""
  }

  isFinished() {
    if (
      this.buffer.length === 0 ||
      this.buffer.indexOf(this.delimiter) === -1
    ) {
      return true
    }
    return false
  }

  push(data) {
    this.buffer += data
  }

  getMessage() {
    const delimiterIndex = this.buffer.indexOf(this.delimiter)
    if (delimiterIndex !== -1) {
      const message = this.buffer.slice(0, delimiterIndex)
      this.buffer = this.buffer.replace(message + this.delimiter, "")
      return message
    }
    return null
  }

  handleData() {
    /**
     * Try to accumulate the buffer with messages
     *
     * If the server isnt sending delimiters for some reason
     * then nothing will ever come back for these requests
     */
    const message = this.getMessage()
    return message
  }
}
```

A quick overview of each method.

1. `push()` concats the data to our stack (which is a string in this case)
2. `handleData()` gets the message and sends it back to the client
3. `getMessage()` checks for a delimiter, pulls out the message and returns it if it exists.
4. `isFinsihed()` just checks if we are done with our buffer for now (if its empty or there is no delimiter)

And now we can incorporate that into our program

```js
let received = new MessageBuffer("\n")
client.on("data", data => {
  received.push(data)
  while (!received.isFinished()) {
    const message = received.handleData()
    console.log(JSON.parse(message))
  }
})
```

Not only is it more readable, but theres more opportunity to use this same type of message handling in other parts of our code.

And if we run it now we get our expected output.

```bash
{ name: 'isaac', age: '28' }
```

## Conclusion

While using a stack to store our data may seem like more work up front, it allows us to create a reusable class component for the rest of our project. It can be extended to both the client and server (which has a similar data handler for messsages coming in from the client). It also lets us write less code in our main function, and makes it more readable for someone who has to work with it later.

I hope this helped you think about some of the complexities or issues that might arise when working with sockets in NodeJS. I think implementation of a buffer stack will help with scalability, as well as provide a useful tool to help learn from.

Thanks for reading!
