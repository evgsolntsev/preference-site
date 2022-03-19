<template>
<div class="grid-container outer">
  <div class="up grid-container">
    <img v-for="(cardInfo, index) in up.cards" :key="index" class="vertical" :src="getImgUrl(cardInfo, '')" :style="getGridColumnStyle(index)">
  </div>
  <div class="left grid-container">
    <img v-for="(cardInfo, index) in left.cards" :key="index" class="horizontal" :src="getImgUrl(cardInfo, 'C')" :style="getGridRowStyle(index)">
  </div>
  <div class="center grid-container">
    <div v-for="(cardInfo, index) in center" :key="index" class="vertical played" :style="getGridColumnStyle(index)">
      <div>{{ cardInfo.player }}</div>
      <img :src="getImgUrl(cardInfo.card, '')">
    </div>
  </div>
  <div class="right grid-container">
    <img v-for="(cardInfo, index) in right.cards" :key="index" class="horizontal" :src="getImgUrl(cardInfo, 'CC')" :style="getGridRowStyle(index)">
  </div>
  <div class="down grid-container">
    <img v-for="(cardInfo, index) in down.cards" :key="index" class="vertical" :src="getImgUrl(cardInfo, '')" :style="getGridColumnStyle(index)">
  </div>
  <div class="players">
    <div> Left player: {{ playerDescription(left) }}</div>
    <div> Up player: {{ playerDescription(up) }}</div>
    <div> Right player: {{ playerDescription(right) }}</div>
    <div> Down player (you): {{ playerDescription(down) }}</div>
  </div>
  <div class="login">
    <template v-if="!logged">
      <div><input v-model="player" type="text" placeholder="login"></div>
      <div><input v-model="password" type="text" placeholder="password"></div>
      <button @click="login"> Submit </button>
    </template>
    <div v-else>Welcome, {{ player }}!</div>
  </div>
</div>
</template>

<script>
import axios from 'axios';
export default {
  name: 'App',
  methods: {
    getGridColumnStyle(index) {
        return 'grid-row: 1; grid-column: '+(index+1)+' / span 2'
    },
    getGridRowStyle(index) {
        return 'grid-column: 1; grid-row: '+(index+1)+' / span 2'
    },
    getImgUrl(cardInfo, prefix) {
        var images = require.context('./assets/', false, /\.png$/);
        return images('./'+prefix+cardInfo.rank+cardInfo.suit+".png");
    },
    playerDescription(side) {
        return side.name
    },
    fetchData() {
      if (this.logged) {
        this.axios.get("http://0.0.0.0:8090/room").then(response => {
          this.room = response.data;
          let playerIndex = -1;
          for (let i = 0; i < response.data.sides.length; i++) {
            if (response.data.sides[i].name == this.player) {
              playerIndex = i;
            }
          }
          if (playerIndex === -1) {
            console.log("player not found: "+response);
            return
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
    login() {
      console.log("trying to log in with '"+this.player+"' and '"+this.password+"'");
      this.axios.post("http://0.0.0.0:8090/login", {
        "login": this.player,
        "password": this.password,
      }).then(response => {
        this.logged = true;
        console.log(response)
      })
    }
  },
  data() {
    let nullSide = [];
    return {
      ready: false,
      room: null,
      down: nullSide,
      up: nullSide,
      left: nullSide,
      right: nullSide,
      center: nullSide,
      player: "",
      password: "",
      logged: false,
      axios: axios.create({
        withCredentials: true
      })
    }
  },
  mounted() {
    this.fetchData()
  },
  created(){
    this.fetchData()
    this.interval = setInterval(() =>{
      this.fetchData()},1000)
  },
  unmounted(){
    clearInterval(this.interval)
  }
}
</script>

<style>
#app {
    font-family: Avenir, Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
    width: 100%;
    height: 100%;
    min-height: 100%;
    background-color: #15A626;
}
html, body {
    height: 100%;
    width: 100%;
    margin: 0;
}

img {
    max-height: 100%;
    max-width: 100%;
    object-fit: contain;
}

.vertical {
    min-height: 100%;
}

.played {
    max-height: 50%;
    min-height: 50%;
}

.horizontal {
    min-width: 100%;
    min-height: 100%;
}

div.grid-container {
    display: grid;
    gap: 0px;
    min-height: 100%;
    max-height: 100%;
}

div.outer {
    grid-template-columns: 1fr 2fr 1fr;
    grid-template-rows: 10fr 20fr 10fr 10fr;
}

.up {
    grid-column: 2;
    grid-row: 1;
    border: 1px solid;
}
.down {
    grid-column: 2;
    grid-row: 3;
    border: 1px solid;
}
.left {
    grid-column: 1;
    grid-row: 1 / 4;
    border: 1px solid;
}
.right {
    grid-column: 3;
    grid-row: 1 / 4;
    border: 1px solid;
}
.center {
    grid-column: 2;
    grid-row: 2;
    border: 1px solid;
}

.players {
    grid-column: 1;
    grid-row: 4;
    border: 1px solid;
}

.login {
    grid-column: 3;
    grid-row: 4;
    border: 1px solid;
}
</style>
