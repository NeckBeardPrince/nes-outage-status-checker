# NES Outage Status Checker

A web app and CLI tool for monitoring Nashville Electric Service power outages.

## Project Structure

- `docs/` - Static web app (GitHub Pages site), NOT documentation
  - `index.html` - Dashboard / landing page (view all active outages, map view, chart view)
  - `monitor.html` - Monitor page (track specific outages by Event ID or address)
  - `all.html` - Redirect stub to `index.html`, kept for pre-1.13.0 bookmarks
  - `how-it-works.html` - Explainer page
  - `version.json` - App version for update detection
- `cmd/` - Go CLI application (Bubble Tea TUI)
- `terraform/` - AWS infrastructure (Lambda, S3, EventBridge, API Gateway)
- `lambda/` - Python Lambda functions
  - `nes_archiver.py` - Archives NES data to S3 every 10 min
  - `health_check.py` - Health check endpoint
- `scripts/` - Utility scripts

## Tech Stack

- **Web**: Vanilla HTML/CSS/JS, Leaflet.js for maps
- **CLI**: Go with Bubble Tea and Lip Gloss
- **Infra**: Terraform, AWS (S3, Lambda, API Gateway, EventBridge)
- **Data**: NES public outage API (no auth required)

## Key Conventions

- Dark theme default, light mode toggle
- All user preferences stored in localStorage
- Version bumps require updating both `docs/version.json` and `CHANGELOG.md`
- Nominatim (OpenStreetMap) used for geocoding - respect 1 req/sec rate limit

## Deployment

- **Web**: Push to `main` auto-deploys to GitHub Pages
- **Domain**: https://nes-outage-checker.com (also on GitHub Pages)
- **Infra**: `cd terraform && terraform apply` (uses S3 backend)

## AWS Details

- Account: 517397653073
- Region: us-east-2
- Profile: personal
- State bucket: astracen-terraform-state
