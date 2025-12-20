package controllers

import (
	"html/template"
	"net/http"
	"webapp/internal/app/middleware"
	"webapp/internal/app/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	userRepo *repositories.UserRepository
	tmpl     *template.Template
}

func NewAuthController(userRepo *repositories.UserRepository, tmpl *template.Template) *AuthController {
	return &AuthController{
		userRepo: userRepo,
		tmpl:     tmpl,
	}
}

func (c *AuthController) ShowLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect if already logged in
	if userID, ok := middleware.GetUserID(r); ok && userID != "" {
		http.Redirect(w, r, "/analyses", http.StatusFound)
		return
	}

	data := map[string]interface{}{
		"Title": "Вход",
	}
	if err := executeTemplate(w, r, c.tmpl, "login", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *AuthController) ShowRegister(w http.ResponseWriter, r *http.Request) {
	// Redirect if already logged in
	if userID, ok := middleware.GetUserID(r); ok && userID != "" {
		http.Redirect(w, r, "/analyses", http.StatusFound)
		return
	}

	data := map[string]interface{}{
		"Title": "Регистрация",
	}
	if err := executeTemplate(w, r, c.tmpl, "register", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		data := map[string]interface{}{
			"Title":   "Вход",
			"Error":   "Необходимо указать имя пользователя и пароль",
			"Username": username,
		}
		executeTemplate(w, r, c.tmpl, "login", data)
		return
	}

	user, err := c.userRepo.GetByUsername(username)
	if err != nil {
		data := map[string]interface{}{
			"Title":   "Вход",
			"Error":   "Неверное имя пользователя или пароль",
			"Username": username,
		}
		executeTemplate(w, r, c.tmpl, "login", data)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		data := map[string]interface{}{
			"Title":   "Вход",
			"Error":   "Неверное имя пользователя или пароль",
			"Username": username,
		}
		executeTemplate(w, r, c.tmpl, "login", data)
		return
	}

	// Create session
	session, err := middleware.GetSession(r)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	session.Values[middleware.UserIDKey()] = user.ID.String()
	session.Values[middleware.UsernameKey()] = user.Username
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/analyses", http.StatusFound)
}

func (c *AuthController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if username == "" || password == "" {
		data := map[string]interface{}{
			"Title":   "Регистрация",
			"Error":   "Необходимо указать имя пользователя и пароль",
			"Username": username,
		}
		executeTemplate(w, r, c.tmpl, "register", data)
		return
	}

	if password != confirmPassword {
		data := map[string]interface{}{
			"Title":   "Регистрация",
			"Error":   "Пароли не совпадают",
			"Username": username,
		}
		executeTemplate(w, r, c.tmpl, "register", data)
		return
	}

	if len(password) < 6 {
		data := map[string]interface{}{
			"Title":   "Регистрация",
			"Error":   "Пароль должен содержать не менее 6 символов",
			"Username": username,
		}
		executeTemplate(w, r, c.tmpl, "register", data)
		return
	}

	// Check if user already exists
	_, err := c.userRepo.GetByUsername(username)
	if err == nil {
		data := map[string]interface{}{
			"Title":   "Регистрация",
			"Error":   "Пользователь с таким именем уже существует",
			"Username": username,
		}
		executeTemplate(w, r, c.tmpl, "register", data)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create user
	user, err := c.userRepo.Create(username, string(hashedPassword))
	if err != nil {
		data := map[string]interface{}{
			"Title":   "Регистрация",
			"Error":   "Не удалось создать пользователя. Попробуйте другое имя.",
			"Username": username,
		}
		executeTemplate(w, r, c.tmpl, "register", data)
		return
	}

	// Create session
	session, err := middleware.GetSession(r)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	session.Values[middleware.UserIDKey()] = user.ID.String()
	session.Values[middleware.UsernameKey()] = user.Username
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/analyses", http.StatusFound)
}

func (c *AuthController) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.GetSession(r)
	if err == nil {
		session.Values = make(map[interface{}]interface{})
		session.Options.MaxAge = -1
		session.Save(r, w)
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

