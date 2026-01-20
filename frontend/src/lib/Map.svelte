<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import maplibregl, { config } from "maplibre-gl";
    import "maplibre-gl/dist/maplibre-gl.css";

    let mapContainer: HTMLElement;
    let map: maplibregl.Map;
    let markers: maplibregl.Marker[] = [];

    const MAP_STYLES = {
        normal: "https://api.maptiler.com/maps/openstreetmap/style.json?key=fB2eDjoDg2nlel5Kw6ym",
        hybrid: "https://api.maptiler.com/maps/hybrid/style.json?key=aUOEn1bA48mz3xc3pL4N",
    };

    // Mock API call - Replace this with your actual URL
    async function fetchLocations() {
        // Simulating a fetch request
        // return await fetch('https://api.example.com/locations').then(res => res.json());
        return [
            {
                id: 1,
                name: "Blocker tengah",
                lng: 110.44053303318597,
                lat: -7.777491824518677,
                desc: "Drone Blocker Markas tengah",
                angleStart: 0,
                config: [
                    { signalCtrl: true, signalGPS: true },
                    { signalCtrl: true, signalGPS: true },
                    { signalCtrl: false, signalGPS: true },
                    { signalCtrl: true, signalGPS: true },
                    { signalCtrl: true, signalGPS: true },
                    { signalCtrl: true, signalGPS: false },
                ],
            },
            {
                id: 2,
                name: "Blocker random",
                lng: 110.4406,
                lat: -7.76,
                desc: "Drone Blocker Markas random",
                angleStart: 0,
                config: [
                    { signalCtrl: false, signalGPS: false },
                    { signalCtrl: false, signalGPS: false },
                    { signalCtrl: false, signalGPS: false },
                    { signalCtrl: true, signalGPS: true },
                    { signalCtrl: true, signalGPS: true },
                    { signalCtrl: true, signalGPS: false },
                ],
            },
        ];
    }

    function addMarkersToMap(data: any[]) {
        data.forEach((loc) => {
            const el = document.createElement("div");
            el.className = "marker-gps";
            el.style.cssText = `
            display: flex;
            align-items: center;
            justify-content: center;
            position: relative;
            `;

            // initial rotation offset
            const baseRotation = loc.angleStart || 0;

            for (let i = 0; i < 6; i++) {
                const angle = i * 60 + baseRotation;
                for (let layer = 0; layer < 2; layer++) {
                    const slice = document.createElement("div");

                    const size = layer === 1 ? 50 : 80;

                    slice.style.cssText = `
                    position: absolute;
                    width: ${size}px; height: ${size}px;
                    background-color: ${layer === 1 ? "green" : "yellow"};
                    border-radius: 50%;
                    border: 1px solid red;
                    pointer-events: none;
                    --angle: ${angle}deg;
                    top: 50%; left: 50%;
                    transform: translate(-50%, -50%) rotate(var(--angle));
                    clip-path: polygon(50% 50%, 100% 20%, 100% 80%);
                    animation: zoom-pulse 2s infinite ease-in-out;
                    `;

                    // hide slice based on config
                    if (loc.config[i].signalCtrl === false && layer === 0) {
                        slice.style.display = "none";
                    }

                    if (loc.config[i].signalGPS === false && layer === 1) {
                        slice.style.display = "none";
                    }

                    el.appendChild(slice);
                }
            }

            // Add the core
            const core = document.createElement("div");
            core.style.cssText =
                "width: 24px; height: 24px; background: red; border: 2px solid white; border-radius: 50%; z-index: 2;";
            el.appendChild(core);

            new maplibregl.Marker({ element: el })
                .setLngLat([loc.lng, loc.lat])
                .addTo(map);
        });
    }

    onMount(async () => {
        // 1. Initialize the map
        map = new maplibregl.Map({
            container: mapContainer,
            style: MAP_STYLES.normal,
            center: [110.44053927286228, -7.777395993083473],
            zoom: 12,
        });
        map.addControl(new maplibregl.NavigationControl(), "top-left");

        // 2. Wait for the map to load its style
        map.on("load", async () => {
            // 3. Fetch the data
            const locationData = await fetchLocations();

            // 4. Add markers
            addMarkersToMap(locationData);
        });
    });

    function switchStyle(styleKey: keyof typeof MAP_STYLES) {
        if (map) map.setStyle(MAP_STYLES[styleKey]);
    }

    onDestroy(() => {
        if (map) map.remove();
    });
</script>

<div class="map-layout">
    <div class="map-buttons">
        <button on:click={() => switchStyle("normal")}>Normal</button>
        <button on:click={() => switchStyle("hybrid")}>Satellite</button>
    </div>
    <div class="map-container" bind:this={mapContainer}></div>
</div>

<style>
    .map-layout {
        height: 100%;
        display: flex;
        flex-direction: column;
        overflow: hidden;
        border: 1px solid #ccc;
    }

    .map-buttons {
        margin-left: 50px;
        margin-top: 10px;
        display: flex;
        gap: 10px;
        background-color: transparent;
        position: absolute;
        z-index: 1;
    }

    button {
        padding: 6px 12px;
        border-radius: 6px;
        border: solid 1px #ccc;
        background-color: white;
    }

    button:active {
        background-color: gray;
    }

    .map-container {
        flex-grow: 1;
    }
</style>
