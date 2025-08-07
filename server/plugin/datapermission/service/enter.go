package service

type ServiceGroup struct {
	DataPermissionService DataPermissionService
}

var ServiceGroupApp = new(ServiceGroup)
