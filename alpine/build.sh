#!/bin/bash
rm  ../dist/terraform-switcher-alpine || true
env ../version
cd .. && docker build  -f alpine/Dockerfile  -t alpine-build .
cd dist && docker run --rm -iv${PWD}:/host-volume alpine-build  sh -s <<EOF
chown -v $(id -u):$(id -g) /app/terraform-switcher
cp -va /app/terraform-switcher  /host-volume/terraform-switcher-alpine
EOF


