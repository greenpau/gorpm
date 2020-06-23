.PHONY: clean build

all: build
	@echo "Done"

clean:
	@echo "OK: clean up completed"

build: clean
	@ cat config.json | jq '.'
	@./build.sh
