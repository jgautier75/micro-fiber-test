package converters

import (
	"micro-fiber-test/pkg/dto/orgs"
	"micro-fiber-test/pkg/model"
)

func ConvertOrgReqToDaoModel(defaultTenantId int64, orgReq orgs.CreateOrgRequest) model.Organization {
	org := model.Organization{}
	org.TenantId = defaultTenantId
	if orgReq.Label != nil {
		org.Label = *orgReq.Label
	}
	if orgReq.Kind != nil {
		org.Type = model.OrganizationType(*orgReq.Kind)
	}
	org.Status = model.OrganizationStatus(orgReq.Status)
	return org
}

func ConvertOrgModelToOrgResp(org model.Organization) orgs.OrganizationResponse {
	return orgs.OrganizationResponse{
		Code:   org.Code,
		Label:  org.Label,
		Status: int(org.Status),
		Kind:   string(org.Type),
	}
}
