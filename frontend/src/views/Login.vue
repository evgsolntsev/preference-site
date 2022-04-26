<template>
    <div :key="componentKey">
      <template v-if="!isLogged()">
        <div><input v-model="player" type="text" placeholder="login"></div>
        <div><input v-model="password" type="text" placeholder="password"></div>
        <button @click="login">Submit</button>
      </template>
      <template v-else>
        <div>Welcome, {{ playerName() }}!</div>
        <button @click="redirectToRoom">To room</button>
        <button @click="logout">Logout</button>
      </template>
    </div>
</template>

<script>
import axios from 'axios';
import VueCookies from 'vue-cookies';
import { useToast } from "vue-toastification";

export default {
  name: 'LoginPage',
  methods: {
    forceRerender() {
      this.componentKey += 1;
    },
    redirectToRoom() {
        this.$router.push("/room");
    },
    isLogged() {
        return VueCookies.isKey("player")
    },
    playerName() {
        return VueCookies.get("player")
    },
    updateLastError(err) {
	const toast = useToast();
	toast.error(err.message);
    },
    login() {
      this.axios.post(this.backend+"/login", {
        "login": this.player,
        "password": this.password,
      }).then(() => {
	VueCookies.set("player", this.player, {expires: "12h"})
        this.forceRerender();
      }).catch(this.updateLastError);
    },
    logout() {
	VueCookies.remove("player")
        this.forceRerender();
    }
  },
  data() {
    return {
      backend: process.env.VUE_APP_HOSTNAME,
      axios: axios.create({
        withCredentials: true
      }),
      player: "",
      password: "",
      componentKey: 0
    }
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
</style>
