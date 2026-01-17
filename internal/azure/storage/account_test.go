package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvironment(t *testing.T) {
	env := NewEnvironment()

	assert.NotNil(t, env)
	assert.NotNil(t, env.Accounts)
	assert.Empty(t, env.Accounts)
}

func TestAddAccount(t *testing.T) {
	env := NewEnvironment()
	account := &Account{
		Name:          "teststorage",
		ResourceGroup: "test-rg",
		Location:      "eastus",
	}

	env.AddAccount(account)

	assert.Len(t, env.Accounts, 1)
	assert.Contains(t, env.Accounts, "teststorage")
	assert.Equal(t, account, env.Accounts["teststorage"])
}

func TestAddMultipleAccounts(t *testing.T) {
	env := NewEnvironment()

	env.AddAccount(&Account{Name: "account1"})
	env.AddAccount(&Account{Name: "account2"})
	env.AddAccount(&Account{Name: "account3"})

	assert.Len(t, env.Accounts, 3)
}

func TestGetAccount(t *testing.T) {
	env := NewEnvironment()
	account := &Account{
		Name:     "myaccount",
		Location: "westus",
	}
	env.AddAccount(account)

	t.Run("existing account", func(t *testing.T) {
		result := env.GetAccount("myaccount")
		assert.NotNil(t, result)
		assert.Equal(t, "myaccount", result.Name)
		assert.Equal(t, "westus", result.Location)
	})

	t.Run("non-existing account", func(t *testing.T) {
		result := env.GetAccount("nonexistent")
		assert.Nil(t, result)
	})
}

func TestListAccounts(t *testing.T) {
	env := NewEnvironment()
	env.AddAccount(&Account{Name: "account1"})
	env.AddAccount(&Account{Name: "account2"})

	accounts := env.ListAccounts()

	assert.Len(t, accounts, 2)
}

func TestListAccountsEmpty(t *testing.T) {
	env := NewEnvironment()

	accounts := env.ListAccounts()

	assert.Empty(t, accounts)
}

func TestGetContainer(t *testing.T) {
	env := NewEnvironment()
	account := &Account{
		Name: "storageaccount",
		Containers: []Container{
			{Name: "container1", PublicAccess: "none"},
			{Name: "container2", PublicAccess: "blob"},
		},
	}
	env.AddAccount(account)

	t.Run("existing container", func(t *testing.T) {
		container := env.GetContainer("storageaccount", "container1")
		assert.NotNil(t, container)
		assert.Equal(t, "container1", container.Name)
		assert.Equal(t, "none", container.PublicAccess)
	})

	t.Run("non-existing container", func(t *testing.T) {
		container := env.GetContainer("storageaccount", "nonexistent")
		assert.Nil(t, container)
	})

	t.Run("non-existing account", func(t *testing.T) {
		container := env.GetContainer("wrongaccount", "container1")
		assert.Nil(t, container)
	})
}

func TestGetBlob(t *testing.T) {
	env := NewEnvironment()
	account := &Account{
		Name: "blobaccount",
		Containers: []Container{
			{
				Name: "data",
				Blobs: []Blob{
					{Name: "file1.txt", Content: "content1"},
					{Name: "file2.txt", Content: "content2"},
				},
			},
		},
	}
	env.AddAccount(account)

	t.Run("existing blob", func(t *testing.T) {
		blob := env.GetBlob("blobaccount", "data", "file1.txt")
		assert.NotNil(t, blob)
		assert.Equal(t, "file1.txt", blob.Name)
		assert.Equal(t, "content1", blob.Content)
	})

	t.Run("second blob", func(t *testing.T) {
		blob := env.GetBlob("blobaccount", "data", "file2.txt")
		assert.NotNil(t, blob)
		assert.Equal(t, "content2", blob.Content)
	})

	t.Run("non-existing blob", func(t *testing.T) {
		blob := env.GetBlob("blobaccount", "data", "nonexistent.txt")
		assert.Nil(t, blob)
	})

	t.Run("non-existing container", func(t *testing.T) {
		blob := env.GetBlob("blobaccount", "wrongcontainer", "file1.txt")
		assert.Nil(t, blob)
	})

	t.Run("non-existing account", func(t *testing.T) {
		blob := env.GetBlob("wrongaccount", "data", "file1.txt")
		assert.Nil(t, blob)
	})
}

func TestContainerIsPubliclyAccessible(t *testing.T) {
	tests := []struct {
		name         string
		publicAccess string
		expected     bool
	}{
		{"none access", "none", false},
		{"private access", "private", false},
		{"empty access", "", false},
		{"blob access", "blob", true},
		{"container access", "container", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := &Container{
				Name:         "test",
				PublicAccess: tt.publicAccess,
			}
			assert.Equal(t, tt.expected, container.IsPubliclyAccessible())
		})
	}
}

func TestAccountFields(t *testing.T) {
	account := &Account{
		Name:                  "fullaccount",
		ResourceGroup:         "my-rg",
		Location:              "eastus2",
		Kind:                  "StorageV2",
		SKU:                   "Standard_LRS",
		AccessTier:            "Hot",
		AllowBlobPublicAccess: true,
		EnableHTTPSOnly:       true,
		MinTLSVersion:         "TLS1_2",
		CreationTime:          "2024-01-01T00:00:00Z",
		ProvisioningState:     "Succeeded",
		PrimaryLocation:       "eastus2",
		StatusOfPrimary:       "available",
		Tags:                  map[string]string{"env": "test"},
	}

	assert.Equal(t, "fullaccount", account.Name)
	assert.Equal(t, "my-rg", account.ResourceGroup)
	assert.Equal(t, "eastus2", account.Location)
	assert.Equal(t, "StorageV2", account.Kind)
	assert.True(t, account.AllowBlobPublicAccess)
	assert.True(t, account.EnableHTTPSOnly)
	assert.Equal(t, "TLS1_2", account.MinTLSVersion)
	assert.Equal(t, "test", account.Tags["env"])
}

func TestBlobFields(t *testing.T) {
	blob := &Blob{
		Name:        "document.pdf",
		Content:     "PDF content here",
		ContentType: "application/pdf",
		Size:        1024,
		Metadata:    map[string]string{"author": "test"},
	}

	assert.Equal(t, "document.pdf", blob.Name)
	assert.Equal(t, "PDF content here", blob.Content)
	assert.Equal(t, "application/pdf", blob.ContentType)
	assert.Equal(t, int64(1024), blob.Size)
	assert.Equal(t, "test", blob.Metadata["author"])
}
