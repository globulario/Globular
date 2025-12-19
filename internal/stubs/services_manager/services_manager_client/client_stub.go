package services_manager_client

// Services_Manager_Client is a lightweight stub for builds lacking the real service manager.
type Services_Manager_Client struct{}

// NewServicesManagerService_Client returns a stub client that reports no actions.
func NewServicesManagerService_Client(_ string, _ string) (*Services_Manager_Client, error) {
	return &Services_Manager_Client{}, nil
}

// GetAllActions returns a placeholder empty action set.
func (c *Services_Manager_Client) GetAllActions() ([]string, error) {
	return []string{}, nil
}
