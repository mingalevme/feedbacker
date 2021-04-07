package feedback
//
//type CustomerMapper interface {
//	CustomerToData(c Customer) *CustomerData
//	DataToCustomer(c *CustomerData) (Customer, error)
//}
//
//type customerMapper struct{}
//
//func (s customerMapper) CustomerToData(c Customer) *CustomerData {
//	if c == nil {
//		return nil
//	}
//	return &CustomerData{
//		Email:          c.GetEmail(),
//		InstallationID: c.GetInstallationID(),
//	}
//}
//
//func (s customerMapper) DataToCustomer(c *CustomerData) (Customer, error) {
//	if c == nil {
//		return nil, nil
//	}
//	return NewCustomer(c.Email, c.InstallationID)
//}
//
//var customerMapperInstance *customerMapper
//
//func GetCustomerMapper() CustomerMapper {
//	if customerMapperInstance == nil {
//		customerMapperInstance = &customerMapper{}
//	}
//	return customerMapperInstance
//}
