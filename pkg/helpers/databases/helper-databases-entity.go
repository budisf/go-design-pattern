package helperDatabases

/*
	|--------------------------------------------------------------------------
	| Entity for needed mapper on to object on table
	|--------------------------------------------------------------------------
	| @noted:
	| {name_struct}Entity
	| please give name struct which according as above.

|
*/
type QueryParamPaginationEntity struct {
	Page    *int    `form:"page"`
	Limit   *int    `form:"limit"`
	Offset  *int    `form:"offset"`
	Search  *string `form:"search"`
	OrderBy *string `form:"order_by`
}
