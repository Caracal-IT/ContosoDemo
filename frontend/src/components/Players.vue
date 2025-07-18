<template>
  <q-card>
    <q-card-section>
      <div class="text-h6">Players</div>
      <q-btn color="primary" icon="add" label="Add Player" @click="showAdd = true" class="q-ml-sm" />
    </q-card-section>
    <q-separator />
    <q-card-section>
      <q-table
        :rows="portal.players"
        :columns="columns"
        row-key="id"
        flat
        dense
        :pagination="{ rowsPerPage: 5 }"
      >
        <template #body-cell-actions="props">
          <q-td>
            <q-btn size="sm" color="primary" icon="edit" flat @click="editPlayer(props.row)" />
            <q-btn size="sm" color="negative" icon="delete" flat @click="removePlayer(props.row)" />
          </q-td>
        </template>
      </q-table>
    </q-card-section>
    <!-- Add/Edit Player Dialog -->
    <q-dialog v-model="showAdd">
      <q-card style="min-width:350px">
        <q-card-section>
          <div class="text-h6">{{ editMode ? 'Edit Player' : 'Add Player' }}</div>
        </q-card-section>
        <q-card-section>
          <q-input v-model="form.name" label="Name" />
          <q-input v-model="form.surname" label="Surname" />
          <q-input v-model.number="form.balance" label="Balance" type="number" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Cancel" v-close-popup @click="resetForm" />
          <q-btn flat label="Save" color="primary" @click="savePlayer" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { usePortalStore } from '../stores/portal'

const portal = usePortalStore()

const showAdd = ref(false)
const editMode = ref(false)
const form = ref({ id: null, name: '', surname: '', balance: 0 })

const columns = [
  { name: 'id', label: 'ID', field: 'id', align: 'left' },
  { name: 'name', label: 'Name', field: 'name', align: 'left' },
  { name: 'surname', label: 'Surname', field: 'surname', align: 'left' },
  { name: 'balance', label: 'Balance', field: 'balance', align: 'right' },
  { name: 'actions', label: 'Actions', field: 'actions', align: 'center' }
]

function resetForm() {
  form.value = { id: null, name: '', surname: '', balance: 0 }
  editMode.value = false
}

function editPlayer(row) {
  form.value = { ...row }
  editMode.value = true
  showAdd.value = true
}

async function savePlayer() {
  if (editMode.value) {
    await portal.updatePlayer(form.value.id, {
      name: form.value.name,
      surname: form.value.surname,
      balance: form.value.balance
    })
  } else {
    await portal.createPlayer({
      name: form.value.name,
      surname: form.value.surname,
      balance: form.value.balance
    })
  }
  showAdd.value = false
  resetForm()
  await portal.fetchPlayers()
}

async function removePlayer(row) {
  await portal.deletePlayer(row.id)
  await portal.fetchPlayers()
}

onMounted(() => {
  portal.fetchPlayers()
})
</script>

