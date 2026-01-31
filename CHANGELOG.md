# Changelog

All notable changes to the NES Outage Status Checker will be documented in this file.

## [1.11.0] - 2026-01-31

*Thank you to [cgoliver](https://github.com/cgoliver) for [PR #9](https://github.com/NeckBeardPrince/nes-outage-status-checker/pull/9)!*

### Added

#### Web App - All Outages Page
- **Chart view** - New view option alongside Grid and Map to visualize trends over time
- **Multiple chart metrics** - Select from Affected Customers, Number of Outages, Active Crews, Customers Being Restored, or Waiting for Crew
- **Historical data collection** - Snapshots stored every 10 minutes with 2-month retention
- **Interactive tooltips** - Hover to see values with metric-specific formatting

### Removed

- **CLI tool deprecated** - The command-line interface has been removed. The web app at [nes-outage-checker.com](https://nes-outage-checker.com) is now the primary way to use this tool.

---

## [1.10.0] - 2026-01-30

*Thank you to [nico22nguyen](https://github.com/nico22nguyen) for [PR #8](https://github.com/NeckBeardPrince/nes-outage-status-checker/pull/8)!*

### Changed

#### Web App - All Outages Page
- **Optimized API interactions** - Filter and sort changes no longer trigger API calls, using cached data instead for faster response on cellular connections

---

## [1.9.0] - 2026-01-30

### Added

#### Web App - All Outages Page
- **Change History modal** - View a log of all changes including new outages, resolved outages, and status changes
- **H keyboard shortcut** - Quick access to History modal

#### Web App - Footer
- **Community acknowledgment** - Added link to GitHub contributors page

---

## [1.8.0] - 2026-01-30

*Thank you to [coleman8er](https://github.com/coleman8er) for [PR #7](https://github.com/NeckBeardPrince/nes-outage-status-checker/pull/7)!*

### Added

#### Web App - Both Pages
- **Data Freshness indicator** - Shows timestamp of last API check and change summary (+X new, -X resolved, X status changes) highlighted when changes occur

#### Web App - Monitor Page
- **System-wide stats in header** - Added Total Outages, Waiting for Crew, and Affected Customers counts to match the All Outages page
- **My Monitors section header** - New section header showing monitor count and total affected customers for tracked events

#### Web App - All Outages Page
- **Sort by Recently Updated** - New sort option to order outages by most recently updated first
- **View on Map modal** - Click "View on Map" on any outage card to open a modal with a larger map showing the outage location, with status-colored border (green for assigned, amber for unassigned) and reverse geocoding to display a human-readable address

---

## [1.7.0] - 2026-01-30

*Thank you to [jjstafford](https://github.com/jjstafford) for [PR #6](https://github.com/NeckBeardPrince/nes-outage-status-checker/pull/6)!*

### Added

#### Web App - Both Pages
- **Map style toggle** - Cycle between Dark, Light, and Aerial (satellite) map views
- **Fullscreen mode** - Expand maps to fullscreen for better visibility
- **Keyboard shortcuts** - S to cycle map style, F for fullscreen
- **Responsive controls** - Map buttons show icons on mobile, text on desktop
- **Style persistence** - Map style preference saved to localStorage

---

## [1.6.0] - 2026-01-30

*Thank you to [nico22nguyen](https://github.com/nico22nguyen) for [PR #5](https://github.com/NeckBeardPrince/nes-outage-status-checker/pull/5)!*

### Added

#### Web App - All Outages Page
- **Affected Customers count** - New stat showing total customers affected across all outages

#### Web App - Both Pages
- **NES transparency warning banner** - Dismissible banner alerting users that NES may have stopped updating crew assignment status

---

## [1.5.0] - 2026-01-29

*Thank you to [coleman8er](https://github.com/coleman8er) for [PR #3](https://github.com/NeckBeardPrince/nes-outage-status-checker/pull/3) and the new features below!*

### Added

#### Web App - Monitor Page
- **Find by Address** - Enter your street address to automatically find nearby outages instead of manually looking up Event IDs
- **Smart matching algorithm** - Weights results by distance AND outage size (larger outages affecting more customers cover wider areas)
- **Top 3 matches** - Shows the 3 best matches with confidence indicators (Best Match, Good Match, Possible)
- **Help modal** - Explains how the address lookup works and why it might occasionally suggest the wrong event
- **Privacy-focused** - Address is sent to OpenStreetMap for geocoding only, no data stored
- **Crews Near Me map** - View NES crews working near your monitored outages

### Changed
- **"Calm Utility" UI refresh** - New warm charcoal color palette, refined typography with DM Serif Display + Outfit fonts, subtle animations and ambient visual effects

---

## [1.4.0] - 2026-01-29

### Added

#### Web App - Both Pages
- **What's New modal** - Popup showing changelog when app updates to a new version
- **Status filter** - Filter events by All, Assigned, or Unassigned status

#### Web App - All Outages Page
- **Sort options** - Sort by Status (assigned first), Oldest (longest duration), or Most Affected
- **Higher filter thresholds** - Added 250+, 500+, 750+, and 1000+ affected filters

---

## [1.3.0] - 2026-01-29

### Added

#### Web App - All Outages Page
- **Map view** - Toggle between grid and map view to see all outages plotted geographically
- **Color-coded markers** - Green for assigned, red for unassigned outages
- **Marker size by impact** - Larger markers indicate more people affected
- **Interactive popups** - Click markers to see outage details
- **Legend** - Visual guide for status colors
- **M keyboard shortcut** - Quick toggle between grid and map views

### Changed
- Moved all controls (refresh, notifications, sound) to top summary bar on both pages

---

## [1.2.0] - 2026-01-28

### Added

#### Web App - Both Pages
- **Light/dark mode toggle** - Switch between dark and light themes with the Theme button (or press T)
- **Colorblind-friendly mode** - Accessible patterns and icons for color vision deficiencies
- **Offline detection** - Banner notification when connection is lost
- **Keyboard shortcuts** - R=Refresh, N=Notifications, T=Theme (C=Copy URL on monitor page)
- **Countdown timer** - Visual countdown to next auto-refresh

#### Web App - Monitor Page
- **Export monitors** - Download monitor configuration as JSON backup
- **Import monitors** - Restore monitors from a JSON file
- **Clear all button** - Remove all monitors with one click
- **Total outage duration** - Combined duration stat across all active monitors

---

## [1.1.0] - 2025-01-28

### Added

#### Web App - Monitor Page
- **Multi-monitor support** - Track up to 4 outage events simultaneously with individual status cards
- **Custom labels** - Name each monitor (e.g., "Home", "Office", "Mom's house") for easy identification
- **Shareable URLs** - Monitor configurations saved in URL (`?monitors=123:Home,456:Office`) for easy sharing
- **Browser notifications** - Get notified when outage status changes or event is resolved (requires permission)
- **Sound alerts** - Optional audio chime when status changes
- **Auto-refresh** - Configurable refresh intervals (1, 5, 10, 15, 30, or 60 minutes)
- **Compact view** - Toggle to show only status badges for a condensed overview
- **Sort by status** - Option to show assigned monitors first
- **Status history** - Track status changes over time for each monitor (last 10 entries)
- **Outage duration** - Display how long each outage has been ongoing
- **Summary bar** - Shows total monitor count and combined affected people
- **Copy URL button** - One-click sharing with toast confirmation
- **Persistent settings** - Compact view, sort preference, and sound settings saved in localStorage
- **Persistent monitors** - Monitor data and labels saved in localStorage
- **Responsive grid layout** - Monitors display side-by-side on wider screens
- **Favicon** - Lightning bolt emoji for easy tab identification
- **Dark theme** - Easy on the eyes during those late-night outages
- **Auto-update detection** - Banner notification when a new version is available

#### Web App - All Outages Page
- **All outages dashboard** - View every active NES outage event in one place
- **Filter by affected** - Filter events by minimum affected people (All, 2+, 10+, 50+, 100+) - defaults to 10+
- **Summary statistics** - Total events, unassigned count, assigned count, and total affected people
- **Browser notifications** - Get notified when any event's status changes or is resolved
- **Sound alerts** - Optional audio chime when status changes
- **Auto-refresh** - Configurable refresh intervals (10, 20, 30, 40, 50, or 60 minutes)
- **Compact view** - Toggle to show condensed event cards
- **Sort by status** - Option to show assigned events first (default on)
- **Responsive grid layout** - Events display in a grid on wider screens

#### CLI
- **Makefile** - Quick build commands (`make build`, `make clean`, `make fmt`, `make vet`)

#### Infrastructure
- **GitHub Pages support** - Static site hosted from `/docs` folder
- **GitHub Actions CI** - Automated build and vet checks on push/PR
- **GitHub Actions Pages deployment** - Automatic deployment when docs change

### Changed
- Renamed binary from `nes-status-checker` to `nes-outage-status-checker`
- Updated README with web app features and installation instructions

### Removed
- View on Map button (non-functional)

---

## [1.0.0] - Initial Release

### Added
- Terminal UI application using Bubble Tea and Lip Gloss
- Real-time outage status monitoring
- Color-coded status display (red for unassigned, flashing green for assigned)
- Auto-refresh every 30 seconds
- Manual refresh with `r` key
- Displays outage details: affected customers, start time, last update, cause

---

*Built with Claude AI*
