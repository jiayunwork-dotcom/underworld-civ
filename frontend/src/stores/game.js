import { writable, derived } from 'svelte/store';

export const gameState = writable(null);
export const playerID = writable(localStorage.getItem('player_id') || '');
export const username = writable(localStorage.getItem('username') || '');
export const currentLayer = writable(0);
export const selectedCell = writable(null);
export const selectedUnit = writable(null);
export const buildMode = writable(null);
export const currentView = writable('menu');
export const gameList = writable([]);
export const operationMode = writable(null);

export const currentPlayer = derived([gameState, playerID], ([$gameState, $playerID]) => {
  if (!$gameState || !$playerID) return null;
  return $gameState.players?.[$playerID] || null;
});

export const prevResources = writable(null);
export const resourceDeltas = writable({});

export const setPlayerID = (id) => {
  playerID.set(id);
  localStorage.setItem('player_id', id);
};

export const setUsername = (name) => {
  username.set(name);
  localStorage.setItem('username', name);
};

export const setOperationMode = (mode) => {
  operationMode.set(mode);
  if (mode !== 'build') {
    buildMode.set(null);
  }
};

export const trackResourceDelta = (current, previous) => {
  if (!previous || !current) return;
  const deltas = {};
  const prevRes = previous.resources || previous;
  const currRes = current.resources || current;
  
  for (const key of ['stone', 'metal', 'glow_mushroom', 'water', 'magic_crystal', 'fossil_fuel']) {
    const prev = prevRes[key] || 0;
    const curr = currRes[key] || 0;
    deltas[key] = curr - prev;
  }
  
  resourceDeltas.set(deltas);
  prevResources.set(current);
};
