#!/usr/bin/bash

DOCKER_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# cd $DOCKER_DIR && docker build -f Dockerfile --tag vault-go:54edf0ba9 ../
# docker scout cves image://vault-go:54edf0ba9

cd $DOCKER_DIR && docker build -f Dockerfile --tag babbili/vault-go:54edf0ba9 ../

# docker run vault-go:54edf0ba9 --name vault-go

docker push babbili/vault-go:54edf0ba9
