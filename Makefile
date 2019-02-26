
.DEFAULT_GOAL := test

export GO111MODULE=on

vars:
	eval $(shell ./envvars.darwin)

flags: vars
ld_flags := "-X github.com/jbvmio/kafkactl/cli/cmd.latestMajor=$(LATESTMAJ) -X github.com/jbvmio/kafkactl/cli/cmd.latestMinor=$(LATESTMIN) -X github.com/jbvmio/kafkactl/cli/cmd.latestPatch=$(LATESTPAT) -X github.com/jbvmio/kafkactl/cli/cmd.release=$(PUBRELEASE) -X github.com/jbvmio/kafkactl/cli/cmd.nextRelease=$(NEXTVER) -X github.com/jbvmio/kafkactl/cli/cmd.revision=$(REV) -X github.com/jbvmio/kafkactl/cli/cmd.buildTime=$(BT) -X github.com/jbvmio/kafkactl/cli/cmd.commitHash=$(GC) -X github.com/jbvmio/kafkactl/cli/cmd.gitVersion=$(FNAME)"

build: flags
	GOOS=darwin ARCH=amd64 go build -ldflags $(ld_flags) -o kafkactl.$(FNAME).darwin
	GOOS=linux ARCH=amd64 go build -ldflags $(ld_flags) -o kafkactl.$(FNAME).linux
	GOOS=darwin ARCH=amd64 go build -ldflags $(ld_flags) -o kafkactl.$(FNAME).exe

exbuild:
	GOOS=darwin ARCH=amd64 go build -ldflags "s -w -X $(LD_FLAGS_STATIC)" -o kafkactl.$(FNAME_STATIC).darwin
	GOOS=linux ARCH=amd64 go build -ldflags "s -w -X $(LD_FLAGS_STATIC)" -o kafkactl.$(FNAME_STATIC).linux
	GOOS=darwin ARCH=amd64 go build -ldflags "s -w -X $(LD_FLAGS_STATIC)" -o kafkactl.$(FNAME_STATIC).exe

clean:
	rm -f kafkactl.$(FNAME).darwin
	rm -f kafkactl.$(FNAME).linux
	rm -f kafkactl.$(FNAME).exe

exclean:
	rm -f kafkactl.$(FNAME_STATIC).darwin
	rm -f kafkactl.$(FNAME_STATIC).linux
	rm -f kafkactl.$(FNAME_STATIC).exe

test: build clean

extest: exbuild exclean

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
