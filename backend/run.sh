#!/usr/bin/env bash
set -e

export JOBTRACK_ADDR=:8080
export JOBTRACK_DB=jobtrack.db
export JOBTRACK_CORS_ORIGINS=http://localhost:3000

go run ./cmd/api
