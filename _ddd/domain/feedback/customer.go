// https://google.github.io/styleguide/jsoncstyleguide.xml

package feedback

type CustomerData struct {
	Email          *string `json:"email,omitempty"`
	InstallationID *string `json:"installationId,omitempty"`
}

type Customer interface {
	GetEmail() *string
	GetInstallationID() *string
	IsEmpty() bool
}

type customer struct {
	CustomerData
}

func NewCustomer(email *string, installationID *string) (Customer, error) {
	return &customer{
		CustomerData{
			Email:          email,
			InstallationID: installationID,
		},
	}, nil
}
//
//func NewCustomerFromData(data CustomerData) (Customer, error) {
//	return NewCustomer(data.Email, data.InstallationID)
//}

func (s *customer) GetEmail() *string {
	return s.Email
}

func (s *customer) GetInstallationID() *string {
	return s.InstallationID
}

func (s *customer) IsEmpty() bool {
	return s.InstallationID == nil && s.Email == nil
}
