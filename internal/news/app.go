package news

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type App struct {
	H          *apphost.Host
	client     *http.Client
	mu         sync.Mutex
	refreshing map[string]bool
	last       map[string]time.Time
}

func New(h *apphost.Host) (*App, error) {
	if e := h.Migrate(1, `CREATE TABLE news_sources(user_id TEXT NOT NULL,id TEXT NOT NULL,value TEXT NOT NULL,PRIMARY KEY(user_id,id));CREATE TABLE news_stories(user_id TEXT NOT NULL,source_id TEXT NOT NULL,id TEXT NOT NULL,value TEXT NOT NULL,PRIMARY KEY(user_id,source_id,id));CREATE TABLE news_saved(user_id TEXT NOT NULL,id TEXT NOT NULL,value TEXT NOT NULL,PRIMARY KEY(user_id,id));CREATE TABLE news_preferences(user_id TEXT PRIMARY KEY,value TEXT NOT NULL);`); e != nil {
		return nil, e
	}
	a := &App{H: h, client: newClient(), refreshing: map[string]bool{}, last: map[string]time.Time{}}
	for pattern, handler := range map[string]http.HandlerFunc{"GET /api/state": a.state, "POST /api/refresh": a.refresh, "POST /api/sources": a.addSource, "DELETE /api/sources/{id}": a.deleteSource, "POST /api/saved/{id}": a.save, "DELETE /api/saved/{id}": a.unsave, "PUT /api/preferences": a.preferences} {
		h.Mux.HandleFunc(pattern, h.Require(handler))
	}
	return a, nil
}
func (a *App) Run(ctx context.Context) error {
	<-ctx.Done()
	a.client.CloseIdleConnections()
	return nil
}
func (a *App) sources(user string) ([]Source, error) {
	var n int
	if e := a.H.DB.QueryRow("SELECT count(*) FROM news_preferences WHERE user_id=?", user).Scan(&n); e != nil {
		return nil, e
	}
	if n == 0 {
		tx, e := a.H.DB.Begin()
		if e != nil {
			return nil, e
		}
		defer tx.Rollback()
		for _, s := range defaults {
			b, _ := json.Marshal(s)
			if _, e = tx.Exec("INSERT OR IGNORE INTO news_sources VALUES(?,?,?)", user, s.ID, string(b)); e != nil {
				return nil, e
			}
		}
		if _, e = tx.Exec("INSERT OR IGNORE INTO news_preferences VALUES(?,?)", user, "{}"); e != nil {
			return nil, e
		}
		if e = tx.Commit(); e != nil {
			return nil, e
		}
	}
	rows, e := a.H.DB.Query("SELECT value FROM news_sources WHERE user_id=? ORDER BY id", user)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []Source{}
	for rows.Next() {
		var v string
		var s Source
		if e = rows.Scan(&v); e != nil {
			return nil, e
		}
		if e = json.Unmarshal([]byte(v), &s); e != nil {
			return nil, e
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
func (a *App) stories(user string, saved bool) ([]Story, error) {
	q := "SELECT value FROM news_stories WHERE user_id=?"
	if saved {
		q = "SELECT value FROM news_saved WHERE user_id=?"
	}
	rows, e := a.H.DB.Query(q, user)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []Story{}
	seen := map[string]bool{}
	for rows.Next() {
		var v string
		var s Story
		if e = rows.Scan(&v); e != nil {
			return nil, e
		}
		if e = json.Unmarshal([]byte(v), &s); e != nil {
			return nil, e
		}
		key := s.SourceID + ":" + s.ID
		if saved {
			key = s.ID
		}
		if !seen[key] {
			seen[key] = true
			out = append(out, s)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Published == out[j].Published {
			return out[i].ID < out[j].ID
		}
		return out[i].Published > out[j].Published
	})
	return out, rows.Err()
}
func (a *App) state(w http.ResponseWriter, r *http.Request) {
	u := apphost.Current(r).User.ID
	s, e := a.sources(u)
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	stories, e := a.stories(u, false)
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	saved, e := a.stories(u, true)
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	var prefs string
	if e = a.H.DB.QueryRow("SELECT value FROM news_preferences WHERE user_id=?", u).Scan(&prefs); e != nil {
		a.H.DBError(w, e)
		return
	}
	apphost.JSON(w, 200, map[string]any{"sources": s, "stories": stories, "saved": saved, "preferences": json.RawMessage(prefs)})
}
func (a *App) refresh(w http.ResponseWriter, r *http.Request) {
	u := apphost.Current(r).User.ID
	a.mu.Lock()
	if a.refreshing[u] || time.Since(a.last[u]) < time.Minute {
		a.mu.Unlock()
		apphost.Problem(w, 429, "Feeds were refreshed recently. Please wait a minute.")
		return
	}
	if len(a.refreshing) >= 4 {
		a.mu.Unlock()
		apphost.Problem(w, 429, "The reader is busy. Try again shortly.")
		return
	}
	a.refreshing[u] = true
	a.last[u] = time.Now()
	a.mu.Unlock()
	defer func() { a.mu.Lock(); delete(a.refreshing, u); a.mu.Unlock() }()
	sources, e := a.sources(u)
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)
	var mu sync.Mutex
	failures := []string{}
	success := 0
	for _, source := range sources {
		wg.Add(1)
		go func(s Source) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-r.Context().Done():
				return
			}
			defer func() { <-sem }()
			b, _, err := fetch(r.Context(), a.client, s.URL, 4<<20)
			var stories []Story
			if err == nil {
				stories, err = parseFeed(b, s)
			}
			if err == nil {
				tx, te := a.H.DB.Begin()
				if te == nil {
					_, te = tx.Exec("DELETE FROM news_stories WHERE user_id=? AND source_id=?", u, s.ID)
					if te == nil {
						for _, story := range stories {
							v, _ := json.Marshal(story)
							if _, te = tx.Exec("INSERT INTO news_stories VALUES(?,?,?,?)", u, s.ID, story.ID, string(v)); te != nil {
								break
							}
						}
					}
					if te == nil {
						te = tx.Commit()
					} else {
						_ = tx.Rollback()
					}
				}
				err = te
			}
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				s.Error = "Could not refresh this feed. Previous stories are kept."
				failures = append(failures, s.Name)
			} else {
				s.Updated = time.Now().UTC().Format(time.RFC3339)
				s.Error = ""
				success++
			}
			v, _ := json.Marshal(s)
			if _, e := a.H.DB.Exec("UPDATE news_sources SET value=? WHERE user_id=? AND id=?", string(v), u, s.ID); e != nil {
				failures = append(failures, s.Name+" (status could not be saved)")
			}
		}(source)
	}
	wg.Wait()
	apphost.JSON(w, 200, map[string]any{"refreshed": success, "failed": failures})
}
func (a *App) addSource(w http.ResponseWriter, r *http.Request) {
	var s Source
	if !apphost.Decode(w, r, &s) {
		return
	}
	s.Name = strings.TrimSpace(s.Name)
	if s.Name == "" || len(s.Name) > 100 || len(s.Region) > 80 || len(s.Category) > 50 {
		apphost.Problem(w, 400, "Give this source a name, region and category")
		return
	}
	if _, e := validURL(s.URL); e != nil {
		apphost.Problem(w, 400, e.Error())
		return
	}
	u := apphost.Current(r).User.ID
	sources, e := a.sources(u)
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	if len(sources) >= 50 {
		apphost.Problem(w, 400, "You can follow up to 50 feeds")
		return
	}
	s.ID = apphost.Hash(s.URL)
	s.Custom = true
	s.Error = ""
	s.Updated = ""
	b, _, e := fetch(r.Context(), a.client, s.URL, 4<<20)
	if e == nil {
		_, e = parseFeed(b, s)
	}
	if e != nil {
		apphost.Problem(w, 400, e.Error())
		return
	}
	v, _ := json.Marshal(s)
	if _, e = a.H.DB.Exec("INSERT OR REPLACE INTO news_sources VALUES(?,?,?)", u, s.ID, string(v)); e != nil {
		a.H.DBError(w, e)
		return
	}
	apphost.JSON(w, 201, s)
}
func (a *App) deleteSource(w http.ResponseWriter, r *http.Request) {
	u := apphost.Current(r).User.ID
	tx, e := a.H.DB.Begin()
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	defer tx.Rollback()
	if _, e = tx.Exec("DELETE FROM news_sources WHERE user_id=? AND id=?", u, r.PathValue("id")); e == nil {
		_, e = tx.Exec("DELETE FROM news_stories WHERE user_id=? AND source_id=?", u, r.PathValue("id"))
	}
	if e == nil {
		e = tx.Commit()
	}
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	apphost.JSON(w, 200, map[string]bool{"ok": true})
}
func (a *App) save(w http.ResponseWriter, r *http.Request) {
	u := apphost.Current(r).User.ID
	id := r.PathValue("id")
	var raw string
	e := a.H.DB.QueryRow("SELECT value FROM news_stories WHERE user_id=? AND id=? LIMIT 1", u, id).Scan(&raw)
	if e != nil {
		e = a.H.DB.QueryRow("SELECT value FROM news_saved WHERE user_id=? AND id=?", u, id).Scan(&raw)
	}
	if e != nil {
		apphost.Problem(w, 404, "This story is no longer in your feeds. Refresh and try another story.")
		return
	}
	var s Story
	if json.Unmarshal([]byte(raw), &s) != nil {
		a.H.DBError(w, errors.New("invalid story"))
		return
	}
	// Try the ordinary public article only when the feed supplies an excerpt.
	// No credentials, scripts, reader proxies or paywall workarounds are used.
	if s.ContentKind == "Feed excerpt" {
		ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
		b, contentType, err := fetch(ctx, a.client, s.URL, 3<<20)
		cancel()
		if err == nil && strings.Contains(contentType, "text/html") {
			text := articleText(b)
			lower := strings.ToLower(text)
			if len(text) > len(s.Text)+400 && !strings.Contains(lower, "subscribe to continue") && !strings.Contains(lower, "sign in to continue") && !strings.Contains(lower, "subscription required") {
				s.Text = text
				s.ContentKind = "Article reading copy"
			}
		}
	}
	s.SavedAt = time.Now().UTC().Format(time.RFC3339)
	b, _ := json.Marshal(s)
	var count int
	if e = a.H.DB.QueryRow("SELECT count(*) FROM news_saved WHERE user_id=?", u).Scan(&count); e != nil {
		a.H.DBError(w, e)
		return
	}
	if count >= 2000 {
		apphost.Problem(w, 400, "Your library has reached 2,000 stories. Remove some before saving more.")
		return
	}
	if _, e = a.H.DB.Exec("INSERT OR REPLACE INTO news_saved VALUES(?,?,?)", u, id, string(b)); e != nil {
		a.H.DBError(w, e)
		return
	}
	apphost.JSON(w, 200, s)
}
func (a *App) unsave(w http.ResponseWriter, r *http.Request) {
	if _, e := a.H.DB.Exec("DELETE FROM news_saved WHERE user_id=? AND id=?", apphost.Current(r).User.ID, r.PathValue("id")); e != nil {
		a.H.DBError(w, e)
		return
	}
	apphost.JSON(w, 200, map[string]bool{"ok": true})
}
func (a *App) preferences(w http.ResponseWriter, r *http.Request) {
	var p struct {
		Following []string `json:"following"`
		Region    string   `json:"region"`
		Theme     string   `json:"theme"`
	}
	if !apphost.Decode(w, r, &p) {
		return
	}
	if len(p.Following) > 50 || len(p.Region) > 80 || (p.Theme != "" && p.Theme != "light" && p.Theme != "dark") {
		apphost.Problem(w, 400, "Invalid preferences")
		return
	}
	b, _ := json.Marshal(p)
	if _, e := a.H.DB.Exec("INSERT OR REPLACE INTO news_preferences VALUES(?,?)", apphost.Current(r).User.ID, string(b)); e != nil {
		a.H.DBError(w, e)
		return
	}
	apphost.JSON(w, 200, p)
}
