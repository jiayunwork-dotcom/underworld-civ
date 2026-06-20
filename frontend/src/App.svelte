<script>
  import { onMount } from 'svelte';
  import { gameState, currentView, gameList, playerID, username, setPlayerID, setUsername } from './stores/game.js';
  import { api } from './utils/api.js';
  import GameMenu from './components/GameMenu.svelte';
  import GameLobby from './components/GameLobby.svelte';
  import GameBoard from './components/GameBoard.svelte';

  onMount(() => {
    if (!$playerID) {
      const id = 'player_' + Math.random().toString(36).substr(2, 9);
      setPlayerID(id);
    }
    if (!$username) {
      const name = '矿工商_' + Math.floor(Math.random() * 10000);
      setUsername(name);
    }
    loadGames();
  });

  async function loadGames() {
    try {
      const games = await api.getGames();
      gameList.set(games);
    } catch (e) {
      console.error('Failed to load games:', e);
    }
  }

  function onGameCreated(data) {
    gameState.set(data.game);
    setPlayerID(data.player_id);
    setUsername(data.username);
    currentView.set('lobby');
  }

  function onGameJoined(data) {
    gameState.set(data.game);
    setPlayerID(data.player_id);
    setUsername(data.username);
    currentView.set('lobby');
  }

  function onGameStarted(state) {
    gameState.set(state);
    currentView.set('game');
  }
</script>

<div class="app">
  {#if $currentView === 'menu'}
    <GameMenu onGameCreated={onGameCreated} onJoinGame={onGameJoined} games={$gameList} refreshGames={loadGames} />
  {:else if $currentView === 'lobby'}
    <GameLobby game={$gameState} playerId={$playerID} onGameStarted={onGameStarted} onBack={() => currentView.set('menu')} />
  {:else if $currentView === 'game'}
    <GameBoard />
  {/if}
</div>

<style>
  .app {
    width: 100vw;
    height: 100vh;
    overflow: hidden;
  }
</style>
