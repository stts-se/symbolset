.PHONY: all
all: linux mac zip

linux: buildinfo main.go converter.go mapper.go meta.go symbolset.go
	GOOS=linux GOARCH=amd64 go build -o symbolset main.go converter.go mapper.go meta.go symbolset.go

mac: buildinfo main.go converter.go mapper.go meta.go symbolset.go
	GOOS=darwin GOARCH=amd64 go build -o symbolset_mac main.go converter.go mapper.go meta.go symbolset.go

buildinfo: 
	@echo -n "" >| buildinfo.txt
	@date "+%Y-%m-%d %H:%M:%S %Z" | sed 's/^/Build timestamp: /' >> buildinfo.txt
	@whoami | sed 's/^/Built by: /' >> buildinfo.txt
	@hash=`git rev-parse HEAD|sed 's/^\(.......\).*/\1/'`; commit=`git rev-parse --abbrev-ref HEAD`; echo "Git commit: $$hash on $$commit" >> buildinfo.txt
	@git describe --tags 2> /dev/null | sed 's/^/Release: /' >> buildinfo.txt

zip: buildinfo linux mac static demo_files
	zip -r -q symbolset.zip symbolset symbolset_mac static/ demo_files/ buildinfo.txt
	rm symbolset
	rm symbolset_mac

clean:
	rm -f symbolset symbolset_mac buildinfo.txt symbolset.zip
