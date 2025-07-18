<template>
  <q-card>
    <q-card-section>
      <div class="text-h6">{{ $t('players') }}</div>
      <q-btn color="primary" icon="add" :label="$t('addPlayer')" @click="showAdd = true" class="q-ml-sm" />
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
        :rows-per-page-label="$t('recordsPerPage')"
        :key="tableKey"
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
          <div class="text-h6">{{ editMode ? $t('editPlayer') : $t('addPlayer') }}</div>
        </q-card-section>
        <q-card-section>
          <q-input v-model="form.name" :label="$t('name')" />
          <q-input v-model="form.surname" :label="$t('surname')" />
          <q-input v-model.number="form.balance" :label="$t('balance')" type="number" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="$t('cancel')" v-close-popup @click="resetForm" />
          <q-btn flat :label="$t('save')" color="primary" @click="savePlayer" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-card>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { usePortalStore } from '../stores/portal'
import { useI18n } from 'vue-i18n'

const portal = usePortalStore()
const { t } = useI18n()

const showAdd = ref(false)
const editMode = ref(false)
const form = ref({ id: null, name: '', surname: '', balance: 0 })
const tableKey = ref(0)

const columns = computed(() => [
  { name: 'id', label: t('id'), field: 'id', align: 'left' },
  { name: 'name', label: t('name'), field: 'name', align: 'left' },
  { name: 'surname', label: t('surname'), field: 'surname', align: 'left' },
  { name: 'balance', label: t('balance'), field: 'balance', align: 'right' },
  { name: 'actions', label: t('actions'), field: 'actions', align: 'center' }
])

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
  let result
  if (editMode.value) {
    result = await portal.updatePlayer(form.value.id, {
      name: form.value.name,
      surname: form.value.surname,
      balance: form.value.balance
    })
  } else {
    result = await portal.createPlayer({
      name: form.value.name,
      surname: form.value.surname,
      balance: form.value.balance
    })
  }
  if (!result?.error) {
    showAdd.value = false
    resetForm()
    await refreshPlayers()
  }
}

async function removePlayer(row) {
  const result = await portal.deletePlayer(row.id)
  // Always refresh after delete, regardless of result (handles 204 No Content)
  await refreshPlayers()
}

async function refreshPlayers() {
  await portal.fetchPlayers()
  tableKey.value++
}

onMounted(() => {
  refreshPlayers()
})
</script>
