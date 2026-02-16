package backend

// Backend is the interface for secret storage providers.
type Backend interface {
	// ListItems returns all item titles from the backend.
	ListItems() ([]string, error)
	// GetSecret fetches the secret value for the given item title.
	GetSecret(title string) (string, error)
}
