package pg

import (
	"fmt"
	"github.com/biter777/countries"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/infrastructure/repository"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/opt"
	"reflect"
	"testing"
	"time"
)

func newMangaForTest(id opt.Optional[string], title, desc, coverUrl string, year uint16, status mangas.Status, region countries.CountryCode) *mangas.Manga {
	temp := mangas.NewManga(title, desc, coverUrl, year, status, region)
	if id.HasValue() {
		temp.Id = *id.Value()
	}
	return &temp // Expected object to be allocated on heap
}

func newVolumeForTest(id opt.Optional[string], mangaId string, number uint32, title, desc string) *mangas.Volume {
	temp := mangas.NewVolume(mangaId, number, title, desc)
	if id.HasValue() {
		temp.Id = *id.Value()
	}
	return &temp
}

func newMangaFavoriteForTest(userId, mangaId string, manga *mangas.Manga, user *users.User) mangas.MangaFavorite {
	temp := mangas.NewFavorite(userId, mangaId)
	temp.Manga = manga
	temp.User = user
	return temp
}

func newMangaHistoryForTest(lastView time.Time, manga *mangas.Manga) mangas.MangaHistory {
	return mangas.MangaHistory{
		LastView: lastView,
		Manga:    manga,
	}
}

func Test_mangaRepository_CreateManga(t *testing.T) {
	type args struct {
		mangas *mangas.Manga
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Normal",
			args: args{
				mangas: newMangaForTest(opt.Null[string](), "One-Punch Man", "something", "https://google.com/images/one_punch_man.png", 2010, mangas.StatusOnGoing, countries.Japan),
			},
			wantErr: false,
		},
		{
			name: "Duplicate Id",
			args: args{
				mangas: newMangaForTest(opt.New("df3be3a1-f02f-4d2e-afe8-83dc61f46839"), "One-Punch Man", "something", "https://google.com/images/one_punch_man.png", 2010, mangas.StatusOnGoing, countries.Japan),
			},
			wantErr: true,
		},
		{
			name: "Duplicate Title",
			args: args{
				mangas: newMangaForTest(opt.Null[string](), "Homeboy", "", "", 2020, mangas.StatusHiatus, countries.China),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tx, err := Db.Begin()
		require.NoError(t, err)
		mangaRepo := NewManga(tx)

		t.Run(tt.name, func(t *testing.T) {
			err := mangaRepo.CreateManga(tt.args.mangas)
			defer func(tx bun.Tx) {
				require.NoError(t, tx.Rollback())
			}(tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateManga() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mangaRepository_CreateVolume(t *testing.T) {
	type args struct {
		volume *mangas.Volume
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Normal",
			args: args{
				volume: newVolumeForTest(opt.Null[string](), "2aa478df-9f0f-4e67-b652-f9b01023eefb", 24, "title", "desc"),
			},
			wantErr: false,
		},
		{
			name: "Duplicate Id",
			args: args{
				volume: newVolumeForTest(opt.New("412be2a3-bd05-49cb-97ba-2748fa3fce7e"), "2aa478df-9f0f-4e67-b652-f9b01023eefb", 1024, "title", "desc"),
			},
			wantErr: true,
		},
		{
			name: "Non-exist Manga",
			args: args{
				volume: newVolumeForTest(opt.Null[string](), uuid.NewString(), 1, "title", "desc"),
			},
			wantErr: true,
		},
		{
			name: "Duplicate Volume Number on Manga",
			args: args{
				volume: newVolumeForTest(opt.Null[string](), "2aa478df-9f0f-4e67-b652-f9b01023eefb", 1, "title", "desc"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tx, err := Db.Begin()
		require.NoError(t, err)
		mangaRepo := NewManga(tx)
		t.Run(tt.name, func(t *testing.T) {
			defer func(tx bun.Tx) {
				require.NoError(t, tx.Rollback())
			}(tx)

			if err := mangaRepo.CreateVolume(tt.args.volume); (err != nil) != tt.wantErr {
				t.Errorf("CreateVolume() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				require.NoError(t, tx.Rollback())
				tx = util.DropError(Db.Begin())
			}
		})
	}
}

func Test_mangaRepository_DeleteVolume(t *testing.T) {
	type args struct {
		mangaId string
		volume  uint32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Normal",
			args: args{
				mangaId: "2aa478df-9f0f-4e67-b652-f9b01023eefb",
				volume:  1,
			},
			wantErr: false,
		},
		{
			name: "Non-existent volume number",
			args: args{
				mangaId: "2aa478df-9f0f-4e67-b652-f9b01023eefb",
				volume:  24,
			},
			wantErr: true,
		},
		{
			name: "Non-existent manga",
			args: args{
				mangaId: uuid.NewString(),
				volume:  2,
			},
			wantErr: true,
		},
		{
			name: "Bad uuid",
			args: args{
				mangaId: "zxczcasdasd",
				volume:  23,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tx, err := Db.Begin()
		require.NoError(t, err)

		mangaRepo := NewManga(tx)
		t.Run(tt.name, func(t *testing.T) {
			defer func(tx2 bun.Tx) {
				require.NoError(t, tx2.Rollback())
			}(tx)

			if err := mangaRepo.DeleteVolume(tt.args.mangaId, tt.args.volume); (err != nil) != tt.wantErr {
				t.Errorf("DeleteVolume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mangaRepository_FindMangaById(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *mangas.Manga
		wantErr bool
	}{
		{
			name: "Normal",
			args: args{
				id: "2aa478df-9f0f-4e67-b652-f9b01023eefb",
			},
			want:    newMangaForTest(opt.New[string]("2aa478df-9f0f-4e67-b652-f9b01023eefb"), "Homeboy", "Maecenas leo odio, condimentum id, luctus nec, molestie sed, justo. Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti.", "", 2013, mangas.StatusDraft, countries.CN),
			wantErr: false,
		},
		{
			name: "Non-existent Manga",
			args: args{
				id: uuid.NewString(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Bad UUID",
			args: args{
				id: "asdasdasd0asdas-cxiasd0asjd",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			mangaRepo := NewManga(Db)
			got, err := mangaRepo.FindMangaById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindMangaById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil && tt.want != nil {
				// Ignore time fields
				got.CreatedAt = tt.want.CreatedAt
				got.UpdatedAt = tt.want.UpdatedAt
				// Ignore relation fields
				got.Genres = tt.want.Genres
				//got.Chapters = tt.want.Chapters
				got.Comments = tt.want.Comments
				got.Ratings = tt.want.Ratings
				got.Translations = tt.want.Translations
				got.Volumes = tt.want.Volumes
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindMangaById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mangaRepository_FindMangaFavorites(t *testing.T) {
	type args struct {
		userId string
		param  repository.QueryParameter
	}
	tests := []struct {
		name    string
		args    args
		want    repository.PagedQueryResult[[]mangas.MangaFavorite]
		wantErr bool
	}{
		{
			name: "User present and get all mangas",
			args: args{
				userId: "c7760836-71e7-4664-99e8-a9503482a296", // Has 18 favorite mangas
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  25,
				},
			},
			want: repository.PagedQueryResult[[]mangas.MangaFavorite]{
				Data: []mangas.MangaFavorite{
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "37dd72b9-93ae-44d5-a30c-4c95d508211c", newMangaForTest(opt.New("37dd72b9-93ae-44d5-a30c-4c95d508211c"), "Aberdeen", "Morbi sem mauris, laoreet ut, rhoncus aliquet, pulvinar sed, nisl. Nunc rhoncus dui vel sem. Sed sagittis. Nam congue, risus semper porta volutpat, quam pede lobortis ligula, sit amet eleifend pede libero quis orci. Nullam molestie nibh in lectus. Pellentesque at nulla. Suspendisse potenti. Cras in purus eu magna vulputate luctus. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.", "", 2008, mangas.StatusOnGoing, countries.ID), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "bddd0fb4-55e8-4eac-978f-540dccfcf23c", newMangaForTest(opt.New("bddd0fb4-55e8-4eac-978f-540dccfcf23c"), "Answer This!", "Nulla tellus. In sagittis dui vel nisl. Duis ac nibh. Fusce lacus purus, aliquet at, feugiat non, pretium quis, lectus. Suspendisse potenti. In eleifend quam a odio.", "", 1987, mangas.StatusHiatus, countries.FR), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "a672026b-10bf-4a9a-83b3-00b4af613533", newMangaForTest(opt.New("a672026b-10bf-4a9a-83b3-00b4af613533"), "Artist, The", "Integer aliquet, massa id lobortis convallis, tortor risus dapibus augue, vel accumsan tellus nisi eu orci. Mauris lacinia sapien quis libero. Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh. In quis justo.", "", 2005, mangas.StatusCompleted, countries.SE), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "79553079-db78-4f40-8da1-db9d8fe7441e", newMangaForTest(opt.New("79553079-db78-4f40-8da1-db9d8fe7441e"), "Bethlehem", "Morbi sem mauris, laoreet ut, rhoncus aliquet, pulvinar sed, nisl. Nunc rhoncus dui vel sem. Sed sagittis. Nam congue, risus semper porta volutpat, quam pede lobortis ligula, sit amet eleifend pede libero quis orci. Nullam molestie nibh in lectus. Pellentesque at nulla. Suspendisse potenti. Cras in purus eu magna vulputate luctus.", "", 1993, mangas.StatusHiatus, countries.LT), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "89dd85c3-c8ce-4da9-802f-2e413dcbb4bb", newMangaForTest(opt.New("89dd85c3-c8ce-4da9-802f-2e413dcbb4bb"), "Butterfly Kiss", "Maecenas rhoncus aliquam lacus. Morbi quis tortor id nulla ultrices aliquet. Maecenas leo odio, condimentum id, luctus nec, molestie sed, justo. Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti. Nullam porttitor lacus at turpis. Donec posuere metus vitae ipsum. Aliquam non mauris. Morbi non lectus. Aliquam sit amet diam in magna bibendum imperdiet. Nullam orci pede, venenatis non, sodales sed, tincidunt eu, felis. Fusce posuere felis sed lacus.", "", 1998, mangas.StatusCompleted, countries.ID), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "d12a4b85-d42e-4a91-a3ce-3f81f7880a13", newMangaForTest(opt.New("d12a4b85-d42e-4a91-a3ce-3f81f7880a13"), "Captain January", "Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh. In quis justo.", "", 2003, mangas.StatusCompleted, countries.PH), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "e73192c5-b5c4-44ff-9b13-04f4e2983010", newMangaForTest(opt.New("e73192c5-b5c4-44ff-9b13-04f4e2983010"), "Common Places (a.k.a. Common Ground) (Lugares comunes)", "Pellentesque eget nunc. Donec quis orci eget orci vehicula condimentum. Curabitur in libero ut massa volutpat convallis. Morbi odio odio, elementum eu, interdum eu, tincidunt in, leo. Maecenas pulvinar lobortis est. Phasellus sit amet erat. Nulla tempus. Vivamus in felis eu sapien cursus vestibulum. Proin eu mi.", "", 2008, mangas.StatusDropped, countries.SI), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "02a8e33f-36f5-4279-92ff-f3b375bd9fdc", newMangaForTest(opt.New("02a8e33f-36f5-4279-92ff-f3b375bd9fdc"), "Cry_Wolf (a.k.a. Cry Wolf)", "Duis aliquam convallis nunc. Proin at turpis a pede posuere nonummy. Integer non velit. Donec diam neque, vestibulum eget, vulputate ut, ultrices vel, augue. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Donec pharetra, magna vestibulum aliquet ultrices, erat tortor sollicitudin mi, sit amet lobortis sapien sapien non mi. Integer ac neque.", "", 2001, mangas.StatusCompleted, countries.TH), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "4584fd42-ae8a-4fc6-be74-6785fe5b25f1", newMangaForTest(opt.New("4584fd42-ae8a-4fc6-be74-6785fe5b25f1"), "Escape Artist, The", "Duis mattis egestas metus. Aenean fermentum. Donec ut mauris eget massa tempor convallis. Nulla neque libero, convallis eget, eleifend luctus, ultricies eu, nibh. Quisque id justo sit amet sapien dignissim vestibulum. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla dapibus dolor vel est. Donec odio justo, sollicitudin ut, suscipit a, feugiat et, eros.", "", 2011, mangas.StatusHiatus, countries.CN), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "1ced9617-1659-49fb-ab7a-b55316630193", newMangaForTest(opt.New("1ced9617-1659-49fb-ab7a-b55316630193"), "Evil That Men Do, The", "Donec ut mauris eget massa tempor convallis. Nulla neque libero, convallis eget, eleifend luctus, ultricies eu, nibh. Quisque id justo sit amet sapien dignissim vestibulum. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla dapibus dolor vel est.", "", 1993, mangas.StatusDropped, countries.MX), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "631f78af-1df6-4b2f-80c6-e4ac1d278f07", newMangaForTest(opt.New("631f78af-1df6-4b2f-80c6-e4ac1d278f07"), "Fifty-Fifty (a.k.a. Schizo) (Shiza)", "Aliquam quis turpis eget elit sodales scelerisque. Mauris sit amet eros. Suspendisse accumsan tortor quis turpis. Sed ante. Vivamus tortor. Duis mattis egestas metus. Aenean fermentum. Donec ut mauris eget massa tempor convallis.", "", 2000, mangas.StatusCompleted, countries.CU), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "45aa60a3-e40b-40be-bd4a-cb10c11bfa85", newMangaForTest(opt.New("45aa60a3-e40b-40be-bd4a-cb10c11bfa85"), "Futurama: Bender's Game", "Donec semper sapien a libero. Nam dui. Proin leo odio, porttitor id, consequat in, consequat ut, nulla. Sed accumsan felis. Ut at dolor quis odio consequat varius.", "", 1980, mangas.StatusCompleted, countries.CN), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "90a685c1-d4ec-4dc4-a649-91f4aaeff24f", newMangaForTest(opt.New("90a685c1-d4ec-4dc4-a649-91f4aaeff24f"), "Greatest, The", "Phasellus id sapien in sapien iaculis congue. Vivamus metus arcu, adipiscing molestie, hendrerit at, vulputate vitae, nisl. Aenean lectus. Pellentesque eget nunc.", "", 1995, mangas.StatusHiatus, countries.PH), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "66bcfab4-1ec8-4a4d-a7a8-c8a730a3822f", newMangaForTest(opt.New("66bcfab4-1ec8-4a4d-a7a8-c8a730a3822f"), "If Looks Could Kill", "Vivamus vestibulum sagittis sapien. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Etiam vel augue. Vestibulum rutrum rutrum neque. Aenean auctor gravida sem. Praesent id massa id nisl venenatis lacinia. Aenean sit amet justo. Morbi ut odio. Cras mi pede, malesuada in, imperdiet et, commodo vulputate, justo. In blandit ultrices enim.", "", 2007, mangas.StatusDraft, countries.AL), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "3f6aa253-bd52-4d6a-a406-6026eb1e9759", newMangaForTest(opt.New("3f6aa253-bd52-4d6a-a406-6026eb1e9759"), "Long Live Death (Viva la muerte)", "Fusce posuere felis sed lacus. Morbi sem mauris, laoreet ut, rhoncus aliquet, pulvinar sed, nisl. Nunc rhoncus dui vel sem. Sed sagittis. Nam congue, risus semper porta volutpat, quam pede lobortis ligula, sit amet eleifend pede libero quis orci. Nullam molestie nibh in lectus. Pellentesque at nulla. Suspendisse potenti. Cras in purus eu magna vulputate luctus.", "", 2011, mangas.StatusCompleted, countries.US), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "04aadb59-0df2-4153-9237-f0a3b606e1c6", newMangaForTest(opt.New("04aadb59-0df2-4153-9237-f0a3b606e1c6"), "Marvin Hamlisch: What He Did for Love", "Aliquam sit amet diam in magna bibendum imperdiet. Nullam orci pede, venenatis non, sodales sed, tincidunt eu, felis. Fusce posuere felis sed lacus. Morbi sem mauris, laoreet ut, rhoncus aliquet, pulvinar sed, nisl. Nunc rhoncus dui vel sem. Sed sagittis. Nam congue, risus semper porta volutpat, quam pede lobortis ligula, sit amet eleifend pede libero quis orci. Nullam molestie nibh in lectus. Pellentesque at nulla. Suspendisse potenti. Cras in purus eu magna vulputate luctus.", "", 2009, mangas.StatusCompleted, countries.CN), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "d5b29b3e-7994-4d36-a43e-4528ff29ba41", newMangaForTest(opt.New("d5b29b3e-7994-4d36-a43e-4528ff29ba41"), "Miracle in Cell No. 7", "Aenean auctor gravida sem. Praesent id massa id nisl venenatis lacinia. Aenean sit amet justo. Morbi ut odio. Cras mi pede, malesuada in, imperdiet et, commodo vulputate, justo. In blandit ultrices enim. Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Proin interdum mauris non ligula pellentesque ultrices. Phasellus id sapien in sapien iaculis congue.", "", 1996, mangas.StatusOnGoing, countries.RU), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "0cd41304-846f-4294-957f-0b7bbbbe3879", newMangaForTest(opt.New("0cd41304-846f-4294-957f-0b7bbbbe3879"), "Promise, The (Versprechen, Das)", "Aliquam quis turpis eget elit sodales scelerisque. Mauris sit amet eros. Suspendisse accumsan tortor quis turpis. Sed ante. Vivamus tortor. Duis mattis egestas metus. Aenean fermentum. Donec ut mauris eget massa tempor convallis. Nulla neque libero, convallis eget, eleifend luctus, ultricies eu, nibh.", "", 2005, mangas.StatusHiatus, countries.ID), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "77f7d74a-8332-4f02-b0a4-5eb4d36d3725", newMangaForTest(opt.New("77f7d74a-8332-4f02-b0a4-5eb4d36d3725"), "Pulse (Kairo)", "Cras non velit nec nisi vulputate nonummy. Maecenas tincidunt lacus at velit. Vivamus vel nulla eget eros elementum pellentesque. Quisque porta volutpat erat. Quisque erat eros, viverra eget, congue eget, semper rutrum, nulla. Nunc purus. Phasellus in felis. Donec semper sapien a libero. Nam dui. Proin leo odio, porttitor id, consequat in, consequat ut, nulla. Sed accumsan felis. Ut at dolor quis odio consequat varius.", "", 2004, mangas.StatusOnGoing, countries.GT), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "79f0b72e-10f1-4505-8fd5-acf07d9296f6", newMangaForTest(opt.New("79f0b72e-10f1-4505-8fd5-acf07d9296f6"), "Rated X: A Journey Through Porn", "Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Vivamus vestibulum sagittis sapien. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Etiam vel augue. Vestibulum rutrum rutrum neque. Aenean auctor gravida sem. Praesent id massa id nisl venenatis lacinia. Aenean sit amet justo. Morbi ut odio. Cras mi pede, malesuada in, imperdiet et, commodo vulputate, justo. In blandit ultrices enim. Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Proin interdum mauris non ligula pellentesque ultrices. Phasellus id sapien in sapien iaculis congue. Vivamus metus arcu, adipiscing molestie, hendrerit at, vulputate vitae, nisl.", "", 2008, mangas.StatusDropped, countries.PT), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "1571d0ec-a7f4-4ba5-8bab-f1737f723e0c", newMangaForTest(opt.New("1571d0ec-a7f4-4ba5-8bab-f1737f723e0c"), "Reel Injun", "Duis mattis egestas metus. Aenean fermentum. Donec ut mauris eget massa tempor convallis. Nulla neque libero, convallis eget, eleifend luctus, ultricies eu, nibh. Quisque id justo sit amet sapien dignissim vestibulum.", "", 1987, mangas.StatusCompleted, countries.CZ), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "8a8a82d7-8dfc-4c22-91db-d0623f229f18", newMangaForTest(opt.New("8a8a82d7-8dfc-4c22-91db-d0623f229f18"), "Ride with the Devil", "Nullam orci pede, venenatis non, sodales sed, tincidunt eu, felis. Fusce posuere felis sed lacus. Morbi sem mauris, laoreet ut, rhoncus aliquet, pulvinar sed, nisl. Nunc rhoncus dui vel sem. Sed sagittis. Nam congue, risus semper porta volutpat, quam pede lobortis ligula, sit amet eleifend pede libero quis orci. Nullam molestie nibh in lectus. Pellentesque at nulla. Suspendisse potenti. Cras in purus eu magna vulputate luctus. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Vivamus vestibulum sagittis sapien. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Etiam vel augue. Vestibulum rutrum rutrum neque.", "", 2006, mangas.StatusOnGoing, countries.BW), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "f1a93136-d653-4134-b3a1-dd00837364d0", newMangaForTest(opt.New("f1a93136-d653-4134-b3a1-dd00837364d0"), "State of Grace", "Quisque id justo sit amet sapien dignissim vestibulum. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla dapibus dolor vel est. Donec odio justo, sollicitudin ut, suscipit a, feugiat et, eros. Vestibulum ac est lacinia nisi venenatis tristique. Fusce congue, diam id ornare imperdiet, sapien urna pretium nisl, ut volutpat sapien arcu sed augue. Aliquam erat volutpat. In congue. Etiam justo. Etiam pretium iaculis justo. In hac habitasse platea dictumst. Etiam faucibus cursus urna. Ut tellus. Nulla ut erat id mauris vulputate elementum.", "", 2007, mangas.StatusDropped, countries.KG), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "8a17c1f7-1dc1-454a-a980-aaf5c118a2ae", newMangaForTest(opt.New("8a17c1f7-1dc1-454a-a980-aaf5c118a2ae"), "Stealing Harvard", "Maecenas leo odio, condimentum id, luctus nec, molestie sed, justo. Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti. Nullam porttitor lacus at turpis.", "", 2005, mangas.StatusHiatus, countries.SE), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "793f2e77-d801-4f0d-885d-916fb37c332a", newMangaForTest(opt.New("793f2e77-d801-4f0d-885d-916fb37c332a"), "Ward, The", "Vestibulum rutrum rutrum neque. Aenean auctor gravida sem. Praesent id massa id nisl venenatis lacinia. Aenean sit amet justo. Morbi ut odio. Cras mi pede, malesuada in, imperdiet et, commodo vulputate, justo. In blandit ultrices enim. Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Proin interdum mauris non ligula pellentesque ultrices. Phasellus id sapien in sapien iaculis congue. Vivamus metus arcu, adipiscing molestie, hendrerit at, vulputate vitae, nisl. Aenean lectus.", "", 1992, mangas.StatusDraft, countries.AR), nil),
				},
				Total: 25,
			},
			wantErr: false,
		},
		{
			name: "User present and get all mangas using 0 limit",
			args: args{
				userId: "c7760836-71e7-4664-99e8-a9503482a296", // Has 18 favorite mangas
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  0,
				},
			},
			want: repository.PagedQueryResult[[]mangas.MangaFavorite]{
				Data: []mangas.MangaFavorite{
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "37dd72b9-93ae-44d5-a30c-4c95d508211c", newMangaForTest(opt.New("37dd72b9-93ae-44d5-a30c-4c95d508211c"), "Aberdeen", "Morbi sem mauris, laoreet ut, rhoncus aliquet, pulvinar sed, nisl. Nunc rhoncus dui vel sem. Sed sagittis. Nam congue, risus semper porta volutpat, quam pede lobortis ligula, sit amet eleifend pede libero quis orci. Nullam molestie nibh in lectus. Pellentesque at nulla. Suspendisse potenti. Cras in purus eu magna vulputate luctus. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.", "", 2008, mangas.StatusOnGoing, countries.ID), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "bddd0fb4-55e8-4eac-978f-540dccfcf23c", newMangaForTest(opt.New("bddd0fb4-55e8-4eac-978f-540dccfcf23c"), "Answer This!", "Nulla tellus. In sagittis dui vel nisl. Duis ac nibh. Fusce lacus purus, aliquet at, feugiat non, pretium quis, lectus. Suspendisse potenti. In eleifend quam a odio.", "", 1987, mangas.StatusHiatus, countries.FR), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "a672026b-10bf-4a9a-83b3-00b4af613533", newMangaForTest(opt.New("a672026b-10bf-4a9a-83b3-00b4af613533"), "Artist, The", "Integer aliquet, massa id lobortis convallis, tortor risus dapibus augue, vel accumsan tellus nisi eu orci. Mauris lacinia sapien quis libero. Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh. In quis justo.", "", 2005, mangas.StatusCompleted, countries.SE), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "79553079-db78-4f40-8da1-db9d8fe7441e", newMangaForTest(opt.New("79553079-db78-4f40-8da1-db9d8fe7441e"), "Bethlehem", "Morbi sem mauris, laoreet ut, rhoncus aliquet, pulvinar sed, nisl. Nunc rhoncus dui vel sem. Sed sagittis. Nam congue, risus semper porta volutpat, quam pede lobortis ligula, sit amet eleifend pede libero quis orci. Nullam molestie nibh in lectus. Pellentesque at nulla. Suspendisse potenti. Cras in purus eu magna vulputate luctus.", "", 1993, mangas.StatusHiatus, countries.LT), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "89dd85c3-c8ce-4da9-802f-2e413dcbb4bb", newMangaForTest(opt.New("89dd85c3-c8ce-4da9-802f-2e413dcbb4bb"), "Butterfly Kiss", "Maecenas rhoncus aliquam lacus. Morbi quis tortor id nulla ultrices aliquet. Maecenas leo odio, condimentum id, luctus nec, molestie sed, justo. Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti. Nullam porttitor lacus at turpis. Donec posuere metus vitae ipsum. Aliquam non mauris. Morbi non lectus. Aliquam sit amet diam in magna bibendum imperdiet. Nullam orci pede, venenatis non, sodales sed, tincidunt eu, felis. Fusce posuere felis sed lacus.", "", 1998, mangas.StatusCompleted, countries.ID), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "d12a4b85-d42e-4a91-a3ce-3f81f7880a13", newMangaForTest(opt.New("d12a4b85-d42e-4a91-a3ce-3f81f7880a13"), "Captain January", "Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh. In quis justo.", "", 2003, mangas.StatusCompleted, countries.PH), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "e73192c5-b5c4-44ff-9b13-04f4e2983010", newMangaForTest(opt.New("e73192c5-b5c4-44ff-9b13-04f4e2983010"), "Common Places (a.k.a. Common Ground) (Lugares comunes)", "Pellentesque eget nunc. Donec quis orci eget orci vehicula condimentum. Curabitur in libero ut massa volutpat convallis. Morbi odio odio, elementum eu, interdum eu, tincidunt in, leo. Maecenas pulvinar lobortis est. Phasellus sit amet erat. Nulla tempus. Vivamus in felis eu sapien cursus vestibulum. Proin eu mi.", "", 2008, mangas.StatusDropped, countries.SI), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "02a8e33f-36f5-4279-92ff-f3b375bd9fdc", newMangaForTest(opt.New("02a8e33f-36f5-4279-92ff-f3b375bd9fdc"), "Cry_Wolf (a.k.a. Cry Wolf)", "Duis aliquam convallis nunc. Proin at turpis a pede posuere nonummy. Integer non velit. Donec diam neque, vestibulum eget, vulputate ut, ultrices vel, augue. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Donec pharetra, magna vestibulum aliquet ultrices, erat tortor sollicitudin mi, sit amet lobortis sapien sapien non mi. Integer ac neque.", "", 2001, mangas.StatusCompleted, countries.TH), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "4584fd42-ae8a-4fc6-be74-6785fe5b25f1", newMangaForTest(opt.New("4584fd42-ae8a-4fc6-be74-6785fe5b25f1"), "Escape Artist, The", "Duis mattis egestas metus. Aenean fermentum. Donec ut mauris eget massa tempor convallis. Nulla neque libero, convallis eget, eleifend luctus, ultricies eu, nibh. Quisque id justo sit amet sapien dignissim vestibulum. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla dapibus dolor vel est. Donec odio justo, sollicitudin ut, suscipit a, feugiat et, eros.", "", 2011, mangas.StatusHiatus, countries.CN), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "1ced9617-1659-49fb-ab7a-b55316630193", newMangaForTest(opt.New("1ced9617-1659-49fb-ab7a-b55316630193"), "Evil That Men Do, The", "Donec ut mauris eget massa tempor convallis. Nulla neque libero, convallis eget, eleifend luctus, ultricies eu, nibh. Quisque id justo sit amet sapien dignissim vestibulum. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla dapibus dolor vel est.", "", 1993, mangas.StatusDropped, countries.MX), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "631f78af-1df6-4b2f-80c6-e4ac1d278f07", newMangaForTest(opt.New("631f78af-1df6-4b2f-80c6-e4ac1d278f07"), "Fifty-Fifty (a.k.a. Schizo) (Shiza)", "Aliquam quis turpis eget elit sodales scelerisque. Mauris sit amet eros. Suspendisse accumsan tortor quis turpis. Sed ante. Vivamus tortor. Duis mattis egestas metus. Aenean fermentum. Donec ut mauris eget massa tempor convallis.", "", 2000, mangas.StatusCompleted, countries.CU), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "45aa60a3-e40b-40be-bd4a-cb10c11bfa85", newMangaForTest(opt.New("45aa60a3-e40b-40be-bd4a-cb10c11bfa85"), "Futurama: Bender's Game", "Donec semper sapien a libero. Nam dui. Proin leo odio, porttitor id, consequat in, consequat ut, nulla. Sed accumsan felis. Ut at dolor quis odio consequat varius.", "", 1980, mangas.StatusCompleted, countries.CN), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "90a685c1-d4ec-4dc4-a649-91f4aaeff24f", newMangaForTest(opt.New("90a685c1-d4ec-4dc4-a649-91f4aaeff24f"), "Greatest, The", "Phasellus id sapien in sapien iaculis congue. Vivamus metus arcu, adipiscing molestie, hendrerit at, vulputate vitae, nisl. Aenean lectus. Pellentesque eget nunc.", "", 1995, mangas.StatusHiatus, countries.PH), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "66bcfab4-1ec8-4a4d-a7a8-c8a730a3822f", newMangaForTest(opt.New("66bcfab4-1ec8-4a4d-a7a8-c8a730a3822f"), "If Looks Could Kill", "Vivamus vestibulum sagittis sapien. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Etiam vel augue. Vestibulum rutrum rutrum neque. Aenean auctor gravida sem. Praesent id massa id nisl venenatis lacinia. Aenean sit amet justo. Morbi ut odio. Cras mi pede, malesuada in, imperdiet et, commodo vulputate, justo. In blandit ultrices enim.", "", 2007, mangas.StatusDraft, countries.AL), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "3f6aa253-bd52-4d6a-a406-6026eb1e9759", newMangaForTest(opt.New("3f6aa253-bd52-4d6a-a406-6026eb1e9759"), "Long Live Death (Viva la muerte)", "Fusce posuere felis sed lacus. Morbi sem mauris, laoreet ut, rhoncus aliquet, pulvinar sed, nisl. Nunc rhoncus dui vel sem. Sed sagittis. Nam congue, risus semper porta volutpat, quam pede lobortis ligula, sit amet eleifend pede libero quis orci. Nullam molestie nibh in lectus. Pellentesque at nulla. Suspendisse potenti. Cras in purus eu magna vulputate luctus.", "", 2011, mangas.StatusCompleted, countries.US), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "04aadb59-0df2-4153-9237-f0a3b606e1c6", newMangaForTest(opt.New("04aadb59-0df2-4153-9237-f0a3b606e1c6"), "Marvin Hamlisch: What He Did for Love", "Aliquam sit amet diam in magna bibendum imperdiet. Nullam orci pede, venenatis non, sodales sed, tincidunt eu, felis. Fusce posuere felis sed lacus. Morbi sem mauris, laoreet ut, rhoncus aliquet, pulvinar sed, nisl. Nunc rhoncus dui vel sem. Sed sagittis. Nam congue, risus semper porta volutpat, quam pede lobortis ligula, sit amet eleifend pede libero quis orci. Nullam molestie nibh in lectus. Pellentesque at nulla. Suspendisse potenti. Cras in purus eu magna vulputate luctus.", "", 2009, mangas.StatusCompleted, countries.CN), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "d5b29b3e-7994-4d36-a43e-4528ff29ba41", newMangaForTest(opt.New("d5b29b3e-7994-4d36-a43e-4528ff29ba41"), "Miracle in Cell No. 7", "Aenean auctor gravida sem. Praesent id massa id nisl venenatis lacinia. Aenean sit amet justo. Morbi ut odio. Cras mi pede, malesuada in, imperdiet et, commodo vulputate, justo. In blandit ultrices enim. Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Proin interdum mauris non ligula pellentesque ultrices. Phasellus id sapien in sapien iaculis congue.", "", 1996, mangas.StatusOnGoing, countries.RU), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "0cd41304-846f-4294-957f-0b7bbbbe3879", newMangaForTest(opt.New("0cd41304-846f-4294-957f-0b7bbbbe3879"), "Promise, The (Versprechen, Das)", "Aliquam quis turpis eget elit sodales scelerisque. Mauris sit amet eros. Suspendisse accumsan tortor quis turpis. Sed ante. Vivamus tortor. Duis mattis egestas metus. Aenean fermentum. Donec ut mauris eget massa tempor convallis. Nulla neque libero, convallis eget, eleifend luctus, ultricies eu, nibh.", "", 2005, mangas.StatusHiatus, countries.ID), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "77f7d74a-8332-4f02-b0a4-5eb4d36d3725", newMangaForTest(opt.New("77f7d74a-8332-4f02-b0a4-5eb4d36d3725"), "Pulse (Kairo)", "Cras non velit nec nisi vulputate nonummy. Maecenas tincidunt lacus at velit. Vivamus vel nulla eget eros elementum pellentesque. Quisque porta volutpat erat. Quisque erat eros, viverra eget, congue eget, semper rutrum, nulla. Nunc purus. Phasellus in felis. Donec semper sapien a libero. Nam dui. Proin leo odio, porttitor id, consequat in, consequat ut, nulla. Sed accumsan felis. Ut at dolor quis odio consequat varius.", "", 2004, mangas.StatusOnGoing, countries.GT), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "79f0b72e-10f1-4505-8fd5-acf07d9296f6", newMangaForTest(opt.New("79f0b72e-10f1-4505-8fd5-acf07d9296f6"), "Rated X: A Journey Through Porn", "Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Vivamus vestibulum sagittis sapien. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Etiam vel augue. Vestibulum rutrum rutrum neque. Aenean auctor gravida sem. Praesent id massa id nisl venenatis lacinia. Aenean sit amet justo. Morbi ut odio. Cras mi pede, malesuada in, imperdiet et, commodo vulputate, justo. In blandit ultrices enim. Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Proin interdum mauris non ligula pellentesque ultrices. Phasellus id sapien in sapien iaculis congue. Vivamus metus arcu, adipiscing molestie, hendrerit at, vulputate vitae, nisl.", "", 2008, mangas.StatusDropped, countries.PT), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "1571d0ec-a7f4-4ba5-8bab-f1737f723e0c", newMangaForTest(opt.New("1571d0ec-a7f4-4ba5-8bab-f1737f723e0c"), "Reel Injun", "Duis mattis egestas metus. Aenean fermentum. Donec ut mauris eget massa tempor convallis. Nulla neque libero, convallis eget, eleifend luctus, ultricies eu, nibh. Quisque id justo sit amet sapien dignissim vestibulum.", "", 1987, mangas.StatusCompleted, countries.CZ), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "8a8a82d7-8dfc-4c22-91db-d0623f229f18", newMangaForTest(opt.New("8a8a82d7-8dfc-4c22-91db-d0623f229f18"), "Ride with the Devil", "Nullam orci pede, venenatis non, sodales sed, tincidunt eu, felis. Fusce posuere felis sed lacus. Morbi sem mauris, laoreet ut, rhoncus aliquet, pulvinar sed, nisl. Nunc rhoncus dui vel sem. Sed sagittis. Nam congue, risus semper porta volutpat, quam pede lobortis ligula, sit amet eleifend pede libero quis orci. Nullam molestie nibh in lectus. Pellentesque at nulla. Suspendisse potenti. Cras in purus eu magna vulputate luctus. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Vivamus vestibulum sagittis sapien. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Etiam vel augue. Vestibulum rutrum rutrum neque.", "", 2006, mangas.StatusOnGoing, countries.BW), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "f1a93136-d653-4134-b3a1-dd00837364d0", newMangaForTest(opt.New("f1a93136-d653-4134-b3a1-dd00837364d0"), "State of Grace", "Quisque id justo sit amet sapien dignissim vestibulum. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla dapibus dolor vel est. Donec odio justo, sollicitudin ut, suscipit a, feugiat et, eros. Vestibulum ac est lacinia nisi venenatis tristique. Fusce congue, diam id ornare imperdiet, sapien urna pretium nisl, ut volutpat sapien arcu sed augue. Aliquam erat volutpat. In congue. Etiam justo. Etiam pretium iaculis justo. In hac habitasse platea dictumst. Etiam faucibus cursus urna. Ut tellus. Nulla ut erat id mauris vulputate elementum.", "", 2007, mangas.StatusDropped, countries.KG), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "8a17c1f7-1dc1-454a-a980-aaf5c118a2ae", newMangaForTest(opt.New("8a17c1f7-1dc1-454a-a980-aaf5c118a2ae"), "Stealing Harvard", "Maecenas leo odio, condimentum id, luctus nec, molestie sed, justo. Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti. Nullam porttitor lacus at turpis.", "", 2005, mangas.StatusHiatus, countries.SE), nil),
					newMangaFavoriteForTest("c7760836-71e7-4664-99e8-a9503482a296", "793f2e77-d801-4f0d-885d-916fb37c332a", newMangaForTest(opt.New("793f2e77-d801-4f0d-885d-916fb37c332a"), "Ward, The", "Vestibulum rutrum rutrum neque. Aenean auctor gravida sem. Praesent id massa id nisl venenatis lacinia. Aenean sit amet justo. Morbi ut odio. Cras mi pede, malesuada in, imperdiet et, commodo vulputate, justo. In blandit ultrices enim. Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Proin interdum mauris non ligula pellentesque ultrices. Phasellus id sapien in sapien iaculis congue. Vivamus metus arcu, adipiscing molestie, hendrerit at, vulputate vitae, nisl. Aenean lectus.", "", 1992, mangas.StatusDraft, countries.AR), nil),
				},
				Total: 25,
			},
			wantErr: false,
		},
		{
			name: "User present and get first half of mangas",
			args: args{
				userId: "dd2166b0-5e62-4b74-b4cb-4be51a5040dc",
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  9,
				},
			},
			want: repository.PagedQueryResult[[]mangas.MangaFavorite]{
				Data: []mangas.MangaFavorite{
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "df3be3a1-f02f-4d2e-afe8-83dc61f46839", newMangaForTest(opt.New("df3be3a1-f02f-4d2e-afe8-83dc61f46839"), "Ace of Hearts", "Integer aliquet, massa id lobortis convallis, tortor risus dapibus augue, vel accumsan tellus nisi eu orci. Mauris lacinia sapien quis libero. Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh.", "", 2012, mangas.StatusDropped, countries.CN), nil),
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "402d9b14-f14c-4982-b71b-1e1885838f6f", newMangaForTest(opt.New("402d9b14-f14c-4982-b71b-1e1885838f6f"), "Adventures of Robin Hood, The", "Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti. Nullam porttitor lacus at turpis. Donec posuere metus vitae ipsum. Aliquam non mauris. Morbi non lectus.", "", 2010, mangas.StatusCompleted, countries.ID), nil),
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "7c634319-937d-4447-9808-3417474309c1", newMangaForTest(opt.New("7c634319-937d-4447-9808-3417474309c1"), "Black Swan, The", "Suspendisse potenti. Nullam porttitor lacus at turpis. Donec posuere metus vitae ipsum. Aliquam non mauris. Morbi non lectus. Aliquam sit amet diam in magna bibendum imperdiet. Nullam orci pede, venenatis non, sodales sed, tincidunt eu, felis.", "", 2004, mangas.StatusCompleted, countries.CA), nil),
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "89dd85c3-c8ce-4da9-802f-2e413dcbb4bb", newMangaForTest(opt.New("89dd85c3-c8ce-4da9-802f-2e413dcbb4bb"), "Butterfly Kiss", "Maecenas rhoncus aliquam lacus. Morbi quis tortor id nulla ultrices aliquet. Maecenas leo odio, condimentum id, luctus nec, molestie sed, justo. Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti. Nullam porttitor lacus at turpis. Donec posuere metus vitae ipsum. Aliquam non mauris. Morbi non lectus. Aliquam sit amet diam in magna bibendum imperdiet. Nullam orci pede, venenatis non, sodales sed, tincidunt eu, felis. Fusce posuere felis sed lacus.", "", 1998, mangas.StatusCompleted, countries.ID), nil),
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "e73192c5-b5c4-44ff-9b13-04f4e2983010", newMangaForTest(opt.New("e73192c5-b5c4-44ff-9b13-04f4e2983010"), "Common Places (a.k.a. Common Ground) (Lugares comunes)", "Pellentesque eget nunc. Donec quis orci eget orci vehicula condimentum. Curabitur in libero ut massa volutpat convallis. Morbi odio odio, elementum eu, interdum eu, tincidunt in, leo. Maecenas pulvinar lobortis est. Phasellus sit amet erat. Nulla tempus. Vivamus in felis eu sapien cursus vestibulum. Proin eu mi.", "", 2008, mangas.StatusDropped, countries.SI), nil),
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "02a8e33f-36f5-4279-92ff-f3b375bd9fdc", newMangaForTest(opt.New("02a8e33f-36f5-4279-92ff-f3b375bd9fdc"), "Cry_Wolf (a.k.a. Cry Wolf)", "Duis aliquam convallis nunc. Proin at turpis a pede posuere nonummy. Integer non velit. Donec diam neque, vestibulum eget, vulputate ut, ultrices vel, augue. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Donec pharetra, magna vestibulum aliquet ultrices, erat tortor sollicitudin mi, sit amet lobortis sapien sapien non mi. Integer ac neque.", "", 2001, mangas.StatusCompleted, countries.TH), nil),
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "4584fd42-ae8a-4fc6-be74-6785fe5b25f1", newMangaForTest(opt.New("4584fd42-ae8a-4fc6-be74-6785fe5b25f1"), "Escape Artist, The", "Duis mattis egestas metus. Aenean fermentum. Donec ut mauris eget massa tempor convallis. Nulla neque libero, convallis eget, eleifend luctus, ultricies eu, nibh. Quisque id justo sit amet sapien dignissim vestibulum. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla dapibus dolor vel est. Donec odio justo, sollicitudin ut, suscipit a, feugiat et, eros.", "", 2011, mangas.StatusHiatus, countries.CN), nil),
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "c5ef679a-6210-4f90-ba05-076f6cb9ec38", newMangaForTest(opt.New("c5ef679a-6210-4f90-ba05-076f6cb9ec38"), "Girl Who Talked to Dolphins, The", "Nulla ac enim. In tempor, turpis nec euismod scelerisque, quam turpis adipiscing lorem, vitae mattis nibh ligula nec sem. Duis aliquam convallis nunc. Proin at turpis a pede posuere nonummy. Integer non velit. Donec diam neque, vestibulum eget, vulputate ut, ultrices vel, augue. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Donec pharetra, magna vestibulum aliquet ultrices, erat tortor sollicitudin mi, sit amet lobortis sapien sapien non mi.", "", 1993, mangas.StatusDraft, countries.GR), nil),
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "dc91dff4-f673-4096-8e4f-c564ea48efd9", newMangaForTest(opt.New("dc91dff4-f673-4096-8e4f-c564ea48efd9"), "Gridlock'd", "Vestibulum ac est lacinia nisi venenatis tristique. Fusce congue, diam id ornare imperdiet, sapien urna pretium nisl, ut volutpat sapien arcu sed augue. Aliquam erat volutpat. In congue. Etiam justo. Etiam pretium iaculis justo. In hac habitasse platea dictumst. Etiam faucibus cursus urna. Ut tellus.", "", 2007, mangas.StatusCompleted, countries.CN), nil),
				},
				Total: 18,
			},
			wantErr: false,
		},
		{
			name: "User present and get second half of mangas",
			args: args{
				userId: "dd2166b0-5e62-4b74-b4cb-4be51a5040dc",
				param: repository.QueryParameter{
					Offset: 13,
					Limit:  9,
				},
			},
			want: repository.PagedQueryResult[[]mangas.MangaFavorite]{
				Data: []mangas.MangaFavorite{
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "dce60990-f315-48b4-bf3a-709fe7a0652a", newMangaForTest(opt.New("dce60990-f315-48b4-bf3a-709fe7a0652a"), "Sniper", "Etiam pretium iaculis justo. In hac habitasse platea dictumst. Etiam faucibus cursus urna. Ut tellus. Nulla ut erat id mauris vulputate elementum. Nullam varius.", "", 1998, mangas.StatusCompleted, countries.CN), nil),
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "fccf9bae-3873-461b-b557-aa8ae6d786e0", newMangaForTest(opt.New("fccf9bae-3873-461b-b557-aa8ae6d786e0"), "The Party", "Praesent blandit. Nam nulla. Integer pede justo, lacinia eget, tincidunt eget, tempus vel, pede. Morbi porttitor lorem id ligula. Suspendisse ornare consequat lectus. In est risus, auctor sed, tristique in, tempus sit amet, sem. Fusce consequat. Nulla nisl. Nunc nisl.", "", 1993, mangas.StatusOnGoing, countries.RW), nil),
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "27a403ad-26bd-47a1-a73d-7bdff4501238", newMangaForTest(opt.New("27a403ad-26bd-47a1-a73d-7bdff4501238"), "Undertow", "Morbi ut odio. Cras mi pede, malesuada in, imperdiet et, commodo vulputate, justo. In blandit ultrices enim. Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Proin interdum mauris non ligula pellentesque ultrices. Phasellus id sapien in sapien iaculis congue. Vivamus metus arcu, adipiscing molestie, hendrerit at, vulputate vitae, nisl. Aenean lectus.", "", 2003, mangas.StatusOnGoing, countries.ID), nil),
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "793f2e77-d801-4f0d-885d-916fb37c332a", newMangaForTest(opt.New("793f2e77-d801-4f0d-885d-916fb37c332a"), "Ward, The", "Vestibulum rutrum rutrum neque. Aenean auctor gravida sem. Praesent id massa id nisl venenatis lacinia. Aenean sit amet justo. Morbi ut odio. Cras mi pede, malesuada in, imperdiet et, commodo vulputate, justo. In blandit ultrices enim. Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Proin interdum mauris non ligula pellentesque ultrices. Phasellus id sapien in sapien iaculis congue. Vivamus metus arcu, adipiscing molestie, hendrerit at, vulputate vitae, nisl. Aenean lectus.", "", 1992, mangas.StatusDraft, countries.AR), nil),
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "823de473-bc35-481d-bdb9-d3adb2127745", newMangaForTest(opt.New("823de473-bc35-481d-bdb9-d3adb2127745"), "Wisdom", "Nunc nisl. Duis bibendum, felis sed interdum venenatis, turpis enim blandit mi, in porttitor pede justo eu massa. Donec dapibus. Duis at velit eu est congue elementum. In hac habitasse platea dictumst. Morbi vestibulum, velit id pretium iaculis, diam erat fermentum justo, nec condimentum neque sapien placerat ante. Nulla justo. Aliquam quis turpis eget elit sodales scelerisque. Mauris sit amet eros.", "", 1995, mangas.StatusOnGoing, countries.GM), nil),
				},
				Total: 18,
			},
			wantErr: false,
		},
		{
			name: "Non-existent user",
			args: args{
				userId: uuid.NewString(),
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  0,
				},
			},
			want:    repository.PagedQueryResult[[]mangas.MangaFavorite]{nil, 0},
			wantErr: true,
		},
		{
			name: "Bad UUID",
			args: args{
				userId: "sadasdoiasd-dasdn012-asd",
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  0,
				},
			},
			want:    repository.PagedQueryResult[[]mangas.MangaFavorite]{nil, 0},
			wantErr: true,
		},
		{
			name: "User present and out of bound limit",
			args: args{
				userId: "dd2166b0-5e62-4b74-b4cb-4be51a5040dc",
				param: repository.QueryParameter{
					Offset: 15,
					Limit:  24,
				},
			},
			want: repository.PagedQueryResult[[]mangas.MangaFavorite]{
				Data: []mangas.MangaFavorite{
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "27a403ad-26bd-47a1-a73d-7bdff4501238", newMangaForTest(opt.New("27a403ad-26bd-47a1-a73d-7bdff4501238"), "Undertow", "Morbi ut odio. Cras mi pede, malesuada in, imperdiet et, commodo vulputate, justo. In blandit ultrices enim. Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Proin interdum mauris non ligula pellentesque ultrices. Phasellus id sapien in sapien iaculis congue. Vivamus metus arcu, adipiscing molestie, hendrerit at, vulputate vitae, nisl. Aenean lectus.", "", 2003, mangas.StatusOnGoing, countries.ID), nil),
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "793f2e77-d801-4f0d-885d-916fb37c332a", newMangaForTest(opt.New("793f2e77-d801-4f0d-885d-916fb37c332a"), "Ward, The", "Vestibulum rutrum rutrum neque. Aenean auctor gravida sem. Praesent id massa id nisl venenatis lacinia. Aenean sit amet justo. Morbi ut odio. Cras mi pede, malesuada in, imperdiet et, commodo vulputate, justo. In blandit ultrices enim. Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Proin interdum mauris non ligula pellentesque ultrices. Phasellus id sapien in sapien iaculis congue. Vivamus metus arcu, adipiscing molestie, hendrerit at, vulputate vitae, nisl. Aenean lectus.", "", 1992, mangas.StatusDraft, countries.AR), nil),
					newMangaFavoriteForTest("dd2166b0-5e62-4b74-b4cb-4be51a5040dc", "823de473-bc35-481d-bdb9-d3adb2127745", newMangaForTest(opt.New("823de473-bc35-481d-bdb9-d3adb2127745"), "Wisdom", "Nunc nisl. Duis bibendum, felis sed interdum venenatis, turpis enim blandit mi, in porttitor pede justo eu massa. Donec dapibus. Duis at velit eu est congue elementum. In hac habitasse platea dictumst. Morbi vestibulum, velit id pretium iaculis, diam erat fermentum justo, nec condimentum neque sapien placerat ante. Nulla justo. Aliquam quis turpis eget elit sodales scelerisque. Mauris sit amet eros.", "", 1995, mangas.StatusOnGoing, countries.GM), nil),
				},
				Total: 18,
			},
			wantErr: false,
		},
		{
			name: "User present and out of bound offset",
			args: args{
				userId: "dd2166b0-5e62-4b74-b4cb-4be51a5040dc\"",
				param: repository.QueryParameter{
					Offset: 18,
					Limit:  9,
				},
			},
			want:    repository.PagedQueryResult[[]mangas.MangaFavorite]{nil, 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mangaRepo := NewManga(Db)
			got, err := mangaRepo.FindMangaFavorites(tt.args.userId, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindMangaFavorites() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got.Data) != len(tt.want.Data) {
				t.Errorf("FindMangaFavorites() error = different length")
				return
			}

			// Ignore time fields
			for i := 0; i < len(got.Data); i++ {
				got.Data[i].CreatedAt = tt.want.Data[i].CreatedAt
				// Ignore user relation
				got.Data[i].User = tt.want.Data[i].User

				// Ignore time fields
				got.Data[i].Manga.CreatedAt = tt.want.Data[i].Manga.CreatedAt
				got.Data[i].Manga.UpdatedAt = tt.want.Data[i].Manga.UpdatedAt
				// Ignore relation fields
				got.Data[i].Manga.Genres = tt.want.Data[i].Manga.Genres
				got.Data[i].Manga.Comments = tt.want.Data[i].Manga.Comments
				got.Data[i].Manga.Ratings = tt.want.Data[i].Manga.Ratings
				got.Data[i].Manga.Translations = tt.want.Data[i].Manga.Translations
				got.Data[i].Manga.Volumes = tt.want.Data[i].Manga.Volumes

				if !reflect.DeepEqual(got.Data[i].Manga, tt.want.Data[i].Manga) {
					t.Errorf("FindMangaFavorites() \ngot = %v, \nwant = %v", got.Data[i].Manga, tt.want.Data[i].Manga)
				}

				// Remove manga to check the outer struct
				got.Data[i].Manga = nil
				tt.want.Data[i].Manga = nil
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindMangaFavorites() got = \n%v, \nwant = \n%v", got, tt.want)
			}

		})

	}
}

