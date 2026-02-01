/**
 * shared-data.js - Shared utilities for NES Outage Checker
 * Handles zip code mapping, reverse geocoding, and area statistics
 */

// Nashville neighborhood zip code mappings
const NASHVILLE_ZIP_CODES = {
    '37201': { name: 'Downtown/Capitol Hill', lat: 36.1627, lon: -86.7816 },
    '37202': { name: 'North Nashville', lat: 36.1950, lon: -86.7810 },
    '37203': { name: 'East Nashville', lat: 36.1633, lon: -86.7520 },
    '37204': { name: 'Germantown/Salemtown', lat: 36.1784, lon: -86.7644 },
    '37205': { name: 'Sylvan Park/West End', lat: 36.1450, lon: -86.8200 },
    '37206': { name: 'Shelby Park/Weaver Park', lat: 36.1533, lon: -86.7180 },
    '37207': { name: 'Inglewood/Parkwood', lat: 36.1317, lon: -86.7420 },
    '37208': { name: 'North Nashville', lat: 36.2100, lon: -86.8050 },
    '37209': { name: 'Hermitage', lat: 36.1083, lon: -86.6580 },
    '37210': { name: 'Antioch', lat: 36.0233, lon: -86.7050 },
    '37211': { name: 'Brentwood', lat: 35.9667, lon: -86.7833 },
    '37212': { name: 'Belmont/The Nations', lat: 36.1350, lon: -86.8550 },
    '37214': { name: 'Southeast Nashville', lat: 36.0733, lon: -86.7050 },
    '37215': { name: 'Belle Meade', lat: 36.1533, lon: -86.9050 },
    '37216': { name: 'East Nashville', lat: 36.1733, lon: -86.7300 },
    '37217': { name: 'Riverside', lat: 36.0933, lon: -86.8700 },
    '37218': { name: 'MetroCenter', lat: 36.1650, lon: -86.8350 },
    '37219': { name: 'Downtown', lat: 36.1600, lon: -86.7750 },
    '37220': { name: 'Green Hills/Buena Vista', lat: 36.1117, lon: -86.8033 },
    '37221': { name: 'Hendersonville', lat: 36.3050, lon: -86.6250 },
    '37222': { name: 'Smyrna/Lavergne', lat: 35.9933, lon: -86.5883 },
    '37224': { name: 'Murfreesboro Pike', lat: 36.0550, lon: -86.6750 },
    '37228': { name: 'Airport/Berry Hill', lat: 36.1250, lon: -86.6880 },
    '37229': { name: 'Goodlettsville', lat: 36.3167, lon: -86.6950 },
    '37230': { name: 'Hermitage/Donelson', lat: 36.0933, lon: -86.6580 },
    '37231': { name: 'Antioch', lat: 36.0433, lon: -86.6980 },
    '37232': { name: 'Madison', lat: 36.1933, lon: -86.7333 },
    '37235': { name: 'Glencliff', lat: 36.0700, lon: -86.8150 },
    '37238': { name: 'Downtown', lat: 36.1600, lon: -86.7700 }
};

/**
 * Convert coordinates to zip code using Nominatim reverse geocoding
 * @param {number} lat - Latitude
 * @param {number} lon - Longitude
 * @returns {Promise<string>} Zip code string
 */
async function reverseGeocodeToZip(lat, lon) {
    try {
        const response = await fetch(
            `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lon}`,
            {
                headers: {
                    'User-Agent': 'NES-Outage-Checker/1.0'
                }
            }
        );
        
        if (!response.ok) return null;
        
        const data = await response.json();
        const postalCode = data.address?.postcode;
        
        return postalCode || null;
    } catch (err) {
        console.error('Reverse geocoding error:', err);
        return null;
    }
}

/**
 * Get neighborhood name from zip code
 * @param {string} zipCode - 5-digit zip code
 * @returns {string} Neighborhood name or zip code if not found
 */
