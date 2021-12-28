#!/bin/sh

# Generated by: tyk-ci/wf-gen
# Generated on: Tue Dec 28 15:00:19 UTC 2021

# Generation commands:
# ./pr.zsh -base release-4.0 -branch aws/TD-664-r4-0 -title TD-664/Aws-cf-templates for release-4.0 -repos tyk
# m4 -E -DxREPO=tyk


if command -V systemctl >/dev/null 2>&1; then
    if [ ! -f /lib/systemd/system/tyk-gateway.service ]; then
        cp /opt/tyk-gateway/install/inits/systemd/system/tyk-gateway.service /lib/systemd/system/tyk-gateway.service
    fi
else
    if [ ! -f /etc/init.d/tyk-gateway ]; then
        cp /opt/tyk-gateway/install/inits/sysv/init.d/tyk-gateway /etc/init.d/tyk-gateway
    fi
fi
