#!/bin/bash

forum="forum"

echo "dokcer: building forum image"
docker build -t $forum .

# echo "docker: running forum container"
# docker run -d --name=$forum -p 8810:8810 --rm $forum

# Запуск Docker контейнера с переменными окружения из файла .env
echo "docker: running forum container"
docker run -d --name $forum -p 8811:8811 --env-file .env --rm $forum