#!/bin/bash
# Run from cron; reads Spaces credentials from same dir: SPACES_KEY.txt, SPACES_SECRET.txt.
# Cron: */10 * * * * /opt/nes-chart/run.sh >> /var/log/nes-chart.log 2>&1
set -e
cd "$(dirname "$0")"
# Space for chart data: nes, nyc3 → https://nes.nyc3.digitaloceanspaces.com/chart-data.json
export SPACES_BUCKET="${SPACES_BUCKET:-nes}"
export SPACES_REGION="${SPACES_REGION:-nyc3}"
if [ -f SPACES_KEY.txt ]; then
  export SPACES_KEY="$(head -1 SPACES_KEY.txt | tr -d '\r\n')"
fi
if [ -f SPACES_SECRET.txt ]; then
  export SPACES_SECRET="$(head -1 SPACES_SECRET.txt | tr -d '\r\n')"
fi
PYTHON="${PYTHON:-python3}"
[ -x venv/bin/python3 ] && PYTHON=venv/bin/python3
exec "$PYTHON" chart_collector.py
