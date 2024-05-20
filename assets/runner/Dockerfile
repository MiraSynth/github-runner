FROM ubuntu:latest

ARG GITHUB_RUNNER_VERSION="2.316.1"
ARG GITHUB_RUNNER_SHASUM="d62de2400eeeacd195db91e2ff011bfb646cd5d85545e81d8f78c436183e09a8"
ARG GITHUB_RUNNER_REPOSITORY="https://github.com"
ARG GITHUB_RUNNER_TOKEN="XXXXXXXXXXXXXXXXXXXXX"
ARG GITHUB_RUNNER_LABELS="ubuntu-latest"
ARG USERNAME=nonroot
ARG USER_UID=1001
ARG USER_GID=$USER_UID

ENV GITHUB_RUNNER_REPOSITORY=$GITHUB_RUNNER_REPOSITORY
ENV GITHUB_RUNNER_TOKEN=$GITHUB_RUNNER_TOKEN
ENV GITHUB_RUNNER_LABELS=$GITHUB_RUNNER_LABELS

RUN apt-get update
RUN apt-get install -y ca-certificates curl tar libdigest-sha-perl libicu-dev

WORKDIR /actions-runner

RUN groupadd --gid $USER_GID $USERNAME
RUN useradd --uid $USER_UID --gid $USER_GID -m $USERNAME

COPY ../../runner.sh ./runner.sh
RUN chmod +x ./runner.sh

RUN chown -R $USERNAME:$USERNAME /actions-runner

USER $USERNAME

RUN curl -o actions-runner-linux-x64-${GITHUB_RUNNER_VERSION}.tar.gz -L https://github.com/actions/runner/releases/download/v${GITHUB_RUNNER_VERSION}/actions-runner-linux-x64-${GITHUB_RUNNER_VERSION}.tar.gz
RUN echo "${GITHUB_RUNNER_SHASUM}  actions-runner-linux-x64-${GITHUB_RUNNER_VERSION}.tar.gz" | shasum -a 256 -c
RUN tar xzf ./actions-runner-linux-x64-${GITHUB_RUNNER_VERSION}.tar.gz

ENTRYPOINT [ "./runner.sh" ]