// https://google.github.io/styleguide/jsoncstyleguide.xml

package main

// import (
// 	"database/sql"
// 	"time"
// )

type Context struct {
	Edition        string `json:"edition,omitempty"`
	AppVersion     string `json:"appVersion,omitempty"`
	AppBuild       string `json:"appBuild,omitempty"`
	OsName         string `json:"osName,omitempty"`
	OsVersion      string `json:"osVersion,omitempty"`
	DeviceBrand    string `json:"deviceBrand,omitempty"`
	DeviceModel    string `json:"deviceModel,omitempty"`
	InstallationId string `json:"installationId,omitempty"`
}

type Feedback struct {
	Id        int      `json:"id,omitempty"`
	Service   string   `json:"service,omitempty"`
	Text      string   `json:"text,omitempty"`
	Email     *string  `json:"email,omitempty"`
	Context   *Context `json:"context,omitempty"`
	CreatedAt int      `json:"created_at,omitempty"`
	UpdatedAt int      `json:"updated_at,omitempty"`
}
