<script>
    // import { dblockerStore, switchSignal } from "./store/dblockerStore";

    const dblockerStore = [
        {
            id: 1,
            name: "Drone Blocker timur",
            serial_numb: "STM-1001",
            latitude: -7.788432918834558,
            longitude: 110.44306970000933,
            desc: "Drone blocker di timur",
            angle_start: 0,
            config: [
                { signal_ctrl: false, signal_gps: false },
                { signal_ctrl: false, signal_gps: false },
                { signal_ctrl: true, signal_gps: true },
                { signal_ctrl: true, signal_gps: true },
                { signal_ctrl: true, signal_gps: true },
                { signal_ctrl: true, signal_gps: true },
            ],
        },
        {
            id: 2,
            name: "Drone Blocker Barat",
            serial_numb: "SCM-1002",
            latitude: -7.786959801210232,
            longitude: 110.4220920495705,
            desc: "Drone blocker di barat",
            angle_start: 0,
            config: [
                { signal_ctrl: true, signal_gps: true },
                { signal_ctrl: true, signal_gps: true },
                { signal_ctrl: true, signal_gps: true },
                { signal_ctrl: true, signal_gps: true },
                { signal_ctrl: true, signal_gps: true },
                { signal_ctrl: true, signal_gps: true },
            ],
        },
    ];
</script>

<div class="list">
    {#each dblockerStore as blocker (blocker.id)}
        <div class="card">
            <div>{blocker.name}</div>
            {#each blocker.config as config, index}
                <div>
                    <div>Sector {index + 1}</div>
                    <div>Block RC: {config.signal_ctrl ? "ON" : "OFF"}</div>
                    <div>Block GPS: {config.signal_gps ? "ON" : "OFF"}</div>
                </div>
                <!-- <div>
                    Sector {index + 1} - Block RC: {config.signal_ctrl
                        ? "ON"
                        : "OFF"}, Block GPS: {config.signal_gps ? "ON" : "OFF"}
                </div> -->
            {/each}
            <div>
                <button>Edit</button>
            </div>
        </div>
    {:else}
        <div class="empty">No DBlockers found</div>
    {/each}
</div>

<style>
    .list {
        display: flex;
        flex-direction: column;
        overflow-y: auto;
        scrollbar-color: var(--border-color-light) var(--background-color-light);
        gap: 6px;
        flex: 1;
        min-height: 0;
    }
    .card {
        margin: 0 2px;
        padding: 6px 8px;
        border: 1px solid var(--border-color-light);
        border-radius: 8px;
        background-color: whitesmoke;
        box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.1);
    }
    .empty {
        text-align: center;
        color: #888;
        margin-top: 2rem;
    }
</style>
