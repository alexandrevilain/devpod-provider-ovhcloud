package ovhcloud

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alexandrevilain/devpod-provider-ovhcloud/pkg/options"
	"github.com/ovh/go-ovh/ovh"
	"golang.org/x/exp/slices"
)

// Client is an OVHCloud API client.
type Client struct {
	client *ovh.Client

	region      string
	serviceName string
}

// NewClient creates a new Client.
func NewClient(auth *options.Authentication) (*Client, error) {
	client, err := ovh.NewClient(
		auth.Endpoint,
		auth.AppKey,
		auth.AppSecret,
		auth.ConsumerKey,
	)
	if err != nil {
		return nil, fmt.Errorf("can't create ovhcloud client: %w", err)
	}

	return &Client{
		client:      client,
		region:      auth.Region,
		serviceName: auth.ServiceName,
	}, nil
}

func sshKeyName(name string) string {
	return fmt.Sprintf("devpod-%s", name)
}

// Init ensures the client can connect to the OVHCloud API.
func (c *Client) Init(ctx context.Context) error {
	_, err := c.listInstances(ctx)
	if err != nil {
		return fmt.Errorf("can't list instances: %w", err)
	}

	return nil
}

// GetInstanceByName retrieves an instance by name.
func (c *Client) GetInstanceByName(ctx context.Context, name string) (*Instance, error) {
	instances, err := c.listInstances(ctx)
	if err != nil {
		return nil, err
	}

	for _, instance := range instances {
		if instance.Name == name && instance.Region == c.region {
			return instance, nil
		}
	}

	return nil, fmt.Errorf("can't find instance '%s'", name)
}

// CreateInstance creates the instance using provided options.
func (c *Client) CreateInstance(ctx context.Context, opts CreateInstanceOptions) error {
	flavor, found, err := c.findFlavorByName(ctx, opts.Flavor)
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("can't find flavor named '%s'", opts.Flavor)
	}

	image, found, err := c.findImageByName(ctx, opts.Image)
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("can't find image named '%s'", opts.Image)
	}

	sskKey, err := c.getOrCreateSSHKey(ctx, sshKeyName(opts.Name), opts.PublicKey)
	if err != nil {
		return err
	}

	req := &CreateInstanceRequest{
		Name:           opts.Name,
		FlavorID:       flavor.ID,
		ImageID:        image.ID,
		MonthlyBilling: false,
		SSHKeyID:       sskKey.ID,
		Region:         c.region,
	}

	_, err = c.createInstance(ctx, req)
	return err
}

func (c *Client) StartInstance(ctx context.Context, name string) error {
	instance, err := c.GetInstanceByName(ctx, name)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/cloud/project/%s/instance/%s/start", c.serviceName, instance.ID)
	return c.client.PostWithContext(ctx, url, nil, nil)
}

func (c *Client) StopInstance(ctx context.Context, name string) error {
	instance, err := c.GetInstanceByName(ctx, name)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/cloud/project/%s/instance/%s/stop", c.serviceName, instance.ID)
	err = c.client.PostWithContext(ctx, url, nil, nil)
	if err != nil {
		return err
	}

	// wait until stopped
	for {
		instance, err := c.getInstanceByID(ctx, instance.ID)
		if err != nil {
			return err
		}

		fmt.Println(instance.Status)

		if instance.Status == InstanceStopped {
			break
		}
		// make sure we don't spam
		time.Sleep(time.Second)
	}

	return nil
}

func (c *Client) GetInstanceStatus(ctx context.Context, name string) (InstanceStatus, error) {
	instance, err := c.GetInstanceByName(ctx, name)
	if err != nil {
		return InstanceUnknown, err
	}

	return instance.Status, nil
}

func (c *Client) DeleteInstance(ctx context.Context, name string) error {
	instance, err := c.GetInstanceByName(ctx, name)
	if err != nil {
		return err
	}

	sshKey, _, err := c.findSSHKeyByName(ctx, sshKeyName(name))
	if err != nil {
		return err
	}

	return errors.Join(
		c.deleteInstance(ctx, instance.ID),
		c.deleteSSHKey(ctx, sshKey.ID),
	)
}

func (c *Client) getInstanceByID(ctx context.Context, id string) (*Instance, error) {
	url := fmt.Sprintf("/cloud/project/%s/instance/%s", c.serviceName, id)
	resp := &Instance{}
	err := c.client.GetWithContext(ctx, url, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) listInstances(ctx context.Context) ([]*Instance, error) {
	url := fmt.Sprintf("/cloud/project/%s/instance?region=%s", c.serviceName, c.region)
	var instances []*Instance
	err := c.client.GetWithContext(ctx, url, &instances)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (c *Client) createInstance(ctx context.Context, req *CreateInstanceRequest) (*Instance, error) {
	resp := &Instance{}
	err := c.client.PostWithContext(ctx, fmt.Sprintf("/cloud/project/%s/instance", c.serviceName), req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) deleteInstance(ctx context.Context, id string) error {
	url := fmt.Sprintf("/cloud/project/%s/instance/%s", c.serviceName, id)
	err := c.client.DeleteWithContext(ctx, url, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) findFlavorByName(ctx context.Context, name string) (*Flavor, bool, error) {
	url := fmt.Sprintf("/cloud/project/%s/flavor", c.serviceName)
	var flavors []*Flavor
	err := c.client.GetWithContext(ctx, url, &flavors)
	if err != nil {
		return nil, false, err
	}
	for _, flavor := range flavors {
		if flavor.Region == c.region && flavor.Name == name {
			return flavor, true, nil
		}
	}

	return nil, false, nil
}

func (c *Client) findImageByName(ctx context.Context, name string) (*Image, bool, error) {
	url := fmt.Sprintf("/cloud/project/%s/image", c.serviceName)
	var images []*Image
	err := c.client.GetWithContext(ctx, url, &images)
	if err != nil {
		return nil, false, err
	}
	for _, image := range images {
		if image.Region == c.region && image.Name == name {
			return image, true, nil
		}
	}

	return nil, false, nil
}

func (c *Client) findSSHKeyByName(ctx context.Context, name string) (*SSHKey, bool, error) {
	url := fmt.Sprintf("/cloud/project/%s/sshkey", c.serviceName)
	var keys []*SSHKey
	err := c.client.GetWithContext(ctx, url, &keys)
	if err != nil {
		return nil, false, err
	}

	for _, key := range keys {
		if slices.Contains(key.Regions, c.region) && key.Name == name {
			return key, true, nil
		}
	}

	return nil, false, nil
}

func (c *Client) getOrCreateSSHKey(ctx context.Context, name string, content []byte) (*SSHKey, error) {
	key, found, err := c.findSSHKeyByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if found {
		return key, nil
	}

	return c.createSSHKey(ctx, &CreateSSHKeyRequest{Name: name, PublicKey: string(content)})
}

func (c *Client) createSSHKey(ctx context.Context, req *CreateSSHKeyRequest) (*SSHKey, error) {
	fmt.Println(req.PublicKey)
	resp := &SSHKey{}
	err := c.client.PostWithContext(ctx, fmt.Sprintf("/cloud/project/%s/sshkey", c.serviceName), req, resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *Client) deleteSSHKey(ctx context.Context, id string) error {
	url := fmt.Sprintf("/cloud/project/%s/sshkey/%s", c.serviceName, id)
	err := c.client.DeleteWithContext(ctx, url, nil)
	if err != nil {
		return err
	}
	return nil
}
