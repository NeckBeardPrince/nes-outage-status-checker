"""
NES Outage Chart Data Collector for DigitalOcean Spaces

Fetches NES outage data every 10 minutes, computes aggregate metrics, and stores
chart-data.json in DigitalOcean Spaces. The web app loads this URL for chart history.

Run via cron every 10 minutes. On server: /opt/nes-chart/run.sh (sources /opt/bibcast/.env for keys)

Environment variables (run.sh sets from SPACES_KEY.txt, SPACES_SECRET.txt in script dir):
  SPACES_KEY       - Spaces access key
  SPACES_SECRET    - Spaces secret
  SPACES_BUCKET   - Space name (default: nes)
  SPACES_REGION   - Region (default: nyc3)
  CHART_OBJECT_KEY - Object key (default: chart-data.json)
"""

import json
import os
import sys
import urllib.request
from datetime import datetime, timezone

try:
    from botocore.exceptions import ClientError
except ImportError:
    ClientError = None

API_URL = "https://utilisocial.io/datacapable/v2/p/NES/map/events"
CHART_OBJECT_KEY = os.environ.get("CHART_OBJECT_KEY", "chart-data.json")


def round_to_10_minutes(dt):
    minutes = (dt.minute // 10) * 10
    return dt.replace(minute=minutes, second=0, microsecond=0)


def fetch_nes_events():
    req = urllib.request.Request(API_URL, headers={"User-Agent": "NES-Chart-Collector/1.0"})
    with urllib.request.urlopen(req, timeout=30) as response:
        return json.loads(response.read().decode("utf-8"))


def compute_metrics(events):
    total_affected = sum(e.get("numPeople") or 0 for e in events)
    total_outages = len(events)
    assigned = [e for e in events if e.get("status") != "Unassigned"]
    unassigned = [e for e in events if e.get("status") == "Unassigned"]
    active_crews = len(assigned)
    waiting_for_crew = len(unassigned)
    customers_being_restored = sum(e.get("numPeople") or 0 for e in assigned)
    now = datetime.now(timezone.utc)
    interval = round_to_10_minutes(now)
    return {
        "timestamp": interval.isoformat().replace("+00:00", "Z"),
        "numPeople": total_affected,
        "eventCount": total_outages,
        "activeCrews": active_crews,
        "customersBeingRestored": customers_being_restored,
        "waitingForCrew": waiting_for_crew,
    }


def get_existing_data(s3_client, bucket, key):
    try:
        resp = s3_client.get_object(Bucket=bucket, Key=key)
        body = resp["Body"].read().decode("utf-8")
        return json.loads(body)
    except Exception as e:
        if ClientError and isinstance(e, ClientError) and e.response.get("Error", {}).get("Code") == "NoSuchKey":
            return []
        print(f"Warning: could not load existing data: {e}", file=sys.stderr)
        return []


def trim_to_two_months(data):
    from datetime import timedelta
    cutoff = datetime.now(timezone.utc) - timedelta(days=60)
    return [d for d in data if datetime.fromisoformat(d["timestamp"].replace("Z", "+00:00")) >= cutoff]


def main():
    bucket = os.environ.get("SPACES_BUCKET")
    region = os.environ.get("SPACES_REGION", "nyc3")
    key_id = os.environ.get("SPACES_KEY")
    secret = os.environ.get("SPACES_SECRET")

    if not all([bucket, key_id, secret]):
        print("Set SPACES_BUCKET, SPACES_REGION, SPACES_KEY, SPACES_SECRET", file=sys.stderr)
        sys.exit(1)

    try:
        import boto3
    except ImportError:
        print("Install boto3: pip install boto3", file=sys.stderr)
        sys.exit(1)

    endpoint = f"https://{region}.digitaloceanspaces.com"
    s3 = boto3.client(
        "s3",
        region_name=region,
        endpoint_url=endpoint,
        aws_access_key_id=key_id,
        aws_secret_access_key=secret,
    )

    events = fetch_nes_events()
    point = compute_metrics(events)
    chart_data = get_existing_data(s3, bucket, CHART_OBJECT_KEY)

    interval_ts = point["timestamp"]
    existing_idx = next((i for i, d in enumerate(chart_data) if d.get("timestamp") == interval_ts), None)
    if existing_idx is not None:
        chart_data[existing_idx] = point
    else:
        chart_data.append(point)

    chart_data.sort(key=lambda d: d["timestamp"])
    chart_data = trim_to_two_months(chart_data)

    s3.put_object(
        Bucket=bucket,
        Key=CHART_OBJECT_KEY,
        Body=json.dumps(chart_data),
        ContentType="application/json",
        ACL="public-read",
    )
    print(f"Updated {CHART_OBJECT_KEY}: {len(chart_data)} points (interval {interval_ts})")


if __name__ == "__main__":
    main()
