package api

import (
    "fmt"
    "io"
    "net/http"
    "os"
)

func UploadCSVHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    err := r.ParseMultipartForm(10 << 20) // max 10MB file
    if err != nil {
        fmt.Println("ParseMultipartForm error:", err)
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    file, _, err := r.FormFile("file") 
    if err != nil {
        fmt.Println("FormFile error:", err)
        http.Error(w, "Failed to get file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    dst, err := os.Create("CurrentBudget.csv")
    if err != nil {
        fmt.Println("Create file error:", err)
        http.Error(w, "Failed to create file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    _, err = io.Copy(dst, file)
    if err != nil {
        fmt.Println("Copy file error:", err)
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("CSV uploaded successfully"))
}
