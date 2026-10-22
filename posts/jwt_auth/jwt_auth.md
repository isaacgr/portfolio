---
title: "How to use JSON Web Tokens for authentication"
created_at: 2019-06-08
description: "Setup JWT for authentication in your web application"
tags:
    - jwt
    - tutorial
    - javascript
    - archive
---

An all too often requirement when building a website is that it include authentication. A user logs in with their name or email and a password, which grants them access to different pages on the website. After some time, or if the user logs out, they are redirected to the login page once again.

I recently needed to do this for a home dashboard page that I built to aggregate some IOT data. I wanted only authorized users to have access to the site.

There are several different ways to acheive this. I could have used firebase with google authentication, or passport, but as often happens I needed a simple solution to a simple problem.

I decided to go with `JSON Web Tokens` (JWT) to authenticate users. JWTs are credentials, which can be used to grant access to resources. Check out more about them [on their webpage](https://jwt.io/). JWT are language agnostic, since they are just JSON. Meaning you can use them with JavaScript, Python, Java, C etc.

The overall concept is pretty simple as I explain below. I will show how to setup the server side to generate and return a token, and how routes can be restricted there. Then I’ll move on to how you can authorize routes within your webapp.

## JWT authentication flow

The flow for using JWT authentication is something like this:

1. Client sends a login request
2. Server checks that the client exists through some means (Im using mongodb in this case, so I check that the user exists and that the password is valid according to the database)
3. If login is valid, then a token is generated for that client and returned
4. Client uses that token when requesting access to different pages of the site

Optionally, the expiry on the token can be set very high or very low, so that after some time the token is no longer valid and the user will have to login again.

## Server side

The server will be waiting for a login request. Once all credentials are confirmed, the token can be generated and returned to the client. The boilerplate for the server side would look something like this. Note that I’m using an express server to handle the post requests.

```js
const jwt = require("jsonwebtoken")

// handle some type of login request
app.post("/api/login", (request, response) => {
  // check that user exists, otherwise return an error
  // check that password for user is valid, otherwise return an error
  // if all good, generate token and return it in the response
  jwt.sign(
    { user },
    // optional secret/private key
    // client must authenticate using this key or the token wil be invalid
    // offers extra layer of security
    process.env.SECRET,
    { expiresIn: "15m" },
    (error, token) => {
      if (error) {
        return response.status(400).json(error["message"])
      }
      return response.json({ token })
    }
  )
})
```

## Client side

The client makes a request to the server. In this case I’m making a `POST` request to `/api/login`.

If the post is successful the client will receive something like the following in the response message.

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7Il9pZCI6IjVjZjRmOWZlNWRlN2I4MzMyNDlmMGU2YiIsInVzZXJuYW1lIjoiaXNhYWMiLCJlbWFpbCI6ImlzYWFjQGVtYWlsLmNvbSIsInNhbHQiOiJiM2U0NmIzNmEyMTFjZTE4ZjFmNDNlYzcwMWFiNzA3NSIsImhhc2giOiI2MTgwMTA2OTdiZjYyZWRhMDgwNWEyMjZhZjVmYmJiZjBkOGQyYjNmYzRlNmYxNGNmOTlkNGM3ZTdhMjZmMTczYTlkODNiOTk3NDJlZjA5NzBlMjg4Mjc0MDcyMGE0OTMwYWIwMjk4ZmE5ZDIxYjM0MmQzMjJkNWI3NjQ4YmYwMmRjNWJlYjdiOTU5NjQzNGMyZjI2ZmM3YTFlMGJmMWEwOGNlYTNkNzk5OTQ2ODMzNjVmYzhiM2Y2NWEwNWFkNzcyNzcyZmIzZjNlZDUyNjJlMGUwMTk1MGI5NjUwOWY0YWFmZmE5MzI0OTlkMWYxYjQ0Mjc2NDc2Yjg3ZjVjYzA3NmQ4M2JkMTg1NzU0NmM2MjRlZWRkYTBkZTNjZjc3MjY4OTgzZjIzNDQ1MGM5MTE1MGI5OGQ0NTg3MDQxYmQzMzc4Mjc4YzI5Y2YzY2Q4N2MxY2U0NThjZjI5ZmMwNmViNGEyNDAzZGYzY2MwZDhjMGY4Y2MxOTlmMjM2OGJjNGY0OWE2NmIxMDcxMDIxYzM3OTAwZmYxNTliYTI0Y2M1OWNkMWUwMzMwN2QzYTVmODNiZTliMWQwNjU2M2FkOWYyOTVlM2FhYjdmMTkxMGJmNGJiZDk3YmE4YmNjMjlmYzU4YjI1ZDA1MTVkZGVkMTU0ZTNiODFkNmUwNTY0YmVkZDBiYzJjMTkxYjgwZjliYWYwYjJlN2FlOTI2YjI0ZmFhMDUyMjBiNjdkZDBlY2NjZWVjMDM4ODE5MDg2YzU1NjhjNjFjNWE1ODQ0MWFiNjhhOWU1MjNjMDJlNDQ4NTMzYmQ1N2VlODdmYTBlZjM4OWQ3YjA3MGQwZmQ0ODNlYTdmOTFiZjcxZDMzOTgxNjM0YmY2MTQ0OTk3YWFhNmYzZjZhNzQ5MjM4NDg0OTMwZTM1NzhiNGJkYTY4NDhjMTFhNmY1Y2MyNjIxOWNmZjkzOWQyNDVjZjZlYzc4MWYzYWNkOGVkYTRkOGEwN2M2N2ZkNDhiM2VkZDA2NTQwNTYwYWVkNzMzNjEzMzBjYWM5NzdiMzNlOTBkYWQzNjRhYzEyYmVkNjllMDhhYmI0MDE0ZTExYWYwMzczNTg0YTQ2ZGFiNjYxOGM0YWU0ZmM2ZTFhODBkMzU1MTlmYmY4ZDRiZTgyNmMwOTg4NWEyOTU0ZjRiZjZkNzQ0MDQxMzI5N2NkMTY5NWNmM2YxNGNlOTlkYjBiYmU2YTZlYjcxNmZhOGJlYWQ2ODU4ZWNhNWVjMGYxMTI1NWVlZDhlNzc1NzA5ZDE1MmY4OWQ3ODM3MDUxMjYxMmI5YTBmMGI1NzUzYmIyNzM4ZDJiYjViMjZjODgwNGNlOWFhMDQwZjA2ODQ0NmI4IiwiY3JlYXRlZEF0Ijoi.wcuPsmHIOeJCeGUhOkZ_xEBAQ1qPS4e0kTKAozyYy4s"
}
```

We could actually pass this very long string along to the [JWT Decoder](https://jwt.io/) and get the decoded payload. We could see info like the username, email, and password (I’ve deleted some sections of it for security purposes :p).

## Server Setup

I’ll run through a quick setup of a server and client, to show how we can use the token to restrict access to the client.

If you dont already, you’ll want to create a folder and install some packages, namely `jsonwebtoken`, `express` and `body-parser` for express, which is a middleware to parse incoming request bodies and add them to the `request.body` property of your app.

```bash
$ mkdir jwt-demo
$ cd jwt-demo/
$ npm init
$ npm install express body-parser jsonwebtoken
```

Create a file to host our server, I’m going to call mine `server.js`. This is where we’ll setup the routes for the server.

Lets add some boilerplate code to get our express server running:

```js
const express = require("express")
const bodyParser = require("body-parser")
const app = express()
const port = 3000

