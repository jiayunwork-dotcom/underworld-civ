<script>
  import { onMount, onDestroy } from 'svelte';
  import { gameState, playerID, currentLayer, selectedCell, buildMode, currentView } from '../stores/game.js';
  import { api, connectWS } from '../utils/api.js';
  import HexMap from './HexMap.svelte';
  import ResourceBar from './ResourceBar.svelte';
  import BuildingPanel from './BuildingPanel.svelte';
  import UnitPanel from './UnitPanel.svelte';
  import DiplomacyPanel from './DiplomacyPanel.svelte';
  import TechPanel from './TechPanel.svelte';
  import LayerSelector from './LayerSelector.svelte';
  import TurnInfo from './TurnInfo.svelte';

  let ws = null;
  let activeTab = 'build';
  let game = null;

  $: {
    game = $gameState;
  }

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
      gameState.set(msg.data);
    }
  }

  async function submitAction(action, data) {
    try {
      await api.submitAction($gameState.id, $playerID, action, data);
    } catch (e) {
      console.error('Action failed:', e);
    }
  }

  function handleCellClick(cell) {
    if ($buildMode) {
      submitAction('build', {
        layer: $currentLayer,
        coord: cell.coord,
        building: $buildMode
      });
      buildMode.set(null);
    } else if (cell.units && cell.units.length > 0 && cell.units[0].owner === $playerID) {
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

  function exitGame() {
    if (confirm('确定要退出游戏吗？')) {
      gameState.set(null);
      currentView.set('menu');
    }
  }
</script>

<div class="game-container">
  <div class="top-bar">
    <ResourceBar player={game?.players?.[$playerID]} />
    <TurnInfo game={game} />
    <button class="btn-exit" on:click={exitGame}>退出</button>
  </div>

  <div class="main-content">
    <div class="left-panel">
      <LayerSelector />
      <div class="tabs">
        <button class={activeTab === 'build' ? 'active' : ''} on:click={() => activeTab = 'build'}>🏗️ 建筑</button>
        <button class={activeTab === 'units' ? 'active' : ''} on:click={() => activeTab = 'units'}>⚔️ 单位</button>
        <button class={activeTab === 'tech' ? 'active' : ''} on:click={() => activeTab = 'tech'}>🔬 科技</button>
        <button class={activeTab === 'diplomacy' ? 'active' : ''} on:click={() => activeTab = 'diplomacy'}>🤝 外交</button>
      </div>

      <div class="panel-content">
        {#if activeTab === 'build'}
          <BuildingPanel onBuild={(type) => buildMode.set(type)} player={game?.players?.[$playerID]} />
        {:else if activeTab === 'units'}
          <UnitPanel player={game?.players?.[$playerID]} />
        {:else if activeTab === 'tech'}
          <TechPanel player={game?.players?.[$playerID]} />
        {:else if activeTab === 'diplomacy'}
          <DiplomacyPanel game={game} playerId={$playerID} />
        {/if}
      </div>
    </div>

    <div class="map-area">
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
      />
    </div>

    <div class="right-panel">
      {#if $selectedCell}
        <div class="cell-info">
          <h3>格子信息</h3>
          <p>坐标: ({$selectedCell.coord.q}, {$selectedCell.coord.r})</p>
          <p>岩石硬度: {$selectedCell.rock_hardness}</p>
          <p>含水量: {Math.round($selectedCell.water_content * 100)}%</p>
          <p>矿物: {$selectedCell.mineral_type}</p>
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
          {#if $selectedCell.is_wall && $selectedCell.owner === $playerID}
            <button class="btn-mine" on:click={() => handleMine($selectedCell)}>⛏️ 挖掘</button>
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

  .building-info, .units-info {
    margin-top: 10px;
    padding-top: 10px;
    border-top: 1px solid #34495e;
  }

  .building-info h4, .units-info h4 {
    color: #3498db;
    font-size: 0.9rem;
    margin-bottom: 8px;
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
