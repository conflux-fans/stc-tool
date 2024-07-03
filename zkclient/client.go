package zkclient

import (
	"context"

	"github.com/mcuadros/go-defaults"
	"github.com/openweb3/go-rpc-provider/interfaces"
	pproviders "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/openweb3/web3go"
	"github.com/openweb3/web3go/providers"
	"github.com/openweb3/web3go/signers"
)

type Client struct {
	*pproviders.MiddlewarableProvider
	// context context.Context
	option *web3go.ClientOption
}

func NewClient(rawurl string) (*Client, error) {
	return NewClientWithOption(rawurl, web3go.ClientOption{})
}

func MustNewClient(rawurl string) *Client {
	c, err := NewClient(rawurl)
	if err != nil {
		panic(err)
	}
	return c
}

func NewClientWithOption(rawurl string, option web3go.ClientOption) (*Client, error) {

	defaults.SetDefaults(&option.Option)

	p, err := pproviders.NewProviderWithOption(rawurl, option.Option)
	if err != nil {
		return nil, err
	}

	if option.SignerManager != nil {
		p = providers.NewSignableProvider(p, option.SignerManager)
	}

	ec := NewClientWithProvider(p)
	ec.option = &option

	return ec, nil
}

func MustNewClientWithOption(rawurl string, option web3go.ClientOption) *Client {
	c, err := NewClientWithOption(rawurl, option)
	if err != nil {
		panic(err)
	}
	return c
}

func NewClientWithProvider(p interfaces.Provider) *Client {
	c := &Client{}
	c.SetProvider(p)
	return c
}

func (c *Client) SetProvider(p interfaces.Provider) {
	if _, ok := p.(*pproviders.MiddlewarableProvider); !ok {
		p = pproviders.NewMiddlewarableProvider(p)
	}

	c.MiddlewarableProvider = p.(*pproviders.MiddlewarableProvider)
}

func (c *Client) Provider() *pproviders.MiddlewarableProvider {
	return c.MiddlewarableProvider
}

// GetSignerManager returns signer manager if exist in option, otherwise return error
func (c *Client) GetSignerManager() (*signers.SignerManager, error) {
	if c.option.SignerManager != nil {
		return c.option.SignerManager, nil
	}
	return nil, web3go.ErrNotFound
}

// func (c *Client) GetProof2(
// 	encoded_vc []byte,
// 	birth_date_threshold uint64,
// 	path_elements []common.Hash,
// 	path_indices []uint,
// ) (val *Proof, err error) {
// 	err = c.CallContext(context.Background(), &val, "zg_generateZkProof",
// 		hex.EncodeToString(encoded_vc),
// 		birth_date_threshold,
// 		path_elements,
// 		path_indices)
// 	return
// }

func (c *Client) GetProof(input *ProveInput) (proof string, err error) {
	err = c.CallContext(context.Background(), &proof, "zg_generateZkProof", input)
	return
}

func (c *Client) Verify(proof string, publicInputs VerifyInput) (val bool, err error) {
	err = c.CallContext(context.Background(), &val, "zg_verifyZkProof", proof, publicInputs)
	return
}
