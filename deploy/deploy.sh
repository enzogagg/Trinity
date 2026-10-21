#!/bin/sh
set -eu
cd "$(dirname "$0")"

docker compose -f compose.infra.yml up -d --wait

CURRENT=$(docker exec trinity-proxy cat /etc/nginx/active/color.conf 2>/dev/null | grep -oE 'blue|green' || true)
if [ "$CURRENT" = "blue" ]; then NEW=green; else NEW=blue; fi
echo "active=${CURRENT:-none} -> deploying $NEW (tag $TAG)"

export COLOR=$NEW
docker compose -p "trinity-$NEW" -f compose.app.yml pull
docker compose -p "trinity-$NEW" -f compose.app.yml up -d --wait

docker exec trinity-proxy sh -c "echo 'set \$color $NEW;' > /etc/nginx/active/color.conf && nginx -s reload"
echo "traffic switched to $NEW"

if [ -n "$CURRENT" ]; then
  COLOR=$CURRENT docker compose -p "trinity-$CURRENT" -f compose.app.yml down
fi
