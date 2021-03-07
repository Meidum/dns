FROM node:10.16.0-stretch AS frontend-build

WORKDIR /frontend
COPY frontend /

RUN npm install -g yarn
RUN yarn install
RUN yarn build

FROM golang:1.15-stretch AS server-build

RUN go get github.com/GeertJohan/go.rice && go get github.com/GeertJohan/go.rice/rice

RUN mkdir -p src/github.com/meidum/dns/frontend/build
WORKDIR src/github.com/meidum/dns

COPY --from=frontend-build build frontend/build
COPY db ./db
COPY records ./records
COPY roles ./roles
COPY users ./users
COPY util ./util
COPY main.go ./main.go
COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go mod download
RUN rice embed-go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /dns -a -installsuffix cgo .

FROM scratch

COPY --from=server-build /dns /dns

ENTRYPOINT ["/dns"]
