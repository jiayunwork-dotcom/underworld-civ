<script>
  import { onMount } from 'svelte';
  import { api } from '../utils/api.js';

  export let player;

  let techs = [];
  let selectedCategory = 'all';

  const categories = [
    { id: 'all', name: '全部' },
    { id: 'mining', name: '采矿' },
    { id: 'industry', name: '工业' },
    { id: 'biology', name: '生物' },
    { id: 'military', name: '军事' },
    { id: 'magic', name: '魔法' },
    { id: 'energy', name: '能源' },
    { id: 'civil', name: '民政' },
    { id: 'exploration', name: '探索' },
    { id: 'economy', name: '经济' },
    { id: 'research', name: '研究' },
    { id: 'defense', name: '防御' }
  ];

  onMount(async () => {
    try {
      techs = await api.getTechs();
    } catch (e) {
      console.error('Failed to load techs:', e);
    }
  });

  function isResearched(techId) {
    return player?.techs?.[techId] || false;
  }

  function canResearch(tech) {
    if (!tech.prerequisites || tech.prerequisites.length === 0) return true;
    return tech.prerequisites.every(p => player?.techs?.[p]);
  }

  function isResearching(techId) {
    return player?.current_research === techId;
  }

  function getFilteredTechs() {
    if (selectedCategory === 'all') return techs;
    return techs.filter(t => t.category === selectedCategory);
  }

  $: filteredTechs = getFilteredTechs();
  $: researchedCount = player?.techs ? Object.keys(player.techs).length : 0;
</script>

