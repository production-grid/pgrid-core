<template>
  <div class="slim-mainpanel">
    <div class="container">
      <div class="slim-pageheader">
        <ol class="breadcrumb slim-breadcrumb">
          <li class="breadcrumb-item" v-if="breadcrumbCategory"><a href="#">{{breadcrumbCategory}}</a></li>
          <li class="breadcrumb-item active" aria-current="page"  v-if="form.id || isSingleton">{{editFormTitle}}</li>
          <li class="breadcrumb-item active" aria-current="page"  v-if="!form.id">{{md.newFormTitle}}</li>
        </ol>
        <h6 class="slim-pagetitle" v-if="!form.id">{{md.newFormTitle}}</h6>
        <h6 class="slim-pagetitle" v-if="form.id || isSingleton">{{editFormTitle}}</h6>
      </div><!-- slim-pageheader -->
      <div class="section-wrapper">
        <div v-if="!form.id">
          <label class="section-title" v-if="!form.id">{{md.newFormTitle}}</label>
          <p class="mg-b-20 mg-sm-b-40" v-if="md.newFormHelp">{{md.newFormHelp}}</p>
       </div>
       <div v-if="form.id  || isSingleton" class="row">
         <div class="col-md-8">
           <label class="section-title" v-if="form.id || isSingleton">{{editFormTitle}}</label>
           <p class="mg-b-20 mg-sm-b-40" v-if="md.editFormHelp">{{md.editFormHelp}}</p>
         </div>

         <div class="col-md-4 mb-2" v-if="hasActionMenu">
           <div id="actionMenu" class="dropdown show float-right">
              <a class="btn btn-primary dropdown-toggle" href="#" role="button" id="dropdownMenuLink" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                ACTIONS
              </a>

              <div class="dropdown-menu dropdown-menu-right" aria-labelledby="dropdownMenuLink">
                <router-link  v-for="result in extraRoutes" class="dropdown-item"  v-bind:key="result.route" :to="{ path: result.route, query: { id: form.id }}">{{result.caption}}</router-link>
              </div>
            </div>
        </div>

       </div>
       <div class="alert alert-outline alert-success" role="alert" v-if="savedVisible">
        <button type="button" class="close" data-dismiss="alert" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
        <strong>Done.</strong> Your changes have been saved.
      </div><!-- alert -->

       <form role="form" v-on:submit="submit">
         <crud-panel
           :globalState="globalState"
           :mode="mode"
           v-model="form"
           :fields="md.fields"
           :fieldsets="md.fieldsets"
           :columnStyle="effectiveColumnStyle"
           :session="session"
         />
           <div v-if="tabsEnabled">
             <div>
               <ul class="nav nav-tabs" role="tablist">
                 <li class="nav-item" v-for="(tab, index) in md.tabs" v-bind:key="tab.id"><a @click="changeTab(index)" v-bind:class="{ active: index == 0 }" class="nav-link font-weight-bold" data-toggle="tab">{{tab.caption}}</a></li>
               </ul>
               <div class="tab-content mg-5">
                 <div class="tab-pane" role="tabpanel" v-bind:class="{ active: index == activeTab }" v-for="(tab, index) in md.tabs" :id="tab.tabId"  v-bind:key="tab.id">
                   <crud-panel
                     :globalState="globalState"
                     :mode="mode"
                     v-model="form"
                     :fields="tab.fields"
                     :fieldsets="tab.fieldsets"
                     :columnStyle="effectiveColumnStyle"
                     :session="session"
                   />
                 </div>
               </div>
             </div>
           </div>
           <div class="mt-4">
             <div class="progress" v-if="processing">
               <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: 100%" aria-valuenow="100" aria-valuemin="0" aria-valuemax="100"></div>
             </div>
             <div v-if="!processing">
               <button type="button" class="btn btn-warning" data-test-id="canceButton" @click="cancelForm()" v-if="listRoute">CANCEL</button>
               <button type="button" class="btn btn-danger mg-l-5" v-if="isDeletable" @click="deleteConfirm()">DELETE</button>
               <button type="submit" class="btn btn-primary mg-l-5 float-right" data-test-id="submitButton">SAVE</button>
             </div>
           </div>
         </form>
       </div><!-- section-wrapper -->
       <!-- /.content -->
       <delete-modal
          :title="deleteDialogTitle"
          v-on:deleteConfirmed="deleteConfirmed"
          id="modal-delete"
          :returnId="form.id"
          :message="deletePrompt">
      </delete-modal>
    </div>
  </div>
