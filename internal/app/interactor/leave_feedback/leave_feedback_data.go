package leave_feedback

import (
	"github.com/mingalevme/feedbacker/pkg/util"
	"github.com/pkg/errors"
)

type LeaveFeedbackData struct {
	App            string  `json:"app" xml:"app" param:"app" query:"app" form:"app"`
	AppVersion     *string `json:"aversion" xml:"aversion" param:"aversion" query:"aversion" form:"aversion"`
	AppBuildNumber *string `json:"bversion" xml:"bversion" param:"bversion" query:"bversion" form:"bversion"`
	Edition        *string `json:"edition" xml:"edition" param:"edition" query:"edition" form:"edition"`
	Body           string  `json:"body" xml:"body" param:"body" query:"body" form:"body"`
	Brand          *string `json:"brand" xml:"brand" param:"brand" query:"brand" form:"brand"`
	Model          *string `json:"model,omitempty" xml:"model" param:"model" query:"model" form:"model"`
	OsName         *string `json:"osname" xml:"osname" param:"osname" query:"osname" form:"osname"`
	OsVersion      *string `json:"osversion" xml:"osversion" param:"osversion" query:"osversion" form:"osversion"`
	Email          *string `json:"email" xml:"email" param:"email" query:"email" form:"email"`
	InstallationID *string `json:"installation_id" xml:"installation_id" param:"installation_id" query:"installation_id" form:"installation_id"`
}

func (s LeaveFeedbackData) Validate() error {
	if util.IsEmptyString(s.App) {
		return errors.Wrap(ErrUnprocessableEntity, "app is required")
	}
	if util.IsEmptyString(s.AppVersion) {
		s.AppVersion = nil
	}
	if util.IsEmptyString(s.AppBuildNumber) {
		s.AppBuildNumber = nil
	}
	if util.IsEmptyString(s.Edition) {
		return errors.Wrap(ErrUnprocessableEntity, "edition is required")
	}
	if util.IsEmptyString(s.Body) {
		return errors.Wrap(ErrUnprocessableEntity, "text (body) is required")
	}
	if util.IsEmptyString(s.Body) {
		return errors.Wrap(ErrUnprocessableEntity, "text (body) is required")
	}
	if util.IsEmptyString(s.Brand) {
		s.Brand = nil
	}
	if util.IsEmptyString(s.Model) {
		s.Model = nil
	}
	if util.IsEmptyString(s.OsName) {
		s.OsName = nil
	}
	if util.IsEmptyString(s.OsVersion) {
		s.OsVersion = nil
	}
	if util.IsEmptyString(s.Email) {
		s.Email = nil
	}
	if util.IsEmptyString(s.InstallationID) {
		s.InstallationID = nil
	}
	return nil
}
