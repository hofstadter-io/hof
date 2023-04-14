# todo


## OCI

### mvp

https://github.com/opencontainers/image-spec/blob/main/manifest.md

- format
  - metadata
	- mod files
	- everything (same hash as git based possible)


- round trip
  - package mod in OCI custom manifest
	- upload
	- download and unpack where expected
	- test that `cue eval` still works


=====================

- integrate into mod cmd
- auth


## Inference

- we want to transparently support git | OCI | (other)
- code to figure out (based on try & response)
- keep config file so we don't have to figure out more than once


## Git

- make use of the codehost code from Go compiler


