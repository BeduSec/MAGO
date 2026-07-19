# Copyright (c) BeduSec. All rights reserved.
#!/bin/bash
echo "Health check:"
curl -s http://localhost:8080/healthz
echo -e "\n\nMetrics:"
curl -s http://localhost:8080/metrics | head -10
echo -e "\n\nTrigger rate limit (default 5 req/s burst):"
for i in {1..6}; do
  curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/
done
echo -e "\n\nReload rules (requires token):"
curl -s -X POST -H "Authorization: Bearer secret-token" http://localhost:8080/v1/rules/reload
echo ""