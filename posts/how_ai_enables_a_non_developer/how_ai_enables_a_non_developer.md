---
title: "AI won't replace your mentor, but it might be the only one available"
description: |
    How AI is helping me get over my writers block and build tools without 
    access to code reviewers, mentor feedback or stake holders
created_at: 2026-05-20
tags: software, ai
---


## I have no one else to take a stake in my code

My work tends more towards product management, "dev-ops" tooling, and product 
verification. We wear several hats at my job, so even the 
Product Manager is spinning up VMs, testing out 
features and bugfixes, taking escalation calls and working with service teams 
to feedback product requests to the developers. Though perhaps he just has a
hard time delegating. And no we're not a startup.

Anyway, over the last several years I have been developing more internal tools for
the product and deployments teams to use. These have evolved from smaller bash
or python scripts into full service packages, installable on our system's OS for
ease of use. They are essentially miniature products in their own right, with
a backend and front end, state management and auth flows. 

Maybe they have less tests than they should...

With these increasingly ambitious projects, I've also found myself trying to
learn more about system design practices, algorithms, and how to structure the
code to help make it more maintainable (even though at this time I'm the only
one that works on it). I copy patterns from our product services; we
utilize a microservice architecture, though most of the services are far from
'micro' at 14 years in. I also attempt to engage our developers when possible to
see if they can provide any feedback or guidance on my decisions. And that is
where I tend to hit a roadblock.

I'd like to be clear that in no way do I feel entitled to anyones help. Everyone
has their own priorities, and unless your work overlaps directly with someone
elses, chances are they wont be able to find the bandwidth to help you with it.

So it becomes hard to find professional input on my problems. Our team is close,
we go to lunch together, go out for drinks, message outside of work, but
the software I'm writing does not excite them. It doesnt tie into their work, 
which is more directly tied to the product as opposed being a "side" project.
So they dont tend to offer a lot of input on it. I ask them about whatever problem I'm
working through, and they do occasionally listen and give some feedback, 
but its not with any *passion*. It cant be.
They're not engaging on the core problem, trying to identify gaps, wanting to
look at the code with me to see what I'm doing or what I can do better. 

Perhaps one of the more annoying aspects of this dynamic is that I get a lot
feedback focusing more on criticizing the decisions made, rather than trying
to understand why those decisions were made and guiding in a more appropriate
direction. Maybe its due to lack of time, its easier to tell someone what they
are doing isnt how you would do it, rather than working through an alternative
solution.

Here are some quotes from a recent interaction:

> this project doesnt matter, so you can do whatever
>
> this seems complicated (with no follow up on how to improve)
>
> im speaking from experience which you dont have so just trust me

The last one being particularly harsh, since I started self-teaching myself
software and putting in the work since 2018.

## Writers Block

Maybe its a lack of confidence, or a lack of experience, but when I find myself working
hard on something with little or no feedback on whether or not the approach I'm taking
is sound, I hit writers block. It becomes impossible to figure out
what that next move should be: 

- How should I approach the problem? 
- What are the best practices when dealing with this pattern? 
- Is this going to be extensible? Maintainable? 

I understand that at a certain point you just have to build it and
make it good later but as someone trying to improve in these areas, getting a
concrete step forward, no matter how minimal, can make a big difference.

> **With AI I have a sound-board**

Another *voice* that can have a discussion around the problem and
consider my context; a stakeholder. 

Its not that I like being told "you're absolutely
right!" (in fact, I use the [caveman skill](https://github.com/JuliusBrussee/caveman)
with Claude to remove the sycophantic responses), but I have something that can
help me get to the next step in the process.

AI will:

- understand what I'm trying to achieve, and relate that what I have done in the past
- provide some back and forth on specific implementation ideas
- not assume the decisions I've made until this point are bad, or misinformed
- provide code review

This makes it a lot easier for me to work through writers block. I am no longer
getting several different opinions on what to do next. I'm not getting flack
for making decisions that seem bad at the surface, but perhaps were necessary given
the project scope. I'm able to get concrete responses on what to do next, how best
to move forward past the current hurdle, and what the possible ways are to handle the
next set of challenges. 

### Asking for help, not answers

Now to be fair, this process does require that I dont ask the AI to solve the
solution for me. That would defeat the purpose. What I often do with claude is enter
`/plan` mode, and prompt it with something like "show me some common methods to
handle this case" or "does it make sense to do x/y/z? what are the caveats to
each option?". This way it doesnt just start flailing wildly and editing the 
codebase, filling my screen with massive multi method diffs.

This also gives me a chance to challenge or question the output. "Why would I
approach it that way vs. this way?", "I think that misses this case" or "That
wont work in my codebase because...". It can then come back with some attempts
at justification, or with corrections.

And now I can spend time trying the methods 
out and doing my own research into the cost/benefit of each. 

