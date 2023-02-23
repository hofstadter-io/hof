---
title: plasmic
type: page
layout: text
draft: true
---


# Plasmic - drag-n-drop that developers can love?

https://www.plasmic.app/ | https://docs.plasmic.app/learn/ | https://github.com/plasmicapp/plasmic

## My Introduction To Plasmic

Saw on [Learning with Jason](https://www.twitch.tv/jlengstorf)

What drew me in

- drag-n-drop with two-way code
- support across many frontend technologies
- they have some amount of open-sourceness (not the DnD UI, understandably)

The main thing I want to do:

- design components and sections for my landing page
- store them in git, have them be part of Hugo build
- figure out a workflow that is going to ... work

We are currently trying to migrate from Webflow to Hugo.
I'm the main editor, I will do more if it's in Hugo.
Plus, Webflow sucks to work with and customize.

- https://hofstadter.io (webflow)
- show WIP site in hugo


## A Bit of Background

We moved the site to Webflow a while ago,
because business team wants nothing to do 
with git, code, or hugo...
but they need to update the website...

Webflow sucks though, and we want something
that is more in the spirit of hof, High Code style.
So drag-n-drop, edit in place for the busienss team,
developers can still user their tools.
Hugo is fast af, highly customizable,
and all of our other sites use it.

### Is this a drag-n-drop tool developers can love?

Ok, I've looked just a bit and slacked with their team a bit more.
Enough to ...

- create an account
- understand how the Plasmic output can move around, and make sure that I can get it into vcs
- to pick out some starting pages
  - https://docs.plasmic.app/learn/rest-quickstart/
	- https://docs.plasmic.app/learn/js-quickstart/
- sounds like we can leverage our bootstrap & font-awesome, but not really sure if / how, no starting places


## Questions Beforehand

It's been almost a week since Plasmic hit my radar.
In that time, I've been brainstorming what I want to
figure out and how I might use or integrate with Plasmic.

- how do I use this with hugo?
	- probably want to eject components as partial templates
	- short-codes / partials for
	- hugo data integration

How do we have a nice preview while designing in Plasmic?

- two-way integration
	- how does this integrate with git? Is this something we manage outside of the platform?
	- make rule to fetch all components, subrules per, can we determine which have changed without fetching?
	- how do we make coded components available in Plasmic
	- how does css/scss work?

- brand assets & style guide, styling generally
	- how do we work with our brand assets and style guide?
	- how do we setup our css, sass, bootstrap, font-awesome
	- Is there a brand guideline template in Plasmic?
	- Can we build one and create a source of truth repo?
	- Design resuable data components in Plasmic, store in repo for import
	- How is the expience of getting styling right compare between code, webflow, and plasmic?

- hosting / serving / dev mode / hot reload
  - definitely want to serve ourselves
	  - site-speed, don't want a downstream dep for page content
		- integrated before site is built with hugo/vuepress
	- switch for dynamic dev mode based on hugo mode / config?
	  - dynamic uses the remote fetch, prod mode uses the fetched version in the source tree
		- can we get hot reload working with remote components?
	- do they have webooks / api?

- how does this relate to _board?
  - does it have an infinite canvas?
	- how does it compare to grapes-js?
	- touch support?
	- can we white label this in?

- integrations
  - what languages do they have examples for (server side)?
	- besides figma, are there other tools our users would integrate with?


Is this drag-n-drop developers can love?


# Taking Plasmic for a Spin

## Getting started plan

- [x] cruise the homepage
- [x] cruise the docs
- try the starting points
  - [x] how to use plasmic (tutorial)
	- developer getting started
- integrate into hugo
- figure out two-way workflow
- try data components
- figure out css / branding
- docs in depth


## The Plasmic Websites

### Homepage

### Docs

### GitHub


## The Plasmic App

### Tutorial

### Developer


## Meeting our goals

```sh
# install
yarn global add @plasmic-app/cli
```

### Integration into existing Hugo site

Making a testimonials section & cards...

Steps:

```sh

```


### Two-way Binding

### Hugo Data and Plasmic Components

### Extras / Other


# Takeways

## Is this drag-n-drop developers can love?

## The Good

## The Bad

## Issues

## Feedback

CLI is slow...

## Questions

- is there a way to codegen anything besides React?
- can we `plasmic sync` Vue or non-framework code?
- how to match css classnames between Plasmic & our existing codebase
- how to minimize css? Is this something you expect a JS bundler to handle?
- lot's of css duplication, excessive html elem spec, how to reduce classes
- how to group {mixins} X {html elems} and give a name

## _board

we need full control over code gen in our _board.

hofmod-[site?] as something equivalent to hugo

hof create should autodiscover context
