// handlers_test.go
package main

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestCreateNoteHandler(t *testing.T) {
    // Создаем тестовый HTTP-запрос с методом POST и форм-данными для новой заметки
    req := httptest.NewRequest("POST", "/create", strings.NewReader("title=TestNote&content=This is a test note"))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Создаем ResponseRecorder для записи ответа
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(CreateNoteHandler)

    // Вызываем обработчик
    handler.ServeHTTP(rr, req)

    // Проверяем статус ответа
    if status := rr.Code; status != http.StatusFound {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
    }

    // Дополнительная проверка: можно извлечь последнюю добавленную заметку из базы данных
    var note Note
    if err := db.Last(&note).Error; err != nil {
        t.Fatalf("could not retrieve note from db: %v", err)
    }

    if note.Title != "TestNote" || note.Content != "This is a test note" {
        t.Errorf("created note does not match: got %v, want %v", note, "TestNote, This is a test note")
    }
}
