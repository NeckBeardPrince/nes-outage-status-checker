# UI Refresh: "Porch Light" (v1.14.0)

## Intent

Make the site feel like a friendly neighborhood tool instead of a utility console. The same data and features, presented warmer and calmer, with a couple of Easter eggs.

## Goals

- New palette, fonts, and page shell shared by `index.html`, `monitor.html`, `how-it-works.html`.
- Shared stylesheet `docs/styles.css` so pages stop drifting.
- Summary bar becomes a plain-language "status sentence" plus big stat cards.
- Toolbar simplified: view tabs visible, filters/sort/refresh in a "Filters" popover, toggles as icon buttons.
- Three Easter eggs (below).

## Non-goals / anti-goals

- No framework, build step, or new JS dependency. Vanilla stays.
- No behavior changes to data fetching, monitoring, notifications, or `localStorage` keys.
- Must not break existing element IDs or `onclick` hooks; JS logic is untouched except for additive UI code.
- Must not loosen the CSP beyond `'self'` + Google Fonts for styles.

## Design

- **Dark (default):** background `#1A1D3A` / `#232750` / `#2D3263`, text `#F4F1E8`, accent porch-light `#FFC857`.
- **Light:** background `#FFF8EC` / `#FFFFFF` / `#F6EEDC`, text `#2A2440`, accent `#E0A21A`.
- **Status:** waiting = coral `#FF7A6B`, crew working = mint `#4FD1A5`.
- **Fonts:** Bricolage Grotesque (display), `Atkinson Hyperlegible` (body), JetBrains Mono (IDs).
- **Shape:** 16px radius cards, soft shadows, generous spacing.
- **Header:** light bulb logo + site name, nav pills, theme toggle.

## Easter eggs

1. The classic cheat code (`↑↑↓↓←→←→BA`) → "Lights out" mode: dark overlay with a flashlight spotlight following the cursor. Esc exits.
2. Click the logo bulb 5 times → flicker, pop, toast "Crew dispatched to your logo."
3. ASCII light bulb in the browser console.

Lives in `docs/eggs.js`, loaded from `'self'` (CSP-safe).

## Approaches considered

- **CSS-variable restyle only:** cheapest, but no new layout. Rejected.
- **Shared stylesheet + new shell, preserve IDs (chosen):** real refresh, JS stays intact.
- **Framework rewrite:** contradicts the vanilla stack. Rejected.

## Left to the builder

Exact spacing, responsive breakpoints, icon glyphs, and popover mechanics.

## Release

Bump `docs/version.json` and in-page `APP_VERSION` to 1.14.0; add CHANGELOG entry. Add `mise.toml` (`serve`, `dev`) for local preview.
