
clean:
	rm -rf ./bin/*;

build: clean
	mkdir -p bin/;
	templ generate
	go build -o bin/htmx-3d-game .;

run: build
	./bin/htmx-3d-game;
