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
    <img v-for="(cardInfo, index) in down.cards" :key="index" :class="'vertical '+isSelected(index)+' '+isHover(index)" :src="getImgUrl(cardInfo, '')" :style="getGridColumnStyle(index)" @mouseover="hovered[index]=true" @mouseleave="hovered[index]=false" @click="selected[index]=!selected[index]">
  </div>
  <div class="players">
    <template v-if="isLogged()">
      <div> Left player: {{ playerDescription(left) }}</div>
      <div> Up player: {{ playerDescription(up) }}</div>
      <div> Right player: {{ playerDescription(right) }}</div>
      <div> Down player (you): {{ playerDescription(down) }}</div>
      <button @click="logout">Logout?</button>
      <button @click="leaveRoom">Leave room?</button>
    </template>
  </div>
  <div class="buttons btn-group">
    <template v-if="isLogged()">
      <button v-for="(buttonInfo, index) in buttons()" :key="index" @click="buttonInfo.Click" :disabled="buttonInfo.IsDisabled()" :style="buttonsStyle()">{{ buttonInfo.Text }}</button>
    </template>
  </div>
  <div class="lastTrick grid-container">
    <img v-for="(cardInfo, index) in lastTrick" :key="index" class="vertical" :src="getImgUrl(cardInfo.card, '')" :style="getGridColumnStyle(index)">
  </div>
</div>
</template>

<script>
import axios from 'axios';
import VueCookies from 'vue-cookies';
import { useToast } from "vue-toastification";

