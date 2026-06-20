<script>
  import { playerID, username } from '../stores/game.js';
  import { api } from '../utils/api.js';

  export let games = [];
  export let onGameCreated;
  export let onJoinGame;
  export let refreshGames;

  let showCreateForm = false;
  let gameName = '地下文明争霸';
  let maxPlayers = 6;
  let selectedRace = 'dwarf';
  let selectedColor = '#e74c3c';
  let races = [];

  const colors = [
    '#e74c3c', '#3498db', '#2ecc71', '#f39c12',
    '#9b59b6', '#1abc9c', '#e67e22', '#34495e'
  ];

  $: {
    loadRaces();
  }

  async function loadRaces() {
    try {
      races = await api.getRaces();
    } catch (e) {
      console.error('Failed to load races:', e);
    }
  }

  async function createGame() {
    try {
      const data = await api.createGame({
        name: gameName,
        max_players: maxPlayers,
        race: selectedRace,
        color: selectedColor
      }, $playerID, $username);
      onGameCreated(data);
    } catch (e) {
      alert('创建游戏失败');
    }
  }

  async function joinGame(gameId) {
    try {
      const data = await api.joinGame(gameId, {
        race: selectedRace,
        color: selectedColor
      }, $playerID, $username);
      onJoinGame(data);
    } catch (e) {
      alert('加入游戏失败');
    }
  }

  function getRaceName(race) {
    const raceMap = {
      dwarf: '矮人',
      mushroom: '蘑菇人',
      elf: '洞穴精灵',
      golem: '石像族',
      zerg: '深渊虫族'
    };
    return raceMap[race] || race;
  }
</script>

