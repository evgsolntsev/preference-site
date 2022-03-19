<template>
  <canvas id="canvas"/>
</template>

<script>
import axios from 'axios';
export default {
  name: 'CanvasRoomTable',
  methods: {
    getImgUrl: function (cardInfo) {
        return require('../assets/'+cardInfo.card.rank+cardInfo.card.suit+".svg")
    },
    fetchData: function() {
      axios.get("http://0.0.0.0:8090/room").then(response => {
        this.room = response.data;
        let playerIndex = 0;
        for (let i = 0; i < response.data.sides.length; i++) {
          if (response.data.sides[i].Name == this.player) {
            playerIndex = i;
          }
        }
        this.down = response.data.sides[playerIndex];
        this.left = response.data.sides[(playerIndex+1)%4];
        this.up = response.data.sides[(playerIndex+2)%4];
        this.right = response.data.sides[(playerIndex+3)%4];
        this.center = response.data.center;
        this.ready = true;
      })
    }
  },
  data() {
    let nullSide = {cards: null};
    return {
      ready: false,
      room: null,
      down: nullSide,
      up: nullSide,
      left: nullSide,
      right: nullSide,
      center: nullSide,
      player: "evgsol"
    }
  },
  created() {
    this.fetchData()
  }
}
</script>

<style>
* {
    padding: 0px;
    margin: 0px;
}
html body {
    width: 100%;
    height: 100%;
}
#canvas {
    border: 3px solid black;
    height: 100%;
    width: 100%;
}
.rotateClockwise {
    -webkit-transform: rotate(90deg);
    -moz-transform: rotate(90deg);
    -o-transform: rotate(90deg);
    -ms-transform: rotate(90deg);
    transform: rotate(90deg);
}
.rotateCounterclockwise {
    -webkit-transform: rotate(270deg);
    -moz-transform: rotate(270deg);
    -o-transform: rotate(270deg);
    -ms-transform: rotate(270deg);
    transform: rotate(270deg);
}
</style>
