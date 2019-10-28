package applications

// TenantLingo Models the way tenants should be presented to the world.
type TenantLingo struct {
	TenantSingular string
	TenantPlural   string
	Types          []TenantType
}

// TenantType models the variations within tenants, if there are any
type TenantType struct {
	Singular string
	Plural   string
}
