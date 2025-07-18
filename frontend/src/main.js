import { createApp } from 'vue'
import App from './App.vue'
import { Quasar } from 'quasar'
import quasarUserOptions from './quasar-user-options'
import { createI18n } from 'vue-i18n'
import { createPinia } from 'pinia'

// Import Quasar CSS and icon fonts
import 'quasar/src/css/index.sass'
import '@quasar/extras/material-icons/material-icons.css'

// Localization messages
const messages = {
  en: {
    title: 'Contoso Quasar App',
    welcome: 'Welcome to Contoso!',
    ping: 'Ping Backend',
    clear: 'Clear',
    players: 'Players',
    addPlayer: 'Add Player',
    editPlayer: 'Edit Player',
    name: 'Name',
    surname: 'Surname',
    balance: 'Balance',
    cancel: 'Cancel',
    save: 'Save',
    pingTab: 'Ping',
    id: 'ID',
    actions: 'Actions',
    recordsPerPage: 'Records per page'
  },
  fr: {
    title: 'Application Quasar Contoso',
    welcome: 'Bienvenue chez Contoso !',
    ping: 'Tester le backend',
    clear: 'Effacer',
    players: 'Joueurs',
    addPlayer: 'Ajouter un joueur',
    editPlayer: 'Modifier le joueur',
    name: 'Prénom',
    surname: 'Nom de famille',
    balance: 'Solde',
    cancel: 'Annuler',
    save: 'Enregistrer',
    pingTab: 'Ping',
    id: 'ID',
    actions: 'Actions',
    recordsPerPage: 'Enregistrements par page'
  }
}

const i18n = createI18n({
  legacy: false, // <-- Add this line
  locale: 'en',
  fallbackLocale: 'en',
  messages
})

const pinia = createPinia()

createApp(App)
  .use(Quasar, quasarUserOptions)
  .use(i18n)
  .use(pinia)
  .mount('#app')
