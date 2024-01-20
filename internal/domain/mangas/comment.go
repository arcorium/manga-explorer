package mangas

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	user_entity "manga-explorer/internal/domain/users"
	"time"
)

type Comment struct {
	bun.BaseModel `bun:"table:comments"`

	Id         string        `bun:",pk,type:uuid"`
	ParentId   string        `bun:",nullzero,default:null,type:uuid"` // To check, whether the comment is in root or replying another comment
	ObjectType CommentObject `bun:",nullzero,notnull"`
	ObjectId   string        `bun:",nullzero,notnull,type:uuid"` // Object can be Manga, Chapter, and even Page based on ObjectType
	UserId     string        `bun:",notnull,type:uuid"`

	Comment  string `bun:",nullzero,notnull,type:text"`
	IsEdited bool   `bun:",default:false"`
	Like     uint64 `bun:",default:0"`
	Dislike  uint64 `bun:",default:0"`

	CreatedAt time.Time `bun:",notnull"`
	UpdatedAt time.Time `bun:",notnull"`

	User          *user_entity.User `bun:"rel:belongs-to,join:user_id=id,on_delete:SET DEFAULT"`
	ParentComment *Comment          `bun:"rel:belongs-to,join:parent_id=id,on_delete:CASCADE"` // Indicate if the comment is replying other comments or not
}

func (c *Comment) ValidateAsReply(parent *Comment) bool {
	return parent.ObjectType == c.ObjectType && parent.ObjectId == c.ObjectId
}

func NewComment(objectType CommentObject, objectId, userId string, comment string, like, dislike uint64) Comment {
	currentTime := time.Now()
	return Comment{
		Id:         uuid.NewString(),
		ObjectType: objectType,
		ObjectId:   objectId,
		UserId:     userId,
		Comment:    comment,
		Like:       like,
		Dislike:    dislike,
		CreatedAt:  currentTime,
		UpdatedAt:  currentTime,
	}
}

func NewReplyComment(parentId string, objectType CommentObject, objectId, userId, comment string, like, dislike uint64) Comment {
	currentTime := time.Now()
	return Comment{
		Id:         uuid.NewString(),
		ParentId:   parentId,
		ObjectType: objectType,
		ObjectId:   objectId,
		UserId:     userId,
		Comment:    comment,
		IsEdited:   false,
		Like:       like,
		Dislike:    dislike,
		CreatedAt:  currentTime,
		UpdatedAt:  currentTime,
	}
}

func NewReplyComment2(parent *Comment, userId, comment string, like, dislike uint64) Comment {
	return NewReplyComment(parent.Id, parent.ObjectType, parent.ObjectId, userId, comment, like, dislike)
}
