#!/usr/bin/env bash

set -e

go vet -tags novet ./...
