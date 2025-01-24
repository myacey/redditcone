#!/bin/sh

set -e

cmd="$@"

>&2 echo "!!!!!!!!!! Check mongo for available !!!!!!!!!!"

until nc -z mongodb 27017; do
    >&2 echo "MongoDB is unavailable - sleeping"
    sleep 5
done

>&2 echo "mongo now available, executing command"

exec $cmd
