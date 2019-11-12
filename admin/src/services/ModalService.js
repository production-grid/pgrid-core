
const raiseModal = (modalId) => {
  $('#' + modalId).modal('show')
}

const closeModal = (modalId) => {
  $('#' + modalId).modal('hide')
}

const alert = (message, type) => {
  if (!type) {
    type = 'success'
  }
  $.notify({
    message: message
  }, {
    type: type,
    delay: 3000,
    animate: {
      exit: 'animated fadeOutUp'
    }
  })
}

export default {
  raiseModal,
  closeModal,
  alert
}
