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

---

A lot of my daily interruptions are questions someone could answer for
themselves if they had the means to.

On any given day I might have to

- Drop what I'm doing to jump on a customer emergency
- Drop what I'm doing to investigate an ongoing problem with a live system or
new deployment
- Drop what I'm doing to answer a question about the product
- Drop what I'm doing to answer a question about the behavior someone
is experiencing onsite, and if it is or is not a bug with the product
    - or if it could be added as a feature!

## The myth of 'job security'

Its common when you're more junior to assume the above 
is actually a good thing. If you're the one everyone is going to, then that 
must mean you're getting pretty important. You want to be the source of truth,
the person with all the insight and know how.

But over time that starts to take a toll. You find it harder to do any work
that helps you excel forward and grow. You find yourself spending too much time
going through the same problems with others, walking through the same issues,
same debugging steps, same architectural forks in the road. 

Worse yet, you start to resent your coworkers for becoming so dependent on you.
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

Thats a great way to stay where you are I suppose. Openly making yourself
the gatekeeper of how to execute a particular process, or the only one with the
know-how on how to troubleshoot or investigate an issue. 

Its pretty obvious that having the mentality of knowledge hoarding for some
perceived employment related gain will eventually lead to or accelerate your
burn-out.

And it doesnt matter at the end of the day. The company doesnt actually consider
these sort of things if they want to let you go. And if you quit, a good
manager should be able to cover your position adequately. And it usually only takes
a few months or so for someone to get up to speed enough. 

---

My mentor was actually very good at this sort of documentation and knowledge
sharing, but he did *so much*
that at the time I recall my boss was quite stressed about the whole thing, unsure
of what to actually do about it (there was a brief period where he would have to
take on the bulk of the work). But
within a couple of months the workload started to move off. Sure, they had to 
offload it to about 3 people but that soon became the norm, and they trained
more people to offload that work and now that position he filled is effectively
just 3 or more different jobs delegated to multiple different teams.

So if he had eagerly hoarded knowledge for 'job security' what good
would that have done? It would have just been making things harder for himself
and others while he was with us, and it would have gotten disseminated (or forgotten)
after he left anyway.

---

## Tooling

A good subset of the questions I deal with can be answered through changelogs. Anything
that asks the question of "is this a known issue?", "is this bug seen at this site
fixed in the version we have?".

Our product is [microservice](https://en.wikipedia.org/wiki/Microservices) based as opposed 
to [monolithic](https://en.wikipedia.org/wiki/Monolithic_architecture), so
a changelog would typically either be for the individual service or for a group
of services that account for a change. I have a tool that I've built up over
many years that scrapes jira, git, confluence and testrail and provides the user
with the list of bugfixes and features between versions.

This tool has been made public within the organization, so anyone with an AD
account can access it. This is the exact tool that I use myself to answer
a lot of these questions, or at least provide a basis for further investigation,
and I've tried to train others to use it as well. Its saved me and others loads
of time since they are no longer waiting on me for an answer and I'm no longer
required to distract myself to answer it for them. Now there is a bit of verification
required, and I do still want people to check with me to make sure things are
correct, but its much better doing it from a point of 'Hey Im seeing this issue onsite, 
it sort of looks like this which is already fixed in this version. Do
you agree?'

## Documentation

Internally we use jira and confluence for this. Confluence is great because its
easy to build out long form docs, TOCs, add rich text and images etc and link to
jira for issue reference. Subjects I use confluence for include troubleshooting
guides, documentation on new features or product manuals. 

A secondary benefit is that its a way for you to work on your reading and writing
skills.

Jira is also nice since it allows you to segment pages into tickets, which come with
a consise key that you can use to point to (ex. MAGSUP-188 links to a confluence page
for NMOS troubleshooting). I do use it primarily for things like release notes 
where each point release gets its own ticket with the description of changes.

## Lunch and Learns

Certainly harder to setup, and not everyones company is going to be able to 
accomodate these. But preparing a 1-2 hour presentation on a topic that people
generally reach out about is definitely more engaging. Discussing Q&A in person
is always easier, and the group setting allows others to hear answers to
questions they might not have thought about in the moment. 

Also, free lunch...

## AI

Yes, companies are using AI tools to ingest internal data and make it
easier for people to get these answers themselves. Our company uses [Glean](https://www.glean.com/)
for this. The worries here are not the same as they are externally, where developers,
artists, musicians etc. are concerned with their work being stolen. In this
case it actually works in our favor. If the documentation is well maintained and
well written, then it helps people ingest it in a more efficient way. Of course,
the AI can hallucinate or not provide adequate context given the domain specific
knowledge that might be required, but it still helps with getting that workload
off of your plate so you can focus on the more important elements of your work.

---

No one will really care at the end of the day who wrote any of the guides or notes
that I've built up. But their effects will be felt, through the autonomy and
action of the people who need the knowledge to do their job. They can be more effective
and will eventually outgrow the need for guides through osmosis. They might even
learn enough to contribute themselves, continuing the cycle for the next group.
That frees up my time, and it ultimately frees up their time as well. Plus theres
a certain satisfaction in knowing you actually *helped teach* someone as opposed
to just telling them what to do.
