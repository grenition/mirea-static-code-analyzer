### Разработка слоя логики базы данных

Слой логики базы данных веб-приложения «Анализатор статического кода» отвечает за хранение, структурирование и безопасную обработку данных, связанных с пользователями, проектами, анализами кода и найденными проблемами. В качестве системы управления базами данных используется PostgreSQL, что обеспечивает высокую производительность, поддержку транзакций, работу с бинарными данными и строгую типизацию данных [19].

База данных состоит из пяти основных таблиц (рисунок 2.5.1):
- users – содержит информацию о зарегистрированных пользователях: имя пользователя, хеш пароля и дату регистрации,
- projects – хранит данные о проектах, которые анализируются: название проекта, идентификатор пользователя и дату создания,
- analysis_runs – содержит информацию о каждом выполненном анализе: статус, тип входных данных, метаданные, результаты анализа и даты создания, начала и завершения,
- issues – хранит найденные проблемы в коде: уровень серьёзности, код правила, сообщение, путь к файлу и номер строки,
- analysis_files – сохраняет содержимое проанализированных файлов в бинарном формате для последующего просмотра.

Рисунок 2.5.1 – Диаграмма таблиц базы данных

Таблицы связаны отношениями «один ко многим» [20]: каждый пользователь может иметь множество проектов, каждый проект может содержать множество анализов, каждый анализ может содержать множество найденных проблем и множество файлов. Связи обеспечиваются внешними ключами user_id, project_id и analysis_id с каскадным удалением, что гарантирует целостность данных и автоматическое удаление связанных записей при удалении родительской сущности.

Для взаимодействия с PostgreSQL в проекте используется драйвер lib/pq, который обеспечивает безопасное выполнение SQL-запросов через параметризованные запросы. Это защищает приложение от SQL-инъекций и обеспечивает корректную работу с различными типами данных, включая UUID и бинарные данные.

Репозиторий AnalysisRepository реализует весь набор операций по работе с анализами. Он предоставляет методы создания нового анализа, обновления его статуса и результатов, получения анализа по идентификатору, получения списка анализов пользователя, а также удаления анализа. Данный репозиторий использует прямое подключение к базе данных через интерфейс *sql.DB. Фрагмент кода репозитория AnalysisRepository представлен на листинге 2.5.1.

Листинг 2.5.1 – Фрагмент кода репозитория AnalysisRepository

```go
package repositories

import (
	"database/sql"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"webapp/internal/app/models"
)

type AnalysisRepository struct {
	db *sql.DB
}

func NewAnalysisRepository(db *sql.DB) *AnalysisRepository {
	return &AnalysisRepository{db: db}
}

func (r *AnalysisRepository) Create(projectID, userID uuid.UUID, inputType string, 
	inputMeta json.RawMessage) (*models.AnalysisRun, error) {
	id := uuid.New()
	now := time.Now()

	_, err := r.db.Exec(
		`INSERT INTO analysis_runs (id, project_id, user_id, status, created_at, input_type, input_meta)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		id, projectID, userID, "Created", now, inputType, inputMeta,
	)
	if err != nil {
		return nil, err
	}

	return &models.AnalysisRun{
		ID:        id,
		ProjectID: projectID,
		UserID:    userID,
		Status:    "Created",
		CreatedAt: now,
		InputType: inputType,
		InputMeta: inputMeta,
	}, nil
}

func (r *AnalysisRepository) ListByUserID(userID uuid.UUID) ([]*models.AnalysisRun, error) {
	rows, err := r.db.Query(
		`SELECT id, project_id, user_id, status, created_at, started_at, finished_at,
		 input_type, input_meta, summary_json
		 FROM analysis_runs WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []*models.AnalysisRun
	for rows.Next() {
		var a models.AnalysisRun
		var startedAt, finishedAt sql.NullTime
		var inputMeta, summaryJSON sql.NullString

		if err := rows.Scan(
			&a.ID, &a.ProjectID, &a.UserID, &a.Status, &a.CreatedAt,
			&startedAt, &finishedAt, &a.InputType, &inputMeta, &summaryJSON,
		); err != nil {
			return nil, err
		}

		if startedAt.Valid {
			a.StartedAt = &startedAt.Time
		}
		if finishedAt.Valid {
			a.FinishedAt = &finishedAt.Time
		}
		if inputMeta.Valid {
			a.InputMeta = json.RawMessage(inputMeta.String)
		}
		if summaryJSON.Valid {
			a.SummaryJSON = json.RawMessage(summaryJSON.String)
		}

		analyses = append(analyses, &a)
	}
	return analyses, rows.Err()
}
```

Репозиторий UserRepository отвечает за работу с пользователями. Он реализует создание нового пользователя с хешированием пароля, получение пользователя по имени и по идентификатору. Данный репозиторий также использует прямое подключение к базе данных. Репозиторий UserRepository представлен на листинге 2.5.2.

Листинг 2.5.2 – Репозиторий UserRepository

```go
package repositories

import (
	"database/sql"
	"time"
	"github.com/google/uuid"
	"webapp/internal/app/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(username, passwordHash string) (*models.User, error) {
	id := uuid.New()
	now := time.Now()

	_, err := r.db.Exec(
		"INSERT INTO users (id, username, password_hash, created_at) VALUES ($1, $2, $3, $4)",
		id, username, passwordHash, now,
	)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    now,
	}, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		"SELECT id, username, password_hash, created_at FROM users WHERE username = $1",
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
```

Репозиторий FileRepository отвечает за работу с файлами, сохранёнными в базе данных. Он реализует пакетное создание файлов для анализа, получение всех файлов анализа и получение конкретного файла по идентификатору. Данный репозиторий использует транзакции для обеспечения целостности данных при пакетной вставке. Фрагмент кода репозитория FileRepository представлен на листинге 2.5.3.

Листинг 2.5.3 – Фрагмент кода репозитория FileRepository

```go
func (r *FileRepository) CreateBatch(analysisID uuid.UUID, files []models.AnalysisFile) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO analysis_files (id, analysis_id, file_path, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()
	for _, file := range files {
		_, err := stmt.Exec(uuid.New(), analysisID, file.FilePath, file.Content, now)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *FileRepository) GetByAnalysisID(analysisID uuid.UUID) ([]*models.AnalysisFile, error) {
	rows, err := r.db.Query(
		`SELECT id, analysis_id, file_path, content, created_at
		 FROM analysis_files WHERE analysis_id = $1 ORDER BY file_path`,
		analysisID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*models.AnalysisFile
	for rows.Next() {
		var f models.AnalysisFile
		if err := rows.Scan(&f.ID, &f.AnalysisID, &f.FilePath, &f.Content, &f.CreatedAt); err != nil {
			return nil, err
		}
		files = append(files, &f)
	}
	return files, rows.Err()
}
```

