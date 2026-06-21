<script>
  import { onMount, onDestroy } from 'svelte';
  import { api } from '../utils/api.js';

  export let game;
  export let playerId;
  export let gameId;

  let timeLeft = 0;
  let interval = null;
  let isReady = false;
  let flashState = false;
  let flashInterval = null;

  $: {
    if (game?.planning_ends_at) {
      updateTimeLeft();
    }
  }

  $: isUrgent = timeLeft <= 10 && timeLeft > 0 && game?.phase === 'planning';

  $: {
    if (isUrgent) {
      startFlash();
    } else {
      stopFlash();
    }
  }

  $: {
    if (game?.players && playerId) {
      const player = game.players[playerId];
      isReady = player?.ready || false;
    }
  }

  function updateTimeLeft() {
    const end = new Date(game.planning_ends_at);
    const now = new Date();
    timeLeft = Math.max(0, Math.floor((end - now) / 1000));
  }

  function formatTime(seconds) {
    const m = Math.floor(seconds / 60);
    const s = seconds % 60;
    return `${m}:${s.toString().padStart(2, '0')}`;
  }

  function startFlash() {
    if (flashInterval) return;
    flashInterval = setInterval(() => {
      flashState = !flashState;
    }, 500);
  }

  function stopFlash() {
    if (flashInterval) {
      clearInterval(flashInterval);
      flashInterval = null;
    }
    flashState = false;
  }

  async function toggleReady() {
    if (!gameId || !playerId) return;
    try {
      const action = isReady ? 'unready' : 'ready';
      await api.submitAction(gameId, playerId, action, {});
    } catch (e) {
      console.error('Failed to toggle ready:', e);
    }
  }

  onMount(() => {
    interval = setInterval(updateTimeLeft, 1000);
  });

  onDestroy(() => {
    if (interval) clearInterval(interval);
    stopFlash();
  });
</script>

<div class="turn-info">
  <div class="turn-header">
    <div class="turn-number">
      <span class="label">回合</span>
      <span class="num">{game?.current_turn || 0}</span>
      <span class="max">/ {game?.max_turns || 50}</span>
    </div>
  </div>

  <div class="phase-section">
    {#if game?.phase === 'planning'}
      <div class="phase-planning">
        <span class="phase-icon">📋</span>
        <span class="phase-text">规划阶段</span>
      </div>
    {:else if game?.phase === 'executing'}
      <div class="phase-executing">
        <span class="phase-icon">⚡</span>
        <span class="phase-text">执行中...</span>
      </div>
    {:else if game?.phase === 'ended'}
      <div class="phase-ended">
        <span class="phase-icon">🏆</span>
        <span class="phase-text">游戏结束</span>
      </div>
    {/if}
  </div>

  {#if game?.phase === 'planning'}
    <div class="timer-section">
      <div class="timer-display" class:urgent={isUrgent} class:flash={flashState && isUrgent}>
        <span class="timer-icon">⏱️</span>
        <span class="timer-value">{formatTime(timeLeft)}</span>
      </div>
      <button class="btn-end-turn {isReady ? 'ready' : ''}" on:click={toggleReady}>
        {isReady ? '✓ 已准备' : '结束回合'}
      </button>
    </div>
  {/if}

  {#if game?.phase === 'ended'}
    <div class="winner">
      <span class="winner-label">胜利者:</span>
      <span class="winner-name">{game?.players?.[game.winner_id]?.username || '未知'}</span>
      <span class="winner-type">
        ({game?.victory_type === 'conquest' ? '征服' : game?.victory_type === 'tech' ? '科技' : '分数'}胜利)
      </span>
    </div>
  {/if}
</div>

<style>
  .turn-info {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    min-width: 180px;
  }

  .turn-header {
    text-align: center;
  }

  .turn-number {
    display: flex;
    align-items: baseline;
    gap: 4px;
  }

  .turn-number .label {
    color: #7f8c8d;
    font-size: 0.8rem;
  }

  .turn-number .num {
    color: #e94560;
    font-weight: bold;
    font-size: 1.4rem;
  }

  .turn-number .max {
    color: #7f8c8d;
    font-size: 0.8rem;
  }

  .phase-section {
    display: flex;
    align-items: center;
  }

  .phase-planning,
  .phase-executing,
  .phase-ended {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .phase-icon {
    font-size: 1rem;
  }

  .phase-text {
    font-weight: bold;
    font-size: 0.9rem;
  }

  .phase-planning .phase-text {
    color: #2ecc71;
  }

  .phase-executing .phase-text {
    color: #f39c12;
  }

  .phase-ended .phase-text {
    color: #9b59b6;
  }

  .timer-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    width: 100%;
  }

  .timer-display {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 4px 12px;
    background: rgba(0, 0, 0, 0.3);
    border-radius: 6px;
    transition: all 0.3s;
  }

  .timer-icon {
    font-size: 0.9rem;
  }

  .timer-value {
    font-weight: bold;
    font-size: 1.1rem;
    color: #bdc3c7;
    font-variant-numeric: tabular-nums;
    min-width: 50px;
    text-align: center;
  }

  .timer-display.urgent {
    background: rgba(231, 76, 60, 0.3);
  }

  .timer-display.urgent .timer-value {
    color: #e74c3c;
  }

  .timer-display.flash .timer-value {
    opacity: 0.5;
  }

  .btn-end-turn {
    width: 100%;
    padding: 8px 16px;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-weight: bold;
    font-size: 0.9rem;
    transition: all 0.2s;
    background: #f39c12;
    color: white;
  }

  .btn-end-turn:hover {
    background: #e67e22;
    transform: translateY(-1px);
  }

  .btn-end-turn.ready {
    background: #27ae60;
  }

  .btn-end-turn.ready:hover {
    background: #2ecc71;
  }

  .winner {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
    padding: 6px 12px;
    background: rgba(241, 196, 15, 0.2);
    border-radius: 6px;
    border: 1px solid #f1c40f;
  }

  .winner-label {
    color: #bdc3c7;
    font-size: 0.75rem;
  }

  .winner-name {
    color: #f1c40f;
    font-weight: bold;
    font-size: 0.95rem;
  }

  .winner-type {
    color: #7f8c8d;
    font-size: 0.7rem;
  }
</style>
