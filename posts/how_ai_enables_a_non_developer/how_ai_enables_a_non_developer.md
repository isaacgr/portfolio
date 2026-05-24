---
title: "AI: Artificial Soundboard"
description: |
    How AI is helping me get over my writers block and build tools without 
    access to code reviewers, mentor feedback or stake holders
created_at: 2026-05-20
tags: software, ai, product management
---

My work tends more towards product management, "dev-ops" tooling, and product 
verification. We tend to wear several hats at my job, so even the 
Product Manager is spinning up VMs, testing out 
features and bugfixes, taking escalation calls and working with service teams 
to feedback product requests to the developers. Though perhaps he just has a
hard time delegating...

Anyway, over the last several years I have been developing more internal tools for
the product and deployments teams to use. These have evolved from smaller bash
or python scripts into full service packages, installable on our system's OS for
ease of use. They are essentially miniature products in their own right, with
a backend and front end, state management and auth flows. 

With these increasingly ambitious projects, I've also found myself trying to
learn more about system design practices, algorithms and how to structure the
code to help make it more maintainable (even though at this time I'm the only
one that works on it). I copy patterns from our product services (we
utilize a microservice architecture, though most of the services are far from
'micro' at 14 years in) and attempt to engage our developers when possible to
see if they can provide any feedback or guidance on my decisions. And that is
where I tend to hit a roadblock.

I'd like to be clear that in no way do I feel entitled to anyones help. Everyone
has their own priorities, and unless your work overlaps directly with someone
elses chances are they wont be able to find the bandwidth to help you with it.

So it becomes hard to find professional input on my problems. Our team is close,
we go to lunch together, go out for drinks, communicate outside of work, but
the software I'm writing does not excite them. It doesnt tie into their work, 
which is more directly tied to the product as opposed being a "side" project.
So they dont offer a lot of input on it. I ask them about whatever problem I'm
working through, and they do give some input. But its not with any *passion*. 
They're not engaging on the core problem, trying to identify gaps, wanting to
look at the code with me to see what I'm doing or what I can do better. 

Perhaps one of the more annoying aspects of this dynamic is that I get a lot
feedback focusing more on criticizing the decisions made, rather than trying
to understand why those decisions were made in the first place.

Here are a couple of quotes from a recent interaction:

    > it seems like too much
    > this seems complicated
    > like its your first or second ever golang project [as a reason why my
    potentially complex workflow was a bad idea]
    > im speaking from experience which you dont have so just trust me


The last one being particularly harsh, since I started self-teaching myself
software in 2018.

This is where AI comes in. Its not that I like being told "you're absolutely
right!" (in fact, I use the [https://github.com/JuliusBrussee/caveman](caveman skill)
with Claude to keep things professional), but with AI I have a
sound-board. Another *thing* that can have a discussion around the problem and
consider my context; a stakeholder. AI will:

- understand what I'm trying to achive, and relate that what I have done in the past
- provide some back and forth on specific implementation ideas
- not assume the decisions I've made until this point are bad, or misinformed
    - however will tell me if what I'm doing is not the best approach
- provide code review

This makes it a lot easier for me to work through writers block. I'm now no longer
getting several different opinions on what to do next. I'm not getting flack
for making decisions that seem bad at the surface, but perhaps were necessary given
the project scope. I'm able to get concrete responses on what to do next, how best
to move forward past the current hurdle, and what the possible ways are to handle the
next set of challenges. 


