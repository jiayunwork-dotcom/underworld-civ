<script>
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';

  export let player;

  let units = [];

  onMount(async () => {
    try {
      units = await api.getUnits();
    } catch (e) {
      console.error('Failed to load units:', e);
    }
  });

  const unitNames = {
    sapper: '工兵',
    infantry: '步兵',
    archer: '弓手',
    siege_ram: '攻城锤',
    iron_guard: '铁甲卫士',
    spore_grenadier: '孢子投弹兵',
    shadow_archer: '暗影射手',
    rock_giant: '岩石巨人',
    burrow_bomber: '钻地爆虫'
  };

  function isRaceSpecific(unit) {
    if (!unit?.race_specific) return true;
    return unit.race_specific === player?.race;
  }

  function getPlayerUnits() {
    if (!player?.units) return [];
    const counts = {};
    for (const unit of player.units) {
      counts[unit.type] = (counts[unit.type] || 0) + 1;
    }
    return counts;
  }

  $: playerUnitCounts = getPlayerUnits();
</script>

<div class="unit-panel">
  <h3>⚔️ 兵种</h3>
  <p class="hint">点击地图上的单位移动或攻击</p>

  <div class="unit-summary">
    <strong>当前部队:</strong>
    {#each Object.keys(playerUnitCounts) as type}
      <span class="unit-count">{unitNames[type] || type}: {playerUnitCounts[type]}</span>
    {/each}
    {#if Object.keys(playerUnitCounts).length === 0}
      <span class="none">暂无单位</span>
    {/if}
  </div>

  <div class="unit-list">
    {#each units as unit}
      {#if isRaceSpecific(unit)}
        <div class="unit-card">
          <div class="unit-header">
            <span class="unit-name">{unitNames[unit.type] || unit.name}</span>
            {#if unit.race_specific}
              <span class="race-tag">独有</span>
            {/if}
          </div>
          <p class="unit-desc">{unit.description}</p>
          <div class="unit-stats">
            <span>❤️ {unit.hp}</span>
            <span>⚔️ {unit.attack}</span>
            <span>🛡️ {unit.defense}</span>
            <span>🎯 {unit.range}</span>
            <span>🏃 {unit.speed}</span>
          </div>
          <div class="unit-cost">
            {#if unit.cost?.stone > 0}
              <span class="cost">🪨{unit.cost.stone}</span>
            {/if}
            {#if unit.cost?.metal > 0}
              <span class="cost">⚙️{unit.cost.metal}</span>
            {/if}
            {#if unit.cost?.glow_mushroom > 0}
              <span class="cost">🍄{unit.cost.glow_mushroom}</span>
            {/if}
            {#if unit.cost?.magic_crystal > 0}
              <span class="cost">💎{unit.cost.magic_crystal}</span>
            {/if}
            {#if unit.cost?.fossil_fuel > 0}
              <span class="cost">🔥{unit.cost.fossil_fuel}</span>
            {/if}
          </div>
        </div>
      {/if}
    {/each}
  </div>
</div>

<style>
  .unit-panel h3 {
    color: #ecf0f1;
    margin-bottom: 8px;
    font-size: 1rem;
  }

  .hint {
    color: #7f8c8d;
    font-size: 0.75rem;
    margin-bottom: 12px;
  }

  .unit-summary {
    background: rgba(0, 0, 0, 0.3);
    padding: 10px;
    border-radius: 6px;
    margin-bottom: 12px;
    font-size: 0.8rem;
    color: #bdc3c7;
  }

  .unit-count {
    display: inline-block;
    margin-right: 8px;
    color: #3498db;
  }

  .none {
    color: #7f8c8d;
    font-style: italic;
  }

  .unit-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .unit-card {
    padding: 10px;
    background: rgba(52, 73, 94, 0.3);
    border-radius: 6px;
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .unit-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 6px;
  }

  .unit-name {
    font-weight: bold;
    color: #ecf0f1;
    font-size: 0.9rem;
  }

  .race-tag {
    background: #9b59b6;
    color: white;
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 0.7rem;
  }

  .unit-desc {
    font-size: 0.75rem;
    color: #bdc3c7;
    margin-bottom: 8px;
  }

  .unit-stats {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-bottom: 8px;
    font-size: 0.75rem;
    color: #f39c12;
  }

  .unit-cost {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
  }

  .cost {
    font-size: 0.7rem;
    color: #e67e22;
    background: rgba(0, 0, 0, 0.3);
    padding: 2px 5px;
    border-radius: 3px;
  }
</style>
