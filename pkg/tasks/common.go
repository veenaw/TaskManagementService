package tasks

import "github.com/gosimple/slug"

func TaskFromRequest(req Task) TaskTable {
	id := req.Id
	if id == "" {
		id = slug.Make(req.Title)
	}
	return TaskTable{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Status:      req.Status,
	}
}

func TaskToResponse(tbl TaskTable) Task {
	id := tbl.Id
	if id == "" {
		id = slug.Make(tbl.Title)
	}
	return Task{
		Id:          id,
		Title:       tbl.Title,
		Description: tbl.Description,
		DueDate:     tbl.DueDate,
		Status:      tbl.Status,
	}
}
