# Add line chart view

**Scope:** Client-only. Chart data is fetched from an external URL; no backend or collector changes in this PR.

## Summary

Adds a **Chart** view alongside Grid and Map so users can see trends over time. The view toggle is now three options: **Grid | Map | Chart**. Only one view is shown at a time (chart no longer appears alongside the map).

### Chart data source

The frontend fetches chart history from a JSON URL: `https://nes.nyc3.digitaloceanspaces.com/chart-data.json`. The app loads it on page load and uses it for the line chart. localStorage is used as a fallback when the fetch fails.

### Chart features

- **Metrics**: Affected Customers, Number of Outages, Active Crews, Customers Being Restored, Waiting for Crew (selectable in the chart dropdown).
- **Data**: Snapshots every 10 minutes (aligned to 10‑minute intervals). Data kept for the last 2 months.
- **UI**: X-axis labels include date and time. Tooltip on hover shows the value with wording that matches the selected metric (e.g. “12 crews”, “1,234 customers”).
- **Y-axis range slider**: Optional slider to set the chart’s Y-axis minimum (0 to max − 10%), so you can zoom into the upper part of the scale when values are clustered high. The chosen value is saved in localStorage.
