package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dinever/dingo/app/model"
	"github.com/dinever/dingo/app/utils"
	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
)


func mockPost() *model.Post {
	p := model.NewPost()
	p.Title = "Welcome to Dingo!"
	p.Slug = "welcome-to-dingo"
	p.Markdown = "sample content"
	p.Html = utils.Markdown2Html(p.Markdown)
	p.Tags = model.GenerateTagsFromCommaString("Welcome, Dingo")
	p.AllowComment = true
	p.Category = ""
	p.CreatedBy = 0
	p.UpdatedBy = 0
	p.IsPublished = false
	p.IsPage = false
	p.Author = &model.User{Id: 0, Name: "Dingo User", Email: "example@example.com"}
	return p
}

func TestPost(t *testing.T) {
	Convey("Initialize database", t, func() {
		testDB := fmt.Sprintf(filepath.Join(os.TempDir(), "ding-testdb-%s"), fmt.Sprintf(time.Now().Format("20060102T150405.000")))
		model.Initialize(testDB, true)

		Convey("When the post is not found", func() {
			ctx := mockContext(nil, "GET", "/someslug/")
			ctx.App.ServeHTTP(ctx.Response, ctx.Request)

			So(ctx.Response.(*httptest.ResponseRecorder).Code, ShouldEqual, 404)
		})

		Convey("When the post is not published yet", func() {
			p := mockPost()
			p.Save()

			ctx := mockContext(nil, "GET", "/welcome-to-dingo/")
			ctx.App.ServeHTTP(ctx.Response, ctx.Request)

			So(ctx.Response.(*httptest.ResponseRecorder).Code, ShouldEqual, 404)
		})

		Reset(func() {
			os.Remove(testDB)
		})
	})
}
