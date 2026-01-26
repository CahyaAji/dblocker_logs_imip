<script lang="ts">
    import { onMount } from "svelte";
    import { dblockerStore, type DBlocker, type DBlockerConfig } from "./store/dblockerStore";

    let editingIds: number[] = $state([]);

    let dblockers: DBlocker[] = $state([]);
    let debounceTimer: ReturnType<typeof setTimeout>;

    function toggleEditMode(blocker: DBlocker) {
        if (editingIds.includes(blocker.id)) {
            // console.log(blocker.id, $state.snapshot(blocker).config);
            updateDBlocker(blocker.id, blocker.config);
            editingIds = editingIds.filter(i => i !== blocker.id);
        } else {
            editingIds = [...editingIds, blocker.id];
        }
    }

    async function readDBlockers() {
        try {
            const res = await fetch("/api/dblockers");
            if (!res.ok) throw new Error("Fetch dblockers failed");
            const json = await res.json();
            // Handle wrapped response { data: [...] } or direct array [...]
            const data: DBlocker[] = Array.isArray(json) ? json : (json.data || []);
            return data;

        } catch (error) {
            console.error("Error fetching dblockers:", error);
            return [];
        }
    }

    async function updateDBlocker(blockerid: number, config: DBlockerConfig[]) {
        const payload = { id: blockerid, config };
        try {
            const res = await fetch(`/api/dblockers/config`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(payload),
            });
            if (!res.ok) throw new Error("Update dblocker failed");
            return true;
        } catch (error) {
            console.error("Error updating dblocker:", error);
            return false;
        }
    }

    $effect(() => {
        const storeData = $dblockerStore;
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
            if (JSON.stringify(dblockers) !== JSON.stringify(storeData)) {
                console.log("DBlockersList: Store data changed, updating local state.");
                dblockers = storeData;
            }
        }, 3000);
        return () => clearTimeout(debounceTimer);
    });

    onMount(async ()=>{
        dblockers = await readDBlockers();
    });

</script>

    <div class="list">
        {#each dblockers as blocker (blocker.id)}
            {@const isEditMode = editingIds.includes(blocker.id)}
            <div class="card">
                <div class="card-header">
                    <div>{blocker.name}</div>
                    {#if isEditMode}
                    <button class="btn-edit" onclick={() => toggleEditMode(blocker)}>Apply</button>
                    {:else}
                    <button class="btn-edit" onclick={() => toggleEditMode(blocker)}>Edit</button>
                    {/if}
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
                                        <input type="checkbox" bind:checked={config.signal_ctrl} disabled={!isEditMode}>
                                        <span class="slider"></span>
                                        </label>
                                    </div>
                                    <div class="control-row">
                                        <div class="control-label">Blcker GPS</div>
                                        <label class="switch">
                                        <input type="checkbox" bind:checked={config.signal_gps} disabled={!isEditMode}>
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
                                        <input type="checkbox" bind:checked={config.signal_ctrl} disabled={!isEditMode}>
                                        <span class="slider"></span>
                                        </label>
                                    </div>
                                    <div class="control-row">
                                        <div class="control-label">Blcker GPS</div>
                                        <label class="switch">
                                        <input type="checkbox" bind:checked={config.signal_gps} disabled={!isEditMode}>
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
        gap: 8px;
        flex: 1;
        min-height: 0;
        padding: 10px 6px;
    }
    .empty {
        text-align: center;
        color: #888;
        margin-top: 2rem;
    }
    .sector {
        display: flex;
        flex-direction: column;
        gap: 2px;
        margin-bottom: 4px;
    }
    .card-content {
        display: flex;
        gap: 10px;
    }
    .col {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 6px;
    }
</style>
