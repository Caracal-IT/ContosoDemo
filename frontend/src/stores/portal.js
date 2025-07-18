import { defineStore } from 'pinia'
import { ref } from 'vue'

export const usePortalStore = defineStore('portal', () => {
  const pingResult = ref('')
  const players = ref([])
  const player = ref(null)

  async function pingBackend() {
    try {
      const res = await fetch('/api/ping')
      const data = await res.json()
      pingResult.value = data.message || JSON.stringify(data)
    } catch (e) {
      pingResult.value = 'Ping failed'
    }
  }

  function clearResult() {
    pingResult.value = ''
  }

  // Player CRUD
  async function fetchPlayers() {
    const res = await fetch('/api/players')
    players.value = await res.json()
  }

  async function fetchPlayer(id) {
    const res = await fetch(`/api/players/${id}`)
    player.value = await res.json()
  }

  async function createPlayer(data) {
    const res = await fetch('/api/players', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    })
    return await res.json()
  }

  async function updatePlayer(id, data) {
    const res = await fetch(`/api/players/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    })
    return await res.json()
  }

  async function deletePlayer(id) {
    const res = await fetch(`/api/players/${id}`, { method: 'DELETE' })
    return await res.json()
  }

  return {
    pingResult,
    pingBackend,
    clearResult,
    players,
    player,
    fetchPlayers,
    fetchPlayer,
    createPlayer,
    updatePlayer,
    deletePlayer
  }
})
