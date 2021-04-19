#==============================================================#
# File      :   Makefile
# Ctime     :   2019-04-13
# Mtime     :   2020-09-17
# Desc      :   Makefile shortcuts
# Path      :   Makefile
# Copyright (C) 2018-2021 Ruohang Feng
#==============================================================#

VERSION=`cat main.go | grep -E 'var Version' | grep -Eo '[0-9.]+'`

###############################################################
# Public objective
###############################################################
build:
	go build -o pigsty

clean:
	rm -rf pigsty

release-darwin: clean
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build  -a -ldflags '-extldflags "-static"' -o pigsty
	upx pigsty
	mv -f pigsty bin/pigsty_v$(VERSION)_darwin-amd64

release-linux: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -a -ldflags '-extldflags "-static"' -o pigsty
	upx pigsty
	mv -f pigsty bin/pigsty_v$(VERSION)_linux-amd64

release: clean release-linux release-darwin # release-windows

test:
	go build -o pigsty
	mv pigsty ~/pigsty/pigsty


install: build
	sudo install -m 0755 pigsty /usr/local/bin/pigsty