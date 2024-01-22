package dto

import "manga-explorer/internal/domain/users/dto"

type CommentResponse struct {
	User     dto.UserResponse `json:"user"`
	Comment  string           `json:"comment"`
	IsEdited bool             `json:"is_edited"`
	Like     uint64           `json:"like"`
	Dislike  uint64           `json:"dislike"`
}

type commentCreatInputParent struct {
	ParentId string `json:"parent_id" binding:"uuid4"`
	UserId   string `json:"-"`
	Comment  string `json:"comment" binding:"required"`
}

func (c commentCreatInputParent) HasParent() bool {
	return len(c.ParentId) != 0
}

type MangaCommentCreateInput struct {
	MangaId string `uri:"manga_id" binding:"required"`
	commentCreatInputParent
}

type ChapterCommentCreateInput struct {
	ChapterId string `uri:"chapter_id" binding:"required"`
	commentCreatInputParent
}

type PageCommentCreateInput struct {
	PageId string `uri:"page_id" binding:"required"`
	commentCreatInputParent
}
