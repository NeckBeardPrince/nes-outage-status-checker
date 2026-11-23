# Porch Light UI Refresh Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restyle and re-lay-out the three site pages with a shared stylesheet, a status sentence, a simpler toolbar, and three Easter eggs, without changing app behavior.

**Architecture:** `docs/styles.css` owns tokens, base, shell, and shared components and is linked after each page's inline `<style>`, so it wins on equal specificity. Each page drops its inline token/base/shell rules and keeps page-specific rules. `docs/ui.js` adds the status sentence (reads existing stat IDs through a `MutationObserver`) and popover behavior. `docs/eggs.js` holds the Easter eggs. The existing page JS stays untouched except for color literal swaps.

**Tech Stack:** Vanilla HTML/CSS/JS, Google Fonts, headless Chrome for screenshot checks.

**Spec:** `design/specs/2026-11-23-ui-refresh-design.md`

## Global Constraints

- Fonts: `Bricolage Grotesque` (display), `Atkinson Hyperlegible Next` (body), `JetBrains Mono` (mono).
- Dark tokens: `#1A1D3A` / `#232750` / `#2D3263`, text `#F4F1E8`, accent `#FFC857`. Light: `#FFF8EC` / `#FFFFFF` / `#F6EEDC`, text `#2A2440`, accent `#E0A21A`. Waiting `#FF7A6B`, active `#4FD1A5`.
- Every existing element `id` and `onclick`/`onchange` attribute must still exist after the change.
- The CSP `<meta>` stays the same (new scripts load from `'self'`).
- Same `localStorage` keys. New keys are allowed only for Easter egg state.
- Version 1.14.0 goes in `docs/version.json`, both `APP_VERSION` constants, and `CHANGELOG.md`.

## Review Focus

1. Light mode: every new component must read the light tokens (check the screenshot in light mode).
2. Narrow screens (375px): the header, stat grid, and toolbar wrap without horizontal scroll.
3. Keyboard shortcuts (R/N/T/M/S/F/H, C on monitor) still fire, and the cheat code listener must not swallow them.
4. Users with `prefers-reduced-motion` get no flicker or spotlight animation.
5. The status sentence must handle 0 outages ("Everyone's got power") and singular counts.

---

### Task 1: Shared stylesheet + index shell

**Files:** Create `docs/styles.css`, `docs/ui.js`. Modify `docs/index.html`.

- [ ] Write `styles.css`: tokens, base, `.site-header`/`.brand`/`.nav`, `.hero` + `.status-sentence`, `.stat-grid`/`.stat`, `.toolbar`, `.popover` (`<details>`), restyled `.icon-btn`, `select`, `.event-card`, `.status-badge`, `.modal*`, banners, footer.
- [ ] Remove inline index rules for tokens, `html`, `body*`, `.container`, header, summary bar, `.icon-btn`, `select`, footer. Swap the font `<link>`. Link `styles.css` after `</style>`, and `ui.js` and `eggs.js` with `defer`.
- [ ] Rebuild the header, hero, stat grid, and toolbar markup with all existing IDs preserved. Filters, sort, and refresh go inside `<details class="popover">`.
- [ ] Map old palette literals to new ones in CSS and JS (map markers, chart colors).
- [ ] `ui.js`: `buildStatusSentence(stats)` returns copy; observe the stat elements and write to `#statusSentence`; close open popovers on outside click or Esc.
- [ ] Verify: serve, screenshot dark + light at 1280px and 375px, check the console for errors, and run `grep` to confirm every original ID is still present.
- [ ] Commit.

### Task 2: Monitor + How It Works + all.html

- [ ] The same shell swap on `monitor.html`: stat IDs `crewsActive`, `customersBeingRestored`, `totalOutages`, `waitingForCrew`, `totalCustomersAffected`. The toolbar puts the data actions (Copy/Export/Import/Clear/Sort) in a "Manage" popover.
- [ ] `how-it-works.html` and `all.html`: tokens, fonts, header shell.
- [ ] Verify with screenshots and the ID check. Commit.

### Task 3: Easter eggs

- [ ] `eggs.js`: the cheat code (`ArrowUp ArrowUp ArrowDown ArrowDown ArrowLeft ArrowRight ArrowLeft ArrowRight b a`) toggles `body.lights-out` with a spotlight that follows the pointer through `--spot-x`/`--spot-y`. Esc exits. Five clicks on `#brandBulb` within 2s trigger `.bulb-pop` and a toast. A console ASCII bulb.
- [ ] Reduced motion: no flicker animation, a static spotlight.
- [ ] Verify: dispatch the key events in headless Chrome and confirm the class toggles. Commit.

### Task 4: Release

- [ ] Bump to 1.14.0, add the CHANGELOG entry, and update the `CLAUDE.md` structure list (styles.css, ui.js, eggs.js, mise.toml). Commit.
