package applications

// TenantLingo Models the way tenants should be presented to the world.
type TenantLingo struct {
	TenantSingularKey     string
	TenantPluralKey       string
	TenantSingularDefault string
	TenantPluralDefault   string
	Types                 []TenantType
}

// TenantType models the variations within tenants, if there are any
type TenantType struct {
	SingularKey     string
	PluralKey       string
	SingularDefault string
	PluralDefault   string
}
