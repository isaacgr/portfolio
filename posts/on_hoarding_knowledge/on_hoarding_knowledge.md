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

A good subset of those questions can be answered through changelogs. Anything
that asks the question of "is this a known issue?", "is this bug seen at this site
fixed in the version we have?".

Our product is [microservice](https://en.wikipedia.org/wiki/Microservices) based as opposed 
to [monolithic](https://en.wikipedia.org/wiki/Monolithic_architecture), so
a changelog would typically either be for the individual service or for a group
of services that account for a change. I have a tool that I've built up over
many years that scrapes jira, git, confluence and testrail and provides the user
with the list of bugfixes and features between versions.

