<script>
  import { onMount, onDestroy, afterUpdate } from 'svelte';
  import { hexToPixel, pixelToHex, getHexCorners, hexDistance, getHexNeighbors, hexKey } from '../utils/hex.js';

  export let game;
  export let layer;
  export let playerId;
  export let selectedCell;
  export let buildMode;
  export let onCellClick;
  export let onMine;
  export let onMove;
  export let onAttack;

  let canvas;
  let ctx;
  let hexSize = 32;
  let offsetX = 100;
  let offsetY = 100;
  let isDragging = false;
  let dragStart = { x: 0, y: 0 };
  let hoveredCell = null;

  const layerColors = [
    { wall: '#3d5a3d', floor: '#2a3f2a', accent: '#7cb87c' },
    { wall: '#5d4e37', floor: '#3d3225', accent: '#b8a07c' },
    { wall: '#4a4a5a', floor: '#2d2d3a', accent: '#8c8ca8' },
    { wall: '#5a2d2d', floor: '#3a1a1a', accent: '#c44e4e' },
    { wall: '#3d3d5a', floor: '#25253a', accent: '#7c7cc4' }
  ];

  const mineralColors = {
    none: 'transparent',
    stone: '#95a5a6',
    iron: '#7f8c8d',
    copper: '#e67e22',
    gold: '#f1c40f',
    glow_mushroom: '#2ecc71',
    magic_crystal: '#9b59b6',
    fossil_fuel: '#1a1a1a'
  };

  $: {
    if (game && layer !== undefined) {
      draw();
    }
  }

  onMount(() => {
    ctx = canvas.getContext('2d');
    draw();

    window.addEventListener('resize', draw);
  });

  onDestroy(() => {
    window.removeEventListener('resize', draw);
  });

  function draw() {
    if (!ctx || !game) return;

    const rect = canvas.getBoundingClientRect();
    canvas.width = rect.width * window.devicePixelRatio;
    canvas.height = rect.height * window.devicePixelRatio;
    ctx.scale(window.devicePixelRatio, window.devicePixelRatio);

    const w = rect.width;
    const h = rect.height;

    ctx.fillStyle = '#0d1421';
    ctx.fillRect(0, 0, w, h);

    const gameLayer = game.map?.layers?.[layer];
    if (!gameLayer) return;

    const cells = Object.values(gameLayer.cells || {});

    for (const cell of cells) {
      drawCell(cell, w, h);
    }

    if (selectedCell && selectedCell.coord) {
      const pos = hexToPixel(selectedCell.coord.q, selectedCell.coord.r, hexSize);
      const x = pos.x + offsetX;
      const y = pos.y + offsetY;

      ctx.strokeStyle = '#f1c40f';
      ctx.lineWidth = 3;
      drawHexagon(x, y, hexSize - 2);
      ctx.stroke();
    }

    if (hoveredCell) {
      const pos = hexToPixel(hoveredCell.q, hoveredCell.r, hexSize);
      const x = pos.x + offsetX;
      const y = pos.y + offsetY;

      ctx.strokeStyle = 'rgba(255, 255, 255, 0.5)';
      ctx.lineWidth = 2;
      drawHexagon(x, y, hexSize - 2);
      ctx.stroke();
    }

    if (buildMode && hoveredCell) {
      const pos = hexToPixel(hoveredCell.q, hoveredCell.r, hexSize);
      const x = pos.x + offsetX;
      const y = pos.y + offsetY;

      ctx.fillStyle = 'rgba(52, 152, 219, 0.4)';
      drawHexagon(x, y, hexSize - 4);
      ctx.fill();
    }
  }

  function drawCell(cell, w, h) {
    if (!cell.coord) return;

    const pos = hexToPixel(cell.coord.q, cell.coord.r, hexSize);
    const x = pos.x + offsetX;
    const y = pos.y + offsetY;

    if (x < -hexSize * 2 || x > w + hexSize * 2 || y < -hexSize * 2 || y > h + hexSize * 2) {
      return;
    }

    const colors = layerColors[layer] || layerColors[0];

    if (cell.is_wall) {
      ctx.fillStyle = colors.wall;
      drawHexagon(x, y, hexSize - 1);
      ctx.fill();

      ctx.fillStyle = 'rgba(0, 0, 0, 0.3)';
      drawHexagon(x, y, hexSize - 4);
      ctx.fill();
    } else {
      ctx.fillStyle = colors.floor;
      drawHexagon(x, y, hexSize - 1);
      ctx.fill();

      if (cell.owner) {
        const player = game.players?.[cell.owner];
        if (player) {
          ctx.fillStyle = player.color + '40';
          drawHexagon(x, y, hexSize - 4);
          ctx.fill();
        }
      }

      if (cell.mineral_type && cell.mineral_type !== 'none') {
        const mColor = mineralColors[cell.mineral_type] || '#fff';
        ctx.fillStyle = mColor;
        ctx.beginPath();
        ctx.arc(x, y, 5, 0, Math.PI * 2);
        ctx.fill();
      }

      if (cell.is_shaft) {
        ctx.fillStyle = '#f39c12';
        ctx.beginPath();
        ctx.arc(x, y, 8, 0, Math.PI * 2);
        ctx.fill();
        ctx.fillStyle = '#fff';
        ctx.font = 'bold 10px sans-serif';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText('↓', x, y);
      }

      if (cell.building) {
        drawBuilding(x, y, cell.building);
      }

      if (cell.units && cell.units.length > 0) {
        drawUnits(x, y, cell.units);
      }

      if (cell.flooded) {
        ctx.fillStyle = 'rgba(52, 152, 219, 0.5)';
        drawHexagon(x, y, hexSize - 4);
        ctx.fill();
      }
    }

    ctx.strokeStyle = 'rgba(0, 0, 0, 0.3)';
    ctx.lineWidth = 1;
    drawHexagon(x, y, hexSize - 1);
    ctx.stroke();
  }

  function drawBuilding(x, y, building) {
    const buildingColors = {
      main_base: '#e94560',
      living_quarters: '#3498db',
      workshop: '#95a5a6',
      fungus_farm: '#2ecc71',
      smelter: '#e67e22',
      academy: '#9b59b6',
      altar: '#f1c40f',
      warehouse: '#7f8c8d',
      watchtower: '#1abc9c',
      wall: '#95a5a6',
      elevator: '#e74c3c',
      forge_shrine: '#e74c3c',
      spore_network: '#2ecc71',
      crystal_tower: '#9b59b6',
      petrified_wall: '#7f8c8d',
      hive_rift: '#e67e22'
    };

    const color = buildingColors[building.type] || '#ecf0f1';

    ctx.fillStyle = color;
    ctx.fillRect(x - 10, y - 10, 20, 20);

    if (!building.completed) {
      ctx.fillStyle = 'rgba(0, 0, 0, 0.5)';
      const progress = building.progress / building.build_time;
      ctx.fillRect(x - 10, y - 10, 20 * (1 - progress), 20);
    }

    ctx.strokeStyle = '#fff';
    ctx.lineWidth = 2;
    ctx.strokeRect(x - 10, y - 10, 20, 20);

    const hpPercent = building.hp / building.max_hp;
    ctx.fillStyle = hpPercent > 0.5 ? '#2ecc71' : hpPercent > 0.25 ? '#f39c12' : '#e74c3c';
    ctx.fillRect(x - 10, y + 12, 20 * hpPercent, 3);
  }

  function drawUnits(x, y, units) {
    const count = units.length;
    const owner = units[0]?.owner;
    const player = game.players?.[owner];
    const color = player?.color || '#fff';

    ctx.fillStyle = color;
    ctx.beginPath();
    ctx.arc(x, y, 8, 0, Math.PI * 2);
    ctx.fill();

    ctx.fillStyle = '#fff';
    ctx.font = 'bold 9px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText(count > 9 ? '9+' : count, x, y);

    ctx.strokeStyle = '#fff';
    ctx.lineWidth = 1.5;
    ctx.beginPath();
    ctx.arc(x, y, 8, 0, Math.PI * 2);
    ctx.stroke();
  }

  function drawHexagon(cx, cy, size) {
    ctx.beginPath();
    for (let i = 0; i < 6; i++) {
      const angle = (60 * i - 30) * Math.PI / 180;
      const px = cx + size * Math.cos(angle);
      const py = cy + size * Math.sin(angle);
      if (i === 0) {
        ctx.moveTo(px, py);
      } else {
        ctx.lineTo(px, py);
      }
    }
    ctx.closePath();
  }

  function handleMouseDown(e) {
    if (e.button === 2 || e.button === 1) {
      isDragging = true;
      dragStart = { x: e.clientX - offsetX, y: e.clientY - offsetY };
    }
  }

  function handleMouseMove(e) {
    const rect = canvas.getBoundingClientRect();
    const mx = e.clientX - rect.left - offsetX;
    const my = e.clientY - rect.top - offsetY;

    const hex = pixelToHex(mx, my, hexSize);

    const gameLayer = game.map?.layers?.[layer];
    if (gameLayer) {
      const key = hexKey(hex.q, hex.r);
      if (gameLayer.cells?.[key]) {
        hoveredCell = hex;
      } else {
        hoveredCell = null;
      }
    }

    if (isDragging) {
      offsetX = e.clientX - dragStart.x;
      offsetY = e.clientY - dragStart.y;
      draw();
    }
  }

  function handleMouseUp(e) {
    isDragging = false;
  }

  function handleClick(e) {
    if (isDragging) return;

    const rect = canvas.getBoundingClientRect();
    const mx = e.clientX - rect.left - offsetX;
    const my = e.clientY - rect.top - offsetY;

    const hex = pixelToHex(mx, my, hexSize);
    const key = hexKey(hex.q, hex.r);

    const gameLayer = game.map?.layers?.[layer];
    if (gameLayer && gameLayer.cells?.[key]) {
      onCellClick(gameLayer.cells[key]);
    }
  }

  function handleWheel(e) {
    e.preventDefault();
    const delta = e.deltaY > 0 ? 0.9 : 1.1;
    hexSize = Math.max(16, Math.min(64, hexSize * delta));
    draw();
  }

  function handleContextMenu(e) {
    e.preventDefault();
  }

  function zoomIn() {
    hexSize = Math.min(64, hexSize * 1.2);
    draw();
  }

  function zoomOut() {
    hexSize = Math.max(16, hexSize / 1.2);
    draw();
  }
</script>

<div class="map-container">
  <canvas
    bind:this={canvas}
    on:mousedown={handleMouseDown}
    on:mousemove={handleMouseMove}
    on:mouseup={handleMouseUp}
    on:mouseleave={handleMouseUp}
    on:click={handleClick}
    on:wheel={handleWheel}
    on:contextmenu={handleContextMenu}
  />

  <div class="zoom-controls">
    <button on:click={zoomIn}>+</button>
    <button on:click={zoomOut}>−</button>
  </div>
</div>

<style>
  .map-container {
    width: 100%;
    height: 100%;
    position: relative;
  }

  canvas {
    width: 100%;
    height: 100%;
    display: block;
    cursor: grab;
  }

  canvas:active {
    cursor: grabbing;
  }

  .zoom-controls {
    position: absolute;
    bottom: 20px;
    right: 20px;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .zoom-controls button {
    width: 36px;
    height: 36px;
    border-radius: 6px;
    border: 1px solid #34495e;
    background: rgba(22, 33, 62, 0.9);
    color: #ecf0f1;
    font-size: 1.2rem;
    cursor: pointer;
  }

  .zoom-controls button:hover {
    background: rgba(52, 73, 94, 0.9);
  }
</style>