export default {
  name: 'RoomPage',
  methods: {
    isLogged() {
        return VueCookies.isKey("player")
    },
    playerName() {
        return VueCookies.get("player")
    },
    buttonsStyle() {
        return 'height: '+ (100/this.buttons().length)+ '%;';
    },
    buttons() {
        let showText = "Show your cards";
        if (this.down.open === true) {
            showText = "Hide your cards";
        }
        let allButtons = [{
            "IsShown": () => (!this.onBuypack && (this.status === 0)),
            "IsDisabled": () => (false),
            "Text": "Open buypack",
            "Click": this.openBuypack
        },  {
            "IsShown": () => (!this.onBuypack && (this.status === 2)),
            "Text": "Take buypack",
            "IsDisabled": () => (false),
            "Click": this.takeBuypack
        }, {
            "IsShown": () => (!this.onBuypack && (this.status === 3) && (this.down.cards.length === 12)),
            "IsDisabled": this.isDropDisabled,
            "Text": "Drop",
            "Click": this.drop
        }, {
            "IsShown": () => (!this.onBuypack && ((this.status === 1) || (this.status === 4))),
            "IsDisabled": this.isMoveDisabled,
            "Text": "Move",
            "Click": this.move
        }, {
            "IsShown": () => ((this.status === 1) || (this.status === 4)),
            "IsDisabled": this.isTakeTrickDisabled,
            "Text": "Take trick",
            "Click": this.takeTrick
        }, {
            "IsShown": () => (this.status !== 5),
            "IsDisabled": () => (false),
            "Text": "Shuffle",
            "Click": this.shuffle
        }, {
            "IsShown": () => (!this.onBuypack && (this.status !== 5)),
            "IsDisabled": () => (false),
            "Text": showText,
            "Click": this.changeVisibility
        }, {
            "IsShown": () => (this.onBuypack && (this.status === 0)),
            "IsDisabled": () => (false),
            "Text": "All pass",
            "Click": this.allPass
        }, {
            "IsShown": () => (this.status === 5),
            "IsDisabled": () => (this.playersCount < 3),
            "Text": "Start",
            "Click": this.start
        }];

        let result = [];
        for (let i = 0; i < allButtons.length; i++) {
            if (allButtons[i].IsShown()) {
                result.push(allButtons[i]);
            }
        }
        return result;
    },
    updateLastError(err) {
	const toast = useToast();
	toast.error(err.message);
    },
    getSelected() {
        var indexes = [];
        for (let i = 0; i < this.selected.length; i++) {
            if (this.selected[i]) {
                indexes.push(i);
            }
        }
        return indexes
    },
    logout() {
        VueCookies.remove("player")
        this.$router.push("/login");
    },
    leaveRoom() {
        this.axios.get(this.backend+"/playerOut").then(() => {
            this.$router.push("/lobby");
        }).catch(this.updateLastError);
    },
    countSelected() {
        return this.getSelected().length
    },
    isDropDisabled() {
        return this.countSelected() != 2
    },
    isMoveDisabled() {
        return this.countSelected() != 1
    },
    isTakeTrickDisabled() {
        return this.center.length < 3
    },
    isHover(index) {
        if (this.hovered[index]) {
            return 'hovered'
        }
        return ''
    },
    isSelected(index) {
        if (this.selected[index]) {
            return 'selected'
        }
        return ''
    },
    getGridColumnStyle(index) {
        return 'grid-row: 1; grid-column: '+(index+1)+' / span 2'
    },
    getGridRowStyle(index) {
        return 'grid-column: 1; grid-row: '+(index+1)+' / span 2'
    },
    getImgUrl(cardInfo, prefix) {
        var images = require.context('../assets/', false, /\.png$/);
        return images('./'+prefix+cardInfo.rank+cardInfo.suit+".png");
    },
    playerDescription(side) {
        return side.name+", "+side.tricks+" tricks"
    },
    fetchData() {
      if (this.isLogged()) {
          this.axios.get(this.backend+"/room").then(response => {
              if (response.data === null) {
                this.$router.push("/lobby");
                return
              }
              this.room = response.data;
              let playerIndex = -1;
              for (let i = 0; i < response.data.sides.length; i++) {
                if (response.data.sides[i].name == this.playerName()) {
                  playerIndex = i;
                }
              }
              if (playerIndex === -1) {
                console.log("player not found: "+response);
                return
              }
              if (this.down.length != response.data.sides[playerIndex].length) {
                this.dropSelected();
              } else {
                for (let i = 0; i < this.down.length; i++) {
                  if (this.down.cards[i] != response.data.sides[playerIndex].cards[i]) {
                    this.dropSelected();
                  }
                }
              }
              this.down = response.data.sides[playerIndex];
              this.left = response.data.sides[(playerIndex+1)%4];
              this.up = response.data.sides[(playerIndex+2)%4];
              this.right = response.data.sides[(playerIndex+3)%4];
              this.center = response.data.center;
              this.status = response.data.status;
              this.onBuypack = (playerIndex == response.data.buypackIndex);
              this.lastTrick = response.data.lastTrick;
              this.playersCount = response.data.playersCount;
          }).catch(this.updateLastError);
      }
    },
    dropSelected() {
        this.selected = [false, false, false, false, false, false, false, false, false, false, false, false]
        this.hovered = [false, false, false, false, false, false, false, false, false, false, false, false]
    },
    updateAll() {
        this.fetchData();
        this.dropSelected();
    },
    shuffle() {
        this.axios.post(this.backend+"/shuffle").then(() => {
            this.updateAll()
        }).catch(this.updateLastError);
    },
    openBuypack() {
        this.axios.post(this.backend+"/openBuypack").then(() => {
            this.updateAll()
        }).catch(this.updateLastError);
    },
    allPass() {
        this.axios.post(this.backend+"/allPass").then(() => {
            this.updateAll()
        }).catch(this.updateLastError);
    },
    takeBuypack() {
        this.axios.post(this.backend+"/takeBuypack").then(() => {
            this.updateAll()
        }).catch(this.updateLastError);
    },
    changeVisibility() {
        this.axios.post(this.backend+"/changeVisibility").then(() => {
            this.updateAll()
        }).catch(this.updateLastError);
    },
    drop() {
        var indexes = this.getSelected();
        this.axios.post(this.backend+"/drop", {"indexes": indexes}).then(() => {
            this.updateAll()
        }).catch(this.updateLastError);
    },
    move() {
        var indexes = this.getSelected();
        this.axios.post(this.backend+"/move", {"index": indexes[0]}).then(() => {
            this.updateAll()
        }).catch(this.updateLastError);
    },
    takeTrick() {
        this.axios.post(this.backend+"/takeTrick").then(() => {
            this.updateAll()
        }).catch(this.updateLastError);
    },
    start() {
        this.axios.post(this.backend+"/roomReady").then(() => {
            this.updateAll()
        }).catch(this.updateLastError);
    }
  },
  data() {
    let nullSide = [];
    return {
      room: null,
      down: nullSide,
      up: nullSide,
      left: nullSide,
      right: nullSide,
      center: nullSide,
      lastTrick: nullSide,
      status: "",
      password: "",
      player: "",
      onBuypack: false,
      selected: [false, false, false, false, false, false, false, false, false, false, false, false],
      hovered: [false, false, false, false, false, false, false, false, false, false, false, false],
      backend: process.env.VUE_APP_HOSTNAME,
      playersCount: 0,
      axios: axios.create({
        withCredentials: true
      })
    }
  },
  mounted() {
    this.fetchData()
  },
  created(){
    if (!this.isLogged()) {
      this.$router.push("/login");
    }
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

.hovered {
    z-index: 100;
}

.selected {
    z-index: 50;
    border-style: outset;
    border: 3px solid blue;
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

.lastTrick {
    grid-column: 3;
    grid-row: 4;
    border: 1px solid;
}

.players {
    grid-column: 1;
    grid-row: 4;
    border: 1px solid;
}

.buttons {
    grid-column: 2;
    grid-row: 4;
    border: 1px solid;
}

.btn-group {
  height: 100%;
}

.btn-group button {
  border: 1px solid black;
  color: black;
  cursor: pointer;
  width: 100%;
  display: block;
}

.btn-group button:not(:last-child) {
  border-bottom: none;
}

.btn-group button:enabled {
  background-color: lightgrey;
}

.btn-group button:disabled {
  background-color: darkgrey;
}

.btn-group button:hover:enabled {
  background-color: #3e8e41;
}

</style>
