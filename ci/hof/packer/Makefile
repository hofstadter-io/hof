packer.gen: debian.cue
	hof gen ./debian.cue \
		-T :images.docker=debian-docker.json \
		-T :images.nerdctl=debian-nerdctl.json \
		-T ':images["nerdctl-rootless"]=debian-nerdctl-rootless.json' \
		-T :images.podman=debian-podman.json

packer.build: packer.gen
	packer build -force debian-docker.json
	packer build -force debian-nerdctl.json
	packer build -force debian-nerdctl-rootless.json
	packer build -force debian-podman.json