function getNeighborhoodName(zipCode) {
    const zip = zipCode?.toString().substring(0, 5);
    return NASHVILLE_ZIP_CODES[zip]?.name || zip || 'Unknown Area';
}

/**
 * Find closest zip code from coordinates
 * @param {number} lat - Latitude
 * @param {number} lon - Longitude
 * @returns {string} Closest zip code
 */
function findClosestZipCode(lat, lon) {
    let closest = null;
    let minDistance = Infinity;
    
    Object.entries(NASHVILLE_ZIP_CODES).forEach(([zip, data]) => {
        const dist = Math.sqrt(
            Math.pow(data.lat - lat, 2) + Math.pow(data.lon - lon, 2)
        );
        if (dist < minDistance) {
            minDistance = dist;
            closest = zip;
        }
    });
    
    return closest;
}

/**
 * Extract or generate zip code from event data
 * Uses reverse geocoding if available, falls back to closest match
 * @param {Object} event - Event object with lat/lon
 * @returns {Promise<string>} Zip code
 */
async function getZipCodeForEvent(event) {
    if (!event) return null;
    
    // If event already has zip code, use it
    if (event.zipCode) return event.zipCode;
    
    if (!event.latitude || !event.longitude) return null;
    
    // Try reverse geocoding first (cached in localStorage to reduce API calls)
    const geocodingCache = JSON.parse(localStorage.getItem('nes-geocoding-cache') || '{}');
    const cacheKey = `${event.latitude.toFixed(4)},${event.longitude.toFixed(4)}`;
    
    if (geocodingCache[cacheKey]) {
        return geocodingCache[cacheKey];
    }
    
    // Rate-limited reverse geocoding (max 1 req/sec per Nominatim ToS)
    const reverseGeoCode = await reverseGeocodeToZip(event.latitude, event.longitude);
    if (reverseGeoCode) {
        geocodingCache[cacheKey] = reverseGeoCode;
        localStorage.setItem('nes-geocoding-cache', JSON.stringify(geocodingCache));
        return reverseGeoCode;
    }
    
    // Fall back to finding closest zip code
    const closest = findClosestZipCode(event.latitude, event.longitude);
    geocodingCache[cacheKey] = closest;
    localStorage.setItem('nes-geocoding-cache', JSON.stringify(geocodingCache));
    return closest;
}

/**
 * Calculate reliability score for an area
 * @param {Object} areaStats - Area statistics object
 * @returns {number} Reliability score (0-100, higher is better)
 */
function calculateReliabilityScore(areaStats) {
    if (!areaStats.outages) return 100;
    
    // Score based on frequency and duration
    const frequencyPenalty = Math.min(areaStats.outages * 5, 50);
    const durationPenalty = Math.min(areaStats.avgDuration * 2, 30);
    const impactPenalty = Math.min((areaStats.totalAffected / 100), 20);
    
    const score = Math.max(0, 100 - frequencyPenalty - durationPenalty - impactPenalty);
    return Math.round(score);
}

/**
 * Compare area to city average
 * @param {Object} areaStats - Area statistics
 * @param {Array} allAreas - All areas for comparison
 * @returns {Object} Comparison metrics
 */
function compareToAverage(areaStats, allAreas) {
    if (allAreas.length === 0) return { factor: 1, rating: 'average' };
    
    const avgOutages = allAreas.reduce((sum, a) => sum + a.outages, 0) / allAreas.length;
    const avgDuration = allAreas.reduce((sum, a) => sum + a.avgDuration, 0) / allAreas.length;
    
    const factor = areaStats.outages / avgOutages;
    let rating = 'average';
    
    if (factor < 0.5) rating = 'excellent';
    else if (factor < 0.8) rating = 'above-average';
    else if (factor > 1.5) rating = 'poor';
    else if (factor > 1.2) rating = 'below-average';
    
    return {
        factor: factor.toFixed(1),
        rating,
        percentageDiff: Math.round((factor - 1) * 100)
    };
}

