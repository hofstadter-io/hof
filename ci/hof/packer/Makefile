packer.files:
	@echo $(FILES)

packer.gen: debian.cue
	hof gen ./debian.cue -T ':images=[]{{ lookup "builders.[0].image_family" . }}.json'

packer.build: packer.gen
	@set -e; for file in $(shell cue eval debian.cue -e 'strings.Join(packer,"\n")' --out text); do \
		echo $$file; \
		packer build -force $$file; \
		echo $$file; \
	done
