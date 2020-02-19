<template>
    <div class="status-bar">
        <div class='digital' id="remain">{{remains}}</div>
        <div @click="reset" id="emoji" class='em' :class='over? "em-sob" : "em-smile"' />
        <div class='digital' id="timer" >{{seconds}}</div>
    </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'

export default {
    name: 'StatusBar',
    methods: {
      ...mapActions(['reset'])
    },
    computed: {
      ...mapGetters(['remains', 'seconds', 'over'])
    }
}
</script>

<style scoped>
@font-face {
    font-family: 'digital-clock-font';
    src: url('/fonts/digital-7.ttf');
}

div.status-bar {
    display: inline-flex;
    justify-content: space-between;
    width: 100%;
}

.digital {
    font-family: 'digital-clock-font';
    font-weight: 900;
    font-size: 18pt;
    color: green;
}

#timer {
    width: 3em;
    text-align: right;
}

#remain {
    width: 3em;
    text-align: left;
}
</style>
