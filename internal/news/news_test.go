package news

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/oneworkspacestudio/my-news/internal/webui"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type transport func(*http.Request) (*http.Response, error)

func (f transport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func response(s, kind string) *http.Response {
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{kind}}, Body: io.NopCloser(strings.NewReader(s))}
}
func fixtureFeed(source string) string {
	titles := []string{"Cities are making more room for people and nature", "A new chapter for clean energy across the world", "Inside the community projects bringing neighbours together", "What the next generation of technology could mean for us", "The researchers looking at our planet in a different light", "A weekend of sport, shared moments and new possibilities"}
	var b strings.Builder
	b.WriteString(`<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/"><channel><title>Synthetic test publisher</title>`)
	for i, title := range titles {
		fmt.Fprintf(&b, `<item><title>%s</title><link>https://public.example.test/%s/story-%d</link><description><![CDATA[<p>Synthetic acceptance story. A closer look at the people, ideas and discoveries shaping everyday life, with context from this publisher.</p>]]></description><pubDate>%s</pubDate>`, title, source, i, time.Now().Add(-time.Duration(i+1)*time.Hour).Format(time.RFC1123Z))
		if i%2 == 1 {
			b.WriteString(`<content:encoded><![CDATA[<p>FULL FEED TEXT. This synthetic test article provides its complete reading text through RSS.</p><p>Offline readers can return to these paragraphs without fetching the publisher again.</p>]]></content:encoded>`)
		}
		b.WriteString(`</item>`)
	}
	b.WriteString(`</channel></rss>`)
	return b.String()
}
func fixtureClient() *http.Client {
	return &http.Client{Transport: transport(func(r *http.Request) (*http.Response, error) {
		if strings.Contains(r.URL.Host, "public.example.test") {
			return response("<html><body><nav>Not part of the article</nav><article><h1>Public reading copy</h1><p>PUBLIC ARTICLE TEXT. "+strings.Repeat("This is synthetic test content supplied by a public article page. It confirms the reading copy survives an offline browser restart. ", 12)+"</p><script>alert('blocked')</script></article></body></html>", "text/html"), nil
		}
		return response(fixtureFeed(apphost.Hash(r.URL.String())[:8]), "application/rss+xml"), nil
	})}
}
func TestFeedParsingAndSanitizing(t *testing.T) {
	s := Source{ID: "x", Name: "Synthetic", URL: "https://example.com/feed", Category: "Science", Region: "Europe"}
	items, e := parseFeed([]byte(fixtureFeed("x")), s)
	if e != nil || len(items) != 6 {
		t.Fatal(e, len(items))
	}
	if items[1].ContentKind != "Feed content" || !strings.Contains(items[1].Text, "FULL FEED TEXT") {
		t.Fatal(items[1])
	}
	atom := `<feed xmlns="http://www.w3.org/2005/Atom"><entry><title>Atom &amp; science</title><link rel="alternate" href="https://example.com/atom"/><summary>Summary</summary><content type="html">&lt;p&gt;Article&lt;/p&gt;</content><updated>2026-11-05T10:00:00Z</updated></entry></feed>`
	items, e = parseFeed([]byte(atom), s)
	if e != nil || len(items) != 1 || items[0].Text != "Article" {
		t.Fatal(e, items)
	}
	if _, e = parseFeed([]byte(`<html><body>Not a feed</body></html>`), s); e == nil {
		t.Fatal("accepted HTML as feed")
	}
	if got := plain(`<p>Hello</p><script>alert('x')</script><style>bad</style><p>World</p>`); got != "Hello\n\nWorld" {
		t.Fatal(got)
	}
	if safeLink("javascript:alert(1)", s.URL) != "" {
		t.Fatal("script URL accepted")
	}
}
func TestNetworkBoundary(t *testing.T) {
	for _, s := range []string{"127.0.0.1", "10.1.1.1", "172.16.0.1", "192.168.1.1", "169.254.169.254", "100.64.0.1", "::1", "::ffff:127.0.0.1", "fc00::1", "2002:7f00:1::", "64:ff9b::7f00:1"} {
		if publicIP(net.ParseIP(s)) {
			t.Fatal("private/reserved address allowed", s)
		}
	}
	for _, s := range []string{"https://user:pass@example.com/feed", "file:///etc/passwd", "https://example.com:8080/feed", "http://127.0.0.1/"} {
		if _, e := validURL(s); e == nil {
			t.Fatal(s)
		}
	}
	if !publicIP(net.ParseIP("8.8.8.8")) {
		t.Fatal("public address blocked")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, _, e := fetch(ctx, newClient(), "http://localhost/", 100); e == nil {
		t.Fatal("localhost fetch escaped guard")
	}
	c := newClient()
	if e := c.CheckRedirect(&http.Request{URL: &url.URL{Scheme: "http", Host: "127.0.0.1", Path: "/"}}, nil); e == nil {
		t.Fatal("redirect accepted private address")
	}
}

func testApp(t *testing.T) (*App, *apphost.Host) {
	t.Helper()
	h, e := apphost.New(apphost.Config{ID: "news", Name: "My News", Data: filepath.Join(t.TempDir(), "private"), Listen: "127.0.0.1:19330", Origin: "http://localhost:19330/news"}, webui.Files())
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { h.Close() })
	if e = h.Bootstrap("reader@example.test", "synthetic-news-test"); e != nil {
		t.Fatal(e)
	}
	a, e := New(h)
	if e != nil {
		t.Fatal(e)
	}
	a.client = fixtureClient()
	return a, h
}
func request(h *apphost.Host, method, path, body string, cookie *http.Cookie, csrf string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "http://localhost:19330/news"+path, strings.NewReader(body))
	r.RemoteAddr = "127.0.0.1:12345"
	r.Header.Set("Origin", "http://localhost:19330")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-CSRF-Token", csrf)
	if cookie != nil {
		r.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	h.Handler().ServeHTTP(w, r)
	return w
}
func login(t *testing.T, h *apphost.Host, email string) (*http.Cookie, string) {
	t.Helper()
	w := request(h, "POST", "/api/login", `{"email":"`+email+`","password":"synthetic-news-test"}`, nil, "")
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	cookie := w.Result().Cookies()[0]
	w = request(h, "GET", "/api/me", "", cookie, "")
	var s apphost.Session
	if e := json.Unmarshal(w.Body.Bytes(), &s); e != nil {
		t.Fatal(e)
	}
	return cookie, s.CSRF
}
func TestNativeWorkflowIsolationAndPersistence(t *testing.T) {
	a, h := testApp(t)
	cookie, csrf := login(t, h, "reader@example.test")
	if w := request(h, "GET", "/api/state", "", nil, ""); w.Code != 401 {
		t.Fatal("anonymous API", w.Code)
	}
	if w := request(h, "POST", "/api/refresh", "", cookie, "invalid"); w.Code != 403 {
		t.Fatal("CSRF", w.Code)
	}
	if w := request(h, "POST", "/api/refresh", "", cookie, csrf); w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	owner := apphost.Identity("local:news", "reader@example.test")
	stories, e := a.stories(owner, false)
	if e != nil || len(stories) != 108 {
		t.Fatal(e, len(stories))
	}
	s := stories[0]
	w := request(h, "POST", "/api/saved/"+s.ID, "", cookie, csrf)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "PUBLIC ARTICLE TEXT") {
		t.Fatal(w.Code, w.Body.String())
	}
	if w = request(h, "POST", "/api/local-users", `{"email":"second@example.test","password":"synthetic-news-test"}`, cookie, csrf); w.Code != 201 {
		t.Fatal(w.Code, w.Body.String())
	}
	other, otherCSRF := login(t, h, "second@example.test")
	w = request(h, "POST", "/api/saved/"+s.ID, "", other, otherCSRF)
	if w.Code != 404 {
		t.Fatal("cross-account story leak", w.Code)
	}
	saved, e := a.stories(owner, true)
	if e != nil || len(saved) != 1 {
		t.Fatal(e, len(saved))
	}
	w = request(h, "POST", "/api/sources", `{"name":"Bad","url":"http://127.0.0.1/feed","region":"Local","category":"Science"}`, cookie, csrf)
	if w.Code != 400 {
		t.Fatal("private source accepted")
	}
	w = request(h, "DELETE", "/api/sources/"+s.SourceID, "", cookie, csrf)
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
	saved, _ = a.stories(owner, true)
	if len(saved) != 1 {
		t.Fatal("source removal lost saved reading copy")
	}
	var stored string
	if e = h.DB.QueryRow("SELECT value FROM news_saved WHERE user_id=? AND id=?", owner, s.ID).Scan(&stored); e != nil || !strings.Contains(stored, "PUBLIC ARTICLE TEXT") {
		t.Fatal("saved reading copy was not persisted", e)
	}
}

