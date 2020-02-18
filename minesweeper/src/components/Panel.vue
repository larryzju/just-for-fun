<template>
    <div>
        <table>
            <tr v-for='row in ms.height' :key='row'>
                <td v-for='col in ms.width' :key='col'>
                    <Block
                        :block='getBlock(row, col)'
                        @click.native.left='onLeftClick(row, col)'
                        @click.native.right.prevent ='onRightClick(row, col)' />
                </td>
            </tr>
        </table>
    </div>
</template>

<script>
import Block from "./Block.vue"

export default {
    name: 'Panel',
    props: ['ms'],
    components: {
        Block,
    },
    methods: {
        blockIndex: function(r, c) {
            return (r-1)*this.ms.width + (c-1);
        },
        getBlock: function(r, c) {
            return this.ms.blocks[this.blockIndex(r, c)]
        },
        onLeftClick: function(r, c) {
            let index = this.blockIndex(r,c);
            return this.ms.openBlock(index);
        },
        onRightClick: function(r, c) {
            let index = this.blockIndex(r,c);
            alert("to be implement, index=" + index);
        }
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
