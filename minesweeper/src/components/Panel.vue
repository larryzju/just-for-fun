<template>
    <div class='panel'>
        <table>
            <tr v-for='row in height' :key='row'>
                <td v-for='col in width' :key='col'>
                    <Block :block='block(row, col)' 
                        @click.native.left='open(row, col)'
                        @click.native.right.prevent ='mark(row, col)' />
                </td>
            </tr>
        </table>
    </div>
</template>

<script>
import { mapGetters } from 'vuex'
import Block from "./Block.vue"

export default {
    name: 'Panel',
    components: {
        Block,
    },
    methods: {
      index(row, col) {
        return this.width * (row-1) + (col-1)
      },
      open(row, col) {
        this.$store.dispatch('open', this.index(row, col))
      },
      mark(row, col) {
        alert("mark "+this.index(row,col))
      },
      block(row, col) {
        return this.blocks[this.index(row,col)];
      }
    },
    computed: {
      ...mapGetters(['height', 'width', 'blocks'])
    }
}
</script>

<style scoped>
table {
    border: 1px solid black;
    border-collapse: collapse
}

#panel {
    width: fit-content;
}

td {
    padding: 0 1px 1px 0;
}
</style>
