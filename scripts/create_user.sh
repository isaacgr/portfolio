#!/bin/bash

set -e

if ! id portfolio; then
    groupadd --system portfolio
    useradd --system --gid portfolio --no-create-home --shell /usr/sbin/nologin portfolio
fi
