package ovhcloud

import (
	"errors"
	"time"
)

// FlavorCapability holds capability for a flavor.
type FlavorCapability struct {
	// Name of the capability
	// Allowed: failoverip┃resize┃snapshot┃volume
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// Flavor holds compute flavor properties.
type Flavor struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
	RAM    int    `json:"ram"`
	Disk   int    `json:"disk"`
	Vcpus  int    `json:"vcpus"`
	Type   string `json:"type"`
	OsType string `json:"osType"`
	// InboundBandwidth is the max capacity of inbound traffic in Mbit/s.
	InboundBandwidth int `json:"inboundBandwidth"`
	// OutboundBandwidth is the max capacity of outbound traffic in Mbit/s.
	OutboundBandwidth int                `json:"outboundBandwidth"`
	Available         bool               `json:"available"`
	Capabilities      []FlavorCapability `json:"capabilities"`
	// Quota is the number instance you can spawn with your actual quota.
	Quota int `json:"quota"`
}

// Image holds compute image properties.
type Image struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Region       string    `json:"region"`
	Visibility   string    `json:"visibility"`
	Type         string    `json:"type"`
	MinDisk      int       `json:"minDisk"`
	MinRAM       int       `json:"minRam"`
	Size         float64   `json:"size"`
	CreationDate time.Time `json:"creationDate"`
	Status       string    `json:"status"`
	// User is the user to connect with.
	User       string   `json:"user"`
	FlavorType string   `json:"flavorType"`
	Tags       []string `json:"tags"`
	PlanCode   string   `json:"planCode"`
}

// IP holds ip adress properties.
type IP struct {
	IP        string `json:"ip"`
	Type      string `json:"type"`
	Version   int    `json:"version"`
	NetworkID string `json:"networkId"`
	GatewayIP string `json:"gatewayIp"`
}

// Instance holds compute instance properties.
type Instance struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	IPAddresses  []IP           `json:"ipAddresses"`
	FlavorID     string         `json:"flavorId"`
	ImageID      string         `json:"imageId"`
	SSHKeyID     string         `json:"sshKeyId"`
	Created      time.Time      `json:"created"`
	Region       string         `json:"region"`
	Status       InstanceStatus `json:"status"`
	OperationIds []string       `json:"operationIds"`
}

// InstanceStatus defines all possibles instance status reported by API.
type InstanceStatus string

const (
	InstanceActive           InstanceStatus = "ACTIVE"
	InstanceBuild            InstanceStatus = "BUILD"
	InstanceBuilding         InstanceStatus = "BUILDING"
	InstanceDeleted          InstanceStatus = "DELETED"
	InstanceDeleting         InstanceStatus = "DELETING"
	InstanceError            InstanceStatus = "ERROR"
	InstanceHardReboot       InstanceStatus = "HARD_REBOOT"
	InstanceMigrating        InstanceStatus = "MIGRATING"
	InstancePassword         InstanceStatus = "PASSWORD"
	InstancePasued           InstanceStatus = "PAUSED"
	InstanceReboot           InstanceStatus = "REBOOT"
	InstanceRebuild          InstanceStatus = "REBUILD"
	InstanceRescue           InstanceStatus = "RESCUE"
	InstanceRescued          InstanceStatus = "RESCUED"
	InstanceRescuing         InstanceStatus = "RESCUING"
	InstanceResize           InstanceStatus = "RESIZE"
	InstanceResized          InstanceStatus = "RESIZED"
	InstanceResuming         InstanceStatus = "RESUMING"
	InstanceRevertResize     InstanceStatus = "REVERT_RESIZE"
	InstanceShelved          InstanceStatus = "SHELVED"
	InstanceShelvedOffLoaded InstanceStatus = "SHELVED_OFFLOADED"
	InstanceSelving          InstanceStatus = "SHELVING"
	InstanceShutOff          InstanceStatus = "SHUTOFF"
	InstanceSnapshotting     InstanceStatus = "SNAPSHOTTING"
	InstanceSoftDeleted      InstanceStatus = "SOFT_DELETED"
	InstanceStopped          InstanceStatus = "STOPPED"
	InstanceSuspended        InstanceStatus = "SUSPENDED"
	InstanceUnknown          InstanceStatus = "UNKNOWN"
	InstanceUnrescuing       InstanceStatus = "UNRESCUING"
	InstanceUnshelving       InstanceStatus = "UNSHELVING"
	InstanceVerifyResize     InstanceStatus = "VERIFY_RESIZE"
)

// CreateInstanceRequest is the payload to create an instance.
type CreateInstanceRequest struct {
	FlavorID       string `json:"flavorId"`
	ImageID        string `json:"imageId"`
	MonthlyBilling bool   `json:"monthlyBilling"`
	Name           string `json:"name"`
	Region         string `json:"region"`
	SSHKeyID       string `json:"sshKeyId"`
	UserData       string `json:"userData"`
}

// CreateSSHKeyRequest is the payload to create an ssh key.
type CreateSSHKeyRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

// SSHKey holds ssh key properties.
type SSHKey struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Regions     []string `json:"regions"`
	FingerPrint string   `json:"fingerPrint"`
	PublicKey   string   `json:"publicKey"`
}

// PublicIP returns provided instance v4 public IP.
func PublicIP(instance *Instance) (string, error) {
	for _, ip := range instance.IPAddresses {
		if ip.Version == 4 {
			return ip.IP, nil
		}
	}
	return "", errors.New("can't find public adress")
}
