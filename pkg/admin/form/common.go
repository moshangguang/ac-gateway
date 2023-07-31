package form

type IdFrom struct {
	Id int64 ` form:"id" binding:"required"`
}
