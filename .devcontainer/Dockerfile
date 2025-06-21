FROM mcr.microsoft.com/devcontainers/go:1.24-bookworm

ARG NODE_VERSION="none"
RUN if [ "${NODE_VERSION}" != "none" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

RUN apt-get update && apt-get install redis -y

ENTRYPOINT ["sleep", "infinity"]
