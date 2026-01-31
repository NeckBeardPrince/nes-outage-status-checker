# Add line chart view (fetches data from DigitalOcean Spaces)

## Summary

Adds a **Chart** view alongside Grid and Map so users can see trends over time. The view toggle is now three options: **Grid | Map | Chart**. Only one view is shown at a time (chart no longer appears alongside the map).

### Chart data source

The frontend **fetches chart history from DigitalOcean Spaces**: `https://nes.nyc3.digitaloceanspaces.com/chart-data.json`. A collector (cron on a server) updates that JSON every 10 minutes. The app loads it on page load and uses it for the line chart. localStorage is used only as a fallback when the fetch fails (e.g. CORS when viewing locally).

### Chart features

- **Metrics**: Affected Customers, Number of Outages, Active Crews, Customers Being Restored, Waiting for Crew (selectable in the chart dropdown).
- **Data**: Snapshots every 10 minutes (aligned to 10‑minute intervals). Data kept for the last 2 months.
- **UI**: X-axis labels include date and time. Tooltip on hover shows the value with wording that matches the selected metric (e.g. “12 crews”, “1,234 customers”).
