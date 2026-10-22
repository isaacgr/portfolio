---
title: "TIL: How to use React to edit an SVG"
created_at: 2019-07-16
description: "How to integrate an SVG as React component"
tags:
  - archive
  - react
skip: true
---

## The problem

I recently came across an issue for a side project at my work.

I needed to setup a sort of sunburst graph with d3 to illustrate resource usage for remote sites. The original idea was to use something like [this](https://bl.ocks.org/kerryrodden/7090426).

![sunburst image](resources/sunburst-image.png)

It looks modern and can be interacted with to provide detailed info for each node.

The only issue is that I have no idea how to use d3. Its something I’ve wanted to get into, and have tried out using libraries such as [vx for react](https://vx-demo.now.sh/gallery), but dont have a fully formed skillset to figure out how to map the SVG, color in the right areas and load data etc.

Thats where our graphics designer came in. He made some incredible mockups for the UI that I could base my design off of, however he was too good and the mockups he did didnt have any analogs in the d3 world.

## The solution

Rather than ask that he adjust the graph to be more like something I could make, I got him to send the SVG file to me with IDs and classes that I could hook into and change using React.

Since this is a React project, I wanted to be able to edit the SVG just the same as I would any other React component. As an example lets look at the SVG markup for this heart.

![heart svg](resources/heart-svg.png)

```html
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <path
    d="M 10,30
       A 20,20 0,0,1 50,30
       A 20,20 0,0,1 90,30
       Q 90,60 50,90
       Q 10,60 10,30 z"
  />
</svg>
```

The `svg` elements wraps the `path` tag. The `d` attribute defines the path to be drawn. The `xmlns` is the namespace, and the `viewBox` defines the position and dimension of the SVG viewport.

Basically, its just HTML. So we can write it up with some JSX. It will error out if we dont make some adjustements, though. We’ll want to remove the `xmlns` attribute since these sort of tags arent recognized by React.

Typically, SVGs have inline styles as well, which should also be removed so we can style the SVG from within our CSS file.

```js
import { Component } from "react"
class Heart extends Component {
  render() {
    return (
      <svg viewBox="0 0 100 100">
        <path
          d="M 10,30
             A 20,20 0,0,1 50,30
             A 20,20 0,0,1 90,30
             Q 90,60 50,90
             Q 10,60 10,30 z"
        />
      </svg>
    )
  }
}
export default Heart
```

And thats it really. We now have a `<Heart/>` component that we can import into our react project. We can assign classes and id’s to the various elements in the SVG, and hook into them with d3 or a seperate file.

Thanks for reading!
