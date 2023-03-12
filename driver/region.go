package driver

import (
	region "ethical-be/modules/v1/utilities/zone/region/handler"
	repos "ethical-be/modules/v1/utilities/zone/region/repository"
	service "ethical-be/modules/v1/utilities/zone/region/services"
)

var (
	RegionsRepository = repos.NewRegionsRepository(DB)
	RegionsService    = service.NewRegionsService(RegionsRepository)
	RegionsHandler    = region.NewRegionsHandler(RegionsService)
)
