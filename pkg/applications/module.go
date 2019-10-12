package applications

//FeatureModule defines the base methods required to define a feature module
type FeatureModule interface {
	ID() string
	BeforeAppInit(*Application) error
	AfterAppInit(*Application) error
	BeforeModuleInit(*Application) error
	AfterModuleInit(*Application) error
	SchemaFiles(*Application) ([]string, error)
	Permissions(*Application) []Permission
	APIRoutes(*Application) ([]APIRoute, error)
}
