# Changelog

All notable changes to the NES Outage Status Checker will be documented in this file.

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
