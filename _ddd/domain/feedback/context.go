// https://google.github.io/styleguide/jsoncstyleguide.xml

package feedback

type Context struct {
	AppVersion  *string `json:"appVersion,omitempty"`
	AppBuild    *string `json:"appBuild,omitempty"`
	OsName      *string `json:"osName,omitempty"`
	OsVersion   *string `json:"osVersion,omitempty"`
	DeviceBrand *string `json:"deviceBrand,omitempty"`
	DeviceModel *string `json:"deviceModel,omitempty"`
}

//type Context interface {
//	GetAppVersion() *string
//	GetAppBuildNumber() *string
//	GetOsName() *string
//	GetOsVersion() *string
//	GetDeviceBrand() *string
//	GetDeviceModel() *string
//}

//func NewContext(appVersion *string, appBuild *string, osName *string, osVersion *string, deviceBrand *string, deviceModel *string) (Context, error) {
//	// validate
//	return &context{
//		ContextData{
//			AppVersion:  appVersion,
//			AppBuild:    appBuild,
//			OsName:      osName,
//			OsVersion:   osVersion,
//			DeviceBrand: deviceBrand,
//			DeviceModel: deviceModel,
//		},
//	}, nil
//}
//
//type context struct {
//	ContextData
//}
//
//func (s *context) GetAppVersion() *string {
//	return s.AppVersion
//}
//
//func (s *context) GetAppBuildNumber() *string {
//	return s.AppBuild
//}
//
//func (s *context) GetOsName() *string {
//	return s.OsName
//}
//
//func (s *context) GetOsVersion() *string {
//	return s.OsVersion
//}
//
//func (s *context) GetDeviceBrand() *string {
//	return s.DeviceBrand
//}
//
//func (s *context) GetDeviceModel() *string {
//	return s.DeviceModel
//}