func Test_mangaRepository_FindMangaHistories(t *testing.T) {
	type args struct {
		userId string
		param  repository.QueryParameter
	}
	tests := []struct {
		name    string
		args    args
		want    repository.PagedQueryResult[[]mangas.MangaHistory]
		wantErr bool
	}{
		{
			name: "User present and get all mangas",
			args: args{
				userId: "dc4402e4-0f88-400a-978e-8bb3880ab063",
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  8,
				},
			},
			want: repository.PagedQueryResult[[]mangas.MangaHistory]{
				Data: []mangas.MangaHistory{
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("1bd31e88-0a22-4db2-a894-06a947b4a311"), "Exit Through the Gift Shop", "In est risus, auctor sed, tristique in, tempus sit amet, sem. Fusce consequat. Nulla nisl. Nunc nisl.", "", 1992, mangas.StatusDraft, countries.PL)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("e1674245-bb91-4382-adca-4b2c38878a89"), "Recipients Live and Die in L.A.", "Suspendisse potenti. In eleifend quam a odio. In hac habitasse platea dictumst. Maecenas ut massa quis augue luctus tincidunt. Nulla mollis molestie lorem. Quisque ut erat. Curabitur gravida nisi at nibh. In hac habitasse platea dictumst. Aliquam augue quam, sollicitudin vitae, consectetuer eget, rutrum at, lorem.", "", 2012, mangas.StatusOnGoing, countries.JP)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("df3be3a1-f02f-4d2e-afe8-83dc61f46839"), "Ace of Hearts", "Integer aliquet, massa id lobortis convallis, tortor risus dapibus augue, vel accumsan tellus nisi eu orci. Mauris lacinia sapien quis libero. Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh.", "", 2012, mangas.StatusDropped, countries.CN)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("2aa478df-9f0f-4e67-b652-f9b01023eefb"), "Homeboy", "Maecenas leo odio, condimentum id, luctus nec, molestie sed, justo. Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti.", "", 2013, mangas.StatusDraft, countries.CN)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("35d1bea2-1a13-45e7-a08c-5d35db26444d"), "Freddy vs. Jason", "Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla dapibus dolor vel est. Donec odio justo, sollicitudin ut, suscipit a, feugiat et, eros. Vestibulum ac est lacinia nisi venenatis tristique. Fusce congue, diam id ornare imperdiet, sapien urna pretium nisl, ut volutpat sapien arcu sed augue. Aliquam erat volutpat. In congue. Etiam justo.", "", 1995, mangas.StatusDraft, countries.MY)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("fc1bea74-5fde-4cf0-a332-c957c914d121"), "Red Riding: 1974", "Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti. Nullam porttitor lacus at turpis. Donec posuere metus vitae ipsum.", "", 2008, mangas.StatusOnGoing, countries.PT)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("62c950be-858b-42f2-8799-a09e49bc8589"), "Satanas", "Nam ultrices, libero non mattis pulvinar, nulla pede ullamcorper augue, a suscipit nulla elit ac nulla. Sed vel enim sit amet nunc viverra dapibus. Nulla suscipit ligula in lacus. Curabitur at ipsum ac tellus semper interdum. Mauris ullamcorper purus sit amet nulla.", "", 2010, mangas.StatusOnGoing, countries.BD)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("b8bd3f1e-36e3-4033-8290-c5e0caaeab6d"), "Count Three and Pray", "Suspendisse potenti. In eleifend quam a odio. In hac habitasse platea dictumst. Maecenas ut massa quis augue luctus tincidunt. Nulla mollis molestie lorem. Quisque ut erat. Curabitur gravida nisi at nibh. In hac habitasse platea dictumst. Aliquam augue quam, sollicitudin vitae, consectetuer eget, rutrum at, lorem.", "", 2003, mangas.StatusOnGoing, countries.LT)),
				},
				Total: 8,
			},
			wantErr: false,
		},
		{
			name: "User present and get all mangas using 0 limit",
			args: args{
				userId: "dc4402e4-0f88-400a-978e-8bb3880ab063",
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  0,
				},
			},
			want: repository.PagedQueryResult[[]mangas.MangaHistory]{
				Data: []mangas.MangaHistory{
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("1bd31e88-0a22-4db2-a894-06a947b4a311"), "Exit Through the Gift Shop", "In est risus, auctor sed, tristique in, tempus sit amet, sem. Fusce consequat. Nulla nisl. Nunc nisl.", "", 1992, mangas.StatusDraft, countries.PL)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("e1674245-bb91-4382-adca-4b2c38878a89"), "Recipients Live and Die in L.A.", "Suspendisse potenti. In eleifend quam a odio. In hac habitasse platea dictumst. Maecenas ut massa quis augue luctus tincidunt. Nulla mollis molestie lorem. Quisque ut erat. Curabitur gravida nisi at nibh. In hac habitasse platea dictumst. Aliquam augue quam, sollicitudin vitae, consectetuer eget, rutrum at, lorem.", "", 2012, mangas.StatusOnGoing, countries.JP)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("df3be3a1-f02f-4d2e-afe8-83dc61f46839"), "Ace of Hearts", "Integer aliquet, massa id lobortis convallis, tortor risus dapibus augue, vel accumsan tellus nisi eu orci. Mauris lacinia sapien quis libero. Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh.", "", 2012, mangas.StatusDropped, countries.CN)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("2aa478df-9f0f-4e67-b652-f9b01023eefb"), "Homeboy", "Maecenas leo odio, condimentum id, luctus nec, molestie sed, justo. Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti.", "", 2013, mangas.StatusDraft, countries.CN)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("35d1bea2-1a13-45e7-a08c-5d35db26444d"), "Freddy vs. Jason", "Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla dapibus dolor vel est. Donec odio justo, sollicitudin ut, suscipit a, feugiat et, eros. Vestibulum ac est lacinia nisi venenatis tristique. Fusce congue, diam id ornare imperdiet, sapien urna pretium nisl, ut volutpat sapien arcu sed augue. Aliquam erat volutpat. In congue. Etiam justo.", "", 1995, mangas.StatusDraft, countries.MY)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("fc1bea74-5fde-4cf0-a332-c957c914d121"), "Red Riding: 1974", "Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti. Nullam porttitor lacus at turpis. Donec posuere metus vitae ipsum.", "", 2008, mangas.StatusOnGoing, countries.PT)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("62c950be-858b-42f2-8799-a09e49bc8589"), "Satanas", "Nam ultrices, libero non mattis pulvinar, nulla pede ullamcorper augue, a suscipit nulla elit ac nulla. Sed vel enim sit amet nunc viverra dapibus. Nulla suscipit ligula in lacus. Curabitur at ipsum ac tellus semper interdum. Mauris ullamcorper purus sit amet nulla.", "", 2010, mangas.StatusOnGoing, countries.BD)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("b8bd3f1e-36e3-4033-8290-c5e0caaeab6d"), "Count Three and Pray", "Suspendisse potenti. In eleifend quam a odio. In hac habitasse platea dictumst. Maecenas ut massa quis augue luctus tincidunt. Nulla mollis molestie lorem. Quisque ut erat. Curabitur gravida nisi at nibh. In hac habitasse platea dictumst. Aliquam augue quam, sollicitudin vitae, consectetuer eget, rutrum at, lorem.", "", 2003, mangas.StatusOnGoing, countries.LT)),
				},
				Total: 8,
			},
			wantErr: false,
		},
		{
			name: "User present and get first half",
			args: args{
				userId: "dc4402e4-0f88-400a-978e-8bb3880ab063",
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  4,
				},
			},
			want: repository.PagedQueryResult[[]mangas.MangaHistory]{
				Data: []mangas.MangaHistory{
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("1bd31e88-0a22-4db2-a894-06a947b4a311"), "Exit Through the Gift Shop", "In est risus, auctor sed, tristique in, tempus sit amet, sem. Fusce consequat. Nulla nisl. Nunc nisl.", "", 1992, mangas.StatusDraft, countries.PL)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("e1674245-bb91-4382-adca-4b2c38878a89"), "Recipients Live and Die in L.A.", "Suspendisse potenti. In eleifend quam a odio. In hac habitasse platea dictumst. Maecenas ut massa quis augue luctus tincidunt. Nulla mollis molestie lorem. Quisque ut erat. Curabitur gravida nisi at nibh. In hac habitasse platea dictumst. Aliquam augue quam, sollicitudin vitae, consectetuer eget, rutrum at, lorem.", "", 2012, mangas.StatusOnGoing, countries.JP)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("df3be3a1-f02f-4d2e-afe8-83dc61f46839"), "Ace of Hearts", "Integer aliquet, massa id lobortis convallis, tortor risus dapibus augue, vel accumsan tellus nisi eu orci. Mauris lacinia sapien quis libero. Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh.", "", 2012, mangas.StatusDropped, countries.CN)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("2aa478df-9f0f-4e67-b652-f9b01023eefb"), "Homeboy", "Maecenas leo odio, condimentum id, luctus nec, molestie sed, justo. Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti.", "", 2013, mangas.StatusDraft, countries.CN)),
				},
				Total: 8,
			},
			wantErr: false,
		},
		{
			name: "User present and get last half",
			args: args{
				userId: "dc4402e4-0f88-400a-978e-8bb3880ab063",
				param: repository.QueryParameter{
					Offset: 3,
					Limit:  4,
				},
			},
			want: repository.PagedQueryResult[[]mangas.MangaHistory]{
				Data: []mangas.MangaHistory{
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("2aa478df-9f0f-4e67-b652-f9b01023eefb"), "Homeboy", "Maecenas leo odio, condimentum id, luctus nec, molestie sed, justo. Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti.", "", 2013, mangas.StatusDraft, countries.CN)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("35d1bea2-1a13-45e7-a08c-5d35db26444d"), "Freddy vs. Jason", "Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla dapibus dolor vel est. Donec odio justo, sollicitudin ut, suscipit a, feugiat et, eros. Vestibulum ac est lacinia nisi venenatis tristique. Fusce congue, diam id ornare imperdiet, sapien urna pretium nisl, ut volutpat sapien arcu sed augue. Aliquam erat volutpat. In congue. Etiam justo.", "", 1995, mangas.StatusDraft, countries.MY)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("fc1bea74-5fde-4cf0-a332-c957c914d121"), "Red Riding: 1974", "Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti. Nullam porttitor lacus at turpis. Donec posuere metus vitae ipsum.", "", 2008, mangas.StatusOnGoing, countries.PT)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("62c950be-858b-42f2-8799-a09e49bc8589"), "Satanas", "Nam ultrices, libero non mattis pulvinar, nulla pede ullamcorper augue, a suscipit nulla elit ac nulla. Sed vel enim sit amet nunc viverra dapibus. Nulla suscipit ligula in lacus. Curabitur at ipsum ac tellus semper interdum. Mauris ullamcorper purus sit amet nulla.", "", 2010, mangas.StatusOnGoing, countries.BD)),
				},
				Total: 8,
			},
			wantErr: false,
		},
		{
			name: "Non-existent user",
			args: args{
				userId: uuid.NewString(),
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  0,
				},
			},
			want: repository.PagedQueryResult[[]mangas.MangaHistory]{
				Data:  nil,
				Total: 0,
			},
			wantErr: true,
		},
		{
			name: "Bad UUID",
			args: args{
				userId: "asdasdubas-abd8219d1",
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  0,
				},
			},
			want: repository.PagedQueryResult[[]mangas.MangaHistory]{
				Data:  nil,
				Total: 0,
			},
			wantErr: true,
		},
		{
			name: "User present and out of bound offset",
			args: args{
				userId: "dc4402e4-0f88-400a-978e-8bb3880ab063",
				param: repository.QueryParameter{
					Offset: 8,
					Limit:  5,
				},
			},
			want: repository.PagedQueryResult[[]mangas.MangaHistory]{
				Data:  nil,
				Total: 8,
			},
			wantErr: true,
		},
		{
			name: "User present and out of bound limit",
			args: args{
				userId: "dc4402e4-0f88-400a-978e-8bb3880ab063",
				param: repository.QueryParameter{
					Offset: 5,
					Limit:  9,
				},
			},
			want: repository.PagedQueryResult[[]mangas.MangaHistory]{
				Data: []mangas.MangaHistory{
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("fc1bea74-5fde-4cf0-a332-c957c914d121"), "Red Riding: 1974", "Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti. Nullam porttitor lacus at turpis. Donec posuere metus vitae ipsum.", "", 2008, mangas.StatusOnGoing, countries.PT)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("62c950be-858b-42f2-8799-a09e49bc8589"), "Satanas", "Nam ultrices, libero non mattis pulvinar, nulla pede ullamcorper augue, a suscipit nulla elit ac nulla. Sed vel enim sit amet nunc viverra dapibus. Nulla suscipit ligula in lacus. Curabitur at ipsum ac tellus semper interdum. Mauris ullamcorper purus sit amet nulla.", "", 2010, mangas.StatusOnGoing, countries.BD)),
					newMangaHistoryForTest(time.Now(), newMangaForTest(opt.New("b8bd3f1e-36e3-4033-8290-c5e0caaeab6d"), "Count Three and Pray", "Suspendisse potenti. In eleifend quam a odio. In hac habitasse platea dictumst. Maecenas ut massa quis augue luctus tincidunt. Nulla mollis molestie lorem. Quisque ut erat. Curabitur gravida nisi at nibh. In hac habitasse platea dictumst. Aliquam augue quam, sollicitudin vitae, consectetuer eget, rutrum at, lorem.", "", 2003, mangas.StatusOnGoing, countries.LT)),
				},
				Total: 8,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mangaRepo := NewManga(Db)
		t.Run(tt.name, func(t *testing.T) {
			got, err := mangaRepo.FindMangaHistories(tt.args.userId, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindMangaHistories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got.Data) != len(tt.want.Data) {
				t.Errorf("FindMangaHistories() error = len doesn't match")
				return
			}

			for i := 0; i < len(got.Data); i++ {
				got.Data[i].LastView = tt.want.Data[i].LastView
				got.Data[i].Manga.UpdatedAt = tt.want.Data[i].Manga.UpdatedAt
				got.Data[i].Manga.CreatedAt = tt.want.Data[i].Manga.CreatedAt
				// Ignore relation fields
				got.Data[i].Manga.Genres = tt.want.Data[i].Manga.Genres
				got.Data[i].Manga.Comments = tt.want.Data[i].Manga.Comments
				got.Data[i].Manga.Ratings = tt.want.Data[i].Manga.Ratings
				got.Data[i].Manga.Translations = tt.want.Data[i].Manga.Translations
				got.Data[i].Manga.Volumes = tt.want.Data[i].Manga.Volumes
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindMangaHistories() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mangaRepository_FindMangasByFilter(t *testing.T) {
	type args struct {
		filter *mangas.SearchFilter
		param  repository.QueryParameter
	}
	tests := []struct {
		name    string
		args    args
		want    repository.PagedQueryResult[[]mangas.Manga]
		wantErr bool
	}{
		{
			name: "Manga exist using title",
			args: args{
				filter: &mangas.SearchFilter{
					Title: "Home",
				},
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  0,
				},
			},
			want: repository.PagedQueryResult[[]mangas.Manga]{
				Data: []mangas.Manga{
					*newMangaForTest(opt.New("2aa478df-9f0f-4e67-b652-f9b01023eefb"), "Homeboy", "Maecenas leo odio, condimentum id, luctus nec, molestie sed, justo. Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti.", "", 2013, mangas.StatusDraft, countries.CN),
					*newMangaForTest(opt.New("8d2c10c0-60d1-4b7d-85c0-dc19eeb5316c"), "Home Alone 3", "Vestibulum ac est lacinia nisi venenatis tristique. Fusce congue, diam id ornare imperdiet, sapien urna pretium nisl, ut volutpat sapien arcu sed augue. Aliquam erat volutpat. In congue. Etiam justo. Etiam pretium iaculis justo. In hac habitasse platea dictumst. Etiam faucibus cursus urna.", "", 2005, mangas.StatusDraft, countries.SY),
					*newMangaForTest(opt.New("67e7e795-4999-4123-93e2-c3a768849cf0"), "Home Fries", "Phasellus sit amet erat. Nulla tempus. Vivamus in felis eu sapien cursus vestibulum. Proin eu mi. Nulla ac enim. In tempor, turpis nec euismod scelerisque, quam turpis adipiscing lorem, vitae mattis nibh ligula nec sem.", "", 1998, mangas.StatusDropped, countries.HR),
				},
				Total: 3,
			},
			wantErr: false,
		},
		{
			name: "Manga exist using title and include genre and AND operation",
			args: args{
				filter: &mangas.SearchFilter{
					Title: "home",
					Genres: common.CriterionOption[string]{
						Include: []string{
							"harem",
							"mecha",
							"martial arts",
						},
						IsAndOperation: true,
					},
				},
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  0,
				},
			},
			want: repository.PagedQueryResult[[]mangas.Manga]{
				Data: []mangas.Manga{
					*newMangaForTest(opt.New("8d2c10c0-60d1-4b7d-85c0-dc19eeb5316c"), "Home Alone 3", "Vestibulum ac est lacinia nisi venenatis tristique. Fusce congue, diam id ornare imperdiet, sapien urna pretium nisl, ut volutpat sapien arcu sed augue. Aliquam erat volutpat. In congue. Etiam justo. Etiam pretium iaculis justo. In hac habitasse platea dictumst. Etiam faucibus cursus urna.", "", 2005, mangas.StatusDraft, countries.SY),
				},
				Total: 1,
			},
			wantErr: false,
		},
		{
			name: "Manga exist using title and exclude genre and OR operation",
			args: args{
				filter: &mangas.SearchFilter{
					Title: "home",
					Genres: common.CriterionOption[string]{
						Include: []string{
							"superpower",
							"mecha",
							"sports",
						},
						Exclude: []string{
							"yuri",
							"harem",
						},
						IsAndOperation: false,
					},
				},
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  0,
				},
			},
			want: repository.PagedQueryResult[[]mangas.Manga]{
				Data: []mangas.Manga{
					*newMangaForTest(opt.New("67e7e795-4999-4123-93e2-c3a768849cf0"), "Home Fries", "Phasellus sit amet erat. Nulla tempus. Vivamus in felis eu sapien cursus vestibulum. Proin eu mi. Nulla ac enim. In tempor, turpis nec euismod scelerisque, quam turpis adipiscing lorem, vitae mattis nibh ligula nec sem.", "", 1998, mangas.StatusDropped, countries.HR),
				},
				Total: 1,
			},
			wantErr: false,
		},
		{
			name: "Manga exist using title and include origin",
			args: args{
				filter: &mangas.SearchFilter{
					Title: "ho",
					Origins: []common.Country{
						common.NewCountry(countries.PL),
						common.NewCountry(countries.GR),
					},
					IsOriginInclude: true,
				},
				param: repository.QueryParameter{
					Offset: 1,
					Limit:  0,
				},
			},
			want: repository.PagedQueryResult[[]mangas.Manga]{
				Data: []mangas.Manga{
					*newMangaForTest(opt.New("c5ef679a-6210-4f90-ba05-076f6cb9ec38"), "Girl Who Talked to Dolphins, The", "Nulla ac enim. In tempor, turpis nec euismod scelerisque, quam turpis adipiscing lorem, vitae mattis nibh ligula nec sem. Duis aliquam convallis nunc. Proin at turpis a pede posuere nonummy. Integer non velit. Donec diam neque, vestibulum eget, vulputate ut, ultrices vel, augue. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Donec pharetra, magna vestibulum aliquet ultrices, erat tortor sollicitudin mi, sit amet lobortis sapien sapien non mi.", "", 1993, mangas.StatusDraft, countries.GR),
				},
				Total: 2,
			},
			wantErr: false,
		},
		{
			name: "Manga exist using title and exclude origin",
			args: args{
				filter: &mangas.SearchFilter{
					Title: "ho",
					Origins: []common.Country{
						common.NewCountry(countries.JP),
						common.NewCountry(countries.ID),
						common.NewCountry(countries.CN),
						common.NewCountry(countries.FR),
						common.NewCountry(countries.RU),
						common.NewCountry(countries.PH),
						common.NewCountry(countries.BR),
						common.NewCountry(countries.PL),
						common.NewCountry(countries.TH),
						common.NewCountry(countries.SE),
					},
					IsOriginInclude: false,
				},
				param: repository.QueryParameter{
					Offset: 2,
					Limit:  3,
				},
			},
			want: repository.PagedQueryResult[[]mangas.Manga]{
				Data: []mangas.Manga{
					*newMangaForTest(opt.New("14156167-632a-4cea-bce4-9b85b847e744"), "Last Holiday", "Suspendisse potenti. Nullam porttitor lacus at turpis. Donec posuere metus vitae ipsum. Aliquam non mauris. Morbi non lectus. Aliquam sit amet diam in magna bibendum imperdiet. Nullam orci pede, venenatis non, sodales sed, tincidunt eu, felis. Fusce posuere felis sed lacus. Morbi sem mauris, laoreet ut, rhoncus aliquet, pulvinar sed, nisl. Nunc rhoncus dui vel sem. Sed sagittis. Nam congue, risus semper porta volutpat, quam pede lobortis ligula, sit amet eleifend pede libero quis orci. Nullam molestie nibh in lectus. Pellentesque at nulla.", "", 1994, mangas.StatusHiatus, countries.CA),
					*newMangaForTest(opt.New("8d2c10c0-60d1-4b7d-85c0-dc19eeb5316c"), "Home Alone 3", "Vestibulum ac est lacinia nisi venenatis tristique. Fusce congue, diam id ornare imperdiet, sapien urna pretium nisl, ut volutpat sapien arcu sed augue. Aliquam erat volutpat. In congue. Etiam justo. Etiam pretium iaculis justo. In hac habitasse platea dictumst. Etiam faucibus cursus urna.", "", 2005, mangas.StatusDraft, countries.SY),
					*newMangaForTest(opt.New("67e7e795-4999-4123-93e2-c3a768849cf0"), "Home Fries", "Phasellus sit amet erat. Nulla tempus. Vivamus in felis eu sapien cursus vestibulum. Proin eu mi. Nulla ac enim. In tempor, turpis nec euismod scelerisque, quam turpis adipiscing lorem, vitae mattis nibh ligula nec sem.", "", 1998, mangas.StatusDropped, countries.HR),
				},
				Total: 6,
			},
			wantErr: false,
		},
		{
			name: "Manga exists using title, AND genre and Include origin",
			args: args{
				filter: &mangas.SearchFilter{
					Title: "HO",
					Genres: common.CriterionOption[string]{
						Include: []string{
							"isekai",
							"school",
						},
						Exclude: []string{
							"harem",
						},
						IsAndOperation: true,
					},
					Origins: []common.Country{
						common.NewCountry(countries.ID),
					},
					IsOriginInclude: true,
				},
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  5,
				},
			},
			want: repository.PagedQueryResult[[]mangas.Manga]{
				Data: []mangas.Manga{
					*newMangaForTest(opt.New("402d9b14-f14c-4982-b71b-1e1885838f6f"), "Adventures of Robin Hood, The", "Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti. Nullam porttitor lacus at turpis. Donec posuere metus vitae ipsum. Aliquam non mauris. Morbi non lectus.", "", 2010, mangas.StatusCompleted, countries.ID),
				},
				Total: 1,
			},
			wantErr: false,
		},
		{
			name: "Manga exists using title, OR genre and Exclude origin",
			args: args{
				filter: &mangas.SearchFilter{
					Title: "HO",
					Genres: common.CriterionOption[string]{
						Include: []string{
							"isekai",
							"school",
						},
						Exclude: []string{
							"harem",
							"ecchi",
						},
						IsAndOperation: false,
					},
					Origins: []common.Country{
						common.NewCountry(countries.PE),
					},
					IsOriginInclude: false,
				},
				param: repository.QueryParameter{
					Offset: 1,
					Limit:  5,
				},
			},
			want: repository.PagedQueryResult[[]mangas.Manga]{
				Data: []mangas.Manga{
					*newMangaForTest(opt.New("49781073-e4cd-4826-9074-3fc39f51bd34"), "Abbott and Costello in Hollywood", "Nam dui. Proin leo odio, porttitor id, consequat in, consequat ut, nulla. Sed accumsan felis. Ut at dolor quis odio consequat varius. Integer ac leo. Pellentesque ultrices mattis odio.", "", 1997, mangas.StatusCompleted, countries.TH),
					*newMangaForTest(opt.New("402d9b14-f14c-4982-b71b-1e1885838f6f"), "Adventures of Robin Hood, The", "Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti. Nullam porttitor lacus at turpis. Donec posuere metus vitae ipsum. Aliquam non mauris. Morbi non lectus.", "", 2010, mangas.StatusCompleted, countries.ID),
					*newMangaForTest(opt.New("67e7e795-4999-4123-93e2-c3a768849cf0"), "Home Fries", "Phasellus sit amet erat. Nulla tempus. Vivamus in felis eu sapien cursus vestibulum. Proin eu mi. Nulla ac enim. In tempor, turpis nec euismod scelerisque, quam turpis adipiscing lorem, vitae mattis nibh ligula nec sem.", "", 1998, mangas.StatusDropped, countries.HR),
				},
				Total: 4,
			},
			wantErr: false,
		},
		{
			name: "Manga not-exist by title",
			args: args{
				filter: &mangas.SearchFilter{
					Title: "iasmdiasdaisd",
				},
			},
			want: repository.PagedQueryResult[[]mangas.Manga]{
				Data:  nil,
				Total: 0,
			},
			wantErr: true,
		},
		{
			name: "Manga not-exists using include genre",
			args: args{
				filter: &mangas.SearchFilter{
					Genres: common.CriterionOption[string]{
						Include: []string{
							"isekai",
							"school",
							"drama",
							"romance",
						},
						IsAndOperation: true,
					},
				},
			},
			want: repository.PagedQueryResult[[]mangas.Manga]{
				Data:  nil,
				Total: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		mangaRepo := NewManga(Db)
		t.Run(tt.name, func(t *testing.T) {
			got, err := mangaRepo.FindMangasByFilter(tt.args.filter, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindMangasByFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got.Data) != len(tt.want.Data) {
				t.Errorf("FindMangasByFilter() error = length doesn't match")
				return
			}

			// Ignore time fields
			for i := 0; i < len(got.Data); i++ {
				got.Data[i].UpdatedAt = tt.want.Data[i].UpdatedAt
				got.Data[i].CreatedAt = tt.want.Data[i].CreatedAt
				// Ignore relation fields
				got.Data[i].Genres = tt.want.Data[i].Genres
				//got.Data[i].Chapters = tt.want.Data[i].Chapters
				got.Data[i].Comments = tt.want.Data[i].Comments
				got.Data[i].Ratings = tt.want.Data[i].Ratings
				got.Data[i].Translations = tt.want.Data[i].Translations
				got.Data[i].Volumes = tt.want.Data[i].Volumes
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindMangasByFilter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mangaRepository_FindMangasById(t *testing.T) {
	type args struct {
		ids []string
	}
	tests := []struct {
		name    string
		args    args
		want    []mangas.Manga
		wantErr bool
	}{
		{
			name: "Normal",
			args: args{
				ids: []string{
					"19382f54-1da7-4cb7-807d-9f6030bb121e",
					"2aa478df-9f0f-4e67-b652-f9b01023eefb",
					"35d1bea2-1a13-45e7-a08c-5d35db26444d",
					"4ab94eda-46ae-4de8-bd3a-734b388e06fc",
					"df3be3a1-f02f-4d2e-afe8-83dc61f46839",
				},
			},
			want: []mangas.Manga{
				*newMangaForTest(opt.New("19382f54-1da7-4cb7-807d-9f6030bb121e"), "David and Bathsheba", "Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla dapibus dolor vel est. Donec odio justo, sollicitudin ut, suscipit a, feugiat et, eros. Vestibulum ac est lacinia nisi venenatis tristique. Fusce congue, diam id ornare imperdiet, sapien urna pretium nisl, ut volutpat sapien arcu sed augue. Aliquam erat volutpat.", "", 1998, mangas.StatusHiatus, countries.BY),
				*newMangaForTest(opt.New("2aa478df-9f0f-4e67-b652-f9b01023eefb"), "Homeboy", "Maecenas leo odio, condimentum id, luctus nec, molestie sed, justo. Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti.", "", 2013, mangas.StatusDraft, countries.CN),
				*newMangaForTest(opt.New("35d1bea2-1a13-45e7-a08c-5d35db26444d"), "Freddy vs. Jason", "Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla dapibus dolor vel est. Donec odio justo, sollicitudin ut, suscipit a, feugiat et, eros. Vestibulum ac est lacinia nisi venenatis tristique. Fusce congue, diam id ornare imperdiet, sapien urna pretium nisl, ut volutpat sapien arcu sed augue. Aliquam erat volutpat. In congue. Etiam justo.", "", 1995, mangas.StatusDraft, countries.MY),
				*newMangaForTest(opt.New("4ab94eda-46ae-4de8-bd3a-734b388e06fc"), "Just Friends", "Duis at velit eu est congue elementum. In hac habitasse platea dictumst. Morbi vestibulum, velit id pretium iaculis, diam erat fermentum justo, nec condimentum neque sapien placerat ante. Nulla justo. Aliquam quis turpis eget elit sodales scelerisque.", "", 2003, mangas.StatusDropped, countries.PT),
				*newMangaForTest(opt.New("df3be3a1-f02f-4d2e-afe8-83dc61f46839"), "Ace of Hearts", "Integer aliquet, massa id lobortis convallis, tortor risus dapibus augue, vel accumsan tellus nisi eu orci. Mauris lacinia sapien quis libero. Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh.", "", 2012, mangas.StatusDropped, countries.CN),
			},
			wantErr: false,
		},
		{
			name: "Bad uuid",
			args: args{
				ids: []string{
					"asdadasdasd",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Bad uuid inside",
			args: args{
				ids: []string{
					"19382f54-1da7-4cb7-807d-9f6030bb121e",
					"2aa478df-9f0f-4e67-b652-f9b01023eefb",
					"asdadasdasd",
					"4ab94eda-46ae-4de8-bd3a-734b388e06fc",
					"df3be3a1-f02f-4d2e-afe8-83dc61f46839",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "all manga not found",
			args: args{
				ids: []string{
					uuid.NewString(),
					uuid.NewString(),
					uuid.NewString(),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "some manga id is not found",
			args: args{
				ids: []string{
					"19382f54-1da7-4cb7-807d-9f6030bb121e",
					"2aa478df-9f0f-4e67-b652-f9b01023eefb",
					uuid.NewString(),
				},
			},
			want: []mangas.Manga{
				*newMangaForTest(opt.New("19382f54-1da7-4cb7-807d-9f6030bb121e"), "David and Bathsheba", "Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla dapibus dolor vel est. Donec odio justo, sollicitudin ut, suscipit a, feugiat et, eros. Vestibulum ac est lacinia nisi venenatis tristique. Fusce congue, diam id ornare imperdiet, sapien urna pretium nisl, ut volutpat sapien arcu sed augue. Aliquam erat volutpat.", "", 1998, mangas.StatusHiatus, countries.BY),
				*newMangaForTest(opt.New("2aa478df-9f0f-4e67-b652-f9b01023eefb"), "Homeboy", "Maecenas leo odio, condimentum id, luctus nec, molestie sed, justo. Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti.", "", 2013, mangas.StatusDraft, countries.CN),
			},
			wantErr: false,
		},
		{
			name: "duplicate manga id",
			args: args{
				ids: []string{
					"19382f54-1da7-4cb7-807d-9f6030bb121e",
					"2aa478df-9f0f-4e67-b652-f9b01023eefb",
					"2aa478df-9f0f-4e67-b652-f9b01023eefb",
				},
			},
			want: []mangas.Manga{
				*newMangaForTest(opt.New("19382f54-1da7-4cb7-807d-9f6030bb121e"), "David and Bathsheba", "Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla dapibus dolor vel est. Donec odio justo, sollicitudin ut, suscipit a, feugiat et, eros. Vestibulum ac est lacinia nisi venenatis tristique. Fusce congue, diam id ornare imperdiet, sapien urna pretium nisl, ut volutpat sapien arcu sed augue. Aliquam erat volutpat.", "", 1998, mangas.StatusHiatus, countries.BY),
				*newMangaForTest(opt.New("2aa478df-9f0f-4e67-b652-f9b01023eefb"), "Homeboy", "Maecenas leo odio, condimentum id, luctus nec, molestie sed, justo. Pellentesque viverra pede ac diam. Cras pellentesque volutpat dui. Maecenas tristique, est et tempus semper, est quam pharetra magna, ac consequat metus sapien ut nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris viverra diam vitae quam. Suspendisse potenti.", "", 2013, mangas.StatusDraft, countries.CN),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		m := NewManga(Db)
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.FindMangasById(tt.args.ids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindMangasById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Len(t, got, len(tt.want), "expected got and want to have same length")

			for i := 0; i < len(got); i++ {
				// Ignore time fields
				got[i].CreatedAt = tt.want[i].CreatedAt
				got[i].UpdatedAt = tt.want[i].UpdatedAt
				// Ignore relation fields
				got[i].Genres = tt.want[i].Genres
				//got[i].Chapters = tt.want[i].Chapters
				got[i].Comments = tt.want[i].Comments
				got[i].Ratings = tt.want[i].Ratings
				got[i].Translations = tt.want[i].Translations
				got[i].Volumes = tt.want[i].Volumes
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindMangasById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mangaRepository_FindRandomMangas(t *testing.T) {
	type args struct {
		limit uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []mangas.Manga
		wantErr bool
	}{
		{
			name: "0 Limit, Take All",
			args: args{
				limit: 0,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "10 Limit",
			args: args{
				limit: 10,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		m := NewManga(Db)
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.FindRandomMangas(tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindRandomMangas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.args.limit == 0 {
				require.GreaterOrEqual(t, len(got), int(tt.args.limit))
			} else {
				require.Len(t, got, int(tt.args.limit))
			}
		})
	}
}

func Test_mangaRepository_ListMangas(t *testing.T) {
	type args struct {
		param repository.QueryParameter
	}
	tests := []struct {
		name    string
		args    args
		want    repository.PagedQueryResult[[]mangas.Manga]
		wantErr bool
	}{
		{
			name: "Offset 0 and Limit 10",
			args: args{
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  10,
				},
			},
			want: repository.PagedQueryResult[[]mangas.Manga]{
				Data:  nil,
				Total: 0,
			},
			wantErr: false,
		},
		{
			name: "Offset 5 and Limit 25",
			args: args{
				param: repository.QueryParameter{
					Offset: 5,
					Limit:  25,
				},
			},
			want:    repository.PagedQueryResult[[]mangas.Manga]{},
			wantErr: false,
		},
		{
			name: "Offset 0 and Limit 0, Take All",
			args: args{
				param: repository.QueryParameter{
					Offset: 0,
					Limit:  0,
				},
			},
			want:    repository.PagedQueryResult[[]mangas.Manga]{},
			wantErr: false,
		},
		{
			name: "Offset 5 and Limit 0, Take All from index 5",
			args: args{
				param: repository.QueryParameter{
					Offset: 5,
					Limit:  0,
				},
			},
			want:    repository.PagedQueryResult[[]mangas.Manga]{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		m := NewManga(Db)
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.ListMangas(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListMangas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.args.param.Limit == 0 {
				require.Len(t, got.Data, int(got.Total-tt.args.param.Offset))
			} else {
				require.Len(t, got.Data, int(tt.args.param.Limit))
			}

		})
	}
}

func Test_mangaRepository_EditManga(t *testing.T) {
	type args struct {
		manga *mangas.Manga
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{},
	}
	for _, tt := range tests {
		tx, err := Db.Begin()
		require.NoError(t, err)
		m := NewManga(tx)

		t.Run(tt.name, func(t *testing.T) {
			defer func(tx bun.Tx) {
				require.NoError(t, tx.Rollback())
			}(tx)

			tt.wantErr(t, m.EditManga(tt.args.manga), fmt.Sprintf("EditManga(%v)", tt.args.manga))
		})
	}
}
