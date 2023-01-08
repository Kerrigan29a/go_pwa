all: run

run: build
	cd build && python3 -m http.server

clean:
	-rm -r build

build: build/app.wasm build/wasm_exec.js ${subst assets,build,$(wildcard assets/*)}
	echo $^

build/app.wasm: app/main.go
	GOOS=js GOARCH=wasm go build -o $@ $<

build/wasm_exec.js:
	cp $(shell go env GOROOT)/misc/wasm/$(notdir $@) $(dir $@)

build/%: assets/%
	cp $< $@
