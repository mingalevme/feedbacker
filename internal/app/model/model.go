// https://google.github.io/styleguide/jsoncstyleguide.xml

package model

import (
	"github.com/mingalevme/feedbacker/pkg/util"
	"time"
)

type Context struct {
	AppVersion  *string `json:"appVersion,omitempty" xml:"AppVersion" param:"app_version" query:"app_version" form:"app_version"`
	AppBuild    *string `json:"appBuild,omitempty" xml:"AppBuild" param:"app_build" query:"app_build" form:"app_build"`
	OsName      *string `json:"osName,omitempty" xml:"OsName" param:"os_name" query:"os_name" form:"os_name"`
	OsVersion   *string `json:"osVersion,omitempty" xml:"OsVersion" param:"os_version" query:"os_version" form:"os_version"`
	DeviceBrand *string `json:"deviceBrand,omitempty" xml:"DeviceBrand" param:"device_brand" query:"device_brand" form:"device_brand"`
	DeviceModel *string `json:"deviceModel,omitempty" xml:"DeviceModel" param:"device_model" query:"device_model" form:"device_model"`
}

type Customer struct {
	Email          *string `json:"email,omitempty" xml:"Email" param:"email" query:"email" form:"email"`
	InstallationID *string `json:"installationId,omitempty" xml:"InstallationId" param:"installation_id" query:"installation_id" form:"installation_id"`
}

func (s Customer) isEmpty() bool {
	return util.IsPointerToEmptyString(s.Email) && util.IsPointerToEmptyString(s.InstallationID)
}

type Feedback struct {
	ID        int       `json:"id" xml:"Id" param:"id" query:"id" form:"id"`
	Service   string    `json:"service" xml:"Service" param:"service" query:"service" form:"service"`
	Edition   *string    `json:"edition" xml:"Edition" param:"edition" query:"edition" form:"edition"`
	Text      string    `json:"text" xml:"Text" param:"text" query:"text" form:"text"`
	Context   *Context   `json:"context,omitempty" xml:"Context" param:"context" query:"context" form:"context"`
	Customer  *Customer `json:"customer,omitempty" xml:"Customer" param:"customer" query:"customer" form:"customer"`
	CreatedAt time.Time `json:"createdAt,omitempty" xml:"CreatedAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" xml:"UpdatedAt"`
}