Of course I often find contradictions to what the AI suggests, which is good. 
Those can lead me down a fun rabbit hole of stdlib docs, stack overflow posts, and 
github issues. (I am not sure why, but theres something quite rewarding about
walking through the [go docs](https://pkg.go.dev/std) and using them to piece
together my own solution).

### Case Study: Designing a Generic Graph API

I was building a Go wrapper around [Apache AGE](https://age.apache.org/age-manual/master/intro/overview.html).
I'm working on a tool that models customer data (things like package versions,
installed hardware, current and outstanding POs etc.) using a graph.
I had a `CreateVertex` method that needed to accept different vertex
types: Roles, Users, Permissions etc. each with their own properties. This is
all fairly out of my depth, and something I've been excited to dive into.

The problem was that my GraphQL models are all strongly typed, but the graph 
layer needed to be generic.

I opened Claude and asked: "What's a good way to pass arbitrary
vertex property fields to CreateVertex?"

#### The Options

It gave me three approaches:

 1. map[string]any
    - simple, pass properties as a map
 2. any with reflection
    - accept any struct, marshal to JSON, then to map
 3. Interface with `ToProperties()`
    - define a contract for vertex inputs

Its recommendation was Option 2, since I was
already using `json.Unmarshal` elsewhere. Fair enough.

#### Drilling Down

But I was unsure: "How do I build the Cypher properties string? Cypher map keys
can't be wrapped in quotes." 

This is where I have to be careful. Instead of searching through Cypher docs and AGE
source code, I got a working `buildPropsString` function complete with
string escaping, nulls, bools, and numbers. Now it did flag considerations I
hadn't thought of, like sanitizing keys from user input, escaping single quotes
and handling nested objects but it still has a bit of a compulsion to actually
provide you with a full on solution.

#### Challenging the Suggestion

I then asked: "Should ToProperties() return a map or a string?"

The answer was maps—keep the model layer unaware of Cypher syntax. That made
sense. Separation of concerns.

But I pushed back: "Is there anything in the Apache AGE driver that does
this already? Where did you get this info?"

This is important. AI will confidently present solutions that may or may not
exist in the ecosystem. It admitted it derived `buildPropsString` from Cypher
syntax rules, not from any library. The AGE Go driver handles parsing results,
not building properties. Good to know that I wasn't missing an obvious stdlib
function.

#### Arriving at the Solution

Finally, I questioned the interface pattern itself: "Is there even a point to
the interface? My [gqlgen](https://gqlgen.com/) models would all need to implement it. 
Can't I just pass the map directly?"

"yes, skip the interface. YAGNI. The repository layer already knows
the model—it can build the map inline" (I had to google what 'yagni' meant).

```go
func (r *Role) Create(
ctx context.Context, 
id string, 
input model.CreateRoleInput
) (*model.Role, error) {
     dest := &model.Role{}
     err := r.graph.CreateVertex(ctx, r.label, map[string]any{
         "id":   id,
         "name": input.Name,
     }, &dest)
 return dest, err
}
```

I know the above looks trivial, but the access patterns and abstractions it took
to get here are things that I am pretty unconfident with. Ultimately the above
method accessed via the repository pattern was excellent boilerplate for me
to move forward with elsewhere.

#### What AI Did Here

 - Gave me three concrete options instead of one "right" answer
 - Provided working code I could test immediately
 - Answered follow-ups without losing context
 - Admitted when information came from syntax rules vs. actual libraries
 - Let me challenge suggestions and course-correct

I didn't ask it to write the solution. I asked it to help me think through
the solution. The final decision, which was to skip the interface and use maps, 
was my own, informed by the tradeoffs we discussed together.

### Code review: still hit or miss

I suppose in my scenario there are two ways to go about AI code review:

1. Once you feel like you've completed the code, you can simply ask it something
like "review this code" and then push that and merge it.

2. You can create a PR, and ask for AI review. Which is effectively the same
thing as #1, but maybe a bit easier to follow, and has more
concrete closure when all of the comments have been resolved.

I tend to opt for #2, as it feels like im working towards an actual "end" 
when I walk through and resolve each comment.

However this is again where inconsistency arises. If I ask for a review more
than once, I find that it tends to contradict previous comments. "Change A
to B" before now suddenly becomes "Change B to A". So I never really do
more than one pass with it this way.

I'll try to scrutinize myself more here as well. Run through all the changes,
and add comments to provide future me with a bit more context, even if they may
be somewhat redundant. When you're the only person maintaining the code this
becomes more valuable as you grind out future changes. Maybe its different for
others, but its easier for me to recall my own intent with a block of code
than someone elses.

Consistency in your codebase will go along way here as well. If I want to ask
AI to just make a change I know is otherwise trivial, it has a very consistent
and well documented method to approach it with. I find this helps to reduce
hallucinations, though recently after asking for some changes I saw that it
had decided to implement its own HTTP client, without any of the security checks
I had caked in, simply because I didnt have an existing method to return the
raw response object.

## A workflow for solo devs

Of course, this is no replacement for actual human engagement. When I can get
input from a colleague its often insightful, well thought out and
coming from years of experience in their craft. I like being grilled
about my decisions and forced to consider the alternatives. These are things
AI is generally poor with, since it leans towards assuming the user is
correct in some way and knows better.

Examples from experienced developers are often much more insightful as well. I
find that I have to spend significant time re-reading and researching an AI
solution to have a better understanding of the pitfalls. With more carefully
considered code commentary I dont have to do that as much, since the behavior is much more
clear.

Most AI coding is framed around replacing your own time. "Write a thing that
does X", "Review this code and apply all recommendations", "Re-write this in rust".
Most AI coding advice is about generating code faster. AI is more valuable as a 
collaborator than a typist.

### tl;dr

1. Use `/plan mode`
2. Ask questions, challenge the output
3. RTFM
4. Limit review passes
5. Do your own research before implementing a solution

I have AI help me work through code because I have no one who is interested in it. 
But I'm sure, especially now, a
lot of people are finding themselves in this position. And I'm also sure there
arent as many who get to enjoy the journey, and are just in it for
the destination. 

We love crafting code and finally have a voice around that 
will take an interest, provide feedback, and guide toward a solution we feel
is complete. And what better way to get past the idea stage 
and start building things that excite you.
