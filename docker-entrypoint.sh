#!/bin/sh
set -e

# Ensure data directory exists and is writable by the easy user
# When Docker bind-mounts a host directory, it may be owned by root,
# making it inaccessible to the non-root 'easy' user (UID 10001).
if [ "$(id -u)" = "0" ]; then
    mkdir -p /etc/easy-proxies/data
    chown -R easy:easy /etc/easy-proxies/data

    # Also ensure config files are readable
    chown easy:easy /etc/easy-proxies/config.yaml 2>/dev/null || true
    chown easy:easy /etc/easy-proxies/nodes.txt 2>/dev/null || true

    # Drop privileges and exec as 'easy' user
    exec gosu easy /usr/local/bin/easy-proxies "$@"
fi

# Already running as non-root (e.g. Kubernetes securityContext)
exec /usr/local/bin/easy-proxies "$@"