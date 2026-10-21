---
title: "Git: Editing commits"
created_at: 2019-07-15
description: "Quick ways to edit commit messages locally or on the remote repo"
tags:
    - archive
    - git
skip: true
---

There are a few cases where you’ll want to edit your last commit.

> 1. You need to change the message, because of a typo or because it was a particularily frustrating day and you ended up being less than professional in it
> 2. You forgot to add a file to that commit before pushing it
> 3. You pushed the commit with the bad message and need to reword it

---

  

For **#1**, you can simply run

```bash
$ git commit --amend
```

This will open an editor for the last commit and allow you to make changes to the message.

  

For **#2** its similar but we first want to add the file that we missed.

```bash
$ git add missed_file
$ git commit --amend
```

This will add the file to the previous commit, and allow you to edit the message if necessary.

  

For **#3** we need to use `git rebase`.

We want to use `git rebase` on the commit hash that we want to edit. You can get the hash quickly from `git log`.

```bash
$ git rebase -i <commit hash you want to change>^
```

This will open your default editor with a list of commits and actions for each one. By default, the action is `pick`. For any commit you wish to change the message, change pick to `reword` then save and quit. Dont forget the caret `^` at the end of the command.

> The caret indicates we want the first parent of that commit object. Read more about it [here](https://stackoverflow.com/questions/2221658/whats-the-difference-between-head-and-head-in-git).

For each commit you’ve chosen to reword, you’ll get an editor to edit the commit message. Change it then save and quit.

Once you’re done editing all the commit messages, you’ll return to the command prompt and have a new tree with the updated messages.

You can now upload them to your repository by running

```bash
$ git push origin --force
```

Thanks for reading!
