package mangas

import (
  "errors"
  "manga-explorer/internal/common"
  "math"
)

var ErrUnknownStatus = errors.New("status unknown")

func NewStatus(val string) (Status, error) {
  switch val {
  case "completed":
    return StatusCompleted, nil
  case "ongoing":
    return StatusOnGoing, nil
  case "drafted":
    return StatusDraft, nil
  case "dropped":
    return StatusDropped, nil
  case "hiatus":
    return StatusHiatus, nil
  default:
    return Status(math.MaxUint8), ErrUnknownStatus
  }
}

const (
  StatusCompleted Status = iota
  StatusOnGoing
  StatusDraft
  StatusDropped
  StatusHiatus
)

type Status uint8

func (s Status) String() string {
  switch s {
  case StatusCompleted:
    return "completed"
  case StatusOnGoing:
    return "ongoing"
  case StatusDraft:
    return "drafted"
  case StatusDropped:
    return "dropped"
  case StatusHiatus:
    return "hiatus"
  default:
    return "unknown"
  }
}

func (s Status) Underlying() uint8 {
  return (uint8)(s)
}

func (s Status) Validate() error {
  val := s.Underlying()
  if val > 4 {
    return ErrUnknownStatus
  }
  return nil
}

// TODO: Move it, it should not be belongs here
type SearchFilter struct {
  Title           string
  Genres          common.CriterionOption[string]
  Origins         []common.Country
  IsOriginInclude bool
}

func (f SearchFilter) HasGenre() bool {
  return f.Genres.HasInclude()
  //|| f.Genres.HasExclude()
}

func (f SearchFilter) HasBothGenre() bool {
  return f.Genres.HasInclude()
  //&& f.Genres.HasExclude()
}

func (f SearchFilter) HasOrigin() bool {
  return len(f.Origins) != 0
}

func (f SearchFilter) HasTitle() bool {
  return len(f.Title) != 0
}

type CommentObject string

const (
  CommentObjectManga   = CommentObject("manga")
  CommentObjectChapter = CommentObject("chapter")
  CommentObjectPage    = CommentObject("page")
)

func (c CommentObject) String() string {
  return c.Underlying()
}

func (c CommentObject) Underlying() string {
  return string(c)
}

func (c CommentObject) Validate() error {
  return nil
}
