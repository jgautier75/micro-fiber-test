package converters

import (
	"micro-fiber-test/pkg/dto/orgs"
	"micro-fiber-test/pkg/model"
)

func ConvertOrgReqToDaoModel(defaultTenantId int64, orgReq orgs.CreateOrgRequest) model.Organization {
	org := model.Organization{}
	org.SetTenantId(defaultTenantId)
	if orgReq.Label != nil {
		org.SetLabel(*orgReq.Label)
	}
	if orgReq.Kind != nil {
		org.SetType(model.OrganizationType(*orgReq.Kind))
	}
	org.SetStatus(model.OrganizationStatus(orgReq.Status))
	return org
}

func ConvertOrgModelToOrgResp(org model.OrganizationInterface) orgs.OrganizationResponse {
	return orgs.OrganizationResponse{
		Code:   org.GetCode(),
		Label:  org.GetLabel(),
		Status: int(org.GetStatus()),
		Kind:   string(org.GetType()),
	}
}
