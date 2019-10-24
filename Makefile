
.DEFAULT_GOAL := test

export GO111MODULE=on

vars:
PUBRELEASE="true"
YTIME="1546300800"
LATEST=$(shell curl -s https://github.com/jbvmio/kafkactl/releases/latest | awk -F '/' '/releases/{print $$8}' | awk -F '"' '{print $$1}')
LATESTMAJ=$(shell echo $(LATEST) | cut -d '.' -f 1)
LATESTMIN=$(shell echo $(LATEST) | cut -d '.' -f 2)
LATESTPAT=$(shell echo $(LATEST) | cut -d '.' -f 3)
NEXTPAT=$(shell echo $(LATESTPAT) + 1 | bc)
NEXTVER=$(shell echo $(LATESTMAJ).$(LATESTMIN).$(NEXTPAT))
BT=$(shell date +%s)
GCT=$(shell git rev-list -1 HEAD --timestamp | awk '{print $$1}')
GC=$(shell git rev-list -1 HEAD --abbrev-commit)
REV=$(shell echo $(GCT) - $(YTIME) | bc)
FNAME=$(shell echo $(NEXTVER))

flags: vars
ld_flags := "-X github.com/jbvmio/kafkactl/cli/cmd.latestMajor=$(LATESTMAJ) -X github.com/jbvmio/kafkactl/cli/cmd.latestMinor=$(LATESTMIN) -X github.com/jbvmio/kafkactl/cli/cmd.latestPatch=$(LATESTPAT) -X github.com/jbvmio/kafkactl/cli/cmd.release=$(PUBRELEASE) -X github.com/jbvmio/kafkactl/cli/cmd.nextRelease=$(NEXTVER) -X github.com/jbvmio/kafkactl/cli/cmd.revision=$(REV) -X github.com/jbvmio/kafkactl/cli/cmd.commitHash=$(GC) -X github.com/jbvmio/kafkactl/cli/cmd.gitVersion=$(FNAME)"

localbuild: flags
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

test: localbuild clean

extest: exbuild exclean

docker:
	GOOS=linux ARCH=amd64 go build -ldflags $(ld_flags) -o /usr/local/bin/kafkactl

exvars:
	eval $(shell ./envvars.darwin)

print: exvars
	@printf "YTIME :\t\t$(YTIME)\n"
	@printf "COMMITTAG :\t$(COMMITTAG)\n"
	@printf "LATESTPAT :\t$(LATESTPAT)\n"
	@printf "NEXTPAT :\t$(NEXTPAT)\n"
	@printf "REV :\t\t$(REV)\n"
	@printf "FNAME :\t\t$(FNAME)\n"
	@printf "PUBRELEASE :\t$(PUBRELEASE)\n"
	@printf "GC :\t\t$(GC)\n"
	@printf "LATESTMAJ :\t$(LATESTMAJ)\n"
	@printf "LATESTMIN :\t$(LATESTMIN)\n"
	@printf "NEXTVER :\t$(NEXTVER)\n"
	@printf "GCT :\t\t$(GCT)\n"

release: exvars
	printf "[ RELEASE $(FNAME) ]\n" > .commit.log
	git log --oneline --decorate >> .commit.log
	git add .
	git commit -m "release $(FNAME)"
	git tag -a $(FNAME) -m "release $(FNAME)"
	git push origin
	git push origin $(FNAME)

