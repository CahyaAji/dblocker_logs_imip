<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import maplibregl from "maplibre-gl";
    import "maplibre-gl/dist/maplibre-gl.css";
    import { settings } from "./configStore";

    let mapContainer: HTMLElement;
    let map: maplibregl.Map | undefined;
    let markers = new Map<number, maplibregl.Marker>();
    let intervalId: number;
    let resizeObserver: ResizeObserver;

    const MAP_STYLES = {
        normal: "https://api.maptiler.com/maps/openstreetmap/style.json?key=fB2eDjoDg2nlel5Kw6ym",
        hybrid: "https://api.maptiler.com/maps/hybrid/style.json?key=aUOEn1bA48mz3xc3pL4N",
    };

    async function fetchLocations() {
        // Mock API
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

    function updateMarkers(data: any[]) {
        if (!map) return;
        const incomingIds = new Set(data.map((loc) => loc.id));

        // Cleanup removed markers
        for (const [id, marker] of markers) {
            if (!incomingIds.has(id)) {
                marker.remove();
                markers.delete(id);
            }
        }

        // Add/Update markers
        data.forEach((loc) => {
            if (markers.has(loc.id)) {
                const existingMarker = markers.get(loc.id);
                existingMarker?.setLngLat([loc.lng, loc.lat]);
            } else {
                const el = createMarkerElement(loc);
                const newMarker = new maplibregl.Marker({ element: el })
                    .setLngLat([loc.lng, loc.lat])
                    .addTo(map!);
                markers.set(loc.id, newMarker);
            }
        });
    }

    function updatePixelScale() {
        if (!map || !mapContainer) return;

        const zoom = map.getZoom();
        const lat = map.getCenter().lat;

        // Math to convert meters to pixels at current zoom level
        const metersPerPixel =
            (156543.03392 * Math.cos((lat * Math.PI) / 180)) /
            Math.pow(2, zoom);

        // 1km Radius = 2000m Diameter
        const diameterInMeters = 2000;

        const pixels = diameterInMeters / metersPerPixel;

        // Update the CSS variable
        mapContainer.style.setProperty("--px-diameter", `${pixels}px`);
    }

    function createMarkerElement(loc: any) {
        const el = document.createElement("div");
        el.className = "marker-gps";
        const baseRotation = loc.angleStart || 0;

        for (let i = 0; i < 6; i++) {
            const angle = i * 60 + baseRotation;
            for (let layer = 0; layer < 2; layer++) {
                if (loc.config[i].signalCtrl === false && layer === 0) continue;
                if (loc.config[i].signalGPS === false && layer === 1) continue;

                // Create 2 ripples per sector for the overlapping effect
                for (let ripple = 0; ripple < 2; ripple++) {
                    const slice = document.createElement("div");
                    slice.className = "radar-slice";

                    slice.style.setProperty("--angle", `${angle}deg`);
                    slice.style.setProperty(
                        "--color",
                        layer === 1 ? "darkgreen" : "yellow",
                    );

                    const scaleWrapper = layer === 1 ? 0.6 : 1.0;
                    slice.style.setProperty(
                        "--scale-factor",
                        `${scaleWrapper}`,
                    );

                    // Delay the second ripple by -1s so it starts halfway through
                    slice.style.animationDelay = `${ripple * -1}s`;

                    el.appendChild(slice);
                }
            }
        }

        const core = document.createElement("div");
        core.className = "marker-core";
        el.appendChild(core);

        return el;
    }

    let debounceTimer: ReturnType<typeof setTimeout>;

    function debounceRender(data: any[]) {
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
            updateMarkers(data);
        }, 100);
    }

    function switchStyle(styleKey: "normal" | "hybrid") {
        if (!map) return;
        $settings.mapStyle = styleKey; // ✨ Update Store
        map.setStyle(MAP_STYLES[styleKey]);
    }

    onMount(async () => {
        // ✨ Initialize using Store value
        map = new maplibregl.Map({
            container: mapContainer,
            style: MAP_STYLES[$settings.mapStyle],
            center: [110.44053927286228, -7.777395993083473],
            zoom: 14,
        });
        map.addControl(new maplibregl.NavigationControl(), "top-left");

        // ✨ Resize Observer: Watches for sidebar changes
        resizeObserver = new ResizeObserver(() => {
            map?.resize();
        });
        resizeObserver.observe(mapContainer);

        map.on("load", async () => {
            updatePixelScale();

            if (map) {
                map.on("move", updatePixelScale);
                map.on("zoom", updatePixelScale);
            }

            const initialData = await fetchLocations();
            updateMarkers(initialData);

            intervalId = setInterval(async () => {
                const data = await fetchLocations();
                debounceRender(data);
            }, 1000);
        });
    });

    onDestroy(() => {
        clearInterval(intervalId);
        resizeObserver?.disconnect(); // ✨ Clean up observer

        markers.forEach((m) => m.remove());
        markers.clear();

        if (map) {
            map.off("move", updatePixelScale);
            map.off("zoom", updatePixelScale);
            map.remove();
        }
    });
</script>

<div class="map-layout">
    <div class="map-buttons">
        <button
            class:active={$settings.mapStyle === "normal"}
            on:click={() => switchStyle("normal")}>Normal</button
        >
        <button
            class:active={$settings.mapStyle === "hybrid"}
            on:click={() => switchStyle("hybrid")}>Satellite</button
        >
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
        position: absolute;
        z-index: 2;
    }

    button {
        padding: 6px 12px;
        border-radius: 6px;
        border: solid 1px #ccc;
        background-color: white;
        cursor: pointer;
    }

    button.active {
        background-color: #333;
        color: white;
        border-color: #333;
    }

    .map-container {
        flex-grow: 1;
    }

    .map-layout :global(.marker-gps) {
        display: flex;
        align-items: center;
        justify-content: center;
        position: relative;
    }

    .map-layout :global(.marker-core) {
        width: 18px;
        height: 18px;
        background: red;
        border: 2px solid white;
        border-radius: 50%;
        z-index: 2;
        position: relative;
    }

    .map-layout :global(.radar-slice) {
        position: absolute;
        width: calc(var(--px-diameter) * var(--scale-factor));
        height: calc(var(--px-diameter) * var(--scale-factor));
        background-color: var(--color);

        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%) rotate(var(--angle));
        clip-path: polygon(50% 50%, 100% 20%, 100% 80%);
        border-radius: 50%;

        pointer-events: none;

        /* Linear makes the ripple overlap smoothly */
        animation: zoom-pulse 2s infinite linear;
    }

    @keyframes zoom-pulse {
        0% {
            transform: translate(-50%, -50%) rotate(var(--angle)) scale(0);
            opacity: 0.8;
        }
        30% {
            transform: translate(-50%, -50%) rotate(var(--angle)) scale(0.3);
            opacity: 1;
        }

        100% {
            transform: translate(-50%, -50%) rotate(var(--angle)) scale(1);
            opacity: 0;
        }
    }
</style>
