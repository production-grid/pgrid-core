import SortService from '@/services/SortService'

export default  {


  defaultRoutePath(session) {
    if (!session.effectivePermissions) {
      return "/login"
    } else {
      return "/hello"
    }
  },

  addRoutes(routes, moduleRoutes) {

    if (!this.navRoutes) {
      this.navRoutes = []
      this.sectionList = []
      this.sectionMap = {}
    }

    for (var i = 0; i < moduleRoutes.length; i++) {
      let route = moduleRoutes[i]
      if (route.nav) {
        let section = this.sectionMap[route.nav.section]
        if (!section) {
          section = {
            id: route.nav.section
          }
          this.sectionList.push(section)
        }
        if (route.nav.sectionIndex) {
          section.index = route.nav.sectionIndex
        }
        if (route.nav.sectionIconClass) {
          section.iconClass = route.nav.sectionIconClass
        }
        if (!section.routes) {
          section.routes = []
        }
        if (!route.nav.index) {
          route.nav.index = section.routes.length
        }
        route.index = route.nav.index
        section.routes.push(route)
        this.sectionMap[section.id] = section
      }
    }

    this.navRoutes = this.navRoutes.concat(moduleRoutes)


    return routes.concat(moduleRoutes)

  },

  isAuthorized(route, session) {

    if (!route.nav) {
      return true
    }

    console.log(route)
    console.log(session.effectivePermissions)

    if (route.nav.permission && !session.effectivePermissions.includes(route.nav.permission)) {
      return false
    }

    var i = 0

    if (route.nav.allRequired && (route.nav.allRequired.length > 0)) {
      for (i = 0; i < route.nav.allRequired.length; i++) {
        if (!session.effectivePermissions.includes(route.nav.allRequired[i])) {
          return false
        }
      }
    }

    if (route.nav.anyRequired && (route.nav.anyRequired.length > 0)) {
      for (i = 0; i < route.nav.anyRequired.length; i++) {
        if (session.effectivePermissions.includes(route.nav.anyRequired[i])) {
          return true
        }
      }

      return false
    }

    return true

  },

  visibleNavRoutes(session) {

    if (session) {

      let results = []
      let keys = Object.keys(this.sectionMap)
      for (var i = 0; i < keys.length; i++) {
        let section = this.sectionMap[keys[i]]
        let sectionCopy = {
          id: section.id,
          index: section.index,
          iconClass: section.iconClass
        }
        let routes = []
        for (var r = 0; r < section.routes.length; r++) {
          let route = section.routes[r]
          if (this.isAuthorized(route, session)) {
            routes.push(route)
          }
        }
        sectionCopy.routes = SortService.quickSortAll(routes)

        results.push(sectionCopy)
      }

      return SortService.quickSortAll(results)
    }

  }

}
