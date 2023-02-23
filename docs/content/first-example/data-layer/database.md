---
title: Database
brief: upgrading the storage
weight: 100
---

Our initial data store was a Go map and overly simplistic.
In this section, we change this to use a relational database.
The change is relatively limited given the nature of the work.
We need to

- update the schema for the database
- add a key and timestamp fields to our types
- update our type library to use SQL
- connect to the database on server startup
- generate a script for starting the database

all of our custom code can stay the same
and our seeding process still works after the swap.
The cool thing is that we can make the database conditional
and keep both implementations.

(show how to migrate template to paritals and add conditions)

We will be using [GORM](https://gorm.io/docs/) ~ "The fantastic ORM library for Golang aims to be developer friendly"

<!--
The full code for this section can be found on GitHub
[code/first-example/more-features](https://github.com/hofstadter-io/hof-docs/tree/main/code/first-example/using-a-database)
-->

