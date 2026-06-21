<script>
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';
  import { gameState, playerID } from '../stores/game.js';

  export let player;
  export let gameId;

  let units = [];
  let opponentPlayers = [];
  let dispatchSpyUnit = null;
  let showDispatchModal = false;

  onMount(async () => {
    try {
      units = await api.getUnits();
      await loadOpponents();
    } catch (e) {
      console.error('Failed to load units:', e);
    }
  });

  async function loadOpponents() {
    if (!gameId || !$playerID) return;
    try {
      const res = await api.getTechTree(gameId, $playerID);
      if (res?.opponent_players) {
        opponentPlayers = res.opponent_players;
      }
    } catch (e) {
      console.error('Failed to load opponents:', e);
    }
  }

  $: game = $gameState;
  $: currentPlayer = game?.players?.[$playerID] || player;

  const unitNames = {
    sapper: '工兵',
    infantry: '步兵',
    archer: '弓手',
    siege_ram: '攻城锤',
    iron_guard: '铁甲卫士',
    spore_grenadier: '孢子投弹兵',
    shadow_archer: '暗影射手',
    rock_giant: '岩石巨人',
    burrow_bomber: '钻地爆虫',
    shadow_scholar: '暗影学者'
  };

  const spyStatusNames = {
    idle: '待命',
    lurking: '潜伏中',
    successful: '窃取成功',
    caught: '被发现'
  };

  function isRaceSpecific(unit) {
    if (!unit?.race_specific) return true;
    return unit.race_specific === player?.race;
  }

  function getPlayerUnits() {
    if (!currentPlayer?.units) return [];
    return currentPlayer.units;
  }

  function getPlayerUnitCounts() {
    const counts = {};
    for (const unit of getPlayerUnits()) {
      counts[unit.type] = (counts[unit.type] || 0) + 1;
    }
    return counts;
  }

  function getShadowScholars() {
    return getPlayerUnits().filter(u => u.type === 'shadow_scholar');
  }

  function getSpyStatusDisplay(unit) {
    if (!unit.spy_status || unit.spy_status === 'idle') return '';
    if (unit.spy_status === 'lurking') {
      return `潜伏中 (${unit.spy_lurking_turns || 0}/3回合)`;
    }
    return spyStatusNames[unit.spy_status] || unit.spy_status;
  }

  function openDispatchModal(unit) {
    dispatchSpyUnit = unit;
    showDispatchModal = true;
  }

  async function dispatchSpy(targetId) {
    if (!dispatchSpyUnit || !targetId) return;
    try {
      await api.dispatchSpy(gameId, $playerID, dispatchSpyUnit.id, targetId);
      showDispatchModal = false;
      dispatchSpyUnit = null;
    } catch (e) {
      alert('派遣间谍失败: ' + (e.message || '未知错误'));
    }
  }

  $: playerUnitCounts = getPlayerUnitCounts();
  $: shadowScholars = getShadowScholars();
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

  {#if shadowScholars.length > 0}
    <div class="spy-section">
      <h4>🕵️ 暗影学者</h4>
      <div class="spy-list">
        {#each shadowScholars as unit}
          <div class="spy-card">
            <div class="spy-info">
              <span class="spy-name">暗影学者</span>
              {#if unit.spy_status && unit.spy_status !== 'idle'}
                <span class="spy-status {unit.spy_status}">
                  {getSpyStatusDisplay(unit)}
                </span>
              {:else}
                <span class="spy-status idle">待命</span>
              {/if}
            </div>
            {#if unit.spy_status === 'idle' || !unit.spy_status}
              <button 
                class="dispatch-btn"
                on:click={() => openDispatchModal(unit)}
                disabled={opponentPlayers.length === 0}
                title={opponentPlayers.length === 0 ? '没有可派遣的目标' : '派遣至对手领地窃取科技'}
              >
                派遣
              </button>
            {/if}
          </div>
        {/each}
      </div>
    </div>
  {/if}

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

  {#if showDispatchModal}
    <div class="modal-overlay" on:click={() => showDispatchModal = false}>
      <div class="modal-content" on:click|stopPropagation>
        <h3>🕵️ 派遣暗影学者</h3>
        <p class="modal-desc">
          派遣暗影学者潜入对手领地，潜伏3回合后有40%概率窃取对手当前研究科技的25%进度。
          <br><br>
          <span class="warning">⚠️ 失败或目标无研究时，间谍将被发现并处决！</span>
        </p>
        <div class="target-list">
          {#each opponentPlayers as opponent}
            <button 
              class="target-btn"
              on:click={() => dispatchSpy(opponent.player_id)}
            >
              <span class="target-color" style="background: {opponent.color}"></span>
              <span class="target-name">{opponent.username}</span>
              <span class="target-race">({opponent.race})</span>
            </button>
          {/each}
          {#if opponentPlayers.length === 0}
            <p class="no-targets">没有可派遣的目标</p>
          {/if}
        </div>
        <button class="close-btn" on:click={() => showDispatchModal = false}>取消</button>
      </div>
    </div>
  {/if}
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

  .spy-section {
    background: rgba(155, 89, 182, 0.1);
    border: 1px solid rgba(155, 89, 182, 0.3);
    border-radius: 6px;
    padding: 10px;
    margin-bottom: 12px;
  }

  .spy-section h4 {
    margin: 0 0 8px 0;
    color: #bb8fce;
    font-size: 0.85rem;
  }

  .spy-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .spy-card {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 6px 8px;
    background: rgba(0, 0, 0, 0.2);
    border-radius: 4px;
    font-size: 0.75rem;
  }

  .spy-info {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .spy-name {
    color: #ecf0f1;
    font-weight: bold;
  }

  .spy-status {
    padding: 2px 6px;
    border-radius: 3px;
    font-size: 0.65rem;
    font-weight: bold;
  }

  .spy-status.idle {
    background: rgba(52, 152, 219, 0.2);
    color: #3498db;
  }

  .spy-status.lurking {
    background: rgba(241, 196, 15, 0.2);
    color: #f1c40f;
    animation: pulse 1.5s ease-in-out infinite;
  }

  .spy-status.successful {
    background: rgba(46, 204, 113, 0.2);
    color: #2ecc71;
  }

  .spy-status.caught {
    background: rgba(231, 76, 60, 0.2);
    color: #e74c3c;
  }

  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.5; }
  }

  .dispatch-btn {
    background: rgba(155, 89, 182, 0.3);
    border: 1px solid rgba(155, 89, 182, 0.5);
    color: #bb8fce;
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 0.7rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .dispatch-btn:hover:not(:disabled) {
    background: rgba(155, 89, 182, 0.5);
  }

  .dispatch-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
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

  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
  }

  .modal-content {
    background: #16213e;
    border: 1px solid #34495e;
    border-radius: 10px;
    padding: 20px;
    min-width: 320px;
    max-width: 90%;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.5);
  }

  .modal-content h3 {
    margin: 0 0 12px 0;
    color: #9b59b6;
    font-size: 1.1rem;
  }

  .modal-desc {
    color: #bdc3c7;
    font-size: 0.85rem;
    line-height: 1.5;
    margin-bottom: 16px;
  }

  .warning {
    color: #e74c3c;
    font-weight: bold;
  }

  .target-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 16px;
  }

  .target-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 12px;
    background: rgba(52, 73, 94, 0.5);
    border: 1px solid #34495e;
    border-radius: 6px;
    color: #ecf0f1;
    font-size: 0.85rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .target-btn:hover {
    background: rgba(155, 89, 182, 0.2);
    border-color: #9b59b6;
  }

  .target-color {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .target-name {
    flex: 1;
    font-weight: bold;
  }

  .target-race {
    color: #7f8c8d;
    font-size: 0.75rem;
  }

  .no-targets {
    color: #7f8c8d;
    font-size: 0.85rem;
    text-align: center;
    padding: 20px;
  }

  .close-btn {
    width: 100%;
    padding: 10px;
    background: rgba(52, 73, 94, 0.5);
    border: 1px solid #34495e;
    border-radius: 6px;
    color: #bdc3c7;
    font-size: 0.85rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .close-btn:hover {
    background: rgba(52, 73, 94, 0.8);
  }
</style>
