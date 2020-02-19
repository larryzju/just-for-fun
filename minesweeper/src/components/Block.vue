<template>
  <div>
    <div class="block" :class='block.status === 0 ? "unknown": "open"'>
        <div class='raw' v-if='block.status === 0' />
        <img v-else-if='isNaN(block.value)' src="explosion.gif" />
        <div v-else :class='style()'> {{block.value || ''}} </div>
    </div>
    </div>
</template>

<script>
export default {
  name: "Block",
  props: ['block'],
  methods: {
    style() {
      return "n"+this.block.value
    },
  }
}
</script>

<style scoped>
.n1 {
    color: blue
}

.n2 {
    color: green
}

.n3 {
    color: red
}

.n4 {
    color: darkblue
}

.n5 {
    color: darkgreen
}

.n6, .n7, .n8 {
    color: darkred
}

img, .block {
    width: 1.5em;
    height: 1.5em;
    font-weight: 900;
    outline: 1px solid black;
    box-shadow: 2px 2px black;
    text-align: center;
    display: inline-block;
    text-shadow: .1em 0 30px;
}

.unknown {
    background-color: #EFCDBA;
}

.open {
    background-color: #FFFFFF;
}

</style>
