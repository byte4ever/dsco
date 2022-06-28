package walker

type PoliciesGetter interface {
	GetPolicies() (constraintLayerPolicies, error)
}