</template>

<script>
import APIService from '@/services/APIService'
import DeleteModal from '../modals/DeleteModal'
import router from '../router'
import modalService from '../services/ModalService'
import CrudPanel from './CrudPanel'

export default {
  name: 'crud-form-content-block',
  props: ['globalState', 'breadcrumbCategory', 'resource', 'id', 'listRoute', 'extraRoutes', 'mode', 'singleton', 'titleField', 'router', "columnStyle", "session"],
  components: {
    'delete-modal': DeleteModal,
    'crud-panel': CrudPanel
  },
  data: () => ({
    md: {},
    form: {id: null},
    savedVisible: false,
    activeTab: 0,
    processing: false
  }),
  computed: {
    effectiveColumnStyle: function () {
      if (!this.columnStyle) {
        return 'col-lg-6'
      }
      return this.columnStyle
    },
    isSingleton: function () {
      return this.singleton
    },
    editFormTitle: function () {
      if (this.isSingleton && this.form[this.titleField]) {
        return this.form[this.titleField]
      } else {
        return this.md.editFormTitle
      }
    },
    hasActionMenu: function () {
      return this.extraRoutes && this.extraRoutes.length > 0
    },
    isDeletable: function () {
      return this.form.id && this.md.deleteEnabled
    },
    deleteDialogTitle: function () {
      return 'Confirm ' + this.md.resourceName + ' Delete'
    },
    deletePrompt: function () {
      if (this.md.deletePrompt) {
        return this.md.deletePrompt
      } else {
        return 'Are you sure you want to delete this ' + this.md.resourceName + '?'
      }
    },
    tabsEnabled: function () {
      return this.md.tabs && this.md.tabs.length > 0
    }
  },
  watch: {
    '$route.query.id' () {
      if (!this.isSingleton) {
        this.id = this.$route.query.id
        this.loadFormData()
      }
    }
  },
  mounted () {
    let self = this
    APIService.get('/api/' + self.resource + '/md', function (response) {
      self.md = response.data
      if (self.id || self.isSingleton) {
        self.loadFormData()
      }
    })
  },
  methods: {
    submit: function (e) {
      e.preventDefault()
      let self = this
      self.processing = true
      APIService.post('/api/' + self.resource, this.form, function (response) {
        self.form = response.data
        self.savedVisible = true
        self.processing = false
        setTimeout(function () {
          self.savedVisible = false
        }, 5000)
      })
    },
    changeTab: function (index) {
      this.activeTab = index
    },
    cancelForm: function () {
      router.push(this.listRoute)
    },
    loadFormData: function () {
      let self = this
      let url = '/api/' + self.resource
      if (!self.isSingleton) {
        url = url + '/' + self.id
      }
      APIService.get(url, function (response) {
        self.form = response.data
        self.$emit('loaded', self.form)
      })
    },
    deleteConfirm: function () {
      modalService.raiseModal('modal-delete')
    },
    deleteConfirmed: function (deleteId) {
      let self = this
      self.processing = true
      APIService.delete('/api/' + self.resource + '/' + deleteId, function (response) {
        self.processing = false
        if (response.data.success) {
          router.push({ path: self.listRoute })
        } else {
          alert(response.data.error)
        }
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

.nav-link {
  cursor: pointer;
}

#actionMenu {
  margin-top: -10px;
}

</style>
