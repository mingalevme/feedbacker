package feedback
//
//type ContextMapper interface {
//	ContextToData(c Context) *ContextData
//	DataToContext(c *ContextData) (Context, error)
//}
//
//type contextMapper struct{}
//
//func (s contextMapper) ContextToData(c Context) *ContextData {
//	if c == nil {
//		return nil
//	}
//	return &ContextData{
//		AppVersion:  c.GetAppVersion(),
//		AppBuild:    c.GetAppBuildNumber(),
//		OsName:      c.GetOsName(),
//		OsVersion:   c.GetOsVersion(),
//		DeviceBrand: c.GetDeviceBrand(),
//		DeviceModel: c.GetDeviceModel(),
//	}
//}
//
//func (s contextMapper) DataToContext(c *ContextData) (Context, error) {
//	return NewContext(c.AppVersion, c.AppBuild, c.OsName, c.OsVersion, c.DeviceBrand, c.DeviceModel)
//}
//
//var contextMapperInstance *contextMapper
//
//func GetContextMapper() ContextMapper {
//	if contextMapperInstance == nil {
//		contextMapperInstance = &contextMapper{}
//	}
//	return contextMapperInstance
//}
