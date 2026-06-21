<script>
  import { onMount, onDestroy } from 'svelte';
  import { gameState, playerID, currentLayer, selectedCell, buildMode, currentView, operationMode, setOperationMode, prevResources, trackResourceDelta } from '../stores/game.js';
  import { api, connectWS } from '../utils/api.js';
  import { hexToPixel } from '../utils/hex.js';
  import HexMap from './HexMap.svelte';
  import ResourceBar from './ResourceBar.svelte';
  import BuildingPanel from './BuildingPanel.svelte';
  import BuildingMenu from './BuildingMenu.svelte';
  import UnitPanel from './UnitPanel.svelte';
  import DiplomacyPanel from './DiplomacyPanel.svelte';
  import TechPanel from './TechPanel.svelte';
  import LayerSelector from './LayerSelector.svelte';
  import TurnInfo from './TurnInfo.svelte';

  let ws = null;
  let activeTab = 'build';
  let game = null;
  let showBuildingMenu = false;
  let buildingMenuCell = null;
  let buildingMenuPosition = { x: 0, y: 0 };

  $: {
    game = $gameState;
  }

  $: currentPlayer = game?.players?.[$playerID] || null;

  onMount(() => {
    if ($gameState?.id) {
      ws = connectWS($gameState.id, $playerID, handleMessage);
    }
  });

  onDestroy(() => {
    if (ws) ws.close();
  });

  function handleMessage(msg) {
    if (msg.type === 'game_state') {
      const oldPlayer = game?.players?.[$playerID];
      gameState.set(msg.data);
      
      if (oldPlayer && msg.data?.players?.[$playerID]) {
        trackResourceDelta(msg.data.players[$playerID].resources, oldPlayer.resources);
        
        setTimeout(() => {
          trackResourceDelta(msg.data.players[$playerID].resources, msg.data.players[$playerID].resources);
        }, 2000);
      }
    }
  }

  async function submitAction(action, data) {
    try {
      await api.submitAction($gameState.id, $playerID, action, data);
    } catch (e) {
      console.error('Action failed:', e);
    }
  }

  function isAdjacentToOwned(cell) {
    if (!cell || !cell.coord) return false;
    const gameLayer = game.map?.layers?.[$currentLayer];
    if (!gameLayer) return false;
    
    const neighbors = getHexNeighbors(cell.coord);
    for (const n of neighbors) {
      const key = `${n.q},${n.r}`;
      const nc = gameLayer.cells?.[key];
      if (nc && nc.owner === $playerID) {
        return true;
      }
    }
    return false;
  }

  function getHexNeighbors(coord) {
    return [
      { q: coord.q + 1, r: coord.r },
      { q: coord.q - 1, r: coord.r },
      { q: coord.q, r: coord.r + 1 },
      { q: coord.q, r: coord.r - 1 },
      { q: coord.q + 1, r: coord.r - 1 },
      { q: coord.q - 1, r: coord.r + 1 }
    ];
  }

  function handleCellClick(cell) {
    if ($operationMode === 'mine') {
      if (cell.is_wall && isAdjacentToOwned(cell) && !cell.mining_owner) {
        submitAction('mine', {
          layer: $currentLayer,
          coord: cell.coord,
          workers: 1
        });
      }
      selectedCell.set(cell);
    } else if ($operationMode !== 'mine' && cell.owner === $playerID && !cell.is_wall && !cell.building && (cell.units?.length || 0) === 0) {
      buildingMenuCell = cell;
      showBuildingMenu = true;
      selectedCell.set(cell);
    } else if ($buildMode) {
      if (!cell.is_wall && cell.owner === $playerID && !cell.building) {
        submitAction('build', {
          layer: $currentLayer,
          coord: cell.coord,
          building: $buildMode
        });
      }
      buildMode.set(null);
      selectedCell.set(cell);
    } else {
      selectedCell.set(cell);
    }
  }

  function handleMine(cell) {
    submitAction('mine', {
      layer: $currentLayer,
      coord: cell.coord,
      workers: 1
    });
  }

  function handleMove(unit, target) {
    submitAction('move', {
      unit_id: unit.id,
      to: target.coord,
      to_layer: $currentLayer
    });
  }

  function handleAttack(unit, target) {
    submitAction('attack', {
      unit_id: unit.id,
      target: target.coord,
      target_layer: $currentLayer
    });
  }

  function selectBuilding(buildingType) {
    if (buildingMenuCell) {
      submitAction('build', {
        layer: $currentLayer,
        coord: buildingMenuCell.coord,
        building: buildingType
      });
    }
    showBuildingMenu = false;
    buildingMenuCell = null;
  }

  function closeBuildingMenu() {
    showBuildingMenu = false;
    buildingMenuCell = null;
  }

  function handleMapClick() {
    if (showBuildingMenu) {
      closeBuildingMenu();
    }
  }

  function setMode(mode) {
    if ($operationMode === mode) {
      setOperationMode(null);
    } else {
      setOperationMode(mode);
      if (mode === 'build') {
        activeTab = 'build';
      }
    }
    closeBuildingMenu();
  }

  function exitGame() {
    if (confirm('确定要退出游戏吗？')) {
      gameState.set(null);
      currentView.set('menu');
    }
  }
