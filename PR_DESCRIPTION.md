# Add line chart view for outage metrics

## Summary

Adds a **Chart** view alongside Grid and Map so users can see trends over time. The view toggle is now three options: **Grid | Map | Chart**. Only one view is shown at a time (chart no longer appears alongside the map).

### Chart features

- **Metrics**: Affected Customers, Number of Outages, Active Crews, Customers Being Restored, Waiting for Crew (selectable in the chart dropdown).
- **Data**: Snapshots are stored every 10 minutes (aligned to 10‑minute intervals) when the page fetches outage data. Data is kept for the last 2 months to limit storage size.
- **UI**: X-axis labels include date and time. Tooltip on hover shows the value with wording that matches the selected metric (e.g. “12 crews”, “1,234 customers”).

Data is stored in the browser’s `localStorage` under the key `nes-outage-chart-data`. **This implies the server will need to store the JSON that collects these data snapshots** if you want chart history to persist across devices or after clearing site data (e.g. a backend endpoint that saves/returns the same structure).

Thank you to the developers for building and maintaining this project.