<div class="menu-container">
  <div class="menu-header">
    <h1>🏔️ 地下文明</h1>
    <p class="subtitle">Underworld Civilization</p>
    <p class="player-info">玩家: {$username}</p>
  </div>

  <div class="menu-content">
    <div class="section">
      <h2>🎮 创建游戏</h2>

      <div class="form-group">
        <label>游戏名称</label>
        <input type="text" bind:value={gameName} placeholder="输入游戏名称">
      </div>

      <div class="form-group">
        <label>玩家数量: {maxPlayers}人</label>
        <input type="range" min="4" max="8" bind:value={maxPlayers}>
      </div>

      <div class="form-group">
        <label>选择种族</label>
        <div class="race-selector">
          {#each races as race}
            <div class="race-option {selectedRace === race.race ? 'selected' : ''}"
                 on:click={() => selectedRace = race.race}>
              <div class="race-name">{race.name}</div>
              <div class="race-desc">{race.description}</div>
            </div>
          {/each}
        </div>
      </div>

      <div class="form-group">
        <label>选择颜色</label>
        <div class="color-picker">
          {#each colors as color}
            <div class="color-option {selectedColor === color ? 'selected' : ''}"
                 style="background: {color}"
                 on:click={() => selectedColor = color}>
            </div>
          {/each}
        </div>
      </div>

      <button class="btn btn-primary" on:click={createGame}>创建游戏</button>
    </div>

    <div class="section">
      <h2>🌍 游戏列表</h2>
      <button class="btn btn-small" on:click={refreshGames}>刷新</button>

      <div class="game-list">
        {#if games.length === 0}
          <p class="empty">暂无游戏，创建一个吧！</p>
        {:else}
          {#each games as game}
            <div class="game-card">
              <div class="game-name">{game.name}</div>
              <div class="game-info">
                <span class="status {game.status}">
                  {game.status === 'waiting' ? '等待中' : game.status === 'playing' ? '进行中' : '已结束'}
                </span>
                <span>{game.player_count}/{game.max_players} 人</span>
              </div>
              <div class="players-list">
                {#each game.players as p}
                  <span class="player-tag" style="color: {p.color}">{p.username} ({getRaceName(p.race)})</span>
                {/each}
              </div>
              {#if game.status === 'waiting' && game.player_count < game.max_players}
                <button class="btn btn-small" on:click={() => joinGame(game.id)}>加入</button>
              {/if}
            </div>
          {/each}
        {/if}
      </div>
    </div>
  </div>
</div>

<style>
  .menu-container {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px;
    overflow-y: auto;
    background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  }

  .menu-header {
    text-align: center;
    margin-bottom: 30px;
  }

  .menu-header h1 {
    font-size: 3rem;
    color: #e94560;
    margin-bottom: 10px;
  }

  .subtitle {
    color: #7f8c8d;
    font-size: 1.2rem;
    margin-bottom: 10px;
  }

  .player-info {
    color: #3498db;
    font-size: 0.9rem;
  }

  .menu-content {
    display: flex;
    gap: 40px;
    max-width: 1200px;
    width: 100%;
  }

  .section {
    flex: 1;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    padding: 24px;
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .section h2 {
    margin-bottom: 20px;
    color: #ecf0f1;
  }

  .form-group {
    margin-bottom: 16px;
  }

  .form-group label {
    display: block;
    margin-bottom: 8px;
    color: #bdc3c7;
    font-size: 0.9rem;
  }

  input[type="text"], input[type="range"] {
    width: 100%;
    padding: 10px 12px;
    border-radius: 6px;
    border: 1px solid #34495e;
    background: #1a1a2e;
    color: #ecf0f1;
    font-size: 1rem;
  }

  .race-selector {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .race-option {
    padding: 12px;
    border-radius: 8px;
    border: 2px solid transparent;
    background: rgba(52, 73, 94, 0.3);
    cursor: pointer;
    transition: all 0.2s;
  }

  .race-option:hover {
    background: rgba(52, 73, 94, 0.5);
  }

  .race-option.selected {
    border-color: #e94560;
    background: rgba(233, 69, 96, 0.2);
  }

  .race-name {
    font-weight: bold;
    color: #ecf0f1;
    margin-bottom: 4px;
  }

  .race-desc {
    font-size: 0.8rem;
    color: #95a5a6;
  }

  .color-picker {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }

  .color-option {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    cursor: pointer;
    border: 3px solid transparent;
    transition: transform 0.2s;
  }

  .color-option:hover {
    transform: scale(1.1);
  }

  .color-option.selected {
    border-color: white;
  }

  .btn {
    padding: 12px 24px;
    border-radius: 8px;
    border: none;
    cursor: pointer;
    font-size: 1rem;
    font-weight: bold;
    transition: all 0.2s;
  }

  .btn-primary {
    background: #e94560;
    color: white;
    width: 100%;
  }

  .btn-primary:hover {
    background: #c73650;
    transform: translateY(-2px);
  }

  .btn-small {
    padding: 6px 16px;
    font-size: 0.85rem;
    background: #3498db;
    color: white;
  }

  .btn-small:hover {
    background: #2980b9;
  }

  .game-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-top: 12px;
    max-height: 400px;
    overflow-y: auto;
  }

  .game-card {
    padding: 16px;
    background: rgba(52, 73, 94, 0.3);
    border-radius: 8px;
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .game-name {
    font-size: 1.1rem;
    font-weight: bold;
    color: #ecf0f1;
    margin-bottom: 8px;
  }

  .game-info {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
    font-size: 0.85rem;
    color: #95a5a6;
  }

  .status {
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 0.75rem;
  }

  .status.waiting {
    background: #27ae60;
    color: white;
  }

  .status.playing {
    background: #f39c12;
    color: white;
  }

  .status.finished {
    background: #95a5a6;
    color: white;
  }

  .player-tag {
    font-size: 0.8rem;
    margin-right: 8px;
  }

  .empty {
    text-align: center;
    color: #7f8c8d;
    padding: 40px 0;
  }

  @media (max-width: 800px) {
    .menu-content {
      flex-direction: column;
    }
  }
</style>
