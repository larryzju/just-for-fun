<template>
    <div class="status-bar">
        <div class="digital" id="remain">{{ms.mines - flags}}</div>
        <div @click="$emit('reset')" id="emoji" class='em' :class='ms.over? "em-sob" : "em-smile"' />
        <div class="digital" id="timer" >{{ms.seconds}}</div>
    </div>
</template>

<script>
export default {
    name: 'StatusBar',
    props: ['ms'],
    data() {
        return {
            flags: 0
        }
    },
    method: {
        reset() {
            alert("reset")
        }
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
