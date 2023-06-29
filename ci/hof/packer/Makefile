packer.gen: debian.cue
	hof gen ./debian.cue -T ':images=[]debian-{{ .tool }}.json'

packer.build: packer.gen
	packer build -force debian-docker.json
	packer build -force debian-nerdctl.json
	packer build -force debian-nerdctl-rootless.json
	packer build -force debian-podman.json
