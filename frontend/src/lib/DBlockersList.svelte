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
                <div class="card-header">
                    <div>{blocker.name}</div>
                    <button class="btn-edit">Edit</button>
                </div>
                <div class="card-content">
                    <div class="col">
                        {#each blocker.config.slice(0, 3) as config, index}
                            <div class="sector">
                                <div class="section-title">Sector {index + 1}</div>
                                <div class="control-group">
                                    <div class="control-row">
                                        <div class="control-label">Blcker RC</div>
                                        <label class="switch">
                                        <input type="checkbox" checked={config.signal_ctrl}>
                                        <span class="slider"></span>
                                        </label>
                                    </div>
                                    <div class="control-row">
                                        <div class="control-label">Blcker GPS</div>
                                        <label class="switch">
                                        <input type="checkbox" checked={config.signal_gps}>
                                        <span class="slider"></span>
                                        </label>
                                    </div>
                                </div>
                            </div>
                        {/each}
                    </div>
                    <div class="col">
                        {#each blocker.config.slice(3, 6) as config, index}
                            <div class="sector">
                                <div class="section-title">Sector {index + 4}</div>
                                <div class="control-group">
                                    <div class="control-row">
                                        <div class="control-label">Blcker RC</div>
                                        <label class="switch">
                                        <input type="checkbox" checked={config.signal_ctrl}>
                                        <span class="slider"></span>
                                        </label>
                                    </div>
                                    <div class="control-row">
                                        <div class="control-label">Blcker GPS</div>
                                        <label class="switch">
                                        <input type="checkbox" checked={config.signal_gps}>
                                        <span class="slider"></span>
                                        </label>
                                    </div>
                                </div>
                            </div>
                        {/each}
                    </div>
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
        scrollbar-color: var(--separator) var(--bg-color);
        gap: 16px;
        flex: 1;
        min-height: 0;
        padding: 16px;
    }
    .empty {
        text-align: center;
        color: #888;
        margin-top: 2rem;
    }
    .sector {
        display: flex;
        flex-direction: column;
        gap: 4px;
        margin-bottom: 6px;
    }
    .card-content {
        display: flex;
        gap: 12px;
    }
    .col {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 6px;
    }
</style>
