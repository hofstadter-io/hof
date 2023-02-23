.PHONY: verify_code verify_code verify_diff
verify: verify_code highlight verify_diff

verify_diff:
	@git diff --exit-code code/

verify_code:
	make -C code all
