
.DEFAULT_GOAL := test

export GO111MODULE=on

PUBRELEASE="false"
LATEST=$(shell curl -s https://github.com/jbvmio/kafkactl/releases/latest | awk -F '/' '/releases/{print $$8}' | awk -F '"' '{print $$1}')
LATESTMAJ=$(shell echo $(LATEST) | cut -d '.' -f 1)
LATESTMIN=$(shell echo $(LATEST) | cut -d '.' -f 2)
LATESTPAT=$(shell echo $(LATEST) | cut -d '.' -f 3)
NEXTPAT=$(shell echo $(LATESTPAT) + 1 | bc)
NEXTVER=$(shell echo $(LATESTMAJ).$(LATESTMIN).$(NEXTPAT))
ifeq ($(shell uname), "Darwin" )
	YTIME=$(shell date -juf "%b %d %Y %T" "Jan 1 2019 00:00:00" +%s)
else
	YTIME=$(shell date -d "Jan 1 2019 00:00:00" +%s)
endif
BT=$(shell date +%s)
GCT=$(shell git rev-list -1 HEAD --timestamp | awk '{print $$1}')
GC=$(shell git rev-list -1 HEAD --abbrev-commit)
REV=$(shell echo $(GCT)-$(YTIME) | bc)
FNAME=$(shell echo $(LATEST)-beta.$(REV))

ld_flags := "-X github.com/jbvmio/kafkactl/cli/cmd.latestMajor=$(LATESTMAJ) -X github.com/jbvmio/kafkactl/cli/cmd.latestMinor=$(LATESTMIN) -X github.com/jbvmio/kafkactl/cli/cmd.latestPatch=$(LATESTPAT) -X github.com/jbvmio/kafkactl/cli/cmd.release=$(PUBRELEASE) -X github.com/jbvmio/kafkactl/cli/cmd.nextRelease=$(NEXTVER) -X github.com/jbvmio/kafkactl/cli/cmd.revision=$(REV) -X github.com/jbvmio/kafkactl/cli/cmd.buildTime=$(BT) -X github.com/jbvmio/kafkactl/cli/cmd.commitHash=$(GC)"

build:
	GOOS=darwin ARCH=amd64 go build -ldflags $(ld_flags) -o kafkactl.$(FNAME).darwin
	GOOS=linux ARCH=amd64 go build -ldflags $(ld_flags) -o kafkactl.$(FNAME).linux
	GOOS=darwin ARCH=amd64 go build -ldflags $(ld_flags) -o kafkactl.$(FNAME).exe

clean:
	rm -f kafkactl.darwin
	rm -f kafkactl.linux
	rm -f kafkactl.exe

test: build clean

docker:
	GOOS=linux ARCH=amd64 go build -ldflags $(ld_flags) -o /usr/local/bin/kafkactl

release:
	printf "[ RELEASE $(FNAME) ]\n" > .commit.log
	git log --oneline --decorate >> .commit.log
	git add .
	git commit -m "release $(FNAME)"
	git tag -a $(FNAME) -m "release $(FNAME)"
	git push origin
	git push origin $(FNAME)
