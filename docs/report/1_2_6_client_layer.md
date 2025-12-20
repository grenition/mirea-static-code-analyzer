### Разработка слоя клиентского представления веб-приложения

Слой клиентского представления веб-приложения «Анализатор статического кода» отвечает за отображение данных, формирование пользовательского интерфейса и обеспечение удобного взаимодействия пользователя с системой. Этот слой реализован с использованием HTML5, CSS3, JavaScript и механизма шаблонов Go, что позволяет отделить внешний вид приложения от серверной логики и структуры базы данных.

Функция executeTemplate играет ключевую роль в клиентском слое проекта. Она принимает имя шаблона и данные, затем подключает необходимые HTML-файлы, формируя итоговую HTML-страницу. Благодаря этому клиентский интерфейс отделён от серверной логики, а каждая страница представлена отдельным файлом .html с разметкой.

Основой клиентского представления является HTML-разметка, определяющая структуру всех страниц приложения. HTML обеспечивает передачу данных на сервер посредством GET и POST [21] запросов через стандартные формы. На уровне разметки применяется встроенная валидация. Обязательные поля помечены атрибутом required, а поля для паролей используют тип password, что предотвращает ввод некорректных данных. Фрагмент кода файла с разметкой login.html представлен на листинге 2.6.1.

Листинг 2.6.1 – Фрагмент кода login.html

```html
<div class="card shadow-sm">
    <div class="card-header bg-primary text-white text-center">
        <h4 class="mb-0">
            <i class="bi bi-box-arrow-in-right me-2"></i>Вход
        </h4>
    </div>
    <div class="card-body">
        {{if .Error}}
        <div class="alert alert-danger" role="alert">
            <i class="bi bi-exclamation-triangle me-2"></i>{{.Error}}
        </div>
        {{end}}
        <form method="POST" action="/login">
            <div class="mb-3">
                <label for="username" class="form-label">Имя пользователя</label>
                <input type="text" class="form-control" id="username" name="username" 
                       value="{{.Username}}" required autofocus>
            </div>
            <div class="mb-3">
                <label for="password" class="form-label">Пароль</label>
                <input type="password" class="form-control" id="password" name="password" required>
            </div>
            <div class="d-grid">
                <button type="submit" class="btn btn-primary">
                    <i class="bi bi-box-arrow-in-right me-2"></i>Войти
                </button>
            </div>
        </form>
    </div>
</div>
```

CSS3 отвечает за визуальное оформление веб-приложения и формирование единых стилей для всех элементов интерфейса. В проекте используется фреймворк Bootstrap 5 для базовых стилей и кастомный файл styles.css для дополнительного оформления. В результате интерфейс приобретает современный и удобочитаемый вид, обеспечивая комфортное взаимодействие пользователя с результатами анализа кода.

Отдельную роль играет JavaScript, который используется в проекте для обеспечения динамического поведения страниц. В частности, JavaScript применяется для фильтрации найденных проблем по уровню серьёзности, отображения модальных окон с содержимым файлов, выполнения AJAX-запросов [22] при анализе кода в реальном времени и подтверждения удаления анализов. Это позволяет обновлять интерфейс без перезагрузки страницы, делая работу приложения быстрее и удобнее. Фрагмент кода из snippet_analyzer.html представлен на листинге 2.6.2.

Листинг 2.6.2 – Фрагмент кода snippet_analyzer.html

```javascript
document.getElementById('autoAnalyze').addEventListener('change', function() {
    const analyzeBtn = document.getElementById('analyzeBtn');
    if (this.checked) {
        analyzeBtn.style.display = 'none';
        analyzeCode();
    } else {
        analyzeBtn.style.display = 'block';
    }
});

document.getElementById('language').addEventListener('change', function() {
    if (document.getElementById('autoAnalyze').checked) {
        setTimeout(analyzeCode, 500);
    }
});

function analyzeCode() {
    const code = document.getElementById('code').value;
    const formData = new FormData();
    formData.append('code', code);
    formData.append('extension', document.getElementById('language').value);
    
    fetch('/api/analyze/snippet', {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        displayResults(data);
    });
}
```

