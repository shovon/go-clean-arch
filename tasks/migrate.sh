#!/bin/bash

# To be used in development, not staging, and especially not prod. Just
# development.

cd $(dirname $(dirname ${BASH_SOURCE[0]}))

PGCONNSTRING_DOCKER=$(echo $PGCONNSTRING | sed 's/localhost/host.docker.internal/')

docker run -i -v "$PWD/db/migrations:/migrations" --network host migrate/migrate -path=/migrations/ -database $PGCONNSTRING_DOCKER $@
