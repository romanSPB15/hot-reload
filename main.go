package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"

	_ "embed"
)

//go:embed script.html
var script []byte

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Host == "localhost:8080"
	},
}

func main() {
	rootDir := flag.String("dir", ".", "Директория с фронтендом")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Ошибка создания watcher:", err)
	}
	defer watcher.Close()

	err = watcher.Add("page.html")
	if err != nil {
		log.Fatal("Ошибка добавления файла в watcher:", err)
	}
	reloadSignal := make(chan struct{})
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					select {
					case reloadSignal <- struct{}{}:
					default:
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Ошибка watcher:", err)
			}
		}
	}()

	fileServer := http.FileServer(http.Dir(*rootDir))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 1. Проверяем, не является ли запрос WebSocket
		if r.URL.Path == "/ws" {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Print("Ошибка апгрейда до WebSocket:", err)
				return
			}
			defer conn.Close()
			for range reloadSignal {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
					return
				}
			}
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		fullPath := filepath.Join(*rootDir, path)

		info, err := os.Stat(fullPath)
		if err != nil || info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".html" && ext != ".htm" {
			fileServer.ServeHTTP(w, r)
			return
		}

		content, err := os.ReadFile(fullPath)
		if err != nil {
			http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
			return
		}

		scriptTag := string(script)
		html := string(content)
		insertPos := strings.LastIndex(html, "</body>")
		if insertPos != -1 {
			html = html[:insertPos] + scriptTag + html[insertPos:]
		} else {

			html += scriptTag
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	})

	srv := &http.Server{
		Addr: ":8080",
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		fmt.Println("Получен сигнал завершения, сервер останавливается...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatal("Ошибка при завершении сервера:", err)
		}
	}()

	fmt.Println("Сервер запущен...")
	fmt.Println("http://localhost:8080/index.html")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
