package application

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/EkaterinaShamanaeva/autorization/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

type app struct {
	ctx   context.Context
	repo  *repository.Repository
	cache map[string]repository.User
}

func (a app) Routes(r *httprouter.Router) {
	r.ServeFiles("/public/*filepath", http.Dir("public"))
	r.GET("/", a.authorized(a.StartPage))
	r.GET("/login", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		a.LoginPage(rw, "")
	})
	// отправка логина и пароля при нажатии кнопки "войти"
	r.POST("/login", a.Login)
}

// middleware
func (a app) authorized(next httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		token, err := readCookie("token", r)
		if err != nil {
			http.Redirect(rw, r, "/login", http.StatusSeeOther)
			return
		}
		if _, ok := a.cache[token]; !ok {
			http.Redirect(rw, r, "/login", http.StatusSeeOther)
			return
		}
		next(rw, r, p)
	}
}

func readCookie(name string, r *http.Request) (value string, err error) {
	if name == "" {
		return value, errors.New("you are trying to read empty cookie")
	}
	cookie, err := r.Cookie(name)
	if err != nil {
		return value, err
	}
	str := cookie.Value
	value, _ = url.QueryUnescape(str)
	return value, err
}

func (a app) StartPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	motivation, err := a.repo.GetRandomMotivation(a.ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	lp := filepath.Join("public", "html", "motivation.html")
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	err = tmpl.ExecuteTemplate(rw, "motivation", motivation)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func NewApp(ctx context.Context, dbpool *pgxpool.Pool) *app {
	return &app{ctx, repository.NewRepository(dbpool), make(map[string]repository.User)}
}

func (a app) LoginPage(rw http.ResponseWriter, message string) {
	lp := filepath.Join("public", "html", "login.html")
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	type answer struct {
		Message string
	}
	data := answer{message}
	err = tmpl.ExecuteTemplate(rw, "login", data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) Login(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	if login == "" || password == "" {
		a.LoginPage(rw, "Необходимо указать логин и пароль!")
		return
	}
	// создаем хеш пароля
	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])
	// авторизация -> получаем юзера
	user, err := a.repo.Login(a.ctx, login, hashedPass)
	if err != nil {
		a.LoginPage(rw, "Вы ввели неверный логин или пароль!")
		return
	}
	// генерируем токен, пишем его в кеш и куки
	time64 := time.Now().Unix()
	timeInt := string(time64)
	token := login + password + timeInt
	hashToken := md5.Sum([]byte(token))
	hashedToken := hex.EncodeToString(hashToken[:])
	// записываем токен и пользователя в хеш
	a.cache[hashedToken] = user
	// определяем время жизни куки
	livingTime := 60 * time.Minute
	expiration := time.Now().Add(livingTime)
	// создаем куку с токеном
	cookie := http.Cookie{Name: "token", Value: url.QueryEscape(hashedToken), Expires: expiration}
	http.SetCookie(rw, &cookie)
	http.Redirect(rw, r, "/", http.StatusSeeOther)
}
