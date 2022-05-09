.PHONY: clean install uninstall run

bin/cdmbar:
	mkdir bin
	go build -o bin/cdmbar github.com/codemicro/bar/cmd/bar

run: bin/cdmbar
	./bin/cdmbar

clean:
	rm -rf bin

install: bin/cdmbar
	cp bin/cdmbar ~/.local/bin

uninstall:
	rm ~/.local/bin/cdmbar
