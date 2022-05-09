.PHONY: clean install uninstall run

install_dir := ~/.local/bin
go_exe := go

bin/cdmbar:
	mkdir -p bin
	$(go_exe) build $(go_args) -o bin/cdmbar github.com/codemicro/bar/cmd/bar

run: bin/cdmbar
	./bin/cdmbar

clean:
	rm -rf bin

install: bin/cdmbar
	cp bin/cdmbar $(install_dir)

uninstall:
	rm $(install_dir)/cdmbar
