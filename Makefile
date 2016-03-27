.PHONY: all get build strip test deb

NAME := ncr
VERSION := 0.0.2
ARCH := x86_64
MAINTAINER := Nick Groenen <nick@groenen.me>
DESCRIPTION := Nagios Check Runner
URL := https://github.com/zoni/nagios-check-runner

all: get build strip

get:
	go get

build:
	go build -ldflags "-X main.version=$(VERSION)" cmd/*.go

strip:
	strip ncr

test:
	go test -race -cover

deb: build strip
	test -e _workarea && rm -rf _workarea || true
	mkdir -p _workarea/opt/ncr _workarea/etc/ncr
	cp config.yml.example _workarea/etc/ncr/ncr.yml
	cp ncr _workarea/opt/ncr/

	mkdir -p _workarea/etc/systemd/system
	cp packaging/systemd/ncr.service _workarea/etc/systemd/system
	fpm -f -t deb -s dir -C _workarea --name "$(NAME)" --version "$(VERSION)" --architecture "$(ARCH)" --maintainer "$(MAINTAINER)" --description "$(DESCRIPTION)" --url "$(URL)" --config-files /etc --deb-compression xz --deb-suggests nagios-plugins --deb-no-default-config-files --after-install packaging/after-install.sh --after-upgrade packaging/after-upgrade.sh --before-remove packaging/before-remove.sh .

	rm -rf _workarea/etc/systemd/
	mkdir -p _workarea/etc/init
	cp packaging/upstart/ncr.conf _workarea/etc/init
	fpm -f -t deb -s dir -C _workarea --name "$(NAME)-upstart" --version "$(VERSION)" --architecture "$(ARCH)" --maintainer "$(MAINTAINER)" --description "$(DESCRIPTION)" --url "$(URL)" --config-files /etc --deb-compression xz --deb-suggests nagios-plugins --deb-no-default-config-files --after-install packaging/after-install.sh --after-upgrade packaging/after-upgrade.sh --before-remove packaging/before-remove.sh .

arm6_deb:
	GOARCH=arm GOARM=6 go build -ldflags "-X main.version=$(VERSION)" cmd/*.go
	test -e _workarea && rm -rf _workarea || true
	mkdir -p _workarea/opt/ncr _workarea/etc/ncr
	cp config.yml.example _workarea/etc/ncr/ncr.yml
	cp ncr _workarea/opt/ncr/

	mkdir -p _workarea/etc/systemd/system
	cp packaging/systemd/ncr.service _workarea/etc/systemd/system
	fpm -f -t deb -s dir -C _workarea --name "$(NAME)" --version "$(VERSION)" --architecture armhf --maintainer "$(MAINTAINER)" --description "$(DESCRIPTION)" --url "$(URL)" --config-files /etc --deb-compression xz --deb-suggests nagios-plugins --deb-no-default-config-files --after-install packaging/after-install.sh --after-upgrade packaging/after-upgrade.sh --before-remove packaging/before-remove.sh .

	rm -rf _workarea/etc/systemd/
	mkdir -p _workarea/etc/init
	cp packaging/upstart/ncr.conf _workarea/etc/init
	fpm -f -t deb -s dir -C _workarea --name "$(NAME)-upstart" --version "$(VERSION)" --architecture armhf --maintainer "$(MAINTAINER)" --description "$(DESCRIPTION)" --url "$(URL)" --config-files /etc --deb-compression xz --deb-suggests nagios-plugins --deb-no-default-config-files --after-install packaging/after-install.sh --after-upgrade packaging/after-upgrade.sh --before-remove packaging/before-remove.sh .
