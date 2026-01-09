if ! wget --no-verbose --tries=1 --timeout=3 --spider http://localhost:8081/healthcheck/live; then
    echo "Liveness check FAILED"
    exit 1
fi

if ! wget --no-verbose --tries=1 --timeout=3 --spider http://localhost:8081/healthcheck/ready; then
    echo "Readiness check FAILED"
    exit 1
fi

exit 0