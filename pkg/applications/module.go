package applications

//FeatureModule defines the base methods required to define a feature module
type FeatureModule interface {
	Name() string
	BeforeAppInit(*Application) error
	AfterAppInit(*Application) error
	BeforeModuleInit(*Application) error
	AfterModuleInit(*Application) error
	SchemaFiles(*Application) ([]string, error)
}
