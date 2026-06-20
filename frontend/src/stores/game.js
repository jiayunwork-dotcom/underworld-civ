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

export const currentPlayer = derived([gameState, playerID], ([$gameState, $playerID]) => {
  if (!$gameState || !$playerID) return null;
  return $gameState.players?.[$playerID] || null;
});

export const setPlayerID = (id) => {
  playerID.set(id);
  localStorage.setItem('player_id', id);
};

export const setUsername = (name) => {
  username.set(name);
  localStorage.setItem('username', name);
};
