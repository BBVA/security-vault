#!/bin/dumb-init /bin/sh
set -e

# Note above that we run dumb-init as PID 1 in order to reap zombie processes
# as well as forward signals to all processes in its session. Normally, sh
# wouldn't do either of these functions so we'd leak zombies as well as do
# unclean termination of all our sub-processes.

# VAULT_CONFIG_DIR isn't exposed as a volume but you can compose additional
# config files in there if you use this image as a base, or use
# VAULT_LOCAL_CONFIG below.
VAULT_CONFIG_DIR=/vault/config

# You can also set the VAULT_LOCAL_CONFIG environment variable to pass some
# Vault configuration JSON without having to bind any volumes.
if [ -n "$VAULT_LOCAL_CONFIG" ]; then
	echo "$VAULT_LOCAL_CONFIG" > "$VAULT_CONFIG_DIR/local.json"
fi

#complete vault configuration

export VAULT_ADDR="http://172.0.0.1:8200"
export VAULT_TOKEN=$VAULT_DEV_ROOT_TOKEN_ID

vault mount pki
vault mount-tune -max-lease-ttl=87600h pki
vault write pki/root/generate/internal common_name=ca.cloudframe.wtf ttl=87600h
vault write pki/config/urls issuing_certificates="http://0.0.0.0:8200/v1/pki/ca" crl_distribution_points="http://0.0.0.0:8200/v1/pki/crl"
vault write pki/roles/cloudframe-dot-wtf allowed_domains="cloudframe.wtf" allow_subdomains="true" max_ttl="72h"
