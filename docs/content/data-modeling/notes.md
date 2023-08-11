---
title: Notes

weight: 100
---

{{<lead>}}
These are some notes that haven't found a home yet or have not been written about.
{{</lead>}}

Data Lenses:

- [Project Cambria](https://www.inkandswitch.com/cambria/) (main inspiration)
- [maybe related?](https://medium.com/javascript-scene/lenses-b85976cb0534)
- https://github.com/grafana/thema and https://www.youtube.com/watch?v=PpoS_ThntEM

Similar:

- [Atlas](https://atlasgo.io/) & [Ent](https://entgo.io/)
- [Prisma](https://www.prisma.io/docs/concepts/components/prisma-schema/data-model)
- [Hasura](https://hasura.io/docs/latest/graphql/core/databases/postgres/schema/index.html)


Questions:

- how to checkpoint with a semver attached while also having a tag in git?
	- Does this need to be done? (see next question)
	- To automate, require clean git workspace and make a commit and tag after the checkpoint is written?
	- What we probably don't want is version skew and/or matrix
- Would it make sense for larger orgs to keep their datamodels in a dedicated repos?
	- This would aid reuse and sharing
	- here we might want to connect checkpoint and git tags
	- otherwise, possibly not

- Single source of Truth
	- https://news.ycombinator.com/item?id=32010699

