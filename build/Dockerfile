FROM golang:1.23 as builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

ARG DIR=/workspace
WORKDIR $DIR


# Copy the go source including go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download


# Copy the rest of the source code
COPY . .

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
COPY --from=builder /workspace/fileDB /usr/local/bin/
