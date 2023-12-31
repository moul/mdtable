GOPKG ?=	moul.io/mdtable
DOCKER_IMAGE ?=	moul/mdtable
GOBINS ?=	.
NPM_PACKAGES ?=	.

include rules.mk

generate: install
	mkdir -p .tmp

	echo 'foo@bar:~$$ mdtable' > .tmp/usage.txt
	(go run . csv -h 2>&1 || true) >> .tmp/usage.txt

	echo 'foo@bar:~$$ cat testdata/file2.csv' > .tmp/usage-csv.txt
	cat testdata/file2.csv >> .tmp/usage-csv.txt
	echo >> .tmp/usage-csv.txt
	echo 'foo@bar:~$$ cat testdata/file2.csv | mdtable csv' >> .tmp/usage-csv.txt
	(cat testdata/file2.csv | go run . csv 2>&1 || true) >> .tmp/usage-csv.txt
	echo >> .tmp/usage-csv.txt
	echo 'foo@bar:~$$ cat testdata/file2.csv | mdtable csv -csv-no-header' >> .tmp/usage-csv.txt
	(cat testdata/file2.csv | go run . csv -csv-no-header 2>&1 || true) >> .tmp/usage-csv.txt
	echo >> .tmp/usage-csv.txt
	echo 'foo@bar:~$$ cat testdata/file2.csv | mdtable csv -md-header="foo\\tbar\\tbaz\\t" -md-body="{{range .Cols}}<<{{.}}>>\\t{{end}}"' >> .tmp/usage-csv.txt
	(cat testdata/file2.csv | go run . csv -md-header="foo\\tbar\\tbaz\\t" -md-body="{{range .Cols}}<<{{.}}>>\\t{{end}}" 2>&1 || true) >> .tmp/usage-csv.txt


	go run github.com/campoy/embedmd -w README.md
	rm -rf .tmp
.PHONY: generate
