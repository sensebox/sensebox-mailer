#!/bin/bash

set -e

function cleanup() {
  echo 'cleanup!'
  sudo docker stop mailhog
  sudo docker rm mailhog
}
trap cleanup EXIT

sudo docker run --name mailhog -p "1025:1025" -p "8025:8025" mailhog/mailhog