// parse application/json requests, we'll be making these to login
app.use(bodyParser.json())

app.get("/", (req, res) => res.send("Homepage"))

app.listen(port, () => console.log(`Server listening on port ${port}!`))
```

We can start our server and confirm its up and running

```bash
$ node server.js
> Example app listening on port 3000!
```

So now that our app is running, we need to tweak it to have a route that users can use to login. We then want to test that the token restricts access to certain routes.

For the purposes of demonstration, we will just hardcode a username and password to accept. The default expiration time is 24 hours, and as shown above you could change this if required.

```js
...
// require the jwt library
const jwt = require('jsonwebtoken')
...

app.post('/api/login', (request, response) =>{
  // get the username and password from the request body
  const {username, password} = request.body
  // verify it somehow
  if (username === "isaac" && password==="jwt"){
    // user is authorized, send them a token
    jwt.sign(
      {username},
      "secretkey",
      (error, token) => {
        if (error){
          // return the error message to the user
          return response.status(400).json(error['message'])
        }
        // all good, return the token to the user
        return response.json({token})
      }
    )
  } else {
    // user is unauthorized, send error
    return response.status(404).send('not authorized')
  }
})

...
```

The login request will be a `POST` to `/api/login` with the username and password in the body. I’m using [Postman](https://www.getpostman.com/) to make the requests. If you havent heard of it, I highly recommend checking it out. Its a useful tool that I use for all my deployments to help test APIs.

POST the username and password in the body to `http://localhost:3000/api/login`

```json
{
  "username": "isaac",
  "password": "jwt"
}
```

Now just as above, you’ll get a token in a response

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImlzYWFjIiwiaWF0IjoxNTU5NzQ4Njk1fQ.qowwebCF0wkKhCY1SSpOInmTrXUKeGC0nhEsrTLbn7U"
}
```

## Authorization with the token

Alright, so we now have a token. We can send this with our requests to get authorization to different routes. We want to send this token in our headers, using the `Bearer Token` authorization header. We just want to add the token string into the header in our request.

The header looks something like this

```js
Authorization: Bearer <token>
```

So when we get the request, the `authorization` header is going to be a string with `Bearer <token>` in it. We will need to strip out the token from the string, and then verify it.

Im going to setup a `verifyToken` middleware for our route. So express is going to run this first before moving on to the request callback. If the token exists in the request then we can do something with it, if not we can reject the request from the user.

```js
// verify the token
app.get("/api/verify", verifyToken, (request, response) => {
  // we have a token, make sure its valid
  jwt.verify(request.token, "secretkey", (error, authData) => {
    // something invalid, return
    if (error) {
      return response.status(403).json({
        message: "unauthorized",
      })
    } else {
      // token is valid, return an object with the authData
      return response.json({
        authorized: true,
        authData,
      })
    }
  })
})

