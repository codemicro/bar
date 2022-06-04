.PHONY: clean install uninstall run

install_dir := ~/.local/bin
go_exe := go

bin/cdmbar:
	mkdir -p bin
	$(go_exe) build $(build_args) -o bin/cdmbar github.com/codemicro/bar/cmd/bar

run: clean bin/cdmbar
	./bin/cdmbar

clean:
	rm -rf bin

install: bin/cdmbar
	mkdir -p $(install_dir)
	cp bin/cdmbar $(install_dir)/cdmbar

uninstall:
	rm $(install_dir)/cdmbar
