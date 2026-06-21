<script>
  import { resourceDeltas } from '../stores/game.js';

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
    return res?.resources?.[key] ?? res?.[key] ?? 0;
  }

  function getProduction(prod, key) {
    return prod?.[key] || 0;
  }

  function isLowResource(player, key) {
    const amount = getResource(player?.resources, key);
    const production = getProduction(player?.production, key);
    return amount < 20 && production <= 0;
  }

  function getDelta(key) {
    return $resourceDeltas[key] || 0;
  }

  function hasPositiveDelta(key) {
    return getDelta(key) > 0;
  }

  function hasNegativeDelta(key) {
    return getDelta(key) < 0;
  }
</script>

<div class="resource-bar">
  {#each resourceInfo as res}
    <div class="resource-item {isLowResource(player, res.key) ? 'warning' : ''}" title={res.name}>
      <span class="icon">{res.icon}</span>
      <span class="value" style="color: {res.color}">
        {getResource(player?.resources, res.key)}
      </span>
      {#if hasPositiveDelta(res.key)}
        <span class="delta positive">
          <span class="arrow">↑</span>
          +{getDelta(res.key)}
        </span>
      {:else if hasNegativeDelta(res.key)}
        <span class="delta negative">
          <span class="arrow">↓</span>
          {getDelta(res.key)}
        </span>
      {:else}
        <span class="production">
          {#if getProduction(player?.production, res.key) > 0}
            <span class="prod-positive">+{getProduction(player?.production, res.key)}</span>
          {:else if getProduction(player?.production, res.key) < 0}
            <span class="prod-negative">{getProduction(player?.production, res.key)}</span>
          {:else}
            <span class="prod-neutral">+0</span>
          {/if}
        </span>
      {/if}
      {#if isLowResource(player, res.key)}
        <span class="warning-icon">⚠️</span>
      {/if}
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
    gap: 12px;
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
    position: relative;
    transition: all 0.3s;
    border: 1px solid transparent;
  }

  .resource-item.warning {
    background: rgba(231, 76, 60, 0.2);
    border-color: #e74c3c;
    animation: pulse-warning 1.5s infinite;
  }

  @keyframes pulse-warning {
    0%, 100% {
      box-shadow: 0 0 0 0 rgba(231, 76, 60, 0.4);
    }
    50% {
      box-shadow: 0 0 8px 2px rgba(231, 76, 60, 0.3);
    }
  }

  .icon {
    font-size: 1.1rem;
  }

  .value {
    font-weight: bold;
    color: #ecf0f1;
    font-variant-numeric: tabular-nums;
    min-width: 30px;
    text-align: right;
  }

  .production {
    font-size: 0.75rem;
  }

  .prod-positive {
    color: #2ecc71;
  }

  .prod-negative {
    color: #e74c3c;
  }

  .prod-neutral {
    color: #7f8c8d;
  }

  .delta {
    font-size: 0.75rem;
    font-weight: bold;
    display: flex;
    align-items: center;
    gap: 2px;
    animation: delta-pop 0.5s ease-out;
  }

  .delta .arrow {
    font-size: 0.65rem;
  }

  .delta.positive {
    color: #2ecc71;
  }

  .delta.negative {
    color: #e74c3c;
  }

  @keyframes delta-pop {
    0% {
      transform: scale(1.3);
      opacity: 0;
    }
    50% {
      transform: scale(1.1);
    }
    100% {
      transform: scale(1);
      opacity: 1;
    }
  }

  .warning-icon {
    font-size: 0.8rem;
    animation: shake 0.5s ease-in-out infinite;
  }

  @keyframes shake {
    0%, 100% { transform: translateX(0); }
    25% { transform: translateX(-1px); }
    75% { transform: translateX(1px); }
  }

  .population .value {
    color: #3498db;
  }

  .research .value {
    color: #9b59b6;
  }
</style>
