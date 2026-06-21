<script>
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';

  export let player;
  export let cell;
  export let onSelect;
  export let onClose;

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

  const buildingIcons = {
    main_base: '🏰',
    living_quarters: '🏠',
    workshop: '🔧',
    fungus_farm: '🍄',
    smelter: '🏭',
    academy: '📚',
    altar: '⛩️',
    warehouse: '📦',
    watchtower: '🗼',
    wall: '🧱',
    elevator: '🛗',
    forge_shrine: '🔨',
    spore_network: '🌐',
    crystal_tower: '💎',
    petrified_wall: '🪨',
    hive_rift: '🐛'
  };

  function handleBuild(building) {
    if (canAfford(building) && isRaceSpecific(building)) {
      onSelect(building.type);
    }
  }
</script>

<div class="building-menu" on:click|stopPropagation>
  <div class="menu-header">
    <h3>🏗️ 选择建筑</h3>
    <button class="btn-close" on:click={onClose}>×</button>
  </div>
  
  <div class="menu-info">
    <span class="coord">坐标: ({cell?.coord?.q}, {cell?.coord?.r})</span>
  </div>

  <div class="building-list">
    {#each buildings as building}
      {#if isRaceSpecific(building) && building.type !== 'main_base'}
        <div class="building-card {!canAfford(building) ? 'disabled' : ''}"
             on:click={() => handleBuild(building)}>
          <div class="building-icon">{buildingIcons[building.type] || '🏠'}</div>
          <div class="building-info">
            <div class="building-name">{buildingNames[building.type] || building.name}</div>
            <div class="building-desc">{building.description}</div>
            <div class="building-meta">
              <span class="build-time">⏱️ {building.build_time}回合</span>
            </div>
            <div class="building-cost">
              {#if building.cost?.stone > 0}
                <span class="cost-item">
                  <span class="cost-icon">🪨</span>
                  <span class:insufficient={(player?.resources?.stone || player?.resources?.resources?.stone || 0) < building.cost.stone}>
                    {building.cost.stone}
                  </span>
                </span>
              {/if}
              {#if building.cost?.metal > 0}
                <span class="cost-item">
                  <span class="cost-icon">⚙️</span>
                  <span class:insufficient={(player?.resources?.metal || player?.resources?.resources?.metal || 0) < building.cost.metal}>
                    {building.cost.metal}
                  </span>
                </span>
              {/if}
              {#if building.cost?.glow_mushroom > 0}
                <span class="cost-item">
                  <span class="cost-icon">🍄</span>
                  <span class:insufficient={(player?.resources?.glow_mushroom || player?.resources?.resources?.glow_mushroom || 0) < building.cost.glow_mushroom}>
                    {building.cost.glow_mushroom}
                  </span>
                </span>
              {/if}
              {#if building.cost?.water > 0}
                <span class="cost-item">
                  <span class="cost-icon">💧</span>
                  <span class:insufficient={(player?.resources?.water || player?.resources?.resources?.water || 0) < building.cost.water}>
                    {building.cost.water}
                  </span>
                </span>
              {/if}
              {#if building.cost?.magic_crystal > 0}
                <span class="cost-item">
                  <span class="cost-icon">💎</span>
                  <span class:insufficient={(player?.resources?.magic_crystal || player?.resources?.resources?.magic_crystal || 0) < building.cost.magic_crystal}>
                    {building.cost.magic_crystal}
                  </span>
                </span>
              {/if}
              {#if building.cost?.fossil_fuel > 0}
                <span class="cost-item">
                  <span class="cost-icon">🔥</span>
                  <span class:insufficient={(player?.resources?.fossil_fuel || player?.resources?.resources?.fossil_fuel || 0) < building.cost.fossil_fuel}>
                    {building.cost.fossil_fuel}
                  </span>
                </span>
              {/if}
            </div>
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
  .building-menu {
    background: #16213e;
    border: 1px solid #2c3e50;
    border-radius: 8px;
    box-shadow: 0 4px 30px rgba(0, 0, 0, 0.6);
    width: 320px;
    max-height: 70vh;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  @keyframes slideIn {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .menu-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 14px;
    background: #0f3460;
    border-bottom: 1px solid #2c3e50;
  }

  .menu-header h3 {
    color: #ecf0f1;
    margin: 0;
    font-size: 0.95rem;
  }

  .btn-close {
    background: none;
    border: none;
    color: #bdc3c7;
    font-size: 1.3rem;
    cursor: pointer;
    padding: 0;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
  }

  .btn-close:hover {
    background: rgba(255, 255, 255, 0.1);
    color: #fff;
  }

  .menu-info {
    padding: 8px 14px;
    background: rgba(0, 0, 0, 0.2);
    border-bottom: 1px solid #2c3e50;
  }

  .coord {
    color: #7f8c8d;
    font-size: 0.8rem;
  }

  .building-list {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .building-card {
    display: flex;
    gap: 10px;
    padding: 8px;
    background: rgba(52, 73, 94, 0.3);
    border-radius: 6px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    cursor: pointer;
    transition: all 0.2s;
    position: relative;
  }

  .building-card:hover:not(.disabled) {
    background: rgba(52, 73, 94, 0.5);
    border-color: #3498db;
    transform: translateX(2px);
  }

  .building-card.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .building-icon {
    font-size: 1.8rem;
    flex-shrink: 0;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.3);
    border-radius: 6px;
  }

  .building-info {
    flex: 1;
    min-width: 0;
  }

  .building-name {
    font-weight: bold;
    color: #ecf0f1;
    font-size: 0.85rem;
    margin-bottom: 2px;
  }

  .building-desc {
    font-size: 0.7rem;
    color: #bdc3c7;
    margin-bottom: 4px;
    line-height: 1.3;
  }

  .building-meta {
    margin-bottom: 4px;
  }

  .build-time {
    font-size: 0.7rem;
    color: #f39c12;
  }

  .building-cost {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .cost-item {
    display: flex;
    align-items: center;
    gap: 2px;
    font-size: 0.7rem;
    color: #bdc3c7;
    background: rgba(0, 0, 0, 0.3);
    padding: 2px 5px;
    border-radius: 3px;
  }

  .cost-icon {
    font-size: 0.75rem;
  }

  .cost-item .insufficient {
    color: #e74c3c;
    font-weight: bold;
  }

  .race-tag {
    position: absolute;
    top: 4px;
    right: 4px;
    padding: 1px 6px;
    background: #9b59b6;
    color: white;
    font-size: 0.6rem;
    border-radius: 3px;
  }
</style>
