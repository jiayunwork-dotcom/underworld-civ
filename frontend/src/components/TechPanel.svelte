<script>
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';
  import { gameState, playerID } from '../stores/game.js';

  export let player;
  export let gameId;

  let techTree = [];
  let hoveredTech = null;
  let hoverPosition = { x: 0, y: 0 };
  let isSettingResearch = false;

  const categoryInfo = {
    military: { name: '军事线', icon: '⚔️', color: '#e74c3c' },
    economy:  { name: '经济线', icon: '💰', color: '#f39c12' },
    mining:   { name: '挖掘线', icon: '⛏️', color: '#3498db' },
    special:  { name: '特殊线', icon: '✨', color: '#9b59b6' }
  };

  const categories = ['military', 'economy', 'mining', 'special'];

  $: game = $gameState;
  $: currentPlayer = game?.players?.[$playerID] || player;

  onMount(async () => {
    if (gameId && $playerID) {
      try {
        const res = await api.getTechTree(gameId, $playerID);
        if (res?.my_tech_tree) {
          techTree = res.my_tech_tree;
        }
      } catch (e) {
        console.error('Failed to load tech tree:', e);
      }
    }
  });

  $: {
    if (game && currentPlayer && techTree.length > 0) {
      for (let t of techTree) {
        t.researched = currentPlayer.techs?.[t.id] || false;
        t.is_current = currentPlayer.current_research === t.id;
        t.progress = (currentPlayer.tech_progresses && currentPlayer.tech_progresses[t.id]) || 0;
        t.can_research = !t.researched && (t.prerequisites?.length === 0 ||
          t.prerequisites.every(p => currentPlayer.techs?.[p]));
      }
    }
  }

  function getTechsByCategory(cat) {
    return techTree.filter(t => t.category === cat).sort((a, b) => a.tier - b.tier);
  }

  function isResearched(techId) {
    return currentPlayer?.techs?.[techId] || false;
  }

  function isCurrentResearch(techId) {
    return currentPlayer?.current_research === techId;
  }

  function getProgress(tech) {
    if (tech.researched) return 100;
    const p = (currentPlayer?.tech_progresses?.[tech.id]) || 0;
    return Math.min(100, (p / tech.cost) * 100);
  }

  function getCurrentProgress() {
    if (!currentPlayer?.current_research) return { progress: 0, cost: 0, name: '' };
    const tech = techTree.find(t => t.id === currentPlayer.current_research);
    if (!tech) return { progress: 0, cost: 0, name: '' };
    const p = currentPlayer.tech_progresses?.[tech.id] || 0;
    return {
      progress: Math.min(100, (p / tech.cost) * 100),
      cost: tech.cost,
      name: tech.name,
      current: p
    };
  }

  function showTooltip(tech, event) {
    const rect = event.currentTarget.getBoundingClientRect();
    const panel = event.currentTarget.closest('.tech-panel');
    const panelRect = panel?.getBoundingClientRect();
    hoveredTech = tech;
    hoverPosition = {
      x: rect.left - (panelRect?.left || 0) + rect.width / 2,
      y: rect.top - (panelRect?.top || 0) - 8
    };
  }

  function hideTooltip() {
    hoveredTech = null;
  }

  async function handleTechClick(tech) {
    if (!tech.can_research || tech.researched || isSettingResearch) return;
    isSettingResearch = true;
    try {
      await api.setResearch(gameId, $playerID, tech.id);
      try {
        const res = await api.getTechTree(gameId, $playerID);
        if (res?.my_tech_tree) {
          techTree = res.my_tech_tree;
        }
      } catch (e) {}
    } catch (e) {
      alert('切换研究目标失败: ' + (e.message || '未知错误'));
    } finally {
      isSettingResearch = false;
    }
  }

  function getStatusClass(tech) {
    if (tech.researched) return 'researched';
    if (tech.is_current) return 'researching';
    if (!tech.can_research) return 'locked';
    return 'available';
  }

  function getTooltipPrereqText(tech) {
    if (!tech.prerequisites || tech.prerequisites.length === 0) return '无';
    return tech.prerequisites.map(p => {
      const prereqTech = techTree.find(t => t.id === p);
      const name = prereqTech?.name || p;
      const done = isResearched(p);
      return `${done ? '✅' : '❌'} ${name}`;
    }).join(', ');
  }

  $: totalCount = techTree.length;
  $: researchedCount = techTree.filter(t => t.researched).length;
  $: curProgress = getCurrentProgress();
</script>

