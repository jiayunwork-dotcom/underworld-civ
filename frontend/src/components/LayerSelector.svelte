<script>
  import { currentLayer } from '../stores/game.js';

  const layers = [
    { id: 0, name: '浅层', color: '#2ecc71', desc: '蘑菇林带' },
    { id: 1, name: '中层', color: '#f39c12', desc: '石灰岩洞窟' },
    { id: 2, name: '深层', color: '#95a5a6', desc: '金属矿脉' },
    { id: 3, name: '熔岩', color: '#e74c3c', desc: '魔晶之地' },
    { id: 4, name: '远古', color: '#9b59b6', desc: '化石遗迹' }
  ];

  function selectLayer(id) {
    currentLayer.set(id);
  }
</script>

<div class="layer-selector">
  <h3>🗺️ 图层</h3>
  <div class="layers">
    {#each layers as layer}
      <button
        class="layer-btn {$currentLayer === layer.id ? 'active' : ''}"
        style="border-color: {layer.color}"
        on:click={() => selectLayer(layer.id)}
      >
        <span class="layer-name" style="color: {layer.color}">{layer.name}</span>
        <span class="layer-desc">{layer.desc}</span>
      </button>
    {/each}
  </div>
</div>

<style>
  .layer-selector {
    padding: 12px;
    border-bottom: 1px solid #2c3e50;
  }

  .layer-selector h3 {
    color: #ecf0f1;
    font-size: 0.95rem;
    margin-bottom: 10px;
  }

  .layers {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .layer-btn {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    padding: 8px 12px;
    background: #1a1a2e;
    border: 2px solid #34495e;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .layer-btn:hover {
    background: #2c3e50;
  }

  .layer-btn.active {
    background: rgba(255, 255, 255, 0.1);
  }

  .layer-name {
    font-weight: bold;
    font-size: 0.9rem;
  }

  .layer-desc {
    font-size: 0.7rem;
    color: #7f8c8d;
  }
</style>
