<script>
  import { api } from '../utils/api.js';

  export let game;
  export let playerId;

  function getPlayers() {
    if (!game?.players) return [];
    return Object.values(game.players).filter(p => p.player_id !== playerId);
  }

  function getRelation(targetId) {
    const player = game?.players?.[playerId];
    return player?.diplomacy?.[targetId] || null;
  }

  function getStatusText(rel) {
    if (!rel) return '中立';
    const statuses = {
      neutral: '中立',
      alliance: '联盟',
      war: '战争',
      trade: '贸易'
    };
    return statuses[rel.status] || rel.status;
  }

  function getStatusColor(rel) {
    if (!rel) return '#95a5a6';
    const colors = {
      neutral: '#95a5a6',
      alliance: '#2ecc71',
      war: '#e74c3c',
      trade: '#3498db'
    };
    return colors[rel.status] || '#95a5a6';
  }

  function getBuildingCount(player) {
    return player.buildings ? Object.keys(player.buildings).length : 0;
  }

  function getTechCount(player) {
    return player.techs ? Object.keys(player.techs).length : 0;
  }

  function formatResources(res) {
    const parts = [];
    if (res?.stone) parts.push(`🪨${res.stone}`);
    if (res?.metal) parts.push(`⚙️${res.metal}`);
    if (res?.glow_mushroom) parts.push(`🍄${res.glow_mushroom}`);
    if (res?.water) parts.push(`💧${res.water}`);
    if (res?.magic_crystal) parts.push(`💎${res.magic_crystal}`);
    if (res?.fossil_fuel) parts.push(`🔥${res.fossil_fuel}`);
    return parts.join(' ') || '无';
  }

  async function proposeAlliance(targetId) {
    try {
      await api.submitAction(game.id, playerId, 'propose_alliance', { to_player: targetId });
    } catch (e) {
      console.error('Failed to propose alliance:', e);
    }
  }

  async function declareWar(targetId) {
    if (!confirm('确定要宣战吗？解除所有条约会有3回合冷却。')) return;
    try {
      await api.submitAction(game.id, playerId, 'declare_war', { to_player: targetId });
    } catch (e) {
      console.error('Failed to declare war:', e);
    }
  }

  async function proposeTrade(targetId) {
    const stone = prompt('给出石材数量:', '0');
    if (stone === null) return;
    try {
      await api.submitAction(game.id, playerId, 'trade_offer', {
        to_player: targetId,
        offer: { stone: parseInt(stone) || 0, metal: 0, glow_mushroom: 0, water: 0, magic_crystal: 0, fossil_fuel: 0 },
        demand: { stone: 0, metal: 10, glow_mushroom: 0, water: 0, magic_crystal: 0, fossil_fuel: 0 }
      });
    } catch (e) {
      console.error('Failed to propose trade:', e);
    }
  }

  async function acceptTrade(tradeId) {
    try {
      await api.submitAction(game.id, playerId, 'accept_trade', { trade_id: tradeId });
    } catch (e) {
      console.error('Failed to accept trade:', e);
    }
  }

  function rejectTrade(tradeId) {
    alert('交易已拒绝');
  }

  const raceNames = {
    dwarf: '矮人',
    mushroom: '蘑菇人',
    elf: '洞穴精灵',
    golem: '石像族',
    zerg: '深渊虫族'
  };

  $: otherPlayers = getPlayers();
</script>

<div class="diplomacy-panel">
  <h3>🤝 外交</h3>

  <div class="player-list">
    {#if otherPlayers.length === 0}
      <p class="empty">暂无其他玩家</p>
    {:else}
      {#each otherPlayers as player}
        {#if !player.eliminated}
          <div class="player-card">
            <div class="player-header">
              <div class="player-avatar" style="background: {player.color}">
                {player.username?.charAt(0).toUpperCase()}
              </div>
              <div class="player-info">
                <div class="player-name">{player.username}</div>
                <div class="player-race">{raceNames[player.race] || player.race}</div>
              </div>
            </div>

            <div class="relation-status" style="color: {getStatusColor(getRelation(player.player_id))}">
              {getStatusText(getRelation(player.player_id))}
            </div>

            <div class="player-stats">
              <span>🏠 {getBuildingCount(player)}</span>
              <span>⚔️ {player.units?.length || 0}</span>
              <span>🔬 {getTechCount(player)}</span>
            </div>

            <div class="diplo-actions">
              <button class="btn-small" on:click={() => proposeAlliance(player.player_id)}>
                🤝 结盟
              </button>
              <button class="btn-small" on:click={() => declareWar(player.player_id)}>
                ⚔️ 宣战
              </button>
              <button class="btn-small" on:click={() => proposeTrade(player.player_id)}>
                💰 贸易
              </button>
            </div>
          </div>
        {/if}
      {/each}
    {/if}
  </div>

  {#if game?.players?.[playerId]?.trade_offers?.length > 0}
    <div class="trade-offers">
      <h4>📨 交易请求</h4>
      {#each game.players[playerId].trade_offers as offer}
        <div class="trade-offer">
          <div class="offer-from">来自: {game.players[offer.from_player]?.username || '未知'}</div>
          <div class="offer-content">
            <div>给出: {formatResources(offer.offer)}</div>
            <div>要求: {formatResources(offer.demand)}</div>
          </div>
          <div class="offer-actions">
            <button class="btn-accept" on:click={() => acceptTrade(offer.id)}>接受</button>
            <button class="btn-reject" on:click={() => rejectTrade(offer.id)}>拒绝</button>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .diplomacy-panel h3 {
    color: #ecf0f1;
    margin-bottom: 12px;
    font-size: 1rem;
  }

  .player-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .empty {
    color: #7f8c8d;
    text-align: center;
    padding: 20px;
    font-style: italic;
  }

  .player-card {
    padding: 12px;
    background: rgba(52, 73, 94, 0.3);
    border-radius: 8px;
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .player-header {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 8px;
  }

  .player-avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-weight: bold;
  }

  .player-info {
    flex: 1;
  }

  .player-name {
    font-weight: bold;
    color: #ecf0f1;
  }

  .player-race {
    font-size: 0.75rem;
    color: #95a5a6;
  }

  .relation-status {
    font-size: 0.8rem;
    font-weight: bold;
    margin-bottom: 8px;
  }

  .player-stats {
    display: flex;
    gap: 10px;
    font-size: 0.75rem;
    color: #bdc3c7;
    margin-bottom: 10px;
  }

  .diplo-actions {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
  }

  .btn-small {
    flex: 1;
    padding: 6px 8px;
    font-size: 0.7rem;
    background: #3498db;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .btn-small:hover {
    background: #2980b9;
  }

  .trade-offers {
    margin-top: 16px;
  }

  .trade-offers h4 {
    color: #f39c12;
    margin-bottom: 10px;
    font-size: 0.9rem;
  }

  .trade-offer {
    padding: 10px;
    background: rgba(243, 156, 18, 0.1);
    border-radius: 6px;
    margin-bottom: 8px;
    border: 1px solid rgba(243, 156, 18, 0.3);
  }

  .offer-from {
    font-size: 0.8rem;
    color: #f39c12;
    margin-bottom: 6px;
  }

  .offer-content {
    font-size: 0.75rem;
    color: #bdc3c7;
    margin-bottom: 8px;
  }

  .offer-actions {
    display: flex;
    gap: 6px;
  }

  .btn-accept {
    flex: 1;
    padding: 6px;
    font-size: 0.75rem;
    background: #2ecc71;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .btn-reject {
    flex: 1;
    padding: 6px;
    font-size: 0.75rem;
    background: #e74c3c;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
</style>
