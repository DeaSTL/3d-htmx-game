
clean:
	rm -rf ./bin/*;

build: clean
	mkdir -p bin/;
	go build -o bin/htmx-3d-game .;

run: build
	./bin/htmx-3d-game;
