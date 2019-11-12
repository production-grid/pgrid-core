<template>
  <div class="pagination-wrapper">
     <nav aria-label="Page navigation">
       <ul class="pagination mg-b-0">
         <li class="page-item" v-if="pagingInfo.currentPage > 1">
           <a class="page-link" @click="newPage(currentPage - 1)" aria-label="Previous">
             <i class="fa fa-angle-left"></i>
           </a>
         </li>
         <li class="page-item"
            v-bind:key="page"
            v-bind:class="{ active: page === currentPage }"
            v-for="page in pages">
            <a class="page-link"  @click="newPage(page)">{{page}}</a>
         </li>
         <li class="page-item" v-if="currentPage < pagingInfo.pages">
           <a class="page-link"  @click="newPage(currentPage + 1)" aria-label="Next">
             <i class="fa fa-angle-right"></i>
           </a>
         </li>
       </ul>
     </nav>
   </div><!-- pagination-wrapper -->
</template>

<script>
export default {
  name: 'paging-widget',
  props: ['globalState', 'pagingInfo'],
  computed: {
    currentPage: function () {
      return this.pagingInfo.currentPage
    },
    pages: function () {
      var results = []
      for (var i = 1 ; i <= this.pagingInfo.pages; i++) {
        results.push(i)
      }
      return results
    }
  },
  methods: {
    isCurrentPage: function(pageIndex) {
      return pageIndex === this.currentPage
    },
    newPage: function(pageIndex) {
      let startIndex = ((pageIndex - 1) * this.pagingInfo.pageSize)
      this.$emit('pageChange', startIndex)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
