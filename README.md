# Notes
Для создания веб-приложений «Заметки» на языке Go мы можем использовать стандартную платформу net/http для маршрутизации и обработки запросов, а также функцию gorm для взаимодействия с базой данных. В качестве базы данных подойдет PostgreSQL. Также мы будем использовать Docker для развертывания.

Вот общий план проекта:
## Структура проекта
```notes_app/
├── main.go              # Точка входа приложения
├── handlers.go          # Обработчики HTTP-запросов
├── models.go            # Определения моделей
├── auth.go              # Обработка аутентификации
├── config.go            # Конфигурация базы данных
├── Dockerfile           # Dockerfile для сборки
├── docker-compose.yml   # Docker Compose для запуска
└── templates/           # HTML-шаблоны
    ├── base.html
    ├── index.html
    ├── note.html
    ├── edit.html
```
### Шаг 1: Определим модели
В файле models.go создаем структуру для заметок и пользователей:
```
// models.go
package main

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `gorm:"unique"`
    Password string
    Notes    []Note
}

type Note struct {
    gorm.Model
    Title     string
    Content   string
    UserID    uint
    ExpiresAt *time.Time // опциональное время истечения
}
```
### Шаг 2: Подключение к базе данных
Создаем config.goподключение к базе данных с использованием GORM:
```
// config.go
package main

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

func InitDB() *gorm.DB {
    dsn := "host=localhost user=postgres password=postgres dbname=notes_app port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect database:", err)
    }
    db.AutoMigrate(&User{}, &Note{})
    return db
}
```
### Шаг 3: Маршруты и обработчики
В файле handlers.goреализуем маршруты для создания, просмотра, редактирования и удаления заметок.
```
// handlers.go
package main

import (
    "net/http"
    "strconv"
    "html/template"
    "gorm.io/gorm"
)

var db *gorm.DB

func InitHandlers(database *gorm.DB) {
    db = database
    http.HandleFunc("/", ListNotesHandler)
    http.HandleFunc("/note", ViewNoteHandler)
    http.HandleFunc("/create", CreateNoteHandler)
    http.HandleFunc("/edit", EditNoteHandler)
    http.HandleFunc("/delete", DeleteNoteHandler)
}

func ListNotesHandler(w http.ResponseWriter, r *http.Request) {
    var notes []Note
    db.Find(&notes)
    tmpl := template.Must(template.ParseFiles("templates/index.html"))
    tmpl.Execute(w, notes)
}

// Обработчик для создания заметки
func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        title := r.FormValue("title")
        content := r.FormValue("content")
        note := Note{Title: title, Content: content, UserID: 1} // для простоты: привязка к user ID 1
        db.Create(&note)
        http.Redirect(w, r, "/", http.StatusFound)
    } else {
        tmpl := template.Must(template.ParseFiles("templates/create.html"))
        tmpl.Execute(w, nil)
    }
}

// Другие обработчики: ViewNoteHandler, EditNoteHandler, DeleteNoteHandler
```
### Шаг 4: Аутентификация
Создана простая авторизация пользователя в файле auth.go. Для хеширования паролей используйте библиотеку bcrypt.
```
// auth.go
package main

import (
    "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```
### Шаг 5: Основной файлmain.go
Настраиваем main.goподключение к базе данных и запускаем сервер.
```
// main.go
package main

import (
    "log"
    "net/http"
)

func main() {
    db := InitDB()
    InitHandlers(db)

    log.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}
```
### Шаг 6: HTML-шаблоны
Шаблоны находятся в каталогах templates/. Пример шаблона index.htmlдля отображения списка заметок:
```
<!-- templates/index.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Notes</title>
</head>
<body>
    <h1>Your Notes</h1>
    <a href="/create">Create a new note</a>
    <ul>
        {{range .}}
            <li><a href="/note?id={{.ID}}">{{.Title}}</a></li>
        {{end}}
    </ul>
</body>
</html>
```
### Шаг 7: Docker и Docker Compose
Dockerfile для создания контейнера с приложением:

```
# Dockerfile
FROM golang:1.20-alpine

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o notes_app

CMD ["./notes_app"]
```
### docker-compose.yml для запуска приложений и данных базы данных:
```
# docker-compose.yml
version: '3'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - db

  db:
    image: postgres:13
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: notes_app
    ports:
      - "5432:5432"
```

