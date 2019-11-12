<template>
  <div class="slim-mainpanel">
    <div class="container">
      <div class="slim-pageheader">
        <ol class="breadcrumb slim-breadcrumb">
          <li class="breadcrumb-item" v-if="breadcrumbCategory">{{breadcrumbCategory}}</li>
          <li class="breadcrumb-item active" aria-current="page">{{md.listPageTitle}}</li>
        </ol>
        <h6 class="slim-pagetitle">{{md.listPageTitle}}</h6>
      </div><!-- slim-pageheader -->
      <div class="section-wrapper">
         <label class="section-title">{{md.listPageTitle}} <router-link v-if="!md.newDisabled" :to="{ path: newLink}" class="float-right btn btn-sm btn-outline-primary"><i class="fa fa-plus"></i> {{md.newFormTitle}}</router-link></label>
         <p class="mg-b-20 mg-sm-b-40">{{md.listPageHelp}}</p>
         <div class="search-panel mg-b-5 row" v-if="searchEnabled">
           <div class="col-md-5">
             <div class="input-group">
              <input type="text" class="form-control" v-on:keyup.13="executeSearch()" placeholder="Search" v-model="searchCriteria" aria-describedby="basic-addon2">
              <div class="input-group-append">
                <a class="input-group-text" id="basic-addon2" @click="executeSearch()"><i  class="fa fa-search"></i></a>
              </div>
            </div>
             <div class="include-deleted-block" v-if="includeDeletedOption"><input type="checkbox" @change="executeSearch()" v-model="includeDeleted"/><span class="text-muted">Include Deleted</span></div>
           </div>
         </div>

         <no-results-alert :results="results" :loading="loading" :caption="emptyListCaption"/>
         <div class="table-wrapper" v-if="hasResults">
           <table class="table table-striped table-bordered display nowrap">
             <thead>
               <tr>
                 <th v-for="column in md.columns" v-bind:key="column.id">{{column.caption}}</th>
               </tr>
             </thead>
             <tbody>
               <tr v-for="result in results"  v-bind:key="result.id" v-bind:data-test-id="result.id" data-test-class="subscriber-row">
                 <td v-for="column in md.columns" v-bind:key="column.id">
                   <router-link :to="clickRoute + '/' + result.id">{{result[column.id]}}</router-link>
                 </td>
               </tr>
             </tbody>
           </table>
         </div><!-- table-wrapper -->
         <div class="pagination-wrapper" v-if="pagingVisible">
            <nav aria-label="Page navigation">
              <ul class="pagination mg-b-0">
                <li class="page-item active"><a class="page-link" href="#">1</a></li>
                <li class="page-item"><a class="page-link" href="#">2</a></li>
                <li class="page-item"><a class="page-link" href="#">3</a></li>
                <li class="page-item"><a class="page-link" href="#">4</a></li>
                <li class="page-item"><a class="page-link" href="#">5</a></li>
                <li class="page-item">
                  <a class="page-link" href="#" aria-label="Next">
                    <i class="fa fa-angle-right"></i>
                  </a>
                </li>
              </ul>
            </nav>
          </div><!-- pagination-wrapper -->
       </div><!-- section-wrapper -->
    </div>
  </div>
</template>

<script>
import APIService from '@/services/APIService'
import NoResultsAlert from '../components/NoResultsAlert'

export default {
  name: 'crud-list',
  props: ['globalState', 'breadcrumbCategory', 'resource', 'newRoute', 'clickRoute', 'searchEnabled', 'includeDeletedOption'],
  components: {
    'no-results-alert': NoResultsAlert
  },
  data: () => ({
    results: [],
    md: {},
    searchCriteria: null,
    includeDeleted: false,
    paging: null,
    loading: true
  }),
  computed: {
    newLink: function () {
      if (this.newRoute) {
        return this.newRoute
      } else {
        return this.clickRoute
      }
    },
    pagingVisible: function () {
      return false
    },
    hasResults: function () {
      return this.results && this.results.length > 0
    },
    emptyListCaption: function () {
      if (this.md.emptyListCaption) {
        return this.md.emptyListCaption
      } else {
        return 'This list contains no data.'
      }
    }
  },
  methods: {
    executeSearch: function () {
      let self = this
      let url = '/api/' + self.resource + '?'
      if (self.searchCriteria) {
        url = url + "q=" + self.searchCriteria + "&"
      }
      if (self.includeDeleted) {
        url = url + "includeDeleted=true"
      }
      APIService.get(url)
        .then(response => {
          self.paging = response.data
          self.results = response.data.results
        })
    }
  },
  mounted () {
    let self = this
    APIService.get('/api/' + self.resource + '/md', function (response) {
      self.md = response.data
      APIService.get('/api/' + self.resource, function (response) {
        self.paging = response.data
        self.results = response.data.visibleResults
        self.loading = false
      })
    })
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

.newButton {
  text-transform: uppercase;
}

.activateLink a { cursor: pointer; cursor: hand; }

</style>
