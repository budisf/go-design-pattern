package driver

import (
	area "ethical-be/modules/v1/utilities/zone/area/handler"
	repos "ethical-be/modules/v1/utilities/zone/area/repository"
	service "ethical-be/modules/v1/utilities/zone/area/services"
)

var (
	AreaRepository = repos.NewAreasRepository(DB)
	AreaService    = service.NewAreasService(AreaRepository)
	AreaHandler    = area.NewAreasHandler(AreaService, RegionsService)
)
