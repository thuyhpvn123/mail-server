package types

type Config interface {
	DnsLink() string
	Version() string
	NodeType() string
	PrivateKey() []byte
	ConnectionAddress() string
}
