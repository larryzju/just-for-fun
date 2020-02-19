import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex);

class Block {
  constructor(value, status) {
    this.value = value
    this.status = status || 0
  }
}

class MineSweeper {
  constructor(width, height, mines) {
    this.width  = width
    this.height = height
    this.mines  = mines
    this.timer  = null,
    this.over   = true,
    this.blocks = []
    this.init()
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

  init() {
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

  open(index) {
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
            g.open(i)
          }
        }
      })
    }
  }
}


export default new Vuex.Store({
  state: new MineSweeper(19, 19, 99),
  mutations: {
    RESET(state) {
      state.init()
    },
    OPEN(state, index) {
      state.open(index)
    }
  },
  actions: {
    init(context, {width, height, mines}) {
      context.commit('INIT', {width, height, mines})
    },
    open(context, index) {
      context.commit('OPEN', index)
    },
    reset(context) {
      context.commit('RESET')
    },
    mark(context, x, y) {
      alert('mark ' + x + ',' + y);
    },
  },
  getters: {
    over:    s => s.over,
    seconds: s => s.seconds,
    mines:   s => s.mines,
    height:  s => s.height,
    width:   s => s.width,
    remains: s => s.mines - s.blocks.filter(b => b.status == 2),
    blocks:  s => s.blocks
  }
})
