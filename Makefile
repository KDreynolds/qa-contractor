.PHONY: build generate serve clean

build: generate
	go run ./cmd/generate/

generate:
	templ generate

serve:
	cd out && python3 -m http.server 8080

clean:
	rm -rf out qa-contractor static-gen cmd/generate/static-gen

dev:
	go run . &
	sleep 1
	@echo "http://localhost:3000"
