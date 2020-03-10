# Use base golang image from Docker Hub
# set Docker image build arguments
# build arguments for user/group configurations
ARG USER=wso2carbon
ARG USER_ID=802
ARG USER_GROUP=wso2
ARG USER_GROUP_ID=802
ARG USER_HOME=/home/${USER}
ARG MOTD="\n\
# create the non-root user and group and set MOTD login message
RUN \
    groupadd --system -g ${USER_GROUP_ID} ${USER_GROUP} \


FROM golang:1.12-alpine AS build
RUN apk add --update --no-cache git
WORKDIR /src/aws-secrets-manager
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /app -v ./cmd/aws-secrets-manager

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=build /app /.
ENTRYPOINT ["/app"]

# set the user and work directory
USER ${USER_ID}
WORKDIR ${USER_HOME}