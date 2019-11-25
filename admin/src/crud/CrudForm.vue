<template>
  <div class="slim-mainpanel">
    <div class="container">
      <page-header :breadcrumbCategory="breadcrumbCategory" :pageTitle="pageTitle"/>
      <div class="section-wrapper">
       <form-header :formTitle="pageTitle" :extraRoutes="extraRoutes" :helpText="helpText" :id="id" />
       <saved-alert :visible="savedVisible"/>

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
           <form-button-panel :processing="processing" @delete="deleteConfirm" :isDeletable="isDeletable" :listRoute="listRoute"/>
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
import PageHeader from '../components/PageHeader'
import FormHeader from '../components/FormHeader'
import SavedAlert from '../components/SavedAlert'
import FormButtonPanel from '../components/FormButtonPanel'

export default {
  name: 'crud-form-content-block',
  props: ['globalState', 'breadcrumbCategory', 'resource', 'id', 'listRoute', 'extraRoutes', 'mode', 'singleton', 'titleField', 'router', "columnStyle", "session", "newForm"],
  components: {
    'delete-modal': DeleteModal,
    'crud-panel': CrudPanel,
    'page-header': PageHeader,
    'form-header': FormHeader,
    'saved-alert': SavedAlert,
    'form-button-panel': FormButtonPanel
  },
  data: () => ({
    md: {},
    form: {id: null},
    savedVisible: false,
    activeTab: 0,
    processing: false
  }),
  computed: {
    pageTitle: function () {
      if (this.form.id || this.isSingleton) {
        return this.editFormTitle
      }
      return this.md.newFormTitle
    },
    helpText: function () {
      if (this.form.id || this.isSingleton) {
        return this.md.editFormHelp
      }
      return this.md.newFormHelp
    },
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
      } else if (self.newForm) {
        self.form = self.newForm
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
