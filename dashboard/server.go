package dashboard

import (
	"fmt"
	"html/template"
	"net/http"
	db "pingo/database"
	"strconv"
	"strings"
)

var Mux *http.ServeMux

func init() {
	Mux = http.NewServeMux()

	// Get latest job logs with ?limit=n (default to 1000)
	Mux.HandleFunc("/latest/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			jobName := strings.TrimPrefix(r.URL.Path, "/latest/")
			if jobName == "" {
				http.Error(w, "no job provided", http.StatusBadRequest)
				return
			}

			var limit int = 1000
			var queryErr error
			if r.URL.Query().Has("limit") {
				limit, queryErr = strconv.Atoi(r.URL.Query().Get("limit"))
				if queryErr != nil {
					http.Error(w, "couldnt parse query parametr 'limit'", http.StatusBadRequest)
					return
				}
			}

			data, err := db.SelectLatestJobLogs(db.DB, jobName, limit)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(data) == 0 {
				http.NotFound(w, r)
				return
			}

			w.WriteHeader(http.StatusOK)
			getLineChart(&data).Render(w)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Get job statistics for status with count
	Mux.HandleFunc("/statistics/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			jobName := strings.TrimPrefix(r.URL.Path, "/statistics/")
			if jobName == "" {
				http.Error(w, "no job provided", http.StatusBadRequest)
				return
			}

			data, err := db.SelectJobsForPieChart(db.DB, jobName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(data) == 0 {
				http.NotFound(w, r)
				return
			}

			w.WriteHeader(http.StatusOK)
			getPieChart(&data).Render(w)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Get job latest status
	Mux.HandleFunc("/status/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			jobName := strings.TrimPrefix(r.URL.Path, "/status/")
			if jobName == "" {
				http.Error(w, "no job provided", http.StatusBadRequest)
				return
			}

			data, err := db.SelectLatestJobLogs(db.DB, jobName, 1)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(data) == 0 {
				http.NotFound(w, r)
				return
			}

			tmpl, err := template.New("").Parse("<div class='status-circle {{.}}'></div>")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, template.HTML(strings.ToLower(data[0].Status)))
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Serve static files
	Mux.Handle("/static/", http.FileServer(http.FS(static)))

	Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Println(db.SelectJobsInfo(db.DB))
			tmpl, err := template.New("").
				Funcs(
					template.FuncMap{
						"mapTarget": mapTarget,
						"mapType":   mapType,
						"jobsInfo":  func() ([]db.DBJob, error) { return db.SelectJobsInfo(db.DB) },
					}).
				ParseFS(static, "static/base.tpl.html", "static/jobs_table.tpl.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			tmpl.ExecuteTemplate(w, "base", struct {
				Title string
			}{
				Title: "Dashboard",
			})
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	Mux.HandleFunc("/jobs/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			jobName := strings.TrimPrefix(r.URL.Path, "/jobs/")
			if jobName == "" {
				http.Error(w, "no job provided", http.StatusBadRequest)
				return
			}
			tmpl, err := template.New("").
				Funcs(
					template.FuncMap{
						"JobName": func() string { return jobName },
					}).
				ParseFS(static, "static/base.tpl.html", "static/job_detail.tpl.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			tmpl.ExecuteTemplate(w, "base", struct {
				Title string
			}{
				Title: jobName,
			})
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

}
