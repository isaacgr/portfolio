---
title: "On hoarding knowledge"
description: |
    Why proprietary knowledge and deep understandings should be well documented
    and shared with all teams
created_at: 2026-10-13
tags: software, product, documentation, rtfm
---

> It is amazing what you can accomplish if you do not care who gets the credit.

*Harry S. Truman*

## On any given day I might have to...

- Drop what I'm doing to jump on a customer emergency
- Drop what I'm doing to investigate an ongoing problem with a live system or
new deployment
- Drop what I'm doing to answer a question about the product
- Drop what I'm doing to answer a question about the behavior someone
is experiencing onsite, and if it is or is not a bug with the product
    - or if can be added as a feature!
- Drop what I'm doing to discuss an ongoing system design or
deployment, possible issues, concerns, gotchas etc.
- Drop what I'm doing to investigate a newly reported issue that will
escalate in a day or two
- Test a new feature or bugfix
- Work on gathering feedback, documenting requirements and working with our
senior developers to fit features, bugfixes or QOL changes into the backlog
- Work on documentation for a featureset, product release notes or [SOPs and MOPs](https://vitralogy.com/blog/what-is-the-difference-between-an-sop-and-a-method-of-procedure-mop/)

The point is a lot of what I'm required to do is provide information about
the product, either directly through simple 'yes we can do that, no we cant do
that and heres why' or through a bit of investigation on my part by looking at
logs, asking questions about system history, doing some in house testing and
finally providing a detailed analysis on why this is something we *can* or *cant*
do, or why its something we *should* or *should not* do.

## The myth of 'job security'

Its common when you're more junior to assume the above 
is actually a good thing. If you're the one everyone is going to, then that 
must mean you're getting pretty important. You want to be the source of truth,
the person with all the insight and know how.

![](./resources/busy.gif)

But over time that starts to take a toll. You find it harder to do any work
that helps you excel forward and grow. You find yourself spending too much time
going through the same problems with others, walking through the same issues,
same debugging steps, same architectural forks in the road. 

Worse yet, you start to resent your coworkers for becoming so dependant on you.
Maybe a lot of people would be less reliant on your input and replies if they
had access to the right tools. Maybe they might be able to own the problem, or
fish themselves for the solution if they had access to the resources and tools
you use. 

But they dont, and because of that they come to you. And they come to
you often, and you think them less competent or capable because they're coming
to you with the same problems requiring the same solutions.

Repeat until you're bitter, sarcastic, rude, short-tempered and burnt out.

Heres a direct quote from a more senior colleague in one of these positions.
This was in response to a suggestion that they take a lot of their troubleshooting
and investigation how-to docs and make them available on a small internal
platform we host for this sort of thing:

> i'll consider tagging some more of my polished private docs
>
> not all of them, need a few in the back pocket for job security

Oh dear...

![](./resources/smaug.jpeg)

## Tooling

A good subset of those questions can be answered through changelogs. Anything
that asks the question of "is this a known issue?", "is this bug seen at this site
fixed in the version we have?".

Our product is [microservice](https://en.wikipedia.org/wiki/Microservices) based as opposed 
to [monolithic](https://en.wikipedia.org/wiki/Monolithic_architecture), so
a changelog would typically either be for the individual service or for a group
of services that account for a change. I have a tool that I've built up over
many years that scrapes jira, git, confluence and testrail and provides the user
with the list of bugfixes and features between versions.

