#! /bin/bash

export HOST_IP=$(ifconfig en0 | awk '/ *inet /{print $2}')
echo "Host IP = $HOST_IP"
docker compose up -d
