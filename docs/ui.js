// Porch Light UI helpers shared by every page.
// Additive only: reads the stat elements the page scripts already update and
// never calls into page code, so page behavior is unchanged.
(function () {
    'use strict';

    // Each page maps its own stat element IDs onto the same four numbers.
    var STAT_IDS = {
        index: { outages: 'totalCount', customers: 'totalCustomersAffected', crews: 'assignedCount', waiting: 'unassignedCount' },
        monitor: { outages: 'totalOutages', customers: 'totalCustomersAffected', crews: 'crewsActive', waiting: 'waitingForCrew' }
    };

    function readNumber(id) {
        var el = document.getElementById(id);
        if (!el) return null;
        var n = parseInt(el.textContent.replace(/[^0-9]/g, ''), 10);
        return Number.isFinite(n) ? n : null;
    }

    function plural(n, one, many) {
        return n === 1 ? one : many;
    }

    function esc(s) {
        return String(s).replace(/[&<>"']/g, function (c) {
            return { '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c];
        });
    }

    // Returns HTML for the hero sentence, or null while data is still loading.
    function buildStatusSentence(stats) {
        if (stats.outages === null || stats.customers === null) return null;
        if (stats.outages === 0) {
            return 'Everyone’s got power right now. <span class="hl-accent">Porch lights on</span> across Nashville.';
        }
        var fmt = function (n) { return esc(n.toLocaleString()); };
        var parts = 'Right now <span class="hl-accent">' + fmt(stats.customers) + '</span> ' +
            plural(stats.customers, 'neighbor is', 'neighbors are') + ' in the dark across ' +
            fmt(stats.outages) + ' ' + plural(stats.outages, 'outage', 'outages') + '.';
        if (stats.crews !== null && stats.crews > 0) {
            parts += ' <span class="hl-active">' + fmt(stats.crews) + ' ' + plural(stats.crews, 'crew is', 'crews are') + '</span> out working on it.';
        } else if (stats.crews === 0) {
            parts += ' Crews haven’t been assigned yet.';
        }
        return parts;
    }

    function initStatusSentence() {
        var target = document.getElementById('statusSentence');
        var page = document.body.dataset.page;
        var ids = STAT_IDS[page];
        if (!target || !ids) return;

        var pending = false;
        var update = function () {
            pending = false;
            var html = buildStatusSentence({
                outages: readNumber(ids.outages),
                customers: readNumber(ids.customers),
                crews: readNumber(ids.crews),
                waiting: readNumber(ids.waiting)
            });
            if (html !== null && target.innerHTML !== html) target.innerHTML = html;
        };
        var schedule = function () {
            if (!pending) { pending = true; requestAnimationFrame(update); }
        };

        var observer = new MutationObserver(schedule);
        Object.keys(ids).forEach(function (key) {
            var el = document.getElementById(ids[key]);
            if (el) observer.observe(el, { childList: true, characterData: true, subtree: true });
        });
        // Only treat the initial values as real data once the first load is done.
        var loading = document.getElementById('loading');
        if (!loading || loading.classList.contains('hidden')) schedule();
        else new MutationObserver(function (_, obs) {
            if (loading.classList.contains('hidden')) { obs.disconnect(); schedule(); }
        }).observe(loading, { attributes: true, attributeFilter: ['class'] });
    }

    function initPopovers() {
        var popovers = function () { return document.querySelectorAll('details.popover[open]'); };
        document.addEventListener('click', function (e) {
            popovers().forEach(function (d) { if (!d.contains(e.target)) d.removeAttribute('open'); });
        });
        document.addEventListener('keydown', function (e) {
            if (e.key !== 'Escape') return;
            popovers().forEach(function (d) {
                d.removeAttribute('open');
                var s = d.querySelector('summary');
                if (s) s.focus();
            });
        });
        // Opening one popover closes the others.
        document.querySelectorAll('details.popover').forEach(function (d) {
            d.addEventListener('toggle', function () {
                if (!d.open) return;
                popovers().forEach(function (o) { if (o !== d) o.removeAttribute('open'); });
            });
        });
    }

    // Shared toast used by Easter eggs; pages keep their own toasts.
    window.porchToast = function (message, ms) {
        var t = document.querySelector('.porch-toast');
        if (!t) {
            t = document.createElement('div');
            t.className = 'porch-toast';
            t.setAttribute('role', 'status');
            document.body.appendChild(t);
        }
        t.textContent = message;
        t.classList.add('show');
        clearTimeout(t._timer);
        t._timer = setTimeout(function () { t.classList.remove('show'); }, ms || 2600);
    };

    window.PorchLight = { buildStatusSentence: buildStatusSentence };

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', function () { initStatusSentence(); initPopovers(); });
    } else {
        initStatusSentence();
        initPopovers();
    }
})();
