GOPKG ?=	moul.io/mdtable
DOCKER_IMAGE ?=	moul/mdtable
GOBINS ?=	.
NPM_PACKAGES ?=	.

include rules.mk

generate: install
	GO111MODULE=off go get github.com/campoy/embedmd
	mkdir -p .tmp
	echo 'foo@bar:~$$ mdtable' > .tmp/usage.txt
	mdtable 2>&1 >> .tmp/usage.txt
	embedmd -w README.md
	rm -rf .tmp
.PHONY: generate
