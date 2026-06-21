const API_BASE = import.meta.env.VITE_API_URL || '/api';
const WS_BASE = import.meta.env.VITE_WS_URL || (location.protocol === 'https:' ? 'wss:' : 'ws:') + '//' + location.host + '/ws';

async function request(url, options = {}) {
  const res = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers
    }
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: `HTTP ${res.status}` }));
    throw new Error(err.error || `请求失败: ${res.status}`);
  }
  return res.json();
}

export const api = {
  async getGames() {
    return request(`${API_BASE}/games`);
  },

  async createGame(data, playerId, userName) {
    return request(`${API_BASE}/games`, {
      method: 'POST',
      body: JSON.stringify({
        ...data,
        player_id: playerId,
        username: userName
      })
    });
  },

  async getGame(gameId, playerId) {
    const url = playerId
      ? `${API_BASE}/games/${gameId}?player_id=${encodeURIComponent(playerId)}`
      : `${API_BASE}/games/${gameId}`;
    return request(url);
  },

  async joinGame(gameId, data, playerId, userName) {
    return request(`${API_BASE}/games/${gameId}/join`, {
      method: 'POST',
      body: JSON.stringify({
        ...data,
        player_id: playerId,
        username: userName
      })
    });
  },

  async startGame(gameId, playerId) {
    return request(`${API_BASE}/games/${gameId}/start`, {
      method: 'POST'
    });
  },

  async submitAction(gameId, playerId, action, data) {
    return request(`${API_BASE}/games/${gameId}/actions`, {
      method: 'POST',
      body: JSON.stringify({ action, data, player_id: playerId })
    });
  },

  async getRaces() {
    return request(`${API_BASE}/races`);
  },

  async getBuildings() {
    return request(`${API_BASE}/buildings`);
  },

  async getUnits() {
    return request(`${API_BASE}/units`);
  },

  async getTechs() {
    return request(`${API_BASE}/techs`);
  },

  async getTechTree(gameId, playerId) {
    const url = playerId
      ? `${API_BASE}/games/${gameId}/tech-tree?player_id=${encodeURIComponent(playerId)}`
      : `${API_BASE}/games/${gameId}/tech-tree`;
    return request(url);
  },

  async setResearch(gameId, playerId, techId) {
    return request(`${API_BASE}/games/${gameId}/tech-tree`, {
      method: 'POST',
      body: JSON.stringify({
        player_id: playerId,
        tech_id: techId
      })
    });
  },

  async blockadeTech(gameId, playerId, targetId, category) {
    return request(`${API_BASE}/games/${gameId}/tech-blockade`, {
      method: 'POST',
      body: JSON.stringify({
        player_id: playerId,
        target_id: targetId,
        category: category
      })
    });
  },

  async dispatchSpy(gameId, playerId, unitId, targetId) {
    return request(`${API_BASE}/games/${gameId}/actions`, {
      method: 'POST',
      body: JSON.stringify({
        action: 'dispatch_spy',
        data: {
          unit_id: unitId,
          target_id: targetId
        },
        player_id: playerId
      })
    });
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
