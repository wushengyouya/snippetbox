package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.alexedwards.net/internal/models"
	"snippetbox.alexedwards.net/internal/validator"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(int64(id))
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{Expires: 365}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// 审查提交的数据是否能转为Form
	if err := r.ParseForm(); err != nil {
		app.serverError(w, err)
		return
	}
	// 判断是否提交过期时间
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// 构建form
	form := &snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}
	// 校验form的字段是否符合
	form.CheckField(validator.NotBlank(form.Title), "title", "标题不能为空")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "标题长度不能超过100字符")
	form.CheckField(validator.NotBlank(form.Content), "content", "内容不能为空")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "过期时间无效")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	// 插入到数据库
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl", data)
}
