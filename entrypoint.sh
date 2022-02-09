#/bin/sh
set -e

ssh-keyscan -H github.com >> /etc/ssh/ssh_known_hosts

exec "$@"