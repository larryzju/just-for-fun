<template>
    <div id="app">
        <status-bar v-on:reset="reset" :ms='game'/>
        <panel v-on:game-over="gameOver" :ms='game'/>
    </div>
</template>

<script>
import StatusBar from "./components/StatusBar.vue"
import Panel from "./components/Panel.vue"


class Block {
  constructor(value, status) {
    this.value = value
    this.status = status || 0
  }
}

class MineSweeper {
  constructor(height, width, mines) {
      this.height = height;
      this.width = width;
      this.mines = mines;
      this.blocks = [];
      this.initBlocks();
  }

  initBlocks() {
    let size = this.width * this.height;
    let m = new Array(size)
    for ( let i = 0; i < size; i++ ) {
      m[i] = i < this.mines ? NaN: 0;
    }

    m.sort(() => Math.random() - 0.5)

    for (let i = 0; i < size; i++ ) {
      if (!isNaN(m[i])) {
        m[i] = this.getAdjoinIndex(i).filter(i => isNaN(m[i])).length
      }
    }

    this.blocks = m.map(v => new Block(v, 0))
    this.over = false;
    this.seconds = 0;
    if (this.timer) {
      clearInterval(this.timer);
      this.timer = null;
    }
  }

  getAdjoinIndex(index) {
    let adjoins = [];
    let c = index % this.width, r = Math.floor(index / this.width);
    for (let x = -1; x <= 1; x++ ) {
      for (let y = -1; y <= 1; y++ ) {
        if (x == 0 && y == 0) {
          continue;
        }
        let nc = c + x;
        let nr = r + y;
        if (nc >= 0 && nc < this.width && nr >= 0 && nr < this.height) {
          adjoins.push(nc + nr*this.width)
        }
      }
    }
    return adjoins
  }

  getAdjoinBlocks(index) {
    return this.getAdjoinIndex(index).map(i => this.blocks[i])
  }

  openBlock(index) {
    let g = this;
    let b = g.blocks[index];

    // if it is the first time to click, start the timer
    if (this.timer == null) {
      this.timer = setInterval(function(){
        g.seconds++;
      }, 1000)
    }

    // TODO open the adjoin blocks if it is opened
    if (b.status > 0) {
      return
    }

    b.status = 1

    // bomb!
    if (isNaN(b.value)) {
      this.over = true;
      this.blocks.forEach(b => b.status=1)
      clearInterval(this.timer)
      this.timer = null;
      return
    }

    // empty block
    if (b.value === 0) {
      let adjoins = this.getAdjoinIndex(index);
      window.console.log(["adjoin of", index, adjoins])
      adjoins.forEach(function(i) {
        let a = g.blocks[i]
        if (a.status == 0) {
          if (a.value > 0) {
            a.status = 1
          }
          if (a.value === 0) {
            g.openBlock(i)
          }
        }
      })
    }
  }
}

export default {
    name: 'app',
    data() {
      return {
        game: new MineSweeper(19, 19, 81)
      }
    },
    components: {
        'status-bar': StatusBar,
        'panel': Panel,
    },
    methods: {
      gameOver() {
        this.game.over = true;
        this.game.blocks.forEach(b => b.status=1)
      },
      reset() {
        this.game.initBlocks();
      }
    }
}
</script>

<style>
@import "https://emoji-css.afeld.me/emoji.css";

#app {
  display: inline-block;
}
</style>
