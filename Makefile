GOPKG ?=	moul.io/mdtable
DOCKER_IMAGE ?=	moul/mdtable
GOBINS ?=	.
NPM_PACKAGES ?=	.

include rules.mk

generate: install
	mkdir -p .tmp
	echo 'foo@bar:~$$ mdtable' > .tmp/usage.txt
	(go run . csv -h 2>&1 || true) >> .tmp/usage.txt
	go run github.com/campoy/embedmd -w README.md
	rm -rf .tmp
.PHONY: generate
