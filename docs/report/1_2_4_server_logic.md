### Разработка слоя серверной логики веб-приложения

Слой серверной логики в веб-приложении «Анализатор статического кода» выполняет ключевую роль. Он обрабатывает входящие HTTP-запросы, реализует бизнес-логику приложения, а также обеспечивает взаимодействие между базой данных и пользовательским интерфейсом. Основу данного слоя составляют контроллеры, каждый из которых отвечает за отдельный функциональный блок системы.

В качестве веб-сервера в приложении используется встроенный HTTP-сервер Go, который обеспечивает приём входящих HTTP-запросов и их обработку. Благодаря этому запросы пользователей проходят через стандартный интерфейс Go net/http, где происходит разбор URL-адреса, обработка параметров запроса и отдача статических файлов, а затем перенаправляются в маршрутизатор приложения.

Центральным элементом серверной логики является маршрутизатор на основе библиотеки gorilla/mux. Он анализирует URL-адрес запроса, сопоставляет его с таблицей маршрутов и вызывает соответствующий метод контроллера. Такой механизм позволяет чётко распределить запросы между различными частями приложения. Фрагмент кода настройки маршрутизатора представлен на листинге 2.4.1.

Листинг 2.4.1 – Фрагмент кода настройки маршрутизатора

```go
r := mux.NewRouter()

// Public routes
r.HandleFunc("/login", authCtrl.ShowLogin).Methods("GET")
r.HandleFunc("/login", authCtrl.HandleLogin).Methods("POST")
r.HandleFunc("/register", authCtrl.ShowRegister).Methods("GET")
r.HandleFunc("/register", authCtrl.HandleRegister).Methods("POST")
r.HandleFunc("/", homeCtrl.Index).Methods("GET")
r.HandleFunc("/about", aboutCtrl.Index).Methods("GET")

// Protected routes
r.HandleFunc("/upload", middleware.RequireAuth(uploadCtrl.ShowUpload)).Methods("GET")
r.HandleFunc("/upload/zip", middleware.RequireAuth(uploadCtrl.HandleZipUpload)).Methods("POST")
r.HandleFunc("/analyses", middleware.RequireAuth(analysesCtrl.List)).Methods("GET")
r.HandleFunc("/analyses/{id}", middleware.RequireAuth(analysesCtrl.Details)).Methods("GET")
```

Контроллер HomeController отвечает за отображение главной страницы приложения. Он передаёт данные в представление без дополнительной обработки, так как главная страница содержит только статическую информацию о функциональности сервиса. Контроллер HomeController представлен на листинге 2.4.2.

Листинг 2.4.2 – Контроллер HomeController

```go
package controllers

import (
	"html/template"
	"net/http"
)

type HomeController struct {
	tmpl *template.Template
}

func NewHomeController(tmpl *template.Template) *HomeController {
	return &HomeController{tmpl: tmpl}
}

func (c *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	if err := executeTemplate(w, r, c.tmpl, "home", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
```

Контроллер AuthController отвечает за авторизацию пользователей. Он реализует механизм регистрации новых пользователей, входа в систему и выхода. Контроллер использует репозиторий UserRepository для работы с базой данных и библиотеку bcrypt для хеширования паролей. При успешной авторизации создаётся сессия, которая хранит идентификатор пользователя. Фрагмент кода контроллера AuthController представлен на листинге 2.4.3.

Листинг 2.4.3 – Фрагмент кода контроллера AuthController

```go
func (c *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := c.userRepo.GetByUsername(username)
	if err != nil {
		// Обработка ошибки
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// Обработка ошибки
		return
	}

	// Создание сессии
	session, err := middleware.GetSession(r)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	session.Values[middleware.UserIDKey()] = user.ID.String()
	session.Values[middleware.UsernameKey()] = user.Username
	session.Save(r, w)

	http.Redirect(w, r, "/analyses", http.StatusFound)
}
```

Контроллер UploadController отвечает за загрузку и анализ кода. Он обрабатывает два типа запросов: загрузку ZIP-архивов с проектами и анализ фрагментов кода через интерактивный анализатор. Контроллер использует сервисы AnalyzerService и StorageService для обработки файлов, создаёт записи в базе данных через репозитории и запускает анализ кода. Фрагмент кода контроллера UploadController представлен на листинге 2.4.4.

Листинг 2.4.4 – Фрагмент кода контроллера UploadController

```go
func (c *UploadController) HandleZipUpload(w http.ResponseWriter, r *http.Request) {
	const maxSize = 20 * 1024 * 1024 // 20 MB
	r.ParseMultipartForm(maxSize)

	projectName := r.FormValue("project_name")
	file, header, err := r.FormFile("zip_file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Получение userID из сессии
	userIDStr, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, _ := uuid.Parse(userIDStr)

	// Создание проекта и анализа
	project, _ := c.projectRepo.Create(projectName, userID)
	analysis, _ := c.analysisRepo.Create(project.ID, userID, "zip", inputMeta)

	// Извлечение файлов из ZIP
	files, _ := c.storageService.ExtractZip(tmpFile.Name(), "")

	// Сохранение файлов в базу данных
	c.fileRepo.CreateBatch(analysis.ID, fileModels)

	// Запуск анализа
	issues, summary := c.analyzerService.Analyze(files, analysis.ID.String())

	// Сохранение результатов
	c.issueRepo.CreateBatch(analysis.ID, issues)
	c.analysisRepo.UpdateSummary(analysis.ID, summaryJSON)

	http.Redirect(w, r, fmt.Sprintf("/analyses/%s", analysis.ID), http.StatusSeeOther)
}
```

Контроллер AnalysesController отвечает за отображение истории анализов и детальной информации о каждом анализе. Он получает список анализов текущего пользователя из базы данных через репозитории, фильтрует их по userID для обеспечения конфиденциальности данных, и передаёт результаты в представление. Контроллер также реализует функциональность удаления анализов. Фрагмент кода контроллера AnalysesController представлен на листинге 2.4.5.

Листинг 2.4.5 – Фрагмент кода контроллера AnalysesController

```go
func (c *AnalysesController) List(w http.ResponseWriter, r *http.Request) {
	// Получение userID из сессии
	userIDStr, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, _ := uuid.Parse(userIDStr)

	// Получение анализов только для текущего пользователя
	analyses, err := c.analysisRepo.ListByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Формирование данных для представления
	var data ListData
	for _, analysis := range analyses {
		project, _ := c.projectRepo.GetByID(analysis.ProjectID)
		data.Analyses = append(data.Analyses, AnalysisWithProject{
			Analysis: analysis,
			Project:  project,
		})
	}

	executeTemplate(w, r, c.tmpl, "analyses", data)
}
```

