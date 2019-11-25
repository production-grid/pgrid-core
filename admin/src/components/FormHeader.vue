<template>
 <div class="row">
   <div class="col-md-8">
     <label class="section-title">{{formTitle}}</label>
     <p class="mg-b-20 mg-sm-b-40" v-if="helpText">{{helpText}}</p>
   </div> <!-- col-md-8 -->

   <div class="col-md-4 mb-2">
     <router-link v-if="newLink" :to="{ path: newLink}" class="float-right btn btn-sm btn-outline-primary newButton"><i class="fa fa-plus"></i> {{newLinkCaption}}</router-link>
     <div id="actionMenu" class="dropdown show float-right" v-if="hasActionMenu">
        <a class="btn btn-primary dropdown-toggle" href="#" role="button" id="dropdownMenuLink" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
          ACTIONS
        </a>
        <div class="dropdown-menu dropdown-menu-right" aria-labelledby="dropdownMenuLink">
          <router-link  v-for="result in extraRoutes" class="dropdown-item"  v-bind:key="result.route" :to="{ path: result.route, query: { id: id }}">{{result.caption}}</router-link>
        </div> <!-- dropdown-menu -->
      </div> <!-- actionMenu -->
    </div> <!-- col-md-4 -->
  </div> <!-- row -->
</template>


<script>

export default {
  name: 'form-header',
  props: ['globalState', 'extraRoutes', 'formTitle', 'helpText', 'id', 'newLink', 'newLinkCaption'],
  components: {
  },
  data: () => ({
  }),
  computed: {
    hasActionMenu: function () {
      return this.extraRoutes && this.extraRoutes.length > 0
    },
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

.newButton {
  text-transform: uppercase;
}

.activateLink a { cursor: pointer; cursor: hand; }

</style>
