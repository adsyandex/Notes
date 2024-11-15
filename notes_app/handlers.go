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
    tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
    tmpl.ExecuteTemplate(w, "base", notes)
}

func ViewNoteHandler(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.URL.Query().Get("id"))
    var note Note
    db.First(&note, id)
    tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/note.html"))
    tmpl.ExecuteTemplate(w, "base", note)
}

func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        title := r.FormValue("title")
        content := r.FormValue("content")
        note := Note{Title: title, Content: content, UserID: 1} // привязка к user ID 1 для примера
        db.Create(&note)
        http.Redirect(w, r, "/", http.StatusFound)
    } else {
        tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/create.html"))
        tmpl.ExecuteTemplate(w, "base", nil)
    }
}

func EditNoteHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        id, _ := strconv.Atoi(r.FormValue("id"))
        var note Note
        db.First(&note, id)
        note.Title = r.FormValue("title")
        note.Content = r.FormValue("content")
        db.Save(&note)
        http.Redirect(w, r, "/", http.StatusFound)
    } else {
        id, _ := strconv.Atoi(r.URL.Query().Get("id"))
        var note Note
        db.First(&note, id)
        tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/edit.html"))
        tmpl.ExecuteTemplate(w, "base", note)
    }
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        id, _ := strconv.Atoi(r.FormValue("id"))
        db.Delete(&Note{}, id)
        http.Redirect(w, r, "/", http.StatusFound)
    }
}
