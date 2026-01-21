<script lang="ts">
  import Map from "./lib/Map.svelte";
  import { settings } from "./lib/store/configStore";

  let isResizing = false;

  const toggleSidebar = () => {
    $settings.sidebarExpanded = !$settings.sidebarExpanded;
  };

  const startResize = (e: MouseEvent) => {
    isResizing = true;
    e.preventDefault();
  };

  const handleMouseMove = (e: MouseEvent) => {
    if (!isResizing) return;
    const newWidth = window.innerWidth - e.clientX;

    // Limits: Min 150px, Max 70% of screen
    if (newWidth > 150 && newWidth < window.innerWidth * 0.7) {
      $settings.sidebarWidth = newWidth;
    }
  };

  const stopResize = () => (isResizing = false);
</script>

<svelte:window on:mousemove={handleMouseMove} on:mouseup={stopResize} />

<div class="app-container">
  <main>
    <div class="map-area">
      <Map />
    </div>

    <div class="sidebar-wrapper">
      {#if $settings.sidebarExpanded}
        <button
          type="button"
          class="resizer"
          class:active={isResizing}
          on:mousedown={startResize}
          aria-label="Resize sidebar"
        ></button>
      {/if}

      <aside
        style={$settings.sidebarExpanded
          ? `width: ${$settings.sidebarWidth}px`
          : "width: 50px"}
        class:resizing={isResizing}
      >
        <div class="sidebar-header">
          <button
            class="hamburger"
            on:click={toggleSidebar}
            aria-label="Toggle Sidebar">☰</button
          >
        </div>
        <div class="sidebar-content">
          {#if $settings.sidebarExpanded}
            <p>Ini menu samping</p>
          {/if}
        </div>
      </aside>
    </div>
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

  main {
    display: flex;
    flex: 1;
    overflow: hidden;
    position: relative;
  }

  .map-area {
    flex: 1;
    background-color: #f0f0f0;
    position: relative;
  }

  .sidebar-wrapper {
    display: flex;
    position: relative;
    z-index: 20;
  }

  aside {
    background-color: white;
    height: 100%;
    display: flex;
    flex-direction: column;
    box-shadow: -2px 0 5px rgba(0, 0, 0, 0.1);
    transition: width 0.3s ease;
  }

  aside.resizing {
    transition: none;
  }

  .sidebar-header {
    display: flex;
    align-items: center;
    justify-content: left;
    padding: 10px;
    border-bottom: 1px solid #eee;
  }

  .hamburger {
    background: transparent;
    border: none;
    font-size: 24px;
    cursor: pointer;
    color: #333;
    padding: 5px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .hamburger:hover {
    background-color: #f0f0f0;
    border-radius: 4px;
  }

  .sidebar-content {
    padding: 1rem;
    overflow-y: auto;
    flex: 1;
  }

  .resizer {
    width: 4px;
    cursor: col-resize;
    background-color: #f0f0f0;
    border-left: 1px solid #ccc;
    transition: background-color 0.2s;
    height: 100%;
    padding: 0;
    border: none;
    background: none;
  }

  .resizer:hover,
  .resizer.active,
  .resizer:focus {
    background-color: #bbb;
  }
</style>
