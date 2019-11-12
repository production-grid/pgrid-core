<template>
     <div class="modal" id="modal-delete">
       <div class="modal-dialog" role="document">
         <div class="modal-content bd-0">
           <div class="modal-header pd-y-20 pd-x-25">
             <h6 class="tx-14 mg-b-0 tx-uppercase tx-inverse tx-bold">{{title}}</h6>
             <button type="button" class="close" data-dismiss="modal" aria-label="Close">
               <span aria-hidden="true">&times;</span>
             </button>
           </div>
           <div class="modal-body pd-25">
             <h5 class="lh-4 mg-b-20 tx-inverse">{{message}}</h5>
           </div>
           <div class="modal-footer">
             <button type="button" class="btn btn-danger"  v-on:click="deleteConfirm">{{confirmCaption}}</button>
             <button type="button" class="btn btn-secondary" data-dismiss="modal">CLOSE</button>
           </div>
         </div>
       </div><!-- modal-dialog -->
     </div><!-- modal -->
</template>

<script>
export default {
  name: 'delete-modal',
  props: ['title', 'message', 'returnId', 'confirmButtonCaption'],
  computed: {
    confirmCaption: function () {
      if (this.confirmButtonCaption) {
        return this.confirmButtonCaption
      } else {
        return "DELETE"
      }
    }
  },
  methods: {
    deleteConfirm: function () {
      $("#modal-delete").modal('hide')
      this.$emit('deleteConfirmed', this.returnId)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
