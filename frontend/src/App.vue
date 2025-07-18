<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-toolbar-title>
          <q-icon name="apps" class="q-mr-sm" />
          {{ $t('title') }}
        </q-toolbar-title>
        <q-btn flat dense round icon="language">
          <q-menu>
            <q-list>
              <q-item clickable v-close-popup @click="setLang('en')">
                <q-item-section>English</q-item-section>
              </q-item>
              <q-item clickable v-close-popup @click="setLang('fr')">
                <q-item-section>Français</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </q-toolbar>
    </q-header>

    <q-page-container>
      <q-page class="q-pa-md flex flex-center bg-grey-2">
        <q-card class="my-card">
          <q-card-section class="bg-primary text-white">
            <div class="text-h5 text-center">{{ $t('welcome') }}</div>
          </q-card-section>
          <q-card-section>
            <div class="q-mb-md text-center">
              <q-btn color="primary" @click="portal.pingBackend" :label="$t('ping')" unelevated rounded />
              <q-btn color="secondary" @click="portal.clearResult" :label="$t('clear')" unelevated rounded class="q-ml-sm" />
            </div>
            <div class="q-mt-md">
              <q-banner v-if="portal.pingResult" dense class="bg-green-2 text-primary text-center">
                {{ portal.pingResult }}
              </q-banner>
            </div>
          </q-card-section>
        </q-card>
      </q-page>
      <q-page class="q-pa-md">
        <Players />
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { usePortalStore } from './stores/portal'
import Players from './components/Players.vue'

const portal = usePortalStore()
const { locale } = useI18n()

function setLang(lang) {
  locale.value = lang
}
</script>

<style>
.my-card {
  min-width: 350px;
  max-width: 400px;
  margin: 40px auto;
  box-shadow: 0 4px 24px rgba(0,0,0,0.08);
  border-radius: 18px;
}
body {
  background: #f5f7fa;
}
</style>
