OUTPUT_DIR = ./builds
GIT_COMMIT = `git rev-parse HEAD | cut -c1-7`
VERSION = 1.0.0
BUILD_OPTIONS = -ldflags "-X main.Version=$(VERSION) -X main.CommitID=$(GIT_COMMIT)"

smtp2mqtt: main.go mail.go mqtt.go version.go Makefile
	dep ensure
	go build ${BUILD_OPTIONS}

.PHONY: all
all: smtp2mqtt


tools:
	go get -u github.com/golang/dep/cmd/dep
	go get github.com/mitchellh/gox

test:
	if [ `go fmt $(go list ./... | grep -v /vendor/) | wc -l` -gt 0 ]; then echo "go fmt error"; exit 1; fi

cross_compile:
	GOARM=5 gox -os="darwin linux freebsd netbsd openbsd" -arch="386 amd64 arm" -osarch="!darwin/arm" -output "${OUTPUT_DIR}/pkg/{{.OS}}_{{.Arch}}/{{.Dir}}"

targz:
	mkdir -p ${OUTPUT_DIR}/dist
	cd ${OUTPUT_DIR}/pkg/; for osarch in *; do (cd $$osarch; tar zcvf ../../dist/smtp2mqtt_${VERSION}_$$osarch.tar.gz ./*); done;

shasums:
	cd ${OUTPUT_DIR}/dist; sha256sum * > ./SHA256SUMS

rel:
	mkdir -p release
	CGO_ENABLED=0 gox -ldflags="-s -w -X main.Version=$(VERSION) -X main.CommitID=$(GIT_COMMIT)" -output="release/{{.Dir}}_{{.OS}}_{{.Arch}}"

release:
	ghr -c ${GIT_COMMIT} --delete --prerelease -u kost -r smtp2mqtt pre-release ${OUTPUT_DIR}/dist
