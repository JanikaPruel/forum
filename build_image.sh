#!/bin/bash

forum="forum"

echo "dokcer: building forum image"
docker build -t $forum .

echo "docker: running forum container"
docker run -d --name=$forum -p 8810:8810 --rm $forum