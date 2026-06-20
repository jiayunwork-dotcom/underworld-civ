<script>
  export let game;

  let timeLeft = 0;
  let interval = null;

  $: {
    if (game?.planning_ends_at) {
      updateTimeLeft();
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

  if (typeof window !== 'undefined') {
    interval = setInterval(updateTimeLeft, 1000);
  }
</script>

<div class="turn-info">
  <div class="turn-number">
    第 <span class="num">{game?.current_turn || 0}</span> / {game?.max_turns || 50} 回合
  </div>
  <div class="phase">
    {#if game?.phase === 'planning'}
      <span class="phase-planning">📋 规划阶段</span>
      <span class="timer">⏱️ {formatTime(timeLeft)}</span>
    {:else if game?.phase === 'executing'}
      <span class="phase-executing">⚡ 执行中...</span>
    {:else if game?.phase === 'ended'}
      <span class="phase-ended">🏆 游戏结束</span>
    {/if}
  </div>
  {#if game?.phase === 'ended'}
    <div class="winner">
      胜利者: {game?.players?.[game.winner_id]?.username || '未知'}
      ({game?.victory_type === 'conquest' ? '征服' : game?.victory_type === 'tech' ? '科技' : '分数'}胜利)
    </div>
  {/if}
</div>

<style>
  .turn-info {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
  }

  .turn-number {
    color: #bdc3c7;
    font-size: 0.85rem;
  }

  .turn-number .num {
    color: #e94560;
    font-weight: bold;
    font-size: 1.1rem;
  }

  .phase {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .phase-planning {
    color: #2ecc71;
    font-weight: bold;
  }

  .phase-executing {
    color: #f39c12;
    font-weight: bold;
  }

  .phase-ended {
    color: #9b59b6;
    font-weight: bold;
  }

  .timer {
    color: #e74c3c;
    font-weight: bold;
    font-size: 1.1rem;
  }

  .winner {
    color: #f1c40f;
    font-size: 0.8rem;
    font-weight: bold;
  }
</style>
