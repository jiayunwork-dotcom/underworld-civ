<script>
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';

  export let onBuild;
  export let player;

  let buildings = [];

  onMount(async () => {
    try {
      buildings = await api.getBuildings();
    } catch (e) {
      console.error('Failed to load buildings:', e);
    }
  });

  function canAfford(building) {
    if (!player?.resources) return false;
    const res = player.resources.resources || player.resources;
    return (res.stone || 0) >= (building.cost?.stone || 0) &&
           (res.metal || 0) >= (building.cost?.metal || 0) &&
           (res.glow_mushroom || 0) >= (building.cost?.glow_mushroom || 0) &&
           (res.water || 0) >= (building.cost?.water || 0) &&
           (res.magic_crystal || 0) >= (building.cost?.magic_crystal || 0) &&
           (res.fossil_fuel || 0) >= (building.cost?.fossil_fuel || 0);
  }

  function isRaceSpecific(building) {
    if (!building?.race_specific) return true;
    return building.race_specific === player?.race;
  }

  const buildingNames = {
    main_base: '主基地',
    living_quarters: '居住窟',
    workshop: '工坊',
    fungus_farm: '农菌场',
    smelter: '冶炼炉',
    academy: '学堂',
    altar: '祭坛',
    warehouse: '仓储洞',
    watchtower: '哨塔',
    wall: '城墙',
    elevator: '电梯',
    forge_shrine: '熔炉圣殿',
    spore_network: '孢子网络',
    crystal_tower: '水晶共鸣塔',
    petrified_wall: '石化长城',
    hive_rift: '虫巢裂隙'
  };
</script>

<div class="building-panel">
  <h3>🏗️ 建筑</h3>
  <p class="hint">选择建筑后点击空地建造</p>

  <div class="building-list">
    {#each buildings as building}
      {#if isRaceSpecific(building) && building.type !== 'main_base'}
        <div class="building-card {!canAfford(building) ? 'disabled' : ''}"
             on:click={() => canAfford(building) && onBuild(building.type)}>
          <div class="building-header">
            <span class="building-name">{buildingNames[building.type] || building.name}</span>
            <span class="build-time">⏱️ {building.build_time}回合</span>
          </div>
          <p class="building-desc">{building.description}</p>
          <div class="building-cost">
            {#if building.cost?.stone > 0}
              <span class="cost">🪨{building.cost.stone}</span>
            {/if}
            {#if building.cost?.metal > 0}
              <span class="cost">⚙️{building.cost.metal}</span>
            {/if}
            {#if building.cost?.glow_mushroom > 0}
              <span class="cost">🍄{building.cost.glow_mushroom}</span>
            {/if}
            {#if building.cost?.water > 0}
              <span class="cost">💧{building.cost.water}</span>
            {/if}
            {#if building.cost?.magic_crystal > 0}
              <span class="cost">💎{building.cost.magic_crystal}</span>
            {/if}
            {#if building.cost?.fossil_fuel > 0}
              <span class="cost">🔥{building.cost.fossil_fuel}</span>
            {/if}
          </div>
          {#if building.race_specific}
            <span class="race-tag">种族独有</span>
          {/if}
        </div>
      {/if}
    {/each}
  </div>
</div>

<style>
  .building-panel h3 {
    color: #ecf0f1;
    margin-bottom: 8px;
    font-size: 1rem;
  }

  .hint {
    color: #7f8c8d;
    font-size: 0.75rem;
    margin-bottom: 12px;
  }

  .building-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .building-card {
    padding: 10px;
    background: rgba(52, 73, 94, 0.3);
    border-radius: 6px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    cursor: pointer;
    transition: all 0.2s;
  }

  .building-card:hover:not(.disabled) {
    background: rgba(52, 73, 94, 0.5);
    border-color: #3498db;
  }

  .building-card.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .building-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 6px;
  }

  .building-name {
    font-weight: bold;
    color: #ecf0f1;
    font-size: 0.9rem;
  }

  .build-time {
    font-size: 0.7rem;
    color: #95a5a6;
  }

  .building-desc {
    font-size: 0.75rem;
    color: #bdc3c7;
    margin-bottom: 8px;
  }

  .building-cost {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .cost {
    font-size: 0.75rem;
    color: #f39c12;
    background: rgba(0, 0, 0, 0.3);
    padding: 2px 6px;
    border-radius: 4px;
  }

  .race-tag {
    display: inline-block;
    margin-top: 6px;
    padding: 2px 8px;
    background: #9b59b6;
    color: white;
    font-size: 0.7rem;
    border-radius: 4px;
  }
</style>
