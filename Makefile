GIT_HASH := $(shell git log --pretty=format:%H -n 1)
GIT_COMMIT := $(shell git describe --always --tags | cut -d'-' -f2)
GIT_TAG := $(shell git describe --always --tags --abbrev=0 | tail -c +2)
GO_VERSION := $(shell go version | cut -d' ' -f3)
VERSION := ${GIT_TAG}.${GIT_COMMIT}
RELEASE := 1
VENDOR := "SKB Kontur"
URL := "https://github.com/moira-alert"
LICENSE := "GPLv3"

.PHONY: test prepare build tar rpm deb

default: test build

prepare:
	go get github.com/kardianos/govendor
	govendor sync
	go get github.com/alecthomas/gometalinter
	gometalinter --install

lint: prepare
	gometalinter ./... --vendor --skip mock --disable=errcheck --disable=gocyclo

test: prepare
	echo 'mode: atomic' > coverage.txt && go list ./... | grep -v "/vendor/" | xargs -n1 -I{} sh -c 'go test -bench=. -covermode=atomic -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> coverage.txt' && rm coverage.tmp

build:
	go build -ldflags "-X main.Version=${VERSION}-${RELEASE} -X main.GoVersion=${GO_VERSION} -X main.GitHash=${GIT_HASH}" -o build/moira github.com/moira-alert/moira-alert/cmd/moira

clean:
	rm -rf build

tar:
	mkdir -p build/root/usr/bin
	mkdir -p build/root/usr/lib/systemd/system
	mkdir -p build/root/etc/logrotate.d
	mkdir -p build/root/etc/moira

	cp build/moira build/root/usr/bin/
	cp pkg/moira.service build/root/usr/lib/systemd/system/moira.service
	cp pkg/logrotate build/root/etc/logrotate.d/moira

	cp pkg/storage-schemas.conf build/root/etc/moira/storage-schemas.conf
	cp pkg/moira.yml build/root/etc/moira/moira.yml

	tar -czvPf build/moira-${VERSION}-${RELEASE}.tar.gz -C build/root .

rpm: tar
	fpm -t rpm \
		-s "tar" \
		--description "Moira" \
		--vendor ${VENDOR} \
		--url ${URL} \
		--license ${LICENSE} \
		--name "moira" \
		--version "${VERSION}" \
		--iteration "${RELEASE}" \
		--config-files "/etc/moira/moira.yml" \
		--config-files "/etc/moira/storage-schemas.conf" \
		--after-install "./pkg/postinst" \
		--depends logrotate \
		-p build \
		build/moira-${VERSION}-${RELEASE}.tar.gz

deb: tar
	fpm -t deb \
		-s "tar" \
		--description "Moira" \
		--vendor ${VENDOR} \
		--url ${URL} \
		--license ${LICENSE} \
		--name "moira" \
		--version "${VERSION}" \
		--iteration "${RELEASE}" \
		--config-files "/etc/moira/moira.yml" \
		--config-files "/etc/moira/storage-schemas.conf" \
		--after-install "./pkg/postinst" \
		--depends logrotate \
		-p build \
		build/moira-${VERSION}-${RELEASE}.tar.gz

docker_build:
	docker build -t kontur/moira .


packages: clean build tar rpm deb
