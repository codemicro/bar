.PHONY: clean install uninstall run

install_dir := ~/.local/bin

bin/cdmbar:
	mkdir bin
	go build -o bin/cdmbar github.com/codemicro/bar/cmd/bar

run: bin/cdmbar
	./bin/cdmbar

clean:
	rm -rf bin

install: bin/cdmbar
	cp bin/cdmbar $(install_dir)

uninstall:
	rm $(install_dir)/cdmbar
