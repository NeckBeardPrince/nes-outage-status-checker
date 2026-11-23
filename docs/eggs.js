// Nothing to see here. Move along. 💡
(function () {
    'use strict';

    var reduceMotion = window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches;

    function toast(msg) {
        if (window.porchToast) window.porchToast(msg);
    }

    // ---- 1. Lights out: ↑ ↑ ↓ ↓ ← → ← → B A ----
    var CODE = ['arrowup', 'arrowup', 'arrowdown', 'arrowdown', 'arrowleft', 'arrowright', 'arrowleft', 'arrowright', 'b', 'a'];
    var progress = 0;
    var hint = null;

    function setSpot(x, y) {
        document.body.style.setProperty('--spot-x', x + 'px');
        document.body.style.setProperty('--spot-y', y + 'px');
    }

    function onPointer(e) {
        var p = e.touches ? e.touches[0] : e;
        setSpot(p.clientX, p.clientY);
    }

    function lightsOut(on) {
        var body = document.body;
        if (on === body.classList.contains('lights-out')) return;
        body.classList.toggle('lights-out', on);
        if (on) {
            setSpot(window.innerWidth / 2, window.innerHeight / 2);
            if (!reduceMotion) {
                window.addEventListener('pointermove', onPointer, { passive: true });
                window.addEventListener('touchmove', onPointer, { passive: true });
            }
            hint = document.createElement('div');
            hint.className = 'lights-out-hint';
            hint.textContent = 'Lights out! Now you know how it feels. Press Esc to restore power.';
            body.appendChild(hint);
        } else {
            window.removeEventListener('pointermove', onPointer);
            window.removeEventListener('touchmove', onPointer);
            if (hint) { hint.remove(); hint = null; }
            toast('Power restored. Crew of one: you. ⚡');
        }
    }

    document.addEventListener('keydown', function (e) {
        // Never interfere with typing in fields or the page's own shortcuts.
        var tag = (e.target && e.target.tagName) || '';
        if (/^(INPUT|TEXTAREA|SELECT)$/.test(tag) || e.target.isContentEditable) return;

        if (e.key === 'Escape' && document.body.classList.contains('lights-out')) {
            lightsOut(false);
            return;
        }
        var key = (e.key || '').toLowerCase();
        if (key === CODE[progress]) {
            progress++;
            if (progress === CODE.length) {
                progress = 0;
                lightsOut(!document.body.classList.contains('lights-out'));
            }
        } else {
            progress = key === CODE[0] ? 1 : 0;
        }
    });

    // ---- 2. Five quick clicks on the logo bulb ----
    var clicks = [];
    var busy = false;
    var LINES = [
        'Crew dispatched to your logo. ETA: 2–4 business seconds.',
        'Bulb replaced. NES would like a word about your clicking.',
        'Outage reported: 1 bulb affected. Crew assigned. 🔧'
    ];
    var lineIdx = 0;

    function initBulb() {
        var bulb = document.getElementById('brandBulb');
        if (!bulb) return;
        bulb.addEventListener('click', function () {
            if (busy) return;
            var now = Date.now();
            clicks = clicks.filter(function (t) { return now - t < 2000; });
            clicks.push(now);
            if (clicks.length < 5) return;
            clicks = [];
            busy = true;
            bulb.classList.add(reduceMotion ? 'popped' : 'flicker');
            setTimeout(function () {
                bulb.classList.remove('flicker');
                bulb.classList.add('popped');
                toast(LINES[lineIdx++ % LINES.length]);
                setTimeout(function () {
                    bulb.classList.remove('popped');
                    busy = false;
                }, 2200);
            }, reduceMotion ? 0 : 900);
        });
    }

    if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', initBulb);
    else initBulb();

    // ---- 3. For the curious who open DevTools ----
    var art = [
        '      .-.',
        '     (   )   Hey, neighbor! 👋',
        '      \\ /',
        '      |=|    You found the console. This site is open source:',
        '      |=|    https://github.com/NeckBeardPrince/nes-outage-status-checker',
        '      \'-\'',
        '             Psst. Try ↑ ↑ ↓ ↓ ← → ← → B A on the page.'
    ].join('\n');
    try {
        console.log('%c' + art, 'color:#FFC857;font-family:monospace;font-size:12px;line-height:1.4');
    } catch (_) { /* console unavailable */ }
})();
