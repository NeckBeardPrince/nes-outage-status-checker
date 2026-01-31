# Add line chart view for outage metrics

## Summary

Adds a **Chart** view alongside Grid and Map so users can see trends over time. The view toggle is now three options: **Grid | Map | Chart**. Only one view is shown at a time (chart no longer appears alongside the map).

### Chart features

- **Metrics**: Affected Customers, Number of Outages, Active Crews, Customers Being Restored, Waiting for Crew (selectable in the chart dropdown).
- **Data**: Snapshots every 10 minutes (aligned to 10‑minute intervals). Data kept for the last 2 months.
- **UI**: X-axis labels include date and time. Tooltip on hover shows the value with wording that matches the selected metric (e.g. “12 crews”, “1,234 customers”).

### Chart data source

The frontend loads chart history from **DigitalOcean Spaces**: `https://nes.nyc3.digitaloceanspaces.com/chart-data.json`. A collector runs every 10 minutes (cron on a server) and updates that JSON; localStorage is used only as a fallback when the server is unreachable (e.g. CORS when viewing locally).
