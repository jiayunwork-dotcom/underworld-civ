const API_BASE = import.meta.env.VITE_API_URL || '/api';
const WS_BASE = import.meta.env.VITE_WS_URL || (location.protocol === 'https:' ? 'wss:' : 'ws:') + '//' + location.host + '/ws';

export const api = {
  async getGames() {
    const res = await fetch(`${API_BASE}/games`);
    return res.json();
  },

  async createGame(data, playerId, userName) {
    const res = await fetch(`${API_BASE}/games`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Player-ID': playerId,
        'X-Username': userName
      },
      body: JSON.stringify(data)
    });
    return res.json();
  },

  async getGame(gameId, playerId) {
    const res = await fetch(`${API_BASE}/games/${gameId}`, {
      headers: { 'X-Player-ID': playerId }
    });
    return res.json();
  },

  async joinGame(gameId, data, playerId, userName) {
    const res = await fetch(`${API_BASE}/games/${gameId}/join`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Player-ID': playerId,
        'X-Username': userName
      },
      body: JSON.stringify(data)
    });
    return res.json();
  },

  async startGame(gameId, playerId) {
    const res = await fetch(`${API_BASE}/games/${gameId}/start`, {
      method: 'POST',
      headers: { 'X-Player-ID': playerId }
    });
    return res.json();
  },

  async submitAction(gameId, playerId, action, data) {
    const res = await fetch(`${API_BASE}/games/${gameId}/actions`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Player-ID': playerId
      },
      body: JSON.stringify({ action, data })
    });
    return res.json();
  },

  async getRaces() {
    const res = await fetch(`${API_BASE}/races`);
    return res.json();
  },

  async getBuildings() {
    const res = await fetch(`${API_BASE}/buildings`);
    return res.json();
  },

  async getUnits() {
    const res = await fetch(`${API_BASE}/units`);
    return res.json();
  },

  async getTechs() {
    const res = await fetch(`${API_BASE}/techs`);
    return res.json();
  }
};

export function connectWS(gameId, playerId, onMessage) {
  const wsUrl = `${WS_BASE}/games/${gameId}?player_id=${encodeURIComponent(playerId)}`;
  const ws = new WebSocket(wsUrl);

  ws.onopen = () => {
    console.log('WebSocket connected');
  };

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      onMessage(data);
    } catch (e) {
      console.error('WS parse error:', e);
    }
  };

  ws.onclose = () => {
    console.log('WebSocket disconnected');
  };

  ws.onerror = (error) => {
    console.error('WebSocket error:', error);
  };

  return ws;
}
