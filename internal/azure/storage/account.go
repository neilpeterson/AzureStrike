package storage

import (
	"time"
)

// Account represents a mocked Azure Storage Account
type Account struct {
	Name                  string            `yaml:"name" json:"name"`
	ResourceGroup         string            `yaml:"resource_group" json:"resourceGroup"`
	Location              string            `yaml:"location" json:"location"`
	Kind                  string            `yaml:"kind" json:"kind"`
	SKU                   string            `yaml:"sku" json:"sku"`
	AccessTier            string            `yaml:"access_tier" json:"accessTier"`
	AllowBlobPublicAccess bool              `yaml:"public_access" json:"allowBlobPublicAccess"`
	EnableHTTPSOnly       bool              `yaml:"https_only" json:"enableHttpsTrafficOnly"`
	MinTLSVersion         string            `yaml:"min_tls_version" json:"minimumTlsVersion"`
	CreationTime          string            `yaml:"creation_time" json:"creationTime"`
	ProvisioningState     string            `yaml:"provisioning_state" json:"provisioningState"`
	PrimaryLocation       string            `yaml:"primary_location" json:"primaryLocation"`
	StatusOfPrimary       string            `yaml:"status_of_primary" json:"statusOfPrimary"`
	Containers            []Container       `yaml:"containers" json:"-"`
	ConnectionStrings     map[string]string `yaml:"-" json:"-"`
	Tags                  map[string]string `yaml:"tags" json:"tags"`
}

// Container represents a blob container within a storage account
type Container struct {
	Name         string `yaml:"name" json:"name"`
	PublicAccess string `yaml:"public_access" json:"publicAccess"` // none, blob, container
	Blobs        []Blob `yaml:"blobs" json:"-"`
}

// Blob represents a blob within a container
type Blob struct {
	Name         string            `yaml:"name" json:"name"`
	Content      string            `yaml:"content" json:"-"`
	ContentType  string            `yaml:"content_type" json:"contentType"`
	Size         int64             `yaml:"size" json:"contentLength"`
	LastModified time.Time         `yaml:"-" json:"lastModified"`
	Metadata     map[string]string `yaml:"metadata" json:"metadata"`
}

// SASToken represents a Shared Access Signature token
type SASToken struct {
	Token         string    `yaml:"token" json:"-"`
	Permissions   string    `yaml:"permissions" json:"-"` // r, w, d, l, a, c
	Expiry        time.Time `yaml:"expiry" json:"-"`
	StartTime     time.Time `yaml:"start_time" json:"-"`
	Services      string    `yaml:"services" json:"-"`       // b, f, q, t
	ResourceTypes string    `yaml:"resource_types" json:"-"` // s, c, o
}

// Environment represents the mocked Azure storage environment
type Environment struct {
	Accounts map[string]*Account
}

// NewEnvironment creates a new storage environment
func NewEnvironment() *Environment {
	return &Environment{
		Accounts: make(map[string]*Account),
	}
}

// AddAccount adds a storage account to the environment
func (e *Environment) AddAccount(account *Account) {
	e.Accounts[account.Name] = account
}

// GetAccount retrieves a storage account by name
func (e *Environment) GetAccount(name string) *Account {
	return e.Accounts[name]
}

// ListAccounts returns all storage accounts
func (e *Environment) ListAccounts() []*Account {
	accounts := make([]*Account, 0, len(e.Accounts))
	for _, acc := range e.Accounts {
		accounts = append(accounts, acc)
	}
	return accounts
}

// GetContainer retrieves a container from an account
func (e *Environment) GetContainer(accountName, containerName string) *Container {
	acc := e.GetAccount(accountName)
	if acc == nil {
		return nil
	}
	for i := range acc.Containers {
		if acc.Containers[i].Name == containerName {
			return &acc.Containers[i]
		}
	}
	return nil
}

// GetBlob retrieves a blob from a container
func (e *Environment) GetBlob(accountName, containerName, blobName string) *Blob {
	container := e.GetContainer(accountName, containerName)
	if container == nil {
		return nil
	}
	for i := range container.Blobs {
		if container.Blobs[i].Name == blobName {
			return &container.Blobs[i]
		}
	}
	return nil
}

// IsPubliclyAccessible checks if a container allows public access
func (c *Container) IsPubliclyAccessible() bool {
	return c.PublicAccess == "blob" || c.PublicAccess == "container"
}
