package ovhcloud

// CreateInstanceOptions holds all parameters to create an instance.
type CreateInstanceOptions struct {
	Name      string
	Flavor    string
	Image     string
	PublicKey []byte
}
