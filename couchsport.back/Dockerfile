FROM golang:1.15-alpine AS setup
WORKDIR /go/app
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ARG TARGETOS
ARG TARGETARCH

FROM setup AS sources
COPY . .

FROM sources AS compile
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o /release/couchsport.back .

FROM sources AS run-dev
COPY ./api/templates .

RUN apk add git curl \
    make openssh-client
RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp -f ./bin/air /bin/air

CMD air -c air.toml

FROM scratch AS release
WORKDIR /
COPY --from=compile /release/couchsport.back .
COPY ./api/templates .

ENTRYPOINT ["./couchsport.back", "--env", "production"]


