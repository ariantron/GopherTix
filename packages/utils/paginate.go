package utils

func Paginate(count int64, page, limit int) (totalPages, offset int) {
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}
	totalPages = int((count + int64(limit) - 1) / int64(limit))
	offset = (page - 1) * limit
	return totalPages, offset
}