<div class="tech-panel">
  <h3>🔬 科技</h3>

  <div class="tech-progress">
    <div class="progress-bar">
      <div class="progress-fill" style="width: {(researchedCount / techs.length * 100) || 0}%"></div>
    </div>
    <span class="progress-text">{researchedCount}/{techs.length}</span>
  </div>

  {#if player?.current_research}
    <div class="current-research">
      <strong>正在研究:</strong>
      <span>{getTechName(player.current_research)}</span>
      <div class="research-progress">
        <div class="rp-fill" style="width: {(player.research_progress / getTechCost(player.current_research) * 100) || 0}%"></div>
      </div>
    </div>
  {/if}

  <div class="category-tabs">
    {#each categories as cat}
      <button class={selectedCategory === cat.id ? 'active' : ''}
              on:click={() => selectedCategory = cat.id}>
        {cat.name}
      </button>
    {/each}
  </div>

  <div class="tech-list">
    {#each filteredTechs as tech}
      <div class="tech-card
                  {isResearched(tech.id) ? 'researched' : ''}
                  {!canResearch(tech) && !isResearched(tech.id) ? 'locked' : ''}
                  {isResearching(tech.id) ? 'researching' : ''}">
        <div class="tech-name">
          {isResearched(tech.id) ? '✅' : canResearch(tech) ? '🔓' : '🔒'}
          {tech.name}
        </div>
        <p class="tech-desc">{tech.description}</p>
        <div class="tech-info">
          <span class="cost">💎 {tech.cost}</span>
          <span class="category">{getCategoryName(tech.category)}</span>
        </div>
        {#if tech.prerequisites && tech.prerequisites.length > 0}
          <div class="prereq">
            需要: {tech.prerequisites.map(p => getTechName(p)).join(', ')}
          </div>
        {/if}
      </div>
    {/each}
  </div>
</div>

<style>
  .tech-panel h3 {
    color: #ecf0f1;
    margin-bottom: 12px;
    font-size: 1rem;
  }

  .tech-progress {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 12px;
  }

  .progress-bar {
    flex: 1;
    height: 8px;
    background: rgba(0, 0, 0, 0.3);
    border-radius: 4px;
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, #9b59b6, #8e44ad);
    transition: width 0.3s;
  }

  .progress-text {
    font-size: 0.8rem;
    color: #9b59b6;
    font-weight: bold;
  }

  .current-research {
    background: rgba(155, 89, 182, 0.2);
    padding: 10px;
    border-radius: 6px;
    margin-bottom: 12px;
    font-size: 0.8rem;
    color: #ecf0f1;
  }

  .research-progress {
    margin-top: 6px;
    height: 6px;
    background: rgba(0, 0, 0, 0.3);
    border-radius: 3px;
    overflow: hidden;
  }

  .rp-fill {
    height: 100%;
    background: #9b59b6;
  }

  .category-tabs {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    margin-bottom: 12px;
  }

  .category-tabs button {
    padding: 4px 8px;
    font-size: 0.7rem;
    background: #1a1a2e;
    color: #bdc3c7;
    border: 1px solid #34495e;
    border-radius: 4px;
    cursor: pointer;
  }

  .category-tabs button.active {
    background: #9b59b6;
    color: white;
    border-color: #9b59b6;
  }

  .tech-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
    max-height: 300px;
    overflow-y: auto;
  }

  .tech-card {
    padding: 10px;
    background: rgba(52, 73, 94, 0.3);
    border-radius: 6px;
    border-left: 3px solid #34495e;
  }

  .tech-card.researched {
    border-left-color: #2ecc71;
    opacity: 0.8;
  }

  .tech-card.researching {
    border-left-color: #9b59b6;
    background: rgba(155, 89, 182, 0.2);
  }

  .tech-card.locked {
    opacity: 0.5;
  }

  .tech-name {
    font-weight: bold;
    color: #ecf0f1;
    font-size: 0.85rem;
    margin-bottom: 4px;
  }

  .tech-desc {
    font-size: 0.75rem;
    color: #bdc3c7;
    margin-bottom: 6px;
  }

  .tech-info {
    display: flex;
    justify-content: space-between;
    font-size: 0.7rem;
  }

  .cost {
    color: #9b59b6;
  }

  .category {
    color: #7f8c8d;
  }

  .prereq {
    margin-top: 4px;
    font-size: 0.65rem;
    color: #e74c3c;
  }
</style>

<script context="module">
  function getTechName(id) {
    const names = {
      basic_mining: '基础采矿',
      advanced_mining: '高级采矿',
      support_beams: '支撑梁技术',
      deep_drilling: '深层钻探',
      basic_smelting: '基础冶炼',
      advanced_smelting: '高级冶炼',
      stonecutting: '石材加工',
      automation: '自动化',
      mushroom_farming: '真菌培育',
      water_recycling: '水循环系统',
      bio_luminescence: '生物发光',
      basic_weapons: '基础武器',
      advanced_weapons: '高级武器',
      armor: '装甲技术',
      tactics: '战术训练',
      siege_craft: '攻城工程',
      magic_crystal_theory: '魔晶理论',
      crystal_resonance: '水晶共振',
      arcane_defense: '奥术防御',
      teleportation: '传送术',
      fossil_extraction: '化石提取',
      steam_power: '蒸汽动力',
      geothermal: '地热能源',
      diplomacy: '外交学',
      logistics: '后勤学',
      urban_planning: '城市规划',
      scouting: '侦察术',
      cartography: '制图学',
      ancient_lore: '远古知识',
      storage_tech: '存储技术',
      trade_routes: '贸易路线',
      banking: '地下银行',
      basic_research: '基础研究',
      advanced_research: '高级研究',
      scientific_method: '科学方法',
      fortification: '要塞化',
      trap_design: '陷阱设计',
      bunker: '地堡系统'
    };
    return names[id] || id;
  }

  function getCategoryName(cat) {
    const names = {
      mining: '采矿',
      industry: '工业',
      biology: '生物',
      military: '军事',
      magic: '魔法',
      energy: '能源',
      civil: '民政',
      exploration: '探索',
      economy: '经济',
      research: '研究',
      defense: '防御'
    };
    return names[cat] || cat;
  }

  function getTechCost(id) {
    return 10;
  }
</script>
