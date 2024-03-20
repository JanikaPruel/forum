package controller

import (
	"sort"

	"forum/internal/model"
)

func (ctl *BaseController) FilterByLikes(posts []*model.Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Likes > posts[j].Likes
	})
}

func (ctl *BaseController) FilterByDislikes(posts []*model.Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Dislikes > posts[j].Dislikes
	})
}
