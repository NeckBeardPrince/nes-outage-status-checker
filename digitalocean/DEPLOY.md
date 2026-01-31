# Deploy NES chart collector to bibcast server

Runs at **/opt/nes-chart** on `root@159.203.69.219`. Uses your existing bibcast access keys from `/opt/bibcast/.env`; cron updates `chart-data.json` in Spaces every 10 minutes. Frontend loads that JSON URL.

## 1. On the server: create folder

```bash
ssh root@159.203.69.219
mkdir -p /opt/nes-chart
```

## 2. From your machine: copy files

From the repo root (where `digitalocean/` lives):

```bash
scp digitalocean/chart_collector.py digitalocean/requirements.txt digitalocean/run.sh root@159.203.69.219:/opt/nes-chart/
```

## 3. On the server: credentials and deps

`run.sh` reads credentials from the same directory (**/opt/nes-chart/**): access key from `SPACES_KEY.txt`, secret from `SPACES_SECRET.txt` (first line of each). It writes to Space **nes** in **nyc3**, so `chart-data.json` is written there.

Then create a venv and install deps (server uses externally-managed Python):

```bash
ssh root@159.203.69.219
cd /opt/nes-chart
python3 -m venv venv
venv/bin/pip install -r requirements.txt
chmod +x run.sh
```

## 4. On the server: add cron (every 10 minutes)

```bash
crontab -e
```

Add:

```
*/10 * * * * /opt/nes-chart/run.sh >> /var/log/nes-chart.log 2>&1
```

(Optional: `touch /var/log/nes-chart.log`.)

## 5. Run once to create the JSON in Spaces

```bash
ssh root@159.203.69.219
/opt/nes-chart/run.sh
```

## 6. URL for the frontend

The script writes to Space **nes** (**nyc3**). The JSON is public at:

**`https://nes.nyc3.digitaloceanspaces.com/chart-data.json`**

This URL is set in **docs/all.html** as `CHART_DATA_URL`. To use a different Space, override `SPACES_BUCKET` / `SPACES_REGION` in `run.sh` and update `CHART_DATA_URL` in the frontend.

## CORS

Ensure the Space has CORS allowed for your app’s origin (e.g. `https://nes-outage-checker.com`, GitHub Pages origin). In DigitalOcean: Space → Settings → CORS, or use the API/CLI to set a rule that allows `GET` from your origins.
