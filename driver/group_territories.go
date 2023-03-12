package driver

import (
	gt "ethical-be/modules/v1/utilities/zone/gt/handler"
	repos "ethical-be/modules/v1/utilities/zone/gt/repository"
	service "ethical-be/modules/v1/utilities/zone/gt/services"
)

var (
	GtRepository = repos.NewGtrespository(DB)
	GtService    = service.NewGtService(GtRepository)
	GtHandler    = gt.NewGtHandler(GtService, AreaService)
)
