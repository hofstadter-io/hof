CUE_FILES  = $(shell find code/ -type f -name '*.cue' ! -path './code/vendor/*')
HTML_FILES = $(patsubst code/%.cue, code/%.html, $(CUE_FILES))

.PHONY: highlight code
highlight-cue: $(HTML_FILES)

.PHONY: codefiles
codefiles: $(CUE_FILES)
	@for f in $(CUE_FILES); do echo $$f; done

.PHONY: htmlfiles cleanhtml
htmlfiles: $(HTML_FILES)
	@for f in $(HTML_FILES); do echo $$f; done
cleanhtml:
	@for f in $(HTML_FILES); do rm $$f; done

code/%.html: code/%.cue
	@echo highlight "$<" as "$@"
	@node hack/highlight.js < "$<" > "$@"

.PHONY: gen-output-examples
gen:
	make -C code/first-example gen
