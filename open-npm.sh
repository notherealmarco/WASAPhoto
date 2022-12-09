#!/usr/bin/env sh

docker run -it --rm -v "$(pwd):/src" -u "$(id -u):$(id -g)" -p 4173:4173 -p 5173:5173 --workdir /src/webui node:lts /bin/bash
