<template>
  <div id="app">
    <status-bar />
    <panel />
  </div>
</template>

<script>
import { mapActions } from 'vuex'

import store from './store'
import StatusBar from "./components/StatusBar.vue"
import Panel from "./components/Panel.vue"

export default {
  name: 'app',
  store,
  components: {
    'status-bar': StatusBar,
    'panel': Panel,
  },
  methods: {
    ...mapActions(['reset'])
  },
}
</script>

<style>
@import "https://emoji-css.afeld.me/emoji.css";

#app {
  display: inline-block;
}
</style>