</script>

<div class="game-container">
  <div class="top-bar">
    <ResourceBar player={currentPlayer} />
    <TurnInfo game={game} playerId={$playerID} gameId={game?.id} />
    <button class="btn-exit" on:click={exitGame}>退出</button>
  </div>

  <div class="main-content">
    <div class="left-panel">
      <LayerSelector />
      
      <div class="operation-modes">
        <button 
          class="mode-btn {$operationMode === 'mine' ? 'active' : ''}"
          on:click={() => setMode('mine')}
          class:disabled={game?.phase !== 'planning'}
        >
          <span class="mode-icon">⛏️</span>
          <span class="mode-text">挖掘</span>
        </button>
        <button 
          class="mode-btn {$operationMode === 'build' ? 'active' : ''}"
          on:click={() => setMode('build')}
          class:disabled={game?.phase !== 'planning'}
        >
          <span class="mode-icon">🏗️</span>
          <span class="mode-text">建造</span>
        </button>
      </div>

      <div class="tabs">
        <button class={activeTab === 'build' ? 'active' : ''} on:click={() => activeTab = 'build'}>🏗️ 建筑</button>
        <button class={activeTab === 'units' ? 'active' : ''} on:click={() => activeTab = 'units'}>⚔️ 单位</button>
        <button class={activeTab === 'tech' ? 'active' : ''} on:click={() => activeTab = 'tech'}>🔬 科技</button>
        <button class={activeTab === 'diplomacy' ? 'active' : ''} on:click={() => activeTab = 'diplomacy'}>🤝 外交</button>
      </div>

      <div class="panel-content">
        {#if activeTab === 'build'}
          <BuildingPanel onBuild={(type) => buildMode.set(type)} player={currentPlayer} />
        {:else if activeTab === 'units'}
          <UnitPanel player={currentPlayer} />
        {:else if activeTab === 'tech'}
          <TechPanel player={currentPlayer} />
        {:else if activeTab === 'diplomacy'}
          <DiplomacyPanel game={game} playerId={$playerID} />
        {/if}
      </div>
    </div>

    <div class="map-area" on:click={handleMapClick}>
      <HexMap
        game={game}
        layer={$currentLayer}
        playerId={$playerID}
        onCellClick={handleCellClick}
        onMine={handleMine}
        onMove={handleMove}
        onAttack={handleAttack}
        selectedCell={$selectedCell}
        buildMode={$buildMode}
        operationMode={$operationMode}
      />

      {#if showBuildingMenu && buildingMenuCell}
        <div class="building-modal-overlay" on:click={closeBuildingMenu}>
          <div class="building-modal" on:click|stopPropagation>
            <BuildingMenu 
              player={currentPlayer}
              cell={buildingMenuCell}
              onSelect={selectBuilding}
              onClose={closeBuildingMenu}
            />
          </div>
        </div>
      {/if}

      {#if $operationMode === 'mine'}
        <div class="mode-hint">
          <span class="hint-icon">💡</span>
          <span>点击与己方领地相邻的岩石格子开始挖掘</span>
        </div>
      {/if}

      {#if $buildMode}
        <div class="mode-hint build-hint">
          <span class="hint-icon">💡</span>
          <span>点击己方领地内的空地放置建筑</span>
          <button class="cancel-btn" on:click={() => buildMode.set(null)}>取消</button>
        </div>
      {/if}
    </div>

    <div class="right-panel">
      {#if $selectedCell}
        <div class="cell-info">
          <h3>格子信息</h3>
          <p>坐标: ({$selectedCell.coord.q}, {$selectedCell.coord.r})</p>
          <p>岩石硬度: {$selectedCell.rock_hardness}</p>
          <p>含水量: {Math.round($selectedCell.water_content * 100)}%</p>
          <p>矿物: {$selectedCell.mineral_type}</p>
          
          {#if $selectedCell.water_content >= 0.5 && $selectedCell.is_wall}
            <p class="warning-text">
              <span>💧</span> 高含水量 - 塌方风险！
            </p>
          {/if}
          
          {#if $selectedCell.mining_owner}
            <div class="mining-info">
              <h4>⛏️ 挖掘中</h4>
              <div class="progress-bar">
                <div 
                  class="progress-fill" 
                  style="width: {($selectedCell.mining_progress / $selectedCell.mining_total) * 100}%"
                ></div>
              </div>
              <p>{$selectedCell.mining_progress} / {$selectedCell.mining_total}</p>
            </div>
          {/if}
          
          {#if $selectedCell.building}
            <div class="building-info">
              <h4>建筑: {$selectedCell.building.type}</h4>
              <p>耐久: {$selectedCell.building.hp}/{$selectedCell.building.max_hp}</p>
              {#if !$selectedCell.building.completed}
                <p>建造中: {Math.round($selectedCell.building.progress / $selectedCell.building.build_time * 100)}%</p>
              {/if}
            </div>
          {/if}
          
          {#if $selectedCell.units && $selectedCell.units.length > 0}
            <div class="units-info">
              <h4>单位 ({$selectedCell.units.length})</h4>
              {#each $selectedCell.units as unit}
                <div class="unit-mini">
                  <span>{unit.type}</span>
                  <span class="hp">{unit.hp}/{unit.max_hp}</span>
                </div>
              {/each}
            </div>
          {/if}
          
          {#if $selectedCell.is_wall && isAdjacentToOwned($selectedCell) && !$selectedCell.mining_owner}
            <button class="btn-mine" on:click={() => handleMine($selectedCell)}>⛏️ 开始挖掘</button>
          {/if}
        </div>
      {/if}

      <div class="event-log">
        <h3>📜 事件日志</h3>
        <div class="events">
          {#each game?.events?.slice(-10) || [] as event}
            <div class="event-item">
              <span class="turn">第{event.turn}回合</span>
              <span class="msg">{event.message}</span>
            </div>
          {/each}
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .game-container {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    background: #1a1a2e;
  }

  .top-bar {
    display: flex;
    align-items: center;
    gap: 20px;
    padding: 10px 20px;
    background: #16213e;
    border-bottom: 1px solid #2c3e50;
  }

  .btn-exit {
    background: #e74c3c;
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 6px;
    cursor: pointer;
  }

  .main-content {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

  .left-panel {
    width: 280px;
    background: #16213e;
    border-right: 1px solid #2c3e50;
    display: flex;
    flex-direction: column;
  }

  .operation-modes {
    display: flex;
    gap: 8px;
    padding: 10px 12px;
    border-bottom: 1px solid #2c3e50;
  }

  .mode-btn {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    padding: 10px 8px;
    background: #1a1a2e;
    color: #bdc3c7;
    border: 2px solid #2c3e50;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .mode-btn:hover:not(.disabled) {
    background: #2c3e50;
    border-color: #3498db;
  }

  .mode-btn.active {
    background: rgba(52, 152, 219, 0.2);
    border-color: #3498db;
    color: #3498db;
  }

  .mode-btn.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .mode-icon {
    font-size: 1.5rem;
  }

  .mode-text {
    font-size: 0.8rem;
    font-weight: bold;
  }

  .tabs {
    display: flex;
    flex-wrap: wrap;
    padding: 8px;
    gap: 4px;
    border-bottom: 1px solid #2c3e50;
  }

  .tabs button {
    flex: 1;
    padding: 8px 4px;
    font-size: 0.75rem;
    background: #1a1a2e;
    color: #bdc3c7;
    border: 1px solid #2c3e50;
    border-radius: 4px;
    cursor: pointer;
  }

  .tabs button.active {
    background: #e94560;
    color: white;
    border-color: #e94560;
  }

  .panel-content {
    flex: 1;
    overflow-y: auto;
    padding: 12px;
  }

  .map-area {
    flex: 1;
    position: relative;
    overflow: hidden;
    background: #0d1421;
  }

  .building-modal-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 200;
    animation: fadeIn 0.2s ease-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  .building-modal {
    animation: scaleIn 0.25s ease-out;
  }

  @keyframes scaleIn {
    from {
      opacity: 0;
      transform: scale(0.9);
    }
    to {
      opacity: 1;
      transform: scale(1);
    }
  }

  .mode-hint {
    position: absolute;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 20px;
    background: rgba(52, 152, 219, 0.9);
    color: white;
    border-radius: 8px;
    font-size: 0.9rem;
    animation: slideUp 0.3s ease-out;
  }

  .mode-hint.build-hint {
    background: rgba(46, 204, 113, 0.9);
  }

  @keyframes slideUp {
    from {
      opacity: 0;
      transform: translateX(-50%) translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateX(-50%) translateY(0);
    }
  }

  .hint-icon {
    font-size: 1.1rem;
  }

  .cancel-btn {
    margin-left: 10px;
    padding: 4px 12px;
    background: rgba(255, 255, 255, 0.2);
    color: white;
    border: 1px solid rgba(255, 255, 255, 0.3);
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.8rem;
  }

  .cancel-btn:hover {
    background: rgba(255, 255, 255, 0.3);
  }

  .right-panel {
    width: 260px;
    background: #16213e;
    border-left: 1px solid #2c3e50;
    padding: 16px;
    overflow-y: auto;
  }

  .cell-info {
    background: rgba(0, 0, 0, 0.3);
    border-radius: 8px;
    padding: 12px;
    margin-bottom: 16px;
  }

  .cell-info h3 {
    color: #ecf0f1;
    margin-bottom: 10px;
    font-size: 1rem;
  }

  .cell-info p {
    color: #bdc3c7;
    font-size: 0.85rem;
    margin-bottom: 4px;
  }

  .warning-text {
    color: #e74c3c !important;
    font-weight: bold;
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .mining-info, .building-info, .units-info {
    margin-top: 10px;
    padding-top: 10px;
    border-top: 1px solid #34495e;
  }

  .mining-info h4, .building-info h4, .units-info h4 {
    color: #f39c12;
    font-size: 0.9rem;
    margin-bottom: 8px;
  }

  .building-info h4 {
    color: #3498db;
  }

  .progress-bar {
    width: 100%;
    height: 8px;
    background: rgba(0, 0, 0, 0.5);
    border-radius: 4px;
    overflow: hidden;
    margin-bottom: 4px;
  }

  .progress-fill {
    height: 100%;
    background: #f39c12;
    border-radius: 4px;
    transition: width 0.3s;
  }

  .unit-mini {
    display: flex;
    justify-content: space-between;
    padding: 4px 0;
    font-size: 0.8rem;
    color: #95a5a6;
  }

  .unit-mini .hp {
    color: #2ecc71;
  }

  .btn-mine {
    width: 100%;
    margin-top: 10px;
    padding: 10px;
    background: #f39c12;
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-weight: bold;
  }

  .btn-mine:hover {
    background: #e67e22;
  }

  .event-log h3 {
    color: #ecf0f1;
    margin-bottom: 10px;
    font-size: 1rem;
  }

  .events {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .event-item {
    font-size: 0.8rem;
    color: #95a5a6;
    padding: 6px;
    background: rgba(0, 0, 0, 0.2);
    border-radius: 4px;
  }

  .event-item .turn {
    color: #f39c12;
    margin-right: 8px;
    font-weight: bold;
  }

  @media (max-width: 1000px) {
    .left-panel, .right-panel {
      width: 200px;
    }
  }
</style>
