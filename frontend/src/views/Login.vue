<template>
  <div class="login-page">
    <div class="form">
        <div :style="registerStyle()">
          <input v-model="playerRegister" type="text" placeholder="name"/>
          <input v-model="passwordRegister" type="password" placeholder="password"/>
          <input v-model="emailRegister" type="text" placeholder="email address"/>
          <button @click="register">create</button>
          <p>Already registered? <a href="#" @click="toLogin">Sign In</a></p>
        </div>
        <div :style="loginStyle()">
          <input v-model="player" type="text" placeholder="username"/>
          <input v-model="password" type="password" placeholder="password"/>
          <button @click="login">login</button>
          <p>Not registered? <a href="#" @click="toRegister">Create an account</a></p>
        </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import VueCookies from 'vue-cookies';
import { useToast } from "vue-toastification";

export default {
  name: 'LoginPage',
  methods: {
    redirect() {
        this.axios.get(this.backend+"/room").then(response => {
            if (response.data === null) {
                this.$router.push("/lobby");
                return
            }
            this.$router.push("/room");
	}).catch(this.updateLastError);
    },
    isLogged() {
        return VueCookies.isKey("player")
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
        VueCookies.set("player", this.player, {expires: "12h"});
        this.redirect()
      }).catch(this.updateLastError);
    },
    register() {
      this.axios.post(this.backend+"/register", {
        "login": this.playerRegister,
        "password": this.passwordRegister,
        "email": this.emailRegister,
      }).then(() => {
        const toast = useToast();
          toast.success("User created");
          this.toLogin()
      }).catch(this.updateLastError);
    },
    logout() {
        VueCookies.remove("player")
        this.$router.push("/login");
    },
    toRegister() {
        this.type = "register"
    },
    toLogin() {
        this.type = "login"
    },
    loginStyle() {
        if (this.type != "login") {
            return "display: none;"
        }
        return ""
    },
    registerStyle() {
        if (this.type != "register") {
            return "display: none;"
        }
        return ""
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
      playerRegister: "",
      passwordRegister: "",
      emailRegister: "",
      type: "login"
    }
  },
  created(){
    if (this.isLogged()) {
      this.redirect()
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

@import url(https://fonts.googleapis.com/css?family=Roboto:300);

.login-page {
  width: 360px;
  padding: 8% 0 0;
  margin: auto;
}
.form {
  position: relative;
  z-index: 1;
  background: #FFFFFF;
  max-width: 360px;
  margin: 0 auto 100px;
  padding: 45px;
  text-align: center;
  box-shadow: 0 0 20px 0 rgba(0, 0, 0, 0.2), 0 5px 5px 0 rgba(0, 0, 0, 0.24);
}
.form input {
  font-family: "Roboto", sans-serif;
  outline: 0;
  background: #f2f2f2;
  width: 100%;
  border: 0;
  margin: 0 0 15px;
  padding: 15px;
  box-sizing: border-box;
  font-size: 14px;
}
.form button {
  font-family: "Roboto", sans-serif;
  text-transform: uppercase;
  outline: 0;
  background: #4CAF50;
  width: 100%;
  border: 0;
  padding: 15px;
  color: #FFFFFF;
  font-size: 14px;
  -webkit-transition: all 0.3 ease;
  transition: all 0.3 ease;
  cursor: pointer;
}
.form button:hover,.form button:active,.form button:focus {
  background: #43A047;
}
.form .message {
  margin: 15px 0 0;
  color: #b3b3b3;
  font-size: 12px;
}
.form .message a {
  color: #4CAF50;
  text-decoration: none;
}
.container {
  position: relative;
  z-index: 1;
  max-width: 300px;
  margin: 0 auto;
}
.container:before, .container:after {
  content: "";
  display: block;
  clear: both;
}
.container .info {
  margin: 50px auto;
  text-align: center;
}
.container .info h1 {
  margin: 0 0 15px;
  padding: 0;
  font-size: 36px;
  font-weight: 300;
  color: #1a1a1a;
}
.container .info span {
  color: #4d4d4d;
  font-size: 12px;
}
.container .info span a {
  color: #000000;
  text-decoration: none;
}
.container .info span .fa {
  color: #EF3B3A;
}
body {
  font-family: "Roboto", sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;      
}
</style>
