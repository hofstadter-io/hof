HUGO_VER   = "0.109.0"

# setup the development workspacce
# missing is install hugo...
# convert to hof/flow?
.PHONY: deps
deps:
	npm install
	npm install broken-link-checker -g
	mkdir tmp
	pushd tmp
	wget https://github.com/gohugoio/hugo/releases/download/v$(HUGO_VER)/hugo_extended_$(HUGO_VER)_Linux-64bit.tar.gz
	tar -xf hugo.tar.gz
	chmod +x hugo
	sudo mv hugo /usr/local/bin/hugo
	popd tmp
	rm -rf tmp
	


