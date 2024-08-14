.PHONY: licenses
licenses:
	rm -rf ./LICENSES
	go-licenses save ./... --save_path=./LICENSES