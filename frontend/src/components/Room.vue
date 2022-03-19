<template>
<div class="row">
  <div class="left column">
    1
  </div>
  <div class="center column">
    2
  </div>
  <div class="right column">
    3
  </div>
</div>
<div class="row">
  <div class="left column">
    4
  </div>
  <div class="center column">
    5
  </div>
  <div class="right column">
    6
  </div>
</div>
<div class="row">
  <div class="left column">
    7
  </div>
  <div class="center column">
    8
  </div>
  <div class="right column">
    9
  </div>
</div>

<!--
  <table>
    <tr class="playable">
      <td>1
      </td>
      <td>
        <img v-for="card in up.cards" :key="card.card.suit+card.card.rank" :src="getImgUrl(card)">
      </td>
      <td>3
      </td>
    </tr>
    <tr>
      <td class="playable">
      </td>
      <td> {{ center }}
      </td>
      <td class="playable">
      </td>
    </tr>
    <tr class="playable">
      <td>7
      </td>
      <td>
        <img v-for="cardInfo in down.cards" :key="cardInfo.card.suit+cardInfo.card.rank" :src="getImgUrl(cardInfo)">
      </td>
      <td>9
      </td>
    </tr>
  </table>
-->
</template>

<script>
import axios from 'axios';
export default {
  name: 'RoomTable',
}
</script>

<style>
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
