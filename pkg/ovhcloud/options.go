package ovhcloud

type CreateInstanceOptions struct {
	Name      string
	Flavor    string
	Image     string
	PublicKey []byte
}
