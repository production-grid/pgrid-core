export default  {

  defaultRoutePath(session) {
    if (!session.effectivePermissions) {
      return "/login"
    } else {
      return "/hello"
    }
  },

  addRoutes(routes, moduleRoutes) {

    return routes.concat(moduleRoutes)

  }

}