/**
 * Export area history to CSV
 * @param {Array} historyData - History array
 * @param {string} zipCode - Optional zip code filter
 * @returns {string} CSV content
 */
function exportToCSV(historyData, zipCode = null) {
    const headers = ['Date', 'Time', 'Duration (hrs)', 'People Affected', 'Status', 'Area', 'Zip Code'];
    const rows = [];
    
    historyData.forEach(entry => {
        if (zipCode && entry.zipCode !== zipCode) return;
        
        const date = new Date(entry.startTime);
        const duration = ((entry.lastUpdatedTime - entry.startTime) / (1000 * 60 * 60)).toFixed(2);
        
        rows.push([
            date.toLocaleDateString(),
            date.toLocaleTimeString(),
            duration,
            entry.numPeople,
            entry.status || 'unknown',
            getNeighborhoodName(entry.zipCode),
            entry.zipCode || 'N/A'
        ]);
    });
    
    const csv = [
        headers.join(','),
        ...rows.map(row => row.map(cell => `"${cell}"`).join(','))
    ].join('\n');
    
    return csv;
}

/**
 * Get hourly distribution of outages (for time-of-day analysis)
 * @param {Array} historyData - History array
 * @returns {Object} Hour -> count mapping
 */
function getHourlyDistribution(historyData) {
    const hourMap = {};
    
    for (let hour = 0; hour < 24; hour++) {
        hourMap[hour] = 0;
    }
    
    historyData.forEach(entry => {
        const date = new Date(entry.startTime);
        const hour = date.getHours();
        hourMap[hour]++;
    });
    
    return hourMap;
}

/**
 * Get monthly summary for reports
 * @param {Array} historyData - History array
 * @returns {Object} Month -> stats mapping
 */
function getMonthlySummary(historyData) {
    const monthMap = {};
    
    historyData.forEach(entry => {
        const date = new Date(entry.startTime);
        const monthKey = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`;
        
        if (!monthMap[monthKey]) {
            monthMap[monthKey] = {
                outages: 0,
                totalDuration: 0,
                totalAffected: 0,
                incidents: []
            };
        }
        
        const duration = (entry.lastUpdatedTime - entry.startTime) / (1000 * 60 * 60);
        monthMap[monthKey].outages++;
        monthMap[monthKey].totalDuration += duration;
        monthMap[monthKey].totalAffected += entry.numPeople;
        monthMap[monthKey].incidents.push(entry);
    });
    
    // Calculate averages
    Object.values(monthMap).forEach(month => {
        month.avgDuration = (month.totalDuration / month.outages).toFixed(2);
        month.avgAffected = Math.round(month.totalAffected / month.outages);
    });
    
    return monthMap;
}

/**
 * Find worst month (by various metrics)
 * @param {Object} monthlySummary - Monthly summary object from getMonthlySummary()
 * @returns {Object} Worst month data with metric
 */
function findWorstMonth(monthlySummary) {
    const months = Object.entries(monthlySummary);
    if (months.length === 0) return null;
    
    const worstByOutages = months.reduce((prev, curr) => 
        curr[1].outages > prev[1].outages ? curr : prev
    );
    
    const worstByDuration = months.reduce((prev, curr) => 
        curr[1].totalDuration > prev[1].totalDuration ? curr : prev
    );
    
    const worstByImpact = months.reduce((prev, curr) => 
        curr[1].totalAffected > prev[1].totalAffected ? curr : prev
    );
    
    return {
        outages: worstByOutages,
        duration: worstByDuration,
        impact: worstByImpact
    };
}

// Export for use in other files
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        NASHVILLE_ZIP_CODES,
        reverseGeocodeToZip,
        getNeighborhoodName,
        findClosestZipCode,
        getZipCodeForEvent,
        calculateReliabilityScore,
        compareToAverage,
        exportToCSV,
        getHourlyDistribution,
        getMonthlySummary,
        findWorstMonth
    };
}