// verifyToken middleware to make sure token exists
function verifyToken(request, response, next) {
  // get authorization header from request
  const bearerHeader = request.headers["authorization"]
  // if it exists, then get the token
  if (typeof bearerHeader !== "undefined") {
    const token = bearerHeader.split(" ")[1]
    // add the token to the request object
    request.token = token
    // move on with the calling function
    next()
  } else {
    // no token, reject
    return response.status(403).json({
      message: "unauthorized",
    })
  }
}
```

If we make a `GET` request to `/api/verify` with the token we got earlier in our header, we should get a response saying that we are authorized.

```json
{
  "authorized": true,
  "authData": {
    "username": "isaac",
    "iat": 1559748695
  }
}
```

Note that the auth data contains the data we used to sign our token earlier. Since we only included the username thats all it returns. **iat** means *issued at* and provides a means to know the date and time a token was given to a user.

We could now use this verification technique to restrict access to certain routes or pages within our application. I’ve showed how this is done using a server and a remote client sending external requests, lets take a look at how it can be done from a webclient.

## Webclient Setup

Working with webclients is a little more involved. We need to authorize the user, store their token, and then check that the token is valid for each path they go to. If it isnt, we need to force them to log back to get a new token by signing in again.

I’m not going to go through all of the steps to setup routes, and pages etc. for a web app. I’m just going to focus on how you can authorize a user, store the token, and then check that tokens validity when a user requests a page.

### The Authorization Handler

The first thing we are going to want to setup is an authorization handler. This will be in charge of all of the above operations.

We will need to install `jwt-decode` to decode the json web token from our webclient.

```bash
$ npm install jwt-decode
```

Then, the boilerplate for the auth handler would look something like this. We are just creating a class with some methods to call to handle logging in, logging out, getting the token, removing the token and validating the token.

```js
import decode from "jwt-decode"
import jwt from "jsonwebtoken"

class AuthHandlder {
  // Get a token from the server
  login(username, password) {
    const data = {
      username,
      password,
    }
    // post request with the login info to the server
    return this.fetch(`/api/login`, {
      method: "POST",
      body: JSON.stringify({
        username,
        password,
      }),
    })
      .then(res => {
        // we will set the token in local storage
        // this is storage persistent until explicitly deleted
        this.setToken(res.token)
        return Promise.resolve(res)
      })
      .catch(error => {
        return Promise.reject(error)
      })
  }
  // Checks if there is a saved token in localstorage and if it's still valid
  loggedIn() {
    // get the token from localstorage
    const token = this.getToken()
    return !!token && this.isTokenValid(token)
  }
  // verify the token is valid
  isTokenValid(token) {
    try {
      return jwt.verify(token, "secretkey", (error, auth) => {
        if (error) {
          return false
        } else {
          return true
        }
      })
    } catch (err) {
      return false
    }
  }
  // saves the users token to localStorage
  setToken(idToken) {
    localStorage.setItem("id_token", idToken)
  }
  // gets the users token from localStorage
  getToken() {
    return localStorage.getItem("id_token")
  }
  // clear the users token and data from localStorage
  logout() {
    localStorage.removeItem("id_token")
  }

  getProfile() {
    // decode the token to get the user info, can be useful
    return decode(this.getToken())
  }

  fetch(url, options) {
    // performs api calls sending the required authentication headers
    const headers = {
      Accept: "application/json",
      "Content-Type": "application/json",
    }

    return fetch(url, {
      headers,
      ...options,
    })
      .then(this._checkStatus)
      .then(response => response.json())
  }

  _checkStatus(response) {
    // raises an error in case response status is not a success
    if (response.status >= 200 && response.status < 300) {
      // Success status lies between 200 to 300
      return response
    } else {
      let message = ""
      return response.json().then(json => {
        message = json["message"]
        var error = new Error(message)
        error.response = response
        throw error
      })
    }
  }
}

// create an instance of the handler and export it to be used by the rest of our application
let Auth = new AuthHandler()

export default Auth
```

So now elsewhere in our application we can check that the user is logged in and allow or restrict access based on that. For example, if a user tries to access a certain page we can redirect them to home if they are not logged in.

```js
import Auth from './auth.js'

...

if(!Auth.loggedIn()){
  // redirect to homepage
} else {
  // do something
}

...
```

It’s actually that easy. I’m using it with `react-router`. So I render all pages through private route components except the login page. If a user is not authenticated I just redirect them.

```js
import React from "react"
import { Route, Redirect } from "react-router-dom"
import Auth from "../functions/AuthHandler"

const PrivateRoute = ({ component: Component, ...rest }) => {
  return (
    <Route
      {...rest}
      component={() => (Auth.loggedIn() ? <Component /> : <Redirect to="/" />)}    />
  )
}

export default PrivateRoute
```

And there we go. Hopefully you’re able to integrate authorization into your app using JWT. They really simplified the setup since it was just a personal page and I didnt need anything crazy.

Thanks for reading!
