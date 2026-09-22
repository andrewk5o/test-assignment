#!/bin/sh
# Patch nginx.conf with the actual runtime DNS address.
# nginx.conf hardcodes 127.0.0.11 (Docker's DNS), but Podman, Kubernetes and
# other runtimes use a different address. This runs before nginx starts via the
# /docker-entrypoint.d/ hook (prefixed 05- to run before all built-in scripts).
set -e

# Extract the first nameserver from the container's resolv.conf.
ns=$(awk '/^nameserver/ { print $2; exit }' /etc/resolv.conf)

if [ -n "$ns" ]; then
    # Replace the hardcoded IP in nginx.conf with the discovered one.
    sed -i "s|resolver 127\.0\.0\.11 |resolver $ns |" /etc/nginx/conf.d/default.conf
    echo "resolver set to $ns"
fi