<div class="tech-panel">
  <div class="panel-header">
    <h3>🔬 科技研究中心</h3>
    <div class="overall-progress">
      <div class="op-bar">
        <div class="op-fill" style="width: {totalCount ? (researchedCount / totalCount * 100) : 0}%"></div>
      </div>
      <span class="op-text">{researchedCount}/{totalCount}</span>
    </div>
  </div>

  {#if currentPlayer?.current_research}
    <div class="current-research-bar">
      <div class="crb-info">
        <span class="crb-label">🔬 研究中:</span>
        <span class="crb-name">{curProgress.name}</span>
        <span class="crb-points">{curProgress.current}/{curProgress.cost} 点</span>
      </div>
      <div class="crb-progress">
        <div class="crb-fill" style="width: {curProgress.progress}%"></div>
      </div>
    </div>
  {/if}

  <div class="tech-grid-container">
    {#each categories as cat}
      <div class="tech-column">
        <div class="column-header" style="border-color: {categoryInfo[cat].color}">
          <span class="ch-icon">{categoryInfo[cat].icon}</span>
          <span class="ch-name">{categoryInfo[cat].name}</span>
        </div>
        <div class="tech-nodes">
          {#each getTechsByCategory(cat) as tech, idx}
            <div
              class="tech-node {getStatusClass(tech)}"
              style="--cat-color: {categoryInfo[cat].color}"
              on:mouseenter={(e) => showTooltip(tech, e)}
              on:mouseleave={hideTooltip}
              on:click={() => handleTechClick(tech)}
              title=""
            >
              <div class="node-tier">{tech.tier}</div>
              <div class="node-content">
                <div class="node-name">{tech.name}</div>
                <div class="node-cost">💎{tech.cost}</div>
              </div>
              <div class="node-progress">
                <div class="np-fill" style="width: {getProgress(tech)}%"></div>
              </div>
              {#if tech.is_current}
                <div class="pulse-ring"></div>
              {/if}
            </div>
            {#if idx < getTechsByCategory(cat).length - 1}
              <div class="connector">
                <div class="connector-line {getTechsByCategory(cat)[idx].researched ? 'done' : ''}"></div>
              </div>
            {/if}
          {/each}
        </div>
      </div>
    {/each}
  </div>

  {#if hoveredTech}
    <div
      class="tech-tooltip"
      style="left: {hoverPosition.x}px; top: {hoverPosition.y}px;"
    >
      <div class="tt-header" style="border-left-color: {categoryInfo[hoveredTech.category]?.color}">
        <span class="tt-category-icon">{categoryInfo[hoveredTech.category]?.icon}</span>
        <span class="tt-name">{hoveredTech.name}</span>
        <span class="tt-tier">T{hoveredTech.tier}</span>
      </div>
      <div class="tt-desc">{hoveredTech.description}</div>
      <div class="tt-stats">
        <div class="tt-stat">
          <span class="tt-s-label">总研究点数:</span>
          <span class="tt-s-value">💎 {hoveredTech.cost}</span>
        </div>
        <div class="tt-stat">
          <span class="tt-s-label">当前累积:</span>
          <span class="tt-s-value">{hoveredTech.researched ? '✅ 已完成' : ((currentPlayer?.tech_progresses?.[hoveredTech.id]) || 0) + ' 点'}</span>
        </div>
        <div class="tt-stat">
          <span class="tt-s-label">进度:</span>
          <span class="tt-s-value">{getProgress(hoveredTech).toFixed(1)}%</span>
        </div>
      </div>
      <div class="tt-progress">
        <div class="ttp-fill" style="width: {getProgress(hoveredTech)}%"></div>
      </div>
      <div class="tt-prereq">
        <span class="tt-p-label">前置科技:</span>
        <span class="tt-p-value">{getTooltipPrereqText(hoveredTech)}</span>
      </div>
      <div class="tt-status">
        {#if hoveredTech.researched}
          <span class="tt-status-tag done">✅ 已研究完成</span>
        {:else if hoveredTech.is_current}
          <span class="tt-status-tag current">⚡ 正在研究</span>
        {:else if hoveredTech.can_research}
          <span class="tt-status-tag available">🔓 可研究 - 点击切换</span>
        {:else}
          <span class="tt-status-tag locked">🔒 未解锁</span>
        {/if}
      </div>
    </div>
  {/if}
</div>

<style>
  .tech-panel {
    position: relative;
    background: rgba(22, 33, 62, 0.95);
    border-radius: 8px;
    padding: 14px;
    color: #ecf0f1;
    max-height: 550px;
    overflow-y: auto;
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
    padding-bottom: 10px;
    border-bottom: 1px solid #34495e;
  }

  .panel-header h3 {
    margin: 0;
    font-size: 1rem;
    color: #ecf0f1;
  }

  .overall-progress {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .op-bar {
    width: 100px;
    height: 8px;
    background: rgba(0,0,0,0.3);
    border-radius: 4px;
    overflow: hidden;
  }

  .op-fill {
    height: 100%;
    background: linear-gradient(90deg, #2ecc71, #27ae60);
    transition: width 0.3s;
  }

  .op-text {
    font-size: 0.75rem;
    color: #2ecc71;
    font-weight: bold;
    min-width: 45px;
  }

  .current-research-bar {
    background: rgba(243, 156, 18, 0.15);
    border: 1px solid rgba(243, 156, 18, 0.4);
    border-radius: 6px;
    padding: 10px;
    margin-bottom: 14px;
  }

  .crb-info {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 6px;
    font-size: 0.8rem;
    flex-wrap: wrap;
  }

  .crb-label {
    color: #f39c12;
    font-weight: bold;
  }

  .crb-name {
    color: #ecf0f1;
    font-weight: bold;
  }

  .crb-points {
    margin-left: auto;
    color: #bdc3c7;
    font-size: 0.75rem;
  }

  .crb-progress {
    height: 6px;
    background: rgba(0,0,0,0.3);
    border-radius: 3px;
    overflow: hidden;
  }

  .crb-fill {
    height: 100%;
    background: linear-gradient(90deg, #f39c12, #e67e22);
    transition: width 0.4s;
    animation: pulseProgress 1.5s ease-in-out infinite;
  }

  @keyframes pulseProgress {
    0%, 100% { filter: brightness(1); }
    50% { filter: brightness(1.3); }
  }

  .tech-grid-container {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 10px;
  }

  .tech-column {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .column-header {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 6px;
    background: rgba(0,0,0,0.3);
    border-radius: 4px;
    border-left: 3px solid;
    margin-bottom: 8px;
    font-size: 0.75rem;
    font-weight: bold;
  }

  .ch-icon {
    font-size: 1rem;
  }

  .ch-name {
    flex: 1;
  }

  .tech-nodes {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .tech-node {
    position: relative;
    padding: 8px 6px;
    border-radius: 6px;
    background: rgba(52, 73, 94, 0.4);
    border: 2px solid #34495e;
    cursor: default;
    transition: all 0.25s;
    overflow: hidden;
  }

  .tech-node.available {
    cursor: pointer;
    border-color: rgba(255,255,255,0.6);
    background: rgba(255,255,255,0.08);
  }

  .tech-node.available:hover {
    transform: translateY(-1px);
    border-color: #fff;
    box-shadow: 0 2px 10px rgba(255,255,255,0.15);
    background: rgba(255,255,255,0.12);
  }

  .tech-node.researched {
    border-color: #2ecc71;
    background: rgba(46, 204, 113, 0.2);
  }

  .tech-node.researching {
    border-color: #f1c40f;
    background: rgba(241, 196, 15, 0.18);
    animation: researchingGlow 1.2s ease-in-out infinite;
  }

  .tech-node.locked {
    opacity: 0.4;
    filter: grayscale(0.6);
    border-color: #2c3e50;
  }

  @keyframes researchingGlow {
    0%, 100% {
      box-shadow: 0 0 5px rgba(241, 196, 15, 0.3),
                  inset 0 0 5px rgba(241, 196, 15, 0.1);
    }
    50% {
      box-shadow: 0 0 15px rgba(241, 196, 15, 0.6),
                  inset 0 0 15px rgba(241, 196, 15, 0.2);
    }
  }

  .pulse-ring {
    position: absolute;
    top: 0; left: 0; right: 0; bottom: 0;
    border-radius: 6px;
    border: 2px solid #f1c40f;
    animation: pulseRing 1.2s ease-out infinite;
    pointer-events: none;
  }

  @keyframes pulseRing {
    0% {
      opacity: 0.8;
      transform: scale(1);
    }
    100% {
      opacity: 0;
      transform: scale(1.15);
    }
  }

  .node-tier {
    position: absolute;
    top: 2px;
    right: 4px;
    font-size: 0.55rem;
    color: #95a5a6;
    font-weight: bold;
  }

  .node-content {
    display: flex;
    flex-direction: column;
    gap: 2px;
    margin-bottom: 4px;
  }

  .node-name {
    font-size: 0.72rem;
    font-weight: bold;
    color: #ecf0f1;
    line-height: 1.2;
    padding-right: 18px;
  }

  .node-cost {
    font-size: 0.6rem;
    color: #9b59b6;
    font-weight: bold;
  }

  .node-progress {
    height: 3px;
    background: rgba(0,0,0,0.4);
    border-radius: 2px;
    overflow: hidden;
  }

  .np-fill {
    height: 100%;
    background: var(--cat-color, #3498db);
    transition: width 0.3s;
  }

  .tech-node.researched .np-fill {
    background: #2ecc71;
  }

  .tech-node.researching .np-fill {
    background: #f1c40f;
  }

  .connector {
    display: flex;
    justify-content: center;
    height: 8px;
    padding: 0;
  }

  .connector-line {
    width: 2px;
    height: 100%;
    background: #34495e;
  }

  .connector-line.done {
    background: #2ecc71;
  }

  .tech-tooltip {
    position: absolute;
    transform: translate(-50%, -100%);
    width: 260px;
    background: #0d1421;
    border: 1px solid #34495e;
    border-radius: 8px;
    padding: 12px;
    box-shadow: 0 8px 24px rgba(0,0,0,0.8),
                0 0 0 1px rgba(155, 89, 182, 0.3);
    z-index: 1000;
    pointer-events: none;
    animation: fadeInTooltip 0.15s ease-out;
  }

  @keyframes fadeInTooltip {
    from { opacity: 0; transform: translate(-50%, calc(-100% - 5px)); }
    to { opacity: 1; transform: translate(-50%, -100%); }
  }

  .tt-header {
    display: flex;
    align-items: center;
    gap: 6px;
    padding-bottom: 8px;
    margin-bottom: 8px;
    border-bottom: 1px solid #34495e;
    border-left: 3px solid;
    padding-left: 8px;
  }

  .tt-category-icon {
    font-size: 1rem;
  }

  .tt-name {
    flex: 1;
    font-weight: bold;
    font-size: 0.9rem;
    color: #ecf0f1;
  }

  .tt-tier {
    font-size: 0.7rem;
    padding: 2px 6px;
    background: rgba(155, 89, 182, 0.3);
    border-radius: 4px;
    color: #bb8fce;
    font-weight: bold;
  }

  .tt-desc {
    font-size: 0.78rem;
    color: #bdc3c7;
    margin-bottom: 10px;
    line-height: 1.4;
  }

  .tt-stats {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-bottom: 8px;
  }

  .tt-stat {
    display: flex;
    justify-content: space-between;
    font-size: 0.72rem;
  }

  .tt-s-label {
    color: #95a5a6;
  }

  .tt-s-value {
    color: #ecf0f1;
    font-weight: bold;
  }

  .tt-progress {
    height: 5px;
    background: rgba(0,0,0,0.4);
    border-radius: 3px;
    overflow: hidden;
    margin-bottom: 10px;
  }

  .ttp-fill {
    height: 100%;
    background: linear-gradient(90deg, #9b59b6, #8e44ad);
    transition: width 0.3s;
  }

  .tt-prereq {
    display: flex;
    gap: 6px;
    font-size: 0.7rem;
    margin-bottom: 8px;
    flex-wrap: wrap;
  }

  .tt-p-label {
    color: #95a5a6;
    flex-shrink: 0;
  }

  .tt-p-value {
    color: #ecf0f1;
    line-height: 1.4;
  }

  .tt-status {
    display: flex;
    justify-content: center;
  }

  .tt-status-tag {
    font-size: 0.75rem;
    font-weight: bold;
    padding: 5px 10px;
    border-radius: 4px;
    text-align: center;
  }

  .tt-status-tag.done {
    background: rgba(46, 204, 113, 0.2);
    color: #2ecc71;
  }

  .tt-status-tag.current {
    background: rgba(241, 196, 15, 0.2);
    color: #f1c40f;
    animation: statusPulse 1.2s ease-in-out infinite;
  }

  .tt-status-tag.available {
    background: rgba(52, 152, 219, 0.2);
    color: #5dade2;
  }

  .tt-status-tag.locked {
    background: rgba(127, 140, 141, 0.2);
    color: #95a5a6;
  }

  @keyframes statusPulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.6; }
  }

  .tech-panel::-webkit-scrollbar {
    width: 6px;
  }
  .tech-panel::-webkit-scrollbar-track {
    background: rgba(0,0,0,0.2);
    border-radius: 3px;
  }
  .tech-panel::-webkit-scrollbar-thumb {
    background: #34495e;
    border-radius: 3px;
  }
</style>
