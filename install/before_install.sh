#!/bin/bash

# Generated by: tyk-ci/wf-gen
# Generated on: Fri  3 Sep 12:21:39 UTC 2021

# Generation commands:
# ./pr.zsh -title [TT-2932] systemd service restore for deb -branch releng/upgrades -repos tyk,tyk-analytics
# m4 -E -DxREPO=tyk


echo "Creating user and group..."
GROUPNAME="tyk"
USERNAME="tyk"

getent group "$GROUPNAME" >/dev/null || groupadd -r "$GROUPNAME"
getent passwd "$USERNAME" >/dev/null || useradd -r -g "$GROUPNAME" -M -s /sbin/nologin -c "Tyk service user" "$USERNAME"


# This stopped being a symlink in PR #3569
if [ -L /opt/tyk-gateway/coprocess/python/proto ]; then
    echo "Removing legacy python protobuf symlink"
    rm /opt/tyk-gateway/coprocess/python/proto
fi
