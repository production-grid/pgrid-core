import axios from 'axios'

export default  {

  get(path, thenHandler, catchHandler, finallyHandler) {

    axios.get(path)
      .then(function (response) {
        thenHandler(response)
      })
      .catch(function (error) {
        if (catchHandler) {
          catchHandler(error)
        }
      })
      .finally (function () {
        if (finallyHandler) {
          finallyHandler()
        }
      })
  },

  post(path, body, thenHandler, catchHandler, finallyHandler) {

    axios.post(path, body)
      .then(function (response) {
        thenHandler(response)
      })
      .catch(function (error) {
        if (catchHandler) {
          catchHandler(error)
        }
      })
      .finally (function () {
        if (finallyHandler) {
          finallyHandler()
        }
      })
  },

  delete(path, thenHandler, catchHandler, finallyHandler) {

    axios.delete(path)
      .then(function (response) {
        thenHandler(response)
      })
      .catch(function (error) {
        if (catchHandler) {
          catchHandler(error)
        }
      })
      .finally (function () {
        if (finallyHandler) {
          finallyHandler()
        }
      })
  }

}
