import { writable } from "svelte/store";

export interface DBlockerConfig {
    signalCtrl: boolean;
    signalGPS: boolean;
}

export interface DBlocker {
    id: number;
    name: string;
    lat: number;
    lon: number;
    serialNumb: string;
    desc: string;
    config: DBlockerConfig[];
}

// CREATE THE STORE ---
// This holds the current state of all blockers
export const dblockerStore = writable<DBlocker[]>([]);

let pollingInteral: number | undefined;

const API_URL = "http://localhost:3003/api/dblockers";

export async function fetchDBlockers() {
    try {

    } catch (err) {
        console.error(err);
    }

}