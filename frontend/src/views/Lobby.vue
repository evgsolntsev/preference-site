<template>
  <div class="grid-containerLobby">
    <div class="leftLobby">
    </div>
    <div class="centerLobby">
      <Vue3EasyDataTable class="tableLobby" :headers="headers" :items="rooms" @click-row="processClickRow" :rows-items="[5, 10]" :rows-per-page="5"/>
      <button class="buttonLobby" @click="create()">Create room</button>
    </div>
    <div class="rightLobby">
      Welcome, {{ playerName() }}!
      <button @click="logout()">Logout?</button>
    </div>
  </div>
</template>

<style>
.grid-containerLobby {
    display: grid;
    gap: 0px;
    border: 0px;
    min-height: 100%;
    max-height: 100%;
    grid-template-columns: 1fr 8fr 1fr;
    grid-template-rows: 1fr 19fr;
}

.leftLobby {
    grid-row: 1;
    grid-column: 1;
    border: 0px;
}
.centerLobby {
    border: 0px;
    grid-row: 1;
    grid-column: 2;
}
.rightLobby {
    border: 0px;
    grid-row: 1;
    grid-column: 3;
    margin-top: 10%;
}
.tableLobby {
    margin-top: 3%;
    max-height: 100%;
}
.buttonLobby {
    margin-top: 3%;
}
</style>

<script>
import axios from 'axios';
import VueCookies from 'vue-cookies';
import Vue3EasyDataTable from 'vue3-easy-data-table';
import 'vue3-easy-data-table/dist/style.css';
import { useToast } from "vue-toastification";

const headers = [
  { text: "ID", value: "id" },
  { text: "PLAYERS", value: "players"},
  { text: "STATUS", value: "status"},
  { text: "", value: "button"},
];

export default {
  name: 'LobbyPage',
  components: { Vue3EasyDataTable },
  methods: {
    isLogged() {
        return VueCookies.isKey("player");
    },
    playerName() {
        return VueCookies.get("player")
    },
    updateLastError(err) {
	const toast = useToast();
	toast.error(err.message);
    },
    fetchData() {
        this.axios.get(this.backend+"/rooms").then(response => {
            this.rooms = response.data;
        }).catch(this.updateLastError);
    },
    create() {
        this.axios.get(this.backend+"/createRoom").then(() => {
            this.$router.push("/room");
        }).catch(this.updateLastError);
    },
    processClickRow(row) {
        this.axios.post(this.backend+"/playerIn", {"roomId": row.id}).then(() => {
            this.$router.push("/room");
        }).catch(this.updateLastError);
    },
    logout() {
        VueCookies.remove("player")
        this.$router.push("/login");
    }
  },
  data() {
    return {
      rooms: [],
      headers: headers,
      backend: process.env.VUE_APP_HOSTNAME,
      axios: axios.create({
        withCredentials: true
      })
    }
  },
  created() {
    if (!this.isLogged()) {
      this.$router.push("/login");
    }
    this.fetchData()
    this.interval = setInterval(() => {
      this.fetchData()}, 5000)
  }
}
</script>
