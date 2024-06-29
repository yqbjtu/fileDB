FROM golang:1.13.10 as builder

ARG DIR=/go/src/fileDB
WORKDIR $DIR

# Copy the go source
COPY . ./
#
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o fileDB $DIR

FROM alpine:3.9
LABEL maintainers="developer"


ARG APK_MIRROR=mirrors.aliyun.com
ARG CodeSource=
ARG CodeBranches=
ARG CodeVersion=

ENV CODE_SOURCE=$CodeSource
ENV CODE_BRANCHES=$CodeBranches
ENV CODE_VERSION=$CodeVersion


LABEL CodeSource=$CodeSource CodeBranches=$CodeBranches CodeVersion=$CodeVersion

WORKDIR /bin/
CMD [ "fileDB" ]

# use hardcode , not %DIR
COPY --from=builder /go/src/gindemo/fileDB /usr/local/bin/
