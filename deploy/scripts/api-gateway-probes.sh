#!/bin/bash

if ! curl -f http://localhost:8080/healthcheck/live; then
    echo "Liveness check FAILED"
    exit 1
fi

if ! curl -f http://localhost:8080/healthcheck/ready; then
    echo "Readiness check FAILED"
    exit 1
fi

exit 0