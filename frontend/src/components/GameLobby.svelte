<script>
  import { onMount, onDestroy } from 'svelte';
  import { api } from '../utils/api.js';
  import { connectWS } from '../utils/api.js';

  export let game;
  export let playerId;
  export let onGameStarted;
  export let onBack;

  let ws = null;
  let players = [];
  let isHost = false;

  $: {
    if (game) {
      players = Object.values(game.players || {});
      isHost = game.players?.[playerId]?.is_host || false;
    }
  }

  onMount(() => {
    if (game?.id) {
      ws = connectWS(game.id, playerId, handleMessage);
    }
  });

  onDestroy(() => {
    if (ws) ws.close();
  });

  function handleMessage(msg) {
    if (msg.type === 'game_state') {
      game = msg.data;
      if (msg.data.status === 'playing' && msg.data.phase !== 'ended') {
        onGameStarted(msg.data);
      }
    }
  }

  async function startGame() {
    try {
      const res = await api.startGame(game.id, playerId);
      if (res.game) {
        onGameStarted(res.game);
      }
    } catch (e) {
      alert('开始游戏失败');
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

  const layerNames = [
    '浅层 - 蘑菇林带',
    '中层 - 石灰岩洞窟',
    '深层 - 金属矿脉',
    '熔岩层 - 魔晶之地',
    '远古层 - 化石遗迹'
  ];
</script>

<div class="lobby-container">
  <div class="lobby-header">
    <button class="btn-back" on:click={onBack}>← 返回</button>
    <h1>{game?.name || '游戏大厅'}</h1>
    <div class="game-status">
      状态: <span class={game?.status}>{game?.status === 'waiting' ? '等待中' : '进行中'}</span>
    </div>
  </div>

  <div class="lobby-content">
    <div class="players-section">
      <h2>👥 玩家列表 ({players.length}/{game?.max_players})</h2>
      <div class="player-list">
        {#each players as player}
          <div class="player-card" style="border-color: {player.color}">
            <div class="player-avatar" style="background: {player.color}">
              {player.username?.charAt(0).toUpperCase()}
            </div>
            <div class="player-info">
              <div class="player-name">{player.username}</div>
              <div class="player-race">{getRaceName(player.race)}</div>
            </div>
            {#if player.is_host}
              <span class="host-badge">房主</span>
            {/if}
          </div>
        {/each}
      </div>
    </div>

    <div class="info-section">
      <h2>📜 游戏规则</h2>
      <div class="rules">
        <p><strong>目标：</strong>发展你的地下文明，征服所有对手</p>
        <p><strong>胜利条件：</strong></p>
        <ul>
          <li>征服胜利：摧毁所有敌方主基地</li>
          <li>科技胜利：率先研究全部科技</li>
          <li>分数胜利：50回合后总分最高</li>
        </ul>
        <p><strong>地图：</strong>5层地下洞穴网络</p>
        <ul>
          {#each layerNames as layer, i}
            <li>{layer}</li>
          {/each}
        </ul>
        <p><strong>回合：</strong>每回合90秒规划阶段</p>
      </div>

      {#if isHost && game?.status === 'waiting'}
        <button class="btn btn-start" on:click={startGame}>
          开始游戏
        </button>
      {:else if !isHost}
        <p class="waiting-text">等待房主开始游戏...</p>
      {/if}
    </div>
  </div>
</div>

<style>
  .lobby-container {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
    padding: 20px;
  }

  .lobby-header {
    display: flex;
    align-items: center;
    gap: 20px;
    margin-bottom: 30px;
  }

  .lobby-header h1 {
    flex: 1;
    color: #e94560;
    font-size: 2rem;
  }

  .btn-back {
    background: #34495e;
    color: white;
    border: none;
    padding: 10px 20px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 1rem;
  }

  .btn-back:hover {
    background: #2c3e50;
  }

  .game-status {
    color: #ecf0f1;
  }

  .game-status .waiting {
    color: #27ae60;
    font-weight: bold;
  }

  .game-status .playing {
    color: #f39c12;
    font-weight: bold;
  }

  .lobby-content {
    display: flex;
    gap: 30px;
    flex: 1;
  }

  .players-section, .info-section {
    flex: 1;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    padding: 24px;
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .players-section h2, .info-section h2 {
    margin-bottom: 20px;
    color: #ecf0f1;
  }

  .player-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .player-card {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px;
    background: rgba(52, 73, 94, 0.3);
    border-radius: 8px;
    border-left: 4px solid #e94560;
  }

  .player-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-weight: bold;
    font-size: 1.2rem;
  }

  .player-info {
    flex: 1;
  }

  .player-name {
    font-weight: bold;
    color: #ecf0f1;
    font-size: 1.1rem;
  }

  .player-race {
    color: #95a5a6;
    font-size: 0.9rem;
  }

  .host-badge {
    background: #f39c12;
    color: white;
    padding: 4px 12px;
    border-radius: 4px;
    font-size: 0.75rem;
    font-weight: bold;
  }

  .rules {
    color: #bdc3c7;
    line-height: 1.8;
  }

  .rules ul {
    margin-left: 20px;
    margin-bottom: 16px;
  }

  .rules p {
    margin-bottom: 8px;
  }

  .btn-start {
    width: 100%;
    padding: 16px;
    background: #e94560;
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 1.2rem;
    font-weight: bold;
    cursor: pointer;
    margin-top: 20px;
    transition: all 0.2s;
  }

  .btn-start:hover {
    background: #c73650;
    transform: translateY(-2px);
  }

  .waiting-text {
    text-align: center;
    color: #7f8c8d;
    margin-top: 20px;
    font-style: italic;
  }

  @media (max-width: 800px) {
    .lobby-content {
      flex-direction: column;
    }
  }
</style>
