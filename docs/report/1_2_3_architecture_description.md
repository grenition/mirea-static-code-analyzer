### Разработка архитектуры веб-приложения на основе паттерна MVC

При создании веб-приложения «Анализатор статического кода» особое внимание уделялось архитектуре системы, поскольку от её продуманности зависит устойчивость приложения и удобство дальнейшего развития. В качестве основной архитектурной модели был выбран шаблон MVC. Этот подход разделяет логику работы веб-приложения на три независимых слоя, что обеспечивает гибкость структуры и уменьшает связанность компонентов (рисунок 2.3).

![Архитектура на основе паттерна MVC](mermaid/2_3_1_architecture_mvc.md)
Рисунок 2.3.1 – Архитектура на основе паттерна MVC

Модель (Model) отвечает за представление структуры данных и их валидацию. В проекте этот слой реализован в виде структур Go, которые описывают сущности приложения: Project, AnalysisRun, Issue, User, AnalysisFile. Каждая модель содержит поля, соответствующие колонкам в базе данных. Модель Project представлена на листинге 2.3.1.

Листинг 2.3.1 — Модель Project

```go
package models

import (
	"time"
	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID
	Name      string
	UserID    uuid.UUID
	CreatedAt time.Time
}
```

Для работы с базой данных используется дополнительный слой репозиториев (Repository), который инкапсулирует логику выполнения SQL-запросов. Репозитории предоставляют методы для создания, чтения, обновления и удаления данных. Класс ProjectRepository представлен на листинге 2.3.2.

Листинг 2.3.2 – Репозиторий ProjectRepository

```go
package repositories

import (
	"database/sql"
	"time"
	"github.com/google/uuid"
	"webapp/internal/app/models"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(name string, userID uuid.UUID) (*models.Project, error) {
	id := uuid.New()
	now := time.Now()
	
	_, err := r.db.Exec(
		"INSERT INTO projects (id, name, user_id, created_at) VALUES ($1, $2, $3, $4)",
		id, name, userID, now,
	)
	if err != nil {
		return nil, err
	}
	
	return &models.Project{
		ID:        id,
		Name:      name,
		UserID:    userID,
		CreatedAt: now,
	}, nil
}
```

Представление (View) реализовано в виде HTML-шаблонов, которые отвечают за отображение информации пользователю. Представления получают только подготовленные данные, не содержат бизнес-логики и занимаются исключительно выводом HTML-разметки с использованием CSS и Bootstrap. В проекте используется функция executeTemplate, отвечающая за подключение шаблонов и передачу им данных. Функция executeTemplate представлена на листинге 2.3.3.

Листинг 2.3.3. – Функция executeTemplate

```go
func executeTemplate(w http.ResponseWriter, r *http.Request, tmpl *template.Template, 
	pageName string, data interface{}) error {
	pageTemplate := template.New("").Funcs(template.FuncMap{
		"uuid": func() string {
			return uuid.New().String()
		},
	})
	
	layoutPath := "./internal/app/views/templates/layout.html"
	pageTemplate, err := pageTemplate.ParseFiles(layoutPath)
	if err != nil {
		return err
	}
	
	if pageName != "" {
		pagePath := filepath.Join("./internal/app/views/templates", pageName+".html")
		pageTemplate, err = pageTemplate.ParseFiles(pagePath)
		if err != nil {
			return err
		}
	}
	
	dataMap := make(map[string]interface{})
	if data != nil {
		if m, ok := data.(map[string]interface{}); ok {
			dataMap = m
		} else {
			dataMap = map[string]interface{}{"Data": data}
		}
	}
	
	if username, ok := middleware.GetUsername(r); ok {
		dataMap["Username"] = username
		dataMap["IsAuthenticated"] = true
	} else {
		dataMap["IsAuthenticated"] = false
	}
	
	return pageTemplate.ExecuteTemplate(w, "layout.html", dataMap)
}
```

Контроллер (Controller) служит связующим звеном между моделью и представлением. Контроллеры принимают входящие HTTP-запросы, запускают соответствующую бизнес-логику через сервисы и репозитории, а затем передают данные в представления. В проекте этот слой реализован в виде структур с методами, которые обрабатывают HTTP-запросы. Каждый контроллер содержит ссылку на шаблоны и необходимые репозитории или сервисы. Контроллер HomeController представлен на листинге 2.3.4.

Листинг 2.3.4 – Контроллер HomeController

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

Дополнительно в проекте используется слой сервисов (Service), который содержит бизнес-логику приложения. Сервисы работают с репозиториями для доступа к данным и выполняют сложные операции, такие как анализ кода, запуск линтеров через Docker и парсинг результатов. Это позволяет вынести бизнес-логику из контроллеров и обеспечить её переиспользование.

