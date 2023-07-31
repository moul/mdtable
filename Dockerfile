# dynamic config
ARG             BUILD_DATE
ARG             VCS_REF
ARG             VERSION

# build
FROM            golang:1.20-alpine as builder
RUN             apk add --no-cache git gcc musl-dev make
ENV             GO111MODULE=on
WORKDIR         /go/src/moul.io/mdtable
COPY            go.* ./
RUN             go mod download
COPY            . ./
RUN             make install

# minimalist runtime
FROM alpine:3.18.0
LABEL           org.label-schema.build-date=$BUILD_DATE \
                org.label-schema.name="mdtable" \
                org.label-schema.description="" \
                org.label-schema.url="https://moul.io/mdtable/" \
                org.label-schema.vcs-ref=$VCS_REF \
                org.label-schema.vcs-url="https://github.com/moul/mdtable" \
                org.label-schema.vendor="Manfred Touron" \
                org.label-schema.version=$VERSION \
                org.label-schema.schema-version="1.0" \
                org.label-schema.cmd="docker run -i -t --rm moul/mdtable" \
                org.label-schema.help="docker exec -it $CONTAINER mdtable --help"
COPY            --from=builder /go/bin/mdtable /bin/
ENTRYPOINT      ["/bin/mdtable"]
#CMD             []
