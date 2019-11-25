<template>
  <div class="login card">
    <div class="card-body">
      <div class="row login-header">
        <div class="col-sm-3">
          <img src="./../../assets/pg-logo.png" class="logo"/>
        </div> <!-- col -->
        <div class="col-sm-9">
          <h4>{{session.applicationName}}</h4>
          <p>{{session.tagline}}</p>
        </div> <!-- col -->
      </div> <!-- row -->
      <div class="alert alert-danger alert-outline" v-if="errorMessage">
        {{errorMessage}}
      </div> <!-- alert -->
      <form @submit="processLogin">
          <div class="form-group">
            <label for="emailInput">E-Mail Address or Mobile Number</label>
            <input type="text" class="form-control" id="emailInput" placeholder="Enter email or mobile number" v-model="loginId">
          </div> <!-- form-group -->
          <div class="form-group">
            <label for="passwordInput">Password</label>
            <input type="password" class="form-control" id="passwordInput" placeholder="Password" v-model="password">
          </div> <!-- form-group -->
          <div class="progress" v-if="loginProcessing">
            <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: 100%" aria-valuenow="100" aria-valuemin="0" aria-valuemax="100"></div>
          </div>
          <button type="submit" class="btn btn-dark btn-block" v-if="!loginProcessing"><i class="fas fa-sign-in-alt"></i> Login</button>
          <div class="text-center mg-t-10"><a href="#">I forgot my password.</a></div>
     </form>
   </div> <!-- card-body -->
 </div> <!-- login card -->
</template>

<script>
import APIService from '@/services/APIService'

export default {
  name: 'HelloWorld',
  props: ['session'],
  data: () => ({
    loginId: null,
    password: null,
    errorMessage: null,
    loginProcessing: false
  }),
  methods: {
    processLogin: function(e) {
      e.preventDefault()
      let self = this
      self.errorMessage = null
      self.loginProcessing = true
      let req = {
        loginId: this.loginId,
        password: this.password
      }
      APIService.post('/api/security/login', req, function (response) {
        self.loginProcessing = false
        let loginResponse = response.data
        if (!loginResponse.success) {
          self.errorMessage = loginResponse.error
        } else {
          self.$emit('sessionInvalidated')
        }
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

@media (max-width: 600px) {
  .login-header {
    text-align: center;
  }
  .login {
    margin-top: 10px !important;
  }
}

form {
  margin-top: 20px;
}

.logo {
  max-width: 100px;
}

.alert {
  margin-top: 20px;
}

.login {
  max-width: 500px;
  margin-top: 200px;
  margin-left: auto;
  margin-right: auto;
}

</style>
