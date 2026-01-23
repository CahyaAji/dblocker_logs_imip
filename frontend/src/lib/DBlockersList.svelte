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
    <!-- {#each $dblockerStore as blocker (blocker.id)} -->
    {#each dblockerStore as blocker (blocker.id)}
        <div class="card">
            <div>{blocker.name}</div>
            <div>{blocker.serial_numb}</div>
            <div>{blocker.latitude}</div>
            <div>{blocker.longitude}</div>
            <div>{blocker.angle_start}</div>
            <div>{blocker.config}</div>
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
        gap: 4px;
    }
    .card {
        margin: 0 2px;
        padding: 4px 8px;
        border: 2px solid brown;
        border-radius: 4px 8px 4px 8px;
    }
    .empty {
        text-align: center;
        color: #888;
        margin-top: 2rem;
    }
</style>
