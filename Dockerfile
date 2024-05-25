# Build section of the docker file
FROM golang:latest as build

ARG APPNAME=githubrunner

WORKDIR /app
COPY . .

RUN go mod tidy && go mod tidy
RUN make build-linux
RUN chmod +x ./build/linux/$APPNAME

# App and server section of the docker file
FROM ubuntu:latest as serve

RUN apt-get update
RUN apt-get install -y ca-certificates

ARG USERNAME=nonroot
ARG USER_UID=1001
ARG USER_GID=$USER_UID

COPY --from=build /app/build/linux/$APPNAME /app/$APPNAME

EXPOSE 3038

RUN groupadd --gid $USER_GID $USERNAME
RUN useradd --uid $USER_UID --gid $USER_GID $USERNAME
RUN passwd -d $USERNAME
RUN chown -R $USERNAME:$USERNAME ./app

USER $USERNAME:$USERNAME

WORKDIR /app
ENTRYPOINT [ "./$APPNAME", "server" ]