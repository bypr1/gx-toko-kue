package form

import (
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type CustomerIdentityPhotoForm struct {
	Link     *string `json:"link" form:"link"`
	MimeType *string `json:"mimeType" form:"mimeType"`
}

type CustomerEmailForm struct {
	ID      *string `json:"id" form:"id" validate:"omitempty"`
	Email   string  `json:"email" form:"email"`
	Primary *bool   `json:"primary" form:"primary"`
	Deleted *bool   `json:"deleted" form:"deleted" validate:"omitempty,required_with=id,boolean"`
}

type CustomerPhoneForm struct {
	ID      *string `json:"id" form:"id" validate:"omitempty"`
	Code    string  `json:"code" form:"code"`
	Phone   string  `json:"phone" form:"phone"`
	Primary *bool   `json:"primary" form:"primary"`
	Deleted *bool   `json:"deleted" form:"deleted" validate:"omitempty,required_with=id,boolean"`
}

type CustomerForm struct {
	Request                 *http.Request             `json:"-" form:"-"`
	FullName                string                    `json:"fullName" form:"fullName"`
	ParentUUID              string                    `json:"parentUUID" form:"parentUUID"`
	GenderId                int                       `json:"genderId" form:"genderId"`
	TypeId                  int                       `json:"typeId" form:"typeId"`
	Emails                  []CustomerEmailForm       `json:"emails" form:"emails" validate:"omitempty,dive"`
	Phones                  []CustomerPhoneForm       `json:"phones" form:"phones" validate:"omitempty,dive"`
	Address                 string                    `json:"address" form:"address"`
	NationalityId           uint                      `json:"nationalityId" form:"nationalityId"`
	SocialStatusId          uint                      `json:"socialStatusId" form:"socialStatusId"`
	IdentityNo              string                    `json:"identityNo" form:"identityNo"`
	BirthDate               string                    `json:"birthDate" form:"birthDate"`
	IdentityPhoto           CustomerIdentityPhotoForm `json:"identityPhoto" form:"identityPhoto"`
	IdentityTypeId          int                       `json:"identityTypeId" form:"identityTypeId"`
	IdentityTypeName        string                    `json:"identityTypeName" form:"identityTypeName"`
	IdentityTypeDescription *string                   `json:"identityTypeDescription" form:"identityTypeDescription"`
	OtherInformation        string                    `json:"otherInformation" form:"otherInformation"`
	Company                 bool                      `json:"company" form:"company"`
	CompanyName             string                    `json:"companyName" form:"companyName"`
	CompanyAddress          string                    `json:"companyAddress" form:"companyAddress"`
	CompanyEmail            string                    `json:"companyEmail" form:"companyEmail"`
	Branches                []BranchForm              `json:"branches" form:"branches" validate:"omitempty,dive"`
	CreatedByUUID           string                    `json:"createdByUUID" form:"createdByUUID"`
	CreatedByName           string                    `json:"createdByName" form:"createdByName"`
}

type BranchForm struct {
	BranchId   uint   `json:"branchId" form:"branchId"`
	BranchName string `json:"branchName" form:"branchName"`
}

func (rule *CustomerForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(rule)
}

func (rule *CustomerForm) APIParse(r *http.Request) {
	rule.Request = r
	form := core.BaseForm{}
	form.FormParse(r, &rule)
}
