<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import maplibregl from "maplibre-gl";
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
                lng: 110.44053927286228,
                lat: -7.777395993083473,
                desc: "Drone Blocker Markas tengah",
                isActive: true,
                ringDegree: 360,
            },
        ];
    }

    function addMarkersToMap(data: any[]) {
        data.forEach((loc) => {
            const el = document.createElement("div");
            el.className = "marker-gps";

            // 1. Pass the degree as a CSS Variable
            el.style.setProperty("--ring-angle", `${loc.ringDegree || 360}deg`);

            // 2. Handle the Active state (either via class or display)
            if (loc.isActive === false) {
                el.classList.add("inactive");
            }

            const popup = new maplibregl.Popup({ offset: 15 }).setHTML(
                `<strong>${loc.name}</strong><p>${loc.desc}</p>`,
            );

            const marker = new maplibregl.Marker({ element: el })
                .setLngLat([loc.lng, loc.lat])
                .setPopup(popup)
                .addTo(map);

            markers.push(marker);
        });
    }

    onMount(async () => {
        // 1. Initialize the map
        map = new maplibregl.Map({
            container: mapContainer,
            style: MAP_STYLES.normal,
            center: [110.44057776801604, -7.777334013872872],
            zoom: 12,
        });

        map.addControl(new maplibregl.NavigationControl(), "top-right");

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
    <div class="header">
        <div class="controls">
            <button on:click={() => switchStyle("normal")}>Normal</button>
            <button on:click={() => switchStyle("hybrid")}>Satellite</button>
        </div>
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
    .header {
        padding: 10px;
        background: transparent;
        display: flex;
        justify-content: space-between;
        align-items: center;
        position: absolute;
        z-index: 1;
    }
    .map-container {
        flex-grow: 1;
    }
    button {
        padding: 5px 15px;
        cursor: pointer;
    }

    :global(.marker-gps) {
        width: 20px;
        height: 20px;
        background-color: red;
        border: 2px solid white;
        border-radius: 50%;
        box-shadow: 0 0 8px rgba(0, 0, 0, 0.4);
        cursor: pointer;
        position: relative;
    }

    /* The pulsing ring/wedge */
    :global(.marker-gps::after) {
        content: "";
        position: absolute;
        /* Center it perfectly */
        top: 50%;
        left: 50%;
        width: 60px;
        height: 60px;
        margin-top: -30px;
        margin-left: -30px;

        border-radius: 50%;

        /* Use conic-gradient to create the partial arc */
        /* --ring-angle is passed from JS */
        background: conic-gradient(red var(--ring-angle), transparent 0deg);

        /* Use a mask to make it a 'ring' instead of a 'pie slice' */
        -webkit-mask-image: radial-gradient(circle, transparent 40%, black 41%);
        mask-image: radial-gradient(circle, transparent 40%, black 41%);

        /* Rotate by -90deg if you want the arc to start from the top center */
        transform: rotate(-90deg) scale(1);
        opacity: 0.4;
        z-index: -1;
    }

    /* Apply animation only if NOT inactive */
    :global(.marker-gps:not(.inactive)::after) {
        animation: pulse 2s infinite;
    }

    /* Hide the ring entirely if inactive (optional) */
    :global(.marker-gps.inactive::after) {
        display: none;
    }

    @keyframes pulse {
        0% {
            transform: rotate(-90deg) scale(1); /* Keep rotation during scale */
            opacity: 0.6;
        }
        100% {
            transform: rotate(-90deg) scale(3);
            opacity: 0;
        }
    }
</style>
