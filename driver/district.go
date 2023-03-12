package driver

import (
	district "ethical-be/modules/v1/utilities/zone/district/handler"
	repos "ethical-be/modules/v1/utilities/zone/district/repository"
	service "ethical-be/modules/v1/utilities/zone/district/services"
)

var (
	DistrictRepository = repos.NewDistrictRepository(DB)
	DistrictService    = service.NewDistrictService(DistrictRepository)
	DistrictHandler    = district.NewDistrictHandler(DistrictService)
)
