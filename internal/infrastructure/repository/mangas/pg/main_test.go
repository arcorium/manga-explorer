package pg

import (
  "manga-explorer/internal/domain/mangas"
  "manga-explorer/internal/util/containers"
  "testing"

  "github.com/uptrace/bun"
  "manga-explorer/database"
  "manga-explorer/internal/app/common"
)

var Conf common.Config
var Db *bun.DB

var MangaParent = mangas.NewComment(mangas.CommentObjectManga, "19382f54-1da7-4cb7-807d-9f6030bb121e", "c7760836-71e7-4664-99e8-a9503482a296", "Hello", 10, 1)                        // Parent
var MangaParent2 = mangas.NewReplyComment(MangaParent.Id, mangas.CommentObjectManga, "19382f54-1da7-4cb7-807d-9f6030bb121e", "dc4402e4-0f88-400a-978e-8bb3880ab063", "Hello2", 10, 1) // Parent
// BUGFIX: Reply comment ObjectId should be on the same object type and id (handled on service)
var MangaComments = []mangas.Comment{
  MangaParent,
  mangas.NewReplyComment(MangaParent.Id, mangas.CommentObjectManga, "19382f54-1da7-4cb7-807d-9f6030bb121e", "4d704d17-8900-45d7-83a0-a10e4a4950d9", "Hi", 5, 1),
  MangaParent2,
  mangas.NewReplyComment(MangaParent2.Id, mangas.CommentObjectManga, "19382f54-1da7-4cb7-807d-9f6030bb121e", "db22a444-41a9-41db-ad8c-cb47759a98a8", "Hi", 5, 1),
}

var ChapterParent = mangas.NewComment(mangas.CommentObjectChapter, "cf7ddaa5-2637-41a8-96ac-af76202302e1", "c7760836-71e7-4664-99e8-a9503482a296", "Hello", 10, 1)
var ChapterParent2 = mangas.NewReplyComment(ChapterParent.Id, mangas.CommentObjectChapter, "cf7ddaa5-2637-41a8-96ac-af76202302e1", "dc4402e4-0f88-400a-978e-8bb3880ab063", "Hello2", 10, 1)

var ChapterComments = []mangas.Comment{
  ChapterParent,
  mangas.NewReplyComment(ChapterParent.Id, mangas.CommentObjectChapter, "cf7ddaa5-2637-41a8-96ac-af76202302e1", "4d704d17-8900-45d7-83a0-a10e4a4950d9", "Hi", 5, 1),
  ChapterParent2,
  mangas.NewReplyComment(ChapterParent2.Id, mangas.CommentObjectChapter, "cf7ddaa5-2637-41a8-96ac-af76202302e1", "db22a444-41a9-41db-ad8c-cb47759a98a8", "Hi", 5, 1),
}

var PageParent = mangas.NewComment(mangas.CommentObjectPage, "04df6dea-668d-40f9-85ff-17eea73a51e7", "c7760836-71e7-4664-99e8-a9503482a296", "Hello", 10, 1)
var PageParent2 = mangas.NewReplyComment(PageParent.Id, mangas.CommentObjectPage, "04df6dea-668d-40f9-85ff-17eea73a51e7", "dc4402e4-0f88-400a-978e-8bb3880ab063", "Hello2", 10, 1)

var PageComments = []mangas.Comment{
  PageParent,
  mangas.NewReplyComment(PageParent.Id, mangas.CommentObjectPage, "04df6dea-668d-40f9-85ff-17eea73a51e7", "4d704d17-8900-45d7-83a0-a10e4a4950d9", "Hi", 5, 1),
  PageParent2,
  mangas.NewReplyComment(PageParent2.Id, mangas.CommentObjectPage, "04df6dea-668d-40f9-85ff-17eea73a51e7", "db22a444-41a9-41db-ad8c-cb47759a98a8", "Hi", 5, 1),
}

func LoadCommentFixtures() bool {
  tx, err := Db.Begin()
  if err != nil {
    return false
  }

  repo := NewComment(tx)

  for _, comment := range containers.CombineSplices(MangaComments, ChapterComments, PageComments) {
    err := repo.CreateComment(&comment)
    if err != nil {
      tx.Rollback()
      return false
    }
  }

  tx.Commit()
  return true
}

func RemoveCommentFixtures() {
  repo := NewComment(Db)

  for _, comment := range []mangas.Comment{MangaParent, ChapterParent, PageParent} {
    err := repo.DeleteComment(comment.Id)
    if err != nil {
      //log.Print(err)
    }
  }
}

func TestMain(m *testing.M) {
  var err error
  Conf, err = common.LoadConfig("test", "../../../../../")
  if err != nil {
    panic(err)
  }
  Db = database.Open(&Conf, false)
  defer database.Close(Db)

  if !LoadCommentFixtures() {
    panic("Failed on load comment fixtures")
  }
  defer RemoveCommentFixtures()

  m.Run()
}
