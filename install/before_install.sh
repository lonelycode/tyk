#!/bin/bash

# Generated by: tyk-ci/wf-gen
# Generated on: Wed  1 Dec 12:39:55 UTC 2021

# Generation commands:
# ./pr.zsh -title sync from templates -branch aws/pkr176 -base aws/pkr176 -repos tyk -p
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
