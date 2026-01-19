<script lang="ts">
  import Map from "./lib/Map.svelte";

  let panelWidth = 300;
  let isVisible = true;
  let isResizing = false;

  const toggleSidebar = () => (isVisible = !isVisible);

  const startResize = (e: MouseEvent) => {
    isResizing = true;
    e.preventDefault();
  };

  const handleMouseMove = (e: MouseEvent) => {
    if (!isResizing) return;
    const newWidth = window.innerWidth - e.clientX;
    if (newWidth > 150 && newWidth < window.innerWidth * 0.7) {
      panelWidth = newWidth;
    }
  };

  const stopResize = () => (isResizing = false);
</script>

<svelte:window on:mousemove={handleMouseMove} on:mouseup={stopResize} />

<div class="app-container">
  <header>
    <div class="logo">Drone Blocker</div>
    <button
      class="hamburger"
      on:click={toggleSidebar}
      aria-label="Toggle Sidebar"
    >
      ☰
    </button>
  </header>

  <main>
    <div class="map-area">
      <Map />
    </div>

    {#if isVisible}
      <button
        type="button"
        class="resizer"
        class:active={isResizing}
        on:mousedown={startResize}
        aria-label="Resize sidebar"
        tabindex="0"
      ></button>

      <aside style="width: {panelWidth}px">
        <div class="sidebar-content">
          <p>Ini menu samping</p>
        </div>
      </aside>
    {/if}
  </main>
</div>

<style>
  .app-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100vw;
    overflow: hidden;
  }

  header {
    height: 40px;
    background-color: lightgray;
    color: black;
    display: flex;
    align-items: center;
    font-size: 16pt;
    justify-content: space-between;
    padding: 0 20px;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
    z-index: 10;
  }

  .hamburger {
    background: transparent;
    border: 1px solid #666;
    color: black;
    font-size: 1.5rem;
    cursor: pointer;
    padding: 0 4px;
    border-radius: 4px;
  }

  .hamburger:hover {
    background: #444;
  }

  main {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  .map-area {
    flex: 1;
    background-color: #f0f0f0;
    position: relative;
  }

  aside {
    background-color: lightgray;
    height: 100%;
    overflow-y: auto;
  }

  .sidebar-content {
    padding: 1rem;
  }

  .resizer {
    width: 8px;
    cursor: col-resize;
    background-color: #eee;
    border-left: 1px solid #ccc;
    border-right: 1px solid #ccc;
    transition: background-color 0.2s;
    border: none;
    padding: 0;
    outline: none;
    height: 100%;
  }

  .resizer:hover,
  .resizer.active,
  .resizer:focus {
    background-color: #bbb;
  }
</style>
