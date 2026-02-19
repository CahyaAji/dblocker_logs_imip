<script lang="ts">
    import { onDestroy } from "svelte";
    import { API_BASE } from "./utils/api";

    const topic = "dbl/#";
    let status: "disconnected" | "connected" | "error" = "disconnected";
    let lastMessage = "";
    let source: EventSource | null = null;

    function connect() {
        disconnect();
        const url = `${API_BASE}/mqtt/stream`;
        source = new EventSource(url);

        source.addEventListener("mqtt", (evt: MessageEvent) => {
            lastMessage = evt.data;
            console.log("MQTT event:", lastMessage);
        });

        source.onerror = () => {
            status = "error";
            disconnect();
        };

        source.onopen = () => {
            status = "connected";
        };
    }

    function disconnect() {
        if (source) {
            source.close();
            source = null;
        }
        status = "disconnected";
    }

    onDestroy(disconnect);
</script>

<div class="box">
    <div class="controls">
        <span class="topic">Topic: {topic}</span>
        <button on:click={connect} disabled={status === "connected"}
            >Connect</button
        >
        <button on:click={disconnect} disabled={status === "disconnected"}
            >Disconnect</button
        >
        <span>({status})</span>
    </div>

    <div class="output">
        <div class="label">Last published message:</div>
        <pre>{lastMessage || "—"}</pre>
    </div>
</div>

<style>
    .box {
        border: 1px solid #ddd;
        border-radius: 6px;
        padding: 12px;
        display: grid;
        gap: 10px;
    }

    .controls {
        display: flex;
        gap: 8px;
        align-items: center;
    }

    input {
        flex: 1;
        padding: 6px 8px;
        border: 1px solid #ccc;
        border-radius: 4px;
    }

    button {
        padding: 6px 10px;
        border: 1px solid #ccc;
        border-radius: 4px;
        background: #fff;
        cursor: pointer;
    }

    button:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .output pre {
        margin: 4px 0 0;
        padding: 8px;
        background: #f7f7f7;
        border: 1px solid #eee;
        border-radius: 4px;
        white-space: pre-wrap;
        word-break: break-word;
    }

    .label {
        font-weight: 600;
    }
</style>
