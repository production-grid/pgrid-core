<template>
  <div class="slim-mainpanel">
    <div class="container">
      <page-header breadcrumbCategory="Security" pageTitle="Invite New User"/>
      <div class="section-wrapper">
        <form-header formTitle="Invite New User" helpText="Invite a new user."/>
        <saved-alert :visible="savedVisible"/>
        <form role="form" v-on:submit="submit">
          <div class="form-layout mg-t-10">
            <div class="row">
              <div class="col-md-4">
                <div class="form-group">
                  <label class="form-control-label font-weight-bold">First Name: <span class="tx-danger">*</span></label>
                  <input  class="form-control" v-model="firstName" required/>
                </div> <!-- form-group -->
               </div> <!-- col-md-4 -->
               <div class="col-md-4">
                 <div class="form-group">
                   <label class="form-control-label font-weight-bold">Last Name: <span class="tx-danger">*</span></label>
                   <input  class="form-control" v-model="lastName" required/>
                 </div> <!-- form-group -->
               </div> <!-- col-md-4 -->
            </div> <!-- row -->
            <div class="row">
              <div class="col-md-8">
                <div class="form-group">
                  <label class="form-control-label font-weight-bold">E-Mail Address: <span class="tx-danger">*</span></label>
                  <input  class="form-control" type="email" v-model="emailAddress" required/>
                </div> <!-- form-group -->
               </div> <!-- col-md-4 -->
            </div> <!-- row -->
            <div class="row">
              <div class="col-md-12">
                <div class="form-group">
                  <label class="form-control-label font-weight-bold">Groups: <span class="tx-danger">*</span></label>
                  <group-input v-model="groups" :session="session" groupSource="/api/security/admin/groups"/>

                </div> <!-- form-group -->
               </div> <!-- col-md-4 -->
            </div> <!-- row -->
          </div> <!-- form-layout -->
          <form-button-panel :processing="processing" listRoute="/admin-users" saveCaption="INVITE"/>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import APIService from '@/services/APIService'
import PageHeader from '../../components/PageHeader'
import FormHeader from '../../components/FormHeader'
import FormButtonPanel from '../../components/FormButtonPanel'
import GroupInput from '../../widgets/GroupInput'

export default {
  name: 'admin-user-new',
  props: ['globalState', 'session'],
  components: {
    'page-header': PageHeader,
    'form-header': FormHeader,
    'form-button-panel': FormButtonPanel,
    'group-input': GroupInput
  },
  data: () => ({
    firstName: null,
    lastName: null,
    emailAddress: null,
    groups: [],
    savedVisible: false,
    processing: false
  }),
  methods: {
    submit: function(e) {
      e.preventDefault()
      let self = this
      self.processing = true
      let req = {
        firstName: this.firstName,
        lastName: this.lastName,
        emailAddress: this.emailAddress,
        groups: this.groups
      }
      APIService.post('/api/security/admin-invite', req, function () {
        self.processing = false
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