// Browser acceptance runs a real native host and database with deterministic
// synthetic publisher responses. No production fetch bypass is compiled in.
func TestBrowserAcceptance(t *testing.T) {
	if os.Getenv("NEWS_BROWSER_TEST") != "1" {
		t.Skip("explicit native browser acceptance")
	}
	a, h := testApp(t)
	listener, e := net.Listen("tcp", "127.0.0.1:19330")
	if e != nil {
		t.Fatal(e)
	}
	server := &http.Server{Handler: h.Handler()}
	go server.Serve(listener)
	defer server.Close()
	_ = a
	cmd := exec.Command("node", "../../web/scripts/browser-verify.mjs")
	cmd.Env = append(os.Environ(), "NEWS_TEST_URL=http://localhost:19330/news/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if e = cmd.Run(); e != nil {
		t.Fatal(e)
	}
}

func TestLivePublisherFeeds(t *testing.T) {
	if os.Getenv("NEWS_LIVE_TEST") != "1" {
		t.Skip("explicit live publisher check")
	}
	client := newClient()
	defer client.CloseIdleConnections()
	type result struct {
		Name    string `json:"name"`
		URL     string `json:"url"`
		Stories int    `json:"stories"`
		Error   string `json:"error,omitempty"`
	}
	results := make(chan result, len(defaults))
	sem := make(chan struct{}, 5)
	for _, source := range defaults {
		go func(s Source) {
			sem <- struct{}{}
			defer func() { <-sem }()
			ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
			defer cancel()
			b, _, e := fetch(ctx, client, s.URL, 4<<20)
			var items []Story
			if e == nil {
				items, e = parseFeed(b, s)
			}
			r := result{Name: s.Name, URL: s.URL, Stories: len(items)}
			if e != nil {
				r.Error = e.Error()
			}
			results <- r
		}(source)
	}
	rows := []result{}
	success := 0
	for range defaults {
		r := <-results
		rows = append(rows, r)
		if r.Error == "" && r.Stories > 0 {
			success++
		}
		t.Log(r.Name, r.Stories, r.Error)
	}
	b, _ := json.MarshalIndent(map[string]any{"checkedAt": time.Now().UTC().Format(time.RFC3339), "success": success, "total": len(defaults), "sources": rows}, "", "  ")
	if e := os.WriteFile("../../docs/live-source-check.json", b, 0600); e != nil {
		t.Fatal(e)
	}
	if success == 0 {
		t.Fatal("no live feeds reached")
	}
}
