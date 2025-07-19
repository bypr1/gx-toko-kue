package form

import (
	"fmt"
	"net/http"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
	core "github.com/globalxtreme/go-core/v2/pkg"
)

type CustomerIdentityPhotoForm struct {
	Link     *string `json:"link"`
	MimeType *string `json:"mimeType" validate:"required"`
}

type CustomerEmailForm struct {
	ID      *string `json:"id" validate:"omitempty"`
	Email   string  `json:"email" validate:"required"`
	Primary *bool   `json:"primary" validate:"required"`
	Deleted *bool   `json:"deleted" validate:"omitempty,required_with=id,boolean"`
}

type CustomerPhoneForm struct {
	ID      *string `json:"id" validate:"omitempty"`
	Code    string  `json:"code" validate:"required"`
	Phone   string  `json:"phone" validate:"required"`
	Primary *bool   `json:"primary" validate:"required"`
	Deleted *bool   `json:"deleted" validate:"omitempty,required_with=id,boolean"`
}

type CustomerForm struct {
	Request                 *http.Request             `json:"-" form:"-"`
	FullName                string                    `json:"fullName" validate:"required"`
	ParentUUID              string                    `json:"parentUUID"`
	GenderId                int                       `json:"genderId"`
	TypeId                  int                       `json:"typeId"`
	Emails                  []CustomerEmailForm       `json:"emails" validate:"omitempty,dive"`
	Phones                  []CustomerPhoneForm       `json:"phones" validate:"omitempty,dive"`
	Address                 string                    `json:"address" validate:"required"`
	NationalityId           uint                      `json:"nationalityId" validate:"required"`
	SocialStatusId          uint                      `json:"socialStatusId"`
	IdentityNo              string                    `json:"identityNo"`
	BirthDate               string                    `json:"birthDate"`
	IdentityPhoto           CustomerIdentityPhotoForm `json:"identityPhoto" validate:"required"`
	IdentityTypeId          int                       `json:"identityTypeId"`
	IdentityTypeName        string                    `json:"identityTypeName"`
	IdentityTypeDescription *string                   `json:"identityTypeDescription"`
	OtherInformation        string                    `json:"otherInformation"`
	Company                 bool                      `json:"primary"`
	CompanyName             string                    `json:"companyName"`
	CompanyAddress          string                    `json:"companyAddress"`
	CompanyEmail            string                    `json:"companyEmail"`
	// Attachments             *[]CustomerAttachmentForm `json:"attachments" validate:"omitempty,dive"`
	Branches      []BranchForm `json:"branches" validate:"required,dive"`
	CreatedByUUID string       `json:"createdByUUID"`
	CreatedByName string       `json:"createdByName"`
}

type BranchForm struct {
	BranchId   uint   `json:"branchId" validate:"required"`
	BranchName string `json:"branchName" validate:"required"`
}

func (rule *CustomerForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(rule)
}

func (rule *CustomerForm) APIParse(r *http.Request) {
	formValue := r.MultipartForm.Value

	rule.Request = r
	rule.FullName = formValue["fullName"][0]
	rule.SocialStatusId = uint(core.ToInt(formValue["socialStatusId"][0]))
	rule.GenderId = core.ToInt(formValue["genderId"][0])
	rule.TypeId = core.ToInt(formValue["typeId"][0])
	rule.Address = formValue["address"][0]
	rule.IdentityNo = formValue["identityNo"][0]
	rule.IdentityTypeId = core.ToInt(formValue["identityTypeId"][0])
	rule.IdentityTypeName = formValue["identityTypeName"][0]
	rule.OtherInformation = formValue["otherInformation"][0]
	// rule.Attachments = SetCustomerAttachments(formValue)
	rule.NationalityId = uint(core.ToInt(formValue["nationalityId"][0]))
	rule.ParentUUID = formValue["parentUUID"][0]
	rule.Company = core.ToBool(formValue["company"][0])

	if values, ok := formValue["identityTypeDescription"]; ok && len(values) > 0 {
		rule.IdentityTypeDescription = &values[0]
	}

	if mimeType, exists := formValue["identityPhoto[mimeType]"]; exists && len(mimeType) > 0 {
		rule.IdentityPhoto.MimeType = &mimeType[0]
	}

	if link, exists := formValue["identityPhoto[link]"]; exists && len(link) > 0 {
		rule.IdentityPhoto.Link = &link[0]
	}

	if rule.Company {
		rule.CompanyName = formValue["companyName"][0]
		rule.CompanyAddress = formValue["companyAddress"][0]
		rule.CompanyEmail = formValue["companyEmail"][0]
	}

	emailLoop := true
	emails := make([]CustomerEmailForm, 0)

	for a := 0; emailLoop; a++ {
		if emailName, ok := formValue[fmt.Sprintf("emails[%d][email]", a)]; ok {
			primary := core.ToBool(formValue[fmt.Sprintf("emails[%d][primary]", a)][0])
			email := CustomerEmailForm{
				Email:   emailName[0],
				Primary: &primary,
			}

			if id, emailOk := formValue[fmt.Sprintf("emails[%d][id]", a)]; emailOk {
				email.ID = &id[0]

				deleted := core.ToBool(formValue[fmt.Sprintf("emails[%d][deleted]", a)][0])
				email.Deleted = &deleted
			}
			emails = append(emails, email)
		} else {
			emailLoop = false
		}
	}
	rule.Emails = emails

	phoneLoop := true
	phones := make([]CustomerPhoneForm, 0)

	for a := 0; phoneLoop; a++ {
		if phoneNumber, ok := formValue[fmt.Sprintf("phones[%d][phone]", a)]; ok {
			primary := core.ToBool(formValue[fmt.Sprintf("phones[%d][primary]", a)][0])
			phone := CustomerPhoneForm{
				Code:    formValue[fmt.Sprintf("phones[%d][code]", a)][0],
				Phone:   phoneNumber[0],
				Primary: &primary,
			}

			if id, emailOk := formValue[fmt.Sprintf("phones[%d][id]", a)]; emailOk {
				phone.ID = &id[0]

				deleted := core.ToBool(formValue[fmt.Sprintf("phones[%d][deleted]", a)][0])
				phone.Deleted = &deleted
			}
			phones = append(phones, phone)
		} else {
			phoneLoop = false
		}
	}
	rule.Phones = phones

	branchLoop := true
	branches := make([]BranchForm, 0)
	for a := 0; branchLoop; a++ {
		if branchId, ok := formValue[fmt.Sprintf("branches[%d][branchId]", a)]; ok {
			branch := BranchForm{
				BranchId:   uint(core.ToInt(branchId[0])),
				BranchName: formValue[fmt.Sprintf("branches[%d][branchName]", a)][0],
			}
			branches = append(branches, branch)
		} else {
			branchLoop = false
		}
	}

	rule.Branches = branches

}
