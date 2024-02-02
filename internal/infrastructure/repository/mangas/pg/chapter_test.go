package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/biter777/countries"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/containers"
	"reflect"
	"testing"
	"time"
)

func createChapterForTest(volumeId, translatorId, title string, lang countries.CountryCode, chapter uint64, id ...string) *mangas.Chapter {
	tmp := mangas.NewChapter(volumeId, translatorId, title, lang, chapter, time.Now())

	if len(id) == 1 {
		tmp.Id = id[0]
	}
	tmp.Language = common.Language(lang.Alpha2())
	return &tmp
}

func createChapterForTest2(volumeId, translatorId, title string, lang string, chapter uint64, id ...string) *mangas.Chapter {
	tmp := mangas.NewChapter(volumeId, translatorId, title, countries.Indonesia, chapter, time.Now())

	tmp.Language = common.Language(lang)
	if len(id) == 1 {
		tmp.Id = id[0]
	}
	return &tmp
}

func createPageForTest(id, chapterId, imageUrl string, page uint16) mangas.Page {
	tmp := mangas.NewPage(chapterId, imageUrl, page)
	tmp.Id = id
	return tmp
}

func Test_chapterRepository_CreateChapter(t *testing.T) {
	type args struct {
		chapter *mangas.Chapter
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Normal",
			args: args{
				chapter: createChapterForTest("92077652-ebc9-413c-8bfc-7f72a60a128c", "4afa29b2-d543-4489-b8ef-93f57781c9f6", "Title", countries.Indonesia, 1),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Duplicate Chapter with different user",
			args: args{
				chapter: createChapterForTest("d2b047b6-e9ca-44da-93ec-e4af7ae56c8e", "4afa29b2-d543-4489-b8ef-93f57781c9f6", "Title", countries.Indonesia, 1),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Duplicate chapter with same user and different language",
			args: args{
				chapter: createChapterForTest("d2b047b6-e9ca-44da-93ec-e4af7ae56c8e", "c7760836-71e7-4664-99e8-a9503482a296", "Title", countries.Indonesia, 5),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Duplicate chapter with same user and language",
			args: args{
				chapter: createChapterForTest("d2b047b6-e9ca-44da-93ec-e4af7ae56c8e", "c7760836-71e7-4664-99e8-a9503482a296", "Title", countries.JP, 1),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "Volume doesn't exists",
			args: args{
				chapter: createChapterForTest(uuid.NewString(), "4afa29b2-d543-4489-b8ef-93f57781c9f6", "Title", countries.Indonesia, 1),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "Bad volume id as uuid",
			args: args{
				chapter: createChapterForTest("asdascbz12309-dsa", "4afa29b2-d543-4489-b8ef-93f57781c9f6", "Title", countries.Indonesia, 1),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "Translator doesn't exists",
			args: args{
				chapter: createChapterForTest("92077652-ebc9-413c-8bfc-7f72a60a128c", uuid.NewString(), "Title", countries.Indonesia, 1),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "Bad translator id as uuid",
			args: args{
				chapter: createChapterForTest("92077652-ebc9-413c-8bfc-7f72a60a128c", "asdsabczjd-123boabs", "Title", countries.Indonesia, 1),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		tx, err := Db.Begin()
		require.NoError(t, err)
		c := NewMangaChapter(tx)

		t.Run(tt.name, func(t *testing.T) {
			defer func(tx bun.Tx) {
				require.NoError(t, tx.Rollback())
			}(tx)

			err := c.CreateChapter(tt.args.chapter)
			if !tt.wantErr(t, err) {
				t.Errorf("CreateChapter(%v)", tt.args.chapter)
				return
			}

			if err != nil {
				return
			}

			got, err := c.FindChapter(tt.args.chapter.Id)
			require.NoError(t, err)
			require.NotNil(t, got)

			// Ignore time
			got.CreatedAt = tt.args.chapter.CreatedAt
			got.UpdatedAt = tt.args.chapter.UpdatedAt
			got.PublishDate = tt.args.chapter.PublishDate
			// Ignore relations
			got.Pages = tt.args.chapter.Pages
			got.Comments = tt.args.chapter.Comments
			got.Volume = tt.args.chapter.Volume
			got.Translator = tt.args.chapter.Translator

			require.Equal(t, tt.args.chapter, got)
		})
	}
}

func Test_chapterRepository_DeleteChapter(t *testing.T) {
	type args struct {
		chapterId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Normal",
			args: args{
				chapterId: "cf7ddaa5-2637-41a8-96ac-af76202302e1",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Chapter doesn't exists",
			args: args{
				chapterId: uuid.NewString(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "Bad chapter id as uuid",
			args: args{
				chapterId: "asdasd123-cxzbi123",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		tx, err := Db.Begin()
		require.NoError(t, err)
		c := NewMangaChapter(tx)

		t.Run(tt.name, func(t *testing.T) {
			defer func(tx bun.Tx) {
				require.NoError(t, tx.Rollback())
			}(tx)
			err := c.DeleteChapter(tt.args.chapterId)
			if !tt.wantErr(t, err) {
				t.Errorf("DeleteChapter(%v)", tt.args.chapterId)
				return
			}

			// Prevent bad driver from transaction
			if err != nil {
				return
			}

			got, err := c.FindChapter(tt.args.chapterId)
			require.Error(t, err)
			require.Nil(t, got)
		})
	}
}

func Test_chapterRepository_DeleteChapterPages(t *testing.T) {
	type args struct {
		chapterId string
		pages     []uint16
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Normal",
			args: args{
				chapterId: "2bf2f231-a352-4528-a02f-13d310bfb6b2",
				pages: []uint16{
					1, 4, 6, 19,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Empty pages to be deleted",
			args: args{
				chapterId: "2bf2f231-a352-4528-a02f-13d310bfb6b2",
				pages:     []uint16{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, sql.ErrNoRows)
			},
		},
		{
			name: "Chapter pages doesn't exists",
			args: args{
				chapterId: "2bf2f231-a352-4528-a02f-13d310bfb6b2",
				pages: []uint16{
					32, 64, 123,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "Chapter pages some exists and some doesn't exists",
			args: args{
				chapterId: "2bf2f231-a352-4528-a02f-13d310bfb6b2",
				pages: []uint16{
					1, 4, 100, 130,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Chapter doesn't exists",
			args: args{
				chapterId: uuid.NewString(),
				pages: []uint16{
					1, 4, 6, 19,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "Bad chapter id as uuid",
			args: args{
				chapterId: "asdadasd-zxcz123",
				pages: []uint16{
					1, 4, 6, 19,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		tx, err := Db.Begin()
		require.NoError(t, err)
		c := NewMangaChapter(tx)

		t.Run(tt.name, func(t *testing.T) {
			defer func(tx bun.Tx) {
				require.NoError(t, tx.Rollback())
			}(tx)
			err := c.DeleteChapterPages(tt.args.chapterId, tt.args.pages)
			if !tt.wantErr(t, err, tt.args.pages) {
				t.Errorf("DeleteChapterPages(%v, %v)", tt.args.chapterId, tt.args.pages)
				return
			}
		})
	}
}

func Test_chapterRepository_EditChapter(t *testing.T) {
	type args struct {
		chapter *mangas.Chapter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Normal",
			args: args{
				chapter: createChapterForTest("d2b047b6-e9ca-44da-93ec-e4af7ae56c8e", "73141bf0-5e64-4f52-acab-098f1efa3fa7", "New Title", countries.Japan, 1, "2bf2f231-a352-4528-a02f-13d310bfb6b2"),
			},
			wantErr: false,
		},
		{
			name: "Edit become duplicate",
			args: args{
				chapter: createChapterForTest("d2b047b6-e9ca-44da-93ec-e4af7ae56c8e", "73141bf0-5e64-4f52-acab-098f1efa3fa7", "New Title", countries.Japan, 1, "84c7a245-bc1b-4b10-aae1-082a2db78ace"),
			},
			wantErr: true,
		},
		{
			name: "Chapter doesn't exists",
			args: args{
				chapter: createChapterForTest("d2b047b6-e9ca-44da-93ec-e4af7ae56c8e", "73141bf0-5e64-4f52-acab-098f1efa3fa7", "New Title", countries.Japan, 1, uuid.NewString()),
			},
			wantErr: true,
		},
		{
			name: "Change some ids into bad uuid",
			args: args{
				chapter: createChapterForTest("asdasbczxfca-dasbdoa", "73141bf0-5e64-4f52-acab-098f1efa3fa7", "New Title", countries.Japan, 1, "2bf2f231-a352-4528-a02f-13d310bfb6b2"),
			},
			wantErr: true,
		},
		{
			name: "Change to non-existent volume",
			args: args{
				chapter: createChapterForTest(uuid.NewString(), "73141bf0-5e64-4f52-acab-098f1efa3fa7", "New Title", countries.Japan, 1, "2bf2f231-a352-4528-a02f-13d310bfb6b2"),
			},
			wantErr: true,
		},
		{
			name: "Change to non-existent translator",
			args: args{
				chapter: createChapterForTest("d2b047b6-e9ca-44da-93ec-e4af7ae56c8e", uuid.NewString(), "New Title", countries.Japan, 1, "2bf2f231-a352-4528-a02f-13d310bfb6b2"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tx, err := Db.Begin()
		require.NoError(t, err)
		c := NewMangaChapter(tx)

		t.Run(tt.name, func(t *testing.T) {
			defer func(tx bun.Tx) {
				require.NoError(t, tx.Rollback())
			}(tx)

			err := c.EditChapter(tt.args.chapter)
			if (err != nil) != tt.wantErr {
				t.Errorf("EditChapter() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			chapters, err := c.FindVolumeDetails(tt.args.chapter.VolumeId)
			require.NoError(t, err)
			res := containers.SliceFilter(chapters, func(chapter *mangas.Chapter) bool {
				return chapter.Id == tt.args.chapter.Id
			})
			require.Len(t, res, 1)

			// Ignore time
			res[0].CreatedAt = tt.args.chapter.CreatedAt
			res[0].UpdatedAt = tt.args.chapter.UpdatedAt
			res[0].PublishDate = tt.args.chapter.PublishDate
			// Ignore relations
			res[0].Translator = tt.args.chapter.Translator
			res[0].Volume = tt.args.chapter.Volume
			res[0].Comments = tt.args.chapter.Comments
			res[0].Pages = tt.args.chapter.Pages

			if !reflect.DeepEqual(&res[0], tt.args.chapter) {
				t.Errorf("EditChapter: expected %v got %v", res[0], *tt.args.chapter)
			}

		})
	}
}

func Test_chapterRepository_FindChapterPages(t *testing.T) {
	type args struct {
		chapterId string
	}
	tests := []struct {
		name    string
		args    args
		want    []mangas.Page
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Normal",
			args: args{
				chapterId: "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5",
			},
			want: []mangas.Page{
				createPageForTest("486d1d5a-28be-4625-aea7-54950362439f", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/220x100.png/cc0000/ffffff", 1),
				createPageForTest("a7189298-a554-4c18-8460-71e7ba1017f0", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/138x100.png/ff4444/ffffff", 2),
				createPageForTest("df8f7a7a-e4ba-40bf-b394-b1f227b4a8e9", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/145x100.png/cc0000/ffffff", 3),
				createPageForTest("a1d5acd0-1416-4447-9703-fafd50953685", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/184x100.png/cc0000/ffffff", 4),
				createPageForTest("4b51f396-90d7-465c-991a-44dd70b24a4c", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/144x100.png/5fa2dd/ffffff", 5),
				createPageForTest("35c913f3-68ec-42f6-ac36-41d1ecea1241", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/144x100.png/ff4444/ffffff", 6),
				createPageForTest("7bded066-3ab4-415d-8524-700258e6b7db", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/140x100.png/ff4444/ffffff", 7),
				createPageForTest("53da784c-39f9-441e-ab95-c4517ef404c8", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/103x100.png/dddddd/000000", 8),
				createPageForTest("23fa68c5-43e3-45ea-ad54-6bb04950db60", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/183x100.png/5fa2dd/ffffff", 9),
				createPageForTest("f06366ae-026c-43c0-a676-902bbbc006f6", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/100x100.png/cc0000/ffffff", 10),
				createPageForTest("7d0188ec-bf77-4596-aaf4-342f90620346", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/234x100.png/5fa2dd/ffffff", 11),
				createPageForTest("a94878b7-0814-434e-bd21-2e3e7cf3f59e", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/217x100.png/ff4444/ffffff", 12),
				createPageForTest("db460163-a680-4e57-9780-e477ac60e0ae", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/156x100.png/cc0000/ffffff", 13),
				createPageForTest("75ad701b-de19-4f88-a67a-927893965e2a", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/157x100.png/cc0000/ffffff", 14),
				createPageForTest("2e20b416-2c83-4130-bd5b-0bc91784b1e2", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/221x100.png/ff4444/ffffff", 15),
				createPageForTest("964db383-8b4c-4771-993c-4fb24c7f09da", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/241x100.png/5fa2dd/ffffff", 16),
				createPageForTest("7317d847-512b-43b1-8f34-5343df8db2a2", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/193x100.png/ff4444/ffffff", 17),
				createPageForTest("25159423-a6e4-4e5b-a90e-cc4d2cc6e437", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/177x100.png/5fa2dd/ffffff", 18),
				createPageForTest("dbbd6dcc-e6a7-4fa2-9746-912c47441b7b", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/200x100.png/dddddd/000000", 19),
				createPageForTest("b5b05e43-33fb-40b8-b7fe-92d95d258802", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/205x100.png/cc0000/ffffff", 20),
				createPageForTest("f31de4a3-bdf2-40c0-99ed-236f5bac5ce6", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/149x100.png/cc0000/ffffff", 21),
				createPageForTest("a494b2a8-ad01-4b9d-914e-5133733d68d6", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/238x100.png/dddddd/000000", 22),
				createPageForTest("de9045b6-95a9-4f46-ab96-a68f2c501aa4", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/100x100.png/dddddd/000000", 23),
				createPageForTest("5417dc5a-699f-4c7a-8055-61e8939c411f", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/201x100.png/dddddd/000000", 24),
				createPageForTest("f09833b1-7607-4d3d-8474-cab35932c0dc", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/237x100.png/cc0000/ffffff", 25),
				createPageForTest("2205a355-51c6-4259-a5ae-a5f66100c7a9", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/164x100.png/cc0000/ffffff", 26),
				createPageForTest("551959e6-db62-4c10-846d-e21d6ca1f3e9", "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "http://dummyimage.com/143x100.png/dddddd/000000", 27),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				if err != nil {
					t.Errorf(i[0].(string))
					return false
				}
				return true
			},
		},
		{
			name: "Chapter doesn't exists",
			args: args{
				chapterId: uuid.NewString(),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				if err == nil {
					t.Errorf(i[0].(string))
					return false
				}
				return true
			},
		},
		{
			name: "Chapter doesn't have pages",
			args: args{
				chapterId: "a63305d3-6d90-4702-afa2-e5844437f344",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				if err == nil {
					t.Errorf(i[0].(string))
					return false
				}
				return true
			},
		},
		{
			name: "Bad chapter id as uuid",
			args: args{
				chapterId: "asdasdas-xzcazsbd123",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				if err == nil {
					t.Errorf(i[0].(string))
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		c := NewMangaChapter(Db)
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.FindChapterDetails(tt.args.chapterId)
			if !tt.wantErr(t, err, fmt.Sprintf("FindChapterDetails(%v)", tt.args.chapterId)) {
				return
			}

			if count := util.NilCount[mangas.Page](got, tt.want); count == 2 {
				return
			} else if count == 1 {
				t.Errorf("FindVolumeDetails(): expected: %v got: %v", tt.want, got)
				return
			}

			require.Len(t, got, len(tt.want))
			for i := 0; i < len(got); i++ {
				// Ignore relations
				got[i].Chapter = tt.want[i].Chapter
			}

			assert.Equalf(t, tt.want, got, "FindChapterDetails(%v)", tt.args.chapterId)
		})
	}
}

func Test_chapterRepository_FindVolumeChapters(t *testing.T) {
	type args struct {
		volumeId string
	}
	tests := []struct {
		name    string
		args    args
		want    []mangas.Chapter
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Normal",
			args: args{
				volumeId: "a06dd728-c7af-472b-b346-6376805c9cd5",
			},
			want: []mangas.Chapter{
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "db22a444-41a9-41db-ad8c-cb47759a98a8", "Milius", "EN", 1, "df3b7788-c4ac-496d-b1ba-015ccbce1361"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "c7760836-71e7-4664-99e8-a9503482a296", "Secret of the Wings", "JP", 1, "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "Sherlock Holmes in Washington", "JP", 1, "772a38d8-4741-490f-ac3d-22da14820c4f"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "73141bf0-5e64-4f52-acab-098f1efa3fa7", "99 francs", "EN", 1, "da807afe-1645-4627-908f-0e0b95c89bc3"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "a11b349d-59eb-4ebf-9276-eeaec5bdeacc", "Satan's Brew (Satansbraten)", "JP", 1, "2b560fcd-b397-4bcd-aa05-2e1bbb677f2a"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "4d704d17-8900-45d7-83a0-a10e4a4950d9", "Love Object", "JP", 1, "e541d656-5ef5-484a-8885-2aeadc01ddce"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "db22a444-41a9-41db-ad8c-cb47759a98a8", "Cast A Deadly Spell", "JP", 1, "15a2b024-5901-47b9-80a9-651ed39bd7bc"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "a11b349d-59eb-4ebf-9276-eeaec5bdeacc", "i hate myself :)", "EN", 1, "390f825e-5a55-42d3-8171-3c0c2b28612c"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "4d704d17-8900-45d7-83a0-a10e4a4950d9", "Watching the Detectives", "EN", 1, "76c8dc51-565e-4fea-988d-6867eeae4416"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "Barefoot Gen 2 (Hadashi no Gen II)", "EN", 1, "2eff8d6f-37ff-4cbc-849e-fd032a763f53"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "73141bf0-5e64-4f52-acab-098f1efa3fa7", "Balseros (Cuban Rafters)", "EN", 2, "51cc9c0c-1828-49f2-944b-4bfebec50663"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "a11b349d-59eb-4ebf-9276-eeaec5bdeacc", "We Don't Live Here Anymore", "JP", 2, "c39dc0a1-a2e2-4243-b6d4-3fbc2da9eef3"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "a11b349d-59eb-4ebf-9276-eeaec5bdeacc", "Ghost World", "EN", 2, "3c44e4c1-198c-4127-96c3-8df5c1ebf68d"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "db22a444-41a9-41db-ad8c-cb47759a98a8", "Total Eclipse", "EN", 2, "61feb7fb-374e-4d43-b997-5d6fb180d28e"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "War of the Dead - Stone's War ", "JP", 2, "df185b48-4869-486a-9bbf-b36568f7df34"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "db22a444-41a9-41db-ad8c-cb47759a98a8", "Conquest", "JP", 2, "6e8c479b-9030-4ec9-aeb4-9d8db0e06ced"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "4d704d17-8900-45d7-83a0-a10e4a4950d9", "Death at a Funeral", "EN", 2, "1ffdcf6d-bc78-478f-a53a-5b2489e1405d"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "4d704d17-8900-45d7-83a0-a10e4a4950d9", "Lost Horizon", "JP", 2, "75687d34-9991-4b9b-b4f8-015f53115e82"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "db22a444-41a9-41db-ad8c-cb47759a98a8", "Beyond (Svinalängorna)", "JP", 3, "ada88c10-4444-4896-ba77-ea09bc8f7b35"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "4d704d17-8900-45d7-83a0-a10e4a4950d9", "Bridget Jones's Diary", "JP", 3, "4c1a24a9-8da9-4bcc-8285-a705f1950369"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "db22a444-41a9-41db-ad8c-cb47759a98a8", "DeadHeads", "EN", 3, "801745a8-9d83-4ac9-8b69-3716e5dd4354"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "4d704d17-8900-45d7-83a0-a10e4a4950d9", "Sovereign's Company", "EN", 3, "96ca3564-4ccf-42d1-a4a4-2d7f0e25ae5d"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "db22a444-41a9-41db-ad8c-cb47759a98a8", "By the Bluest of Seas (U samogo sinego morya)", "EN", 4, "a63305d3-6d90-4702-afa2-e5844437f344"),
				*createChapterForTest2("a06dd728-c7af-472b-b346-6376805c9cd5", "db22a444-41a9-41db-ad8c-cb47759a98a8", "Ballad of a Soldier (Ballada o soldate)", "EN", 5, "642a3df0-f330-4476-851c-1704e814226f"),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Volume doesn't exists",
			args: args{
				volumeId: uuid.NewString(),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, sql.ErrNoRows)
			},
		},
		{
			name: "Volume doesn't have chapters",
			args: args{
				volumeId: "8178a559-7565-46b0-a4ff-8c9f2162a7c7",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "Bad volume id as uuid",
			args: args{
				volumeId: "asdadsa0zxc23-zcabsd",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		c := NewMangaChapter(Db)
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.FindVolumeDetails(tt.args.volumeId)
			if !tt.wantErr(t, err) {
				t.Errorf("FindVolumeDetails(%v)", tt.args.volumeId)
				return
			}

			if count := util.NilCount[mangas.Chapter](got, tt.want); count == 2 {
				return
			} else if count == 1 {
				t.Errorf("FindVolumeDetails(): expected: %v got: %v", tt.want, got)
				return
			}

			require.Len(t, got, len(tt.want))
			for i := 0; i < len(got); i++ {
				// Ignore time
				got[i].CreatedAt = tt.want[i].CreatedAt
				got[i].UpdatedAt = tt.want[i].UpdatedAt
				got[i].PublishDate = tt.want[i].PublishDate
				// Ignore relations
				got[i].Translator = tt.want[i].Translator
				got[i].Volume = tt.want[i].Volume
				got[i].Comments = tt.want[i].Comments
				got[i].Pages = tt.want[i].Pages
			}
			assert.Equalf(t, tt.want, got, "FindVolumeDetails(%v)", tt.args.volumeId)
		})
	}
}

func Test_chapterRepository_InsertChapterPage(t *testing.T) {
	type args struct {
		pages []mangas.Page
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Normal",
			args: args{
				pages: []mangas.Page{
					createPageForTest(uuid.NewString(), "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "something9.jpg", 28),
					createPageForTest(uuid.NewString(), "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "something8.jpg", 29),
					createPageForTest(uuid.NewString(), "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "something6.jpg", 30),
					createPageForTest(uuid.NewString(), "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "something3.jpg", 31),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Page duplicate",
			args: args{
				pages: []mangas.Page{
					createPageForTest(uuid.NewString(), "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "sadiabsd.jpg", 27),
					createPageForTest(uuid.NewString(), "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "sadiabsd.jpg", 21),
					createPageForTest(uuid.NewString(), "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "sadiabsd.jpg", 22),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "Different pages in one chapter using same image",
			args: args{
				pages: []mangas.Page{
					createPageForTest(uuid.NewString(), "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "sadiabsd.jpg", 41),
					createPageForTest(uuid.NewString(), "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "sadiabsd.jpg", 42),
					createPageForTest(uuid.NewString(), "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "sadiabsd.jpg", 43),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Different chapter using same image",
			args: args{
				pages: []mangas.Page{
					createPageForTest(uuid.NewString(), "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "sadiabsd.jpg", 41),
					createPageForTest(uuid.NewString(), "df3b7788-c4ac-496d-b1ba-015ccbce1361", "sadiabsd.jpg", 41),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Chapter doesn't exists",
			args: args{
				pages: []mangas.Page{
					createPageForTest(uuid.NewString(), uuid.NewString(), "sadiabsd.jpg", 41),
					createPageForTest(uuid.NewString(), uuid.NewString(), "sadiabsd.jpg", 41),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "Bad chapter id as uuid",
			args: args{
				pages: []mangas.Page{
					createPageForTest(uuid.NewString(), "asdsabcnai2", "sadiabsd.jpg", 41),
					createPageForTest(uuid.NewString(), "asdadsad123-zxc", "sadiabsd.jpg", 41),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "Mixed",
			args: args{
				pages: []mangas.Page{
					createPageForTest(uuid.NewString(), "d2d71f7a-1f84-4dd0-8775-fb31aa723fd5", "sadiabsd.jpg", 41),
					createPageForTest(uuid.NewString(), "asdsabcnai2", "sadiabsd.jpg", 41),
					createPageForTest(uuid.NewString(), uuid.NewString(), "sadiabsd.jpg", 41),
					createPageForTest(uuid.NewString(), "df3b7788-c4ac-496d-b1ba-015ccbce1361", "sadiabsd.jpg", 41),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		tx, err := Db.Begin()
		require.NoError(t, err)
		c := NewMangaChapter(tx)

		t.Run(tt.name, func(t *testing.T) {
			defer func(tx2 bun.Tx) {
				require.NoError(t, tx2.Rollback())
			}(tx)
			err := c.InsertChapterPages(tt.args.pages)
			if !tt.wantErr(t, err) {
				t.Errorf(fmt.Sprintf("InsertChapterPages(%v)", tt.args.pages))
				return
			}
		})
	}
}

func Test_chapterRepository_FindChapter(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *mangas.Chapter
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Normal",
			args: args{
				id: "2bf2f231-a352-4528-a02f-13d310bfb6b2",
			},
			want: createChapterForTest("d2b047b6-e9ca-44da-93ec-e4af7ae56c8e", "73141bf0-5e64-4f52-acab-098f1efa3fa7", "Leaving (Partir)", countries.Japan, 1, "2bf2f231-a352-4528-a02f-13d310bfb6b2"),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Chapter doesn't exists",
			args: args{
				id: uuid.NewString(),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, sql.ErrNoRows)
			},
		},
		{
			name: "Bad chapter id as uuid",
			args: args{
				id: "asdabc-123boas9",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		c := NewMangaChapter(Db)
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.FindChapter(tt.args.id)
			if !tt.wantErr(t, err) {
				t.Errorf("FindChapter(%v)", tt.args.id)
				return
			}

			if count := util.NilCount[mangas.Chapter](got, tt.want); count == 2 || err != nil {
				return
			} else if count == 1 {
				t.Errorf("FindChapterComments(): expected: %v got: %v", tt.want, got)
				return
			}

			got.CreatedAt = tt.want.CreatedAt
			got.UpdatedAt = tt.want.UpdatedAt
			got.PublishDate = tt.want.PublishDate
			// Ignore relations
			got.Translator = tt.want.Translator
			got.Volume = tt.want.Volume
			got.Comments = tt.want.Comments
			got.Pages = tt.want.Pages

			assert.Equalf(t, tt.want, got, "FindChapter(%v)", tt.args.id)
		})
	}
}
