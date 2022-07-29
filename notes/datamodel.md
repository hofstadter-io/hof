# datamodel

next* version datamodel upgrades

- [ ] load datamodel history into generators
- [ ] diff3 tests, can we get an easy fix in? (timebox this)
- [ ] add thema mode & support, can we bidirectionally convert between formats?
- [ ] support using raw cue values with constraints as valid input and markup for the code engine
- [ ] support translation methods between models / versions to be defined by user, parse refs in cue to something consumable by code gen
- [ ] we're more focused on datamodel++ (the extra config/metadata needed for codegen of applications), than using history at runtime (?) (tradeoffs)
- [ ] implement views / relns into info system
- [ ] mock (faux) data (https://www.getsynth.com/docs/getting_started/core-concepts)
- [ ] maybe we just need to add lacuna to our schemas? do they represent the cases where we cannot automate or the user wants to be explicit?
- [ ] support with @rename() / @move() for fields (& models?), use `><` in the diff algo, helper #Defs for users in their CUE
- [ ] commit should take a message, we might want to support a message or description field in several places, t.b.d.

can we use the concepts from hof dm and thema
to support code gen mods interoperability?

first, generators need an explicit schema / datamodel version, indep of their release version
this is like k8s resources having apiVersion or docker-compose / OpenAPI spec version when writing

examples:

- how to make client cli tools support newer API versions, and the client-side filler data
- even if the server is not Go, can it send CUE to a Go CLI such that these can be pushed dynamically
- what about a single CLI talking to many copies of the API, like kubectl -> clusters


---

- https://github.com/Qovery/Replibyte
- https://docs.open-metadata.org/metadata-standard/schemas/overview

