<script>
  export let player;

  const resourceInfo = [
    { key: 'stone', name: '石材', icon: '🪨', color: '#95a5a6' },
    { key: 'metal', name: '金属', icon: '⚙️', color: '#7f8c8d' },
    { key: 'glow_mushroom', name: '发光菌', icon: '🍄', color: '#2ecc71' },
    { key: 'water', name: '地下水', icon: '💧', color: '#3498db' },
    { key: 'magic_crystal', name: '魔晶', icon: '💎', color: '#9b59b6' },
    { key: 'fossil_fuel', name: '化石燃料', icon: '🔥', color: '#e74c3c' }
  ];

  function getResource(res, key) {
    return res?.resources?.[key] || 0;
  }

  function getProduction(prod, key) {
    return prod?.[key] || 0;
  }
</script>

<div class="resource-bar">
  {#each resourceInfo as res}
    <div class="resource-item" title={res.name}>
      <span class="icon">{res.icon}</span>
      <span class="value" style="color: {res.color}">
        {getResource(player?.resources, res.key)}
      </span>
      <span class="production">
        (+{getProduction(player?.production, res.key)})
      </span>
    </div>
  {/each}

  <div class="resource-item population">
    <span class="icon">👥</span>
    <span class="value">{player?.population || 0}/{player?.population_cap || 0}</span>
  </div>

  <div class="resource-item research">
    <span class="icon">🔬</span>
    <span class="value">{player?.research_points || 0}</span>
  </div>
</div>

<style>
  .resource-bar {
    display: flex;
    gap: 16px;
    flex: 1;
    flex-wrap: wrap;
  }

  .resource-item {
    display: flex;
    align-items: center;
    gap: 6px;
    background: rgba(0, 0, 0, 0.3);
    padding: 6px 12px;
    border-radius: 6px;
    font-size: 0.9rem;
  }

  .icon {
    font-size: 1.1rem;
  }

  .value {
    font-weight: bold;
    color: #ecf0f1;
  }

  .production {
    font-size: 0.75rem;
    color: #2ecc71;
  }

  .population .value {
    color: #3498db;
  }

  .research .value {
    color: #9b59b6;
  }
</style>
