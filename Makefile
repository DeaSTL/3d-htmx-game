
clean:
	rm -rf ./bin/*;

build: clean
	mkdir -p bin/;
	templ generate
	go build -o bin/main .;

run: build
	./bin/main;
