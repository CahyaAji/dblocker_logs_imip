import { writable } from 'svelte/store';

export interface DBlockerConfig {
    signalCtrl: boolean;
    signalGPS: boolean;
}

export interface DBlocker {
    id: number;
    name: string;
    lat: number;
    lng: number;
    desc: string;
    angleStart: number;
    config: DBlockerConfig[]; 
}

// --- STORE ---
export const dblockerStore = writable<DBlocker[]>([]);

// --- CONFIG ---
const API_BASE = "http://localhost:3003/api";
let pollingInterval: number | undefined;


// --- READ DATA (GET) ---
export async function fetchDBlockers() {
    try {
        const res = await fetch(`${API_BASE}/dblockers`);
        if (!res.ok) throw new Error("Fetch dblockers failed");
        
        const data: DBlocker[] = await res.json();
        dblockerStore.set(data);
    } catch (err) {
        console.error("Polling Error:", err);
    }
}

export function startPolling(intervalMs = 1000) {
    fetchDBlockers();
    stopPolling();
    pollingInterval = setInterval(fetchDBlockers, intervalMs);
}

export function stopPolling() {
    if (pollingInterval) clearInterval(pollingInterval);
}


// CHANGE DATA (POST)
// This is the function you call when user clicks a button
export async function switchSignal(
    blockerId: number, 
    sectorIdx: number, 
    type: 'signalCtrl' | 'signalGPS', 
    newValue: boolean
) {
    // A. Optimistic Update: Update UI *immediately* so it feels fast
    dblockerStore.update(items => items.map(b => {
        if (b.id !== blockerId) return b;
        
        // Create deep copy of config to trigger Svelte update
        const newConfig = [...b.config];
        newConfig[sectorIdx] = { ...newConfig[sectorIdx], [type]: newValue };
        
        return { ...b, config: newConfig };
    }));

    // B. Send Request to Server
    try {
        const payload = {
            id: blockerId,
            sector: sectorIdx,
            type: type,   // "signalCtrl" or "signalGPS"
            value: newValue // true or false
        };

        const res = await fetch(`${API_BASE}/dblockers/switch`, {
            method: 'POST', // or 'PUT'
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!res.ok) throw new Error("Update failed");

        // Optional: If server returns the new full object, update store again here
        // const updatedBlocker = await res.json();
        // dblockerStore.update(...)

    } catch (err) {
        console.error("Failed to switch signal:", err);
        
        // C. Rollback: If server failed, flip the switch back!
        dblockerStore.update(items => items.map(b => {
            if (b.id !== blockerId) return b;
            const newConfig = [...b.config];
            newConfig[sectorIdx] = { ...newConfig[sectorIdx], [type]: !newValue }; // Revert
            return { ...b, config: newConfig };
        }));
        
        alert("Failed to update signal. Check connection.");
    }
}