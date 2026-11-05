# My News

Your world, in perspective. A personal news reader inspired by Google News,
with its own identity and a library you can read without an internet connection.

- Browse headlines from 18 RSS feeds across nine world regions.
- Explore Business, Technology, Science, Health, Sports, Entertainment and Environment.
- Search stories and publishers, follow topics, or add a local RSS/Atom source.
- Save a reading copy, then reopen it offline on the same device.
- Read in a quiet text view with adjustable type and light/dark themes.

My News uses BBC, Guardian and Al Jazeera feeds as starting points. It is not
affiliated with Google or these publishers and does not include every publication
or location. Add a publisher's feed and give it a city/country to expand coverage.
Publisher attribution and links remain attached to every reading copy.

## Offline reading

Use the download button on a story. My News saves the feed text and, for excerpts,
tries to extract a reading copy from the ordinary public article page. A label
identifies **Article reading copy**, **Feed content** or **Feed excerpt**. Extraction
can omit text; paywalled, unavailable or script-rendered articles may remain excerpts.
Images, video and the publisher's interactive website are not downloaded.

The confirmation means the reading text was written and read back from this device's
storage. Visit **Saved & offline** while disconnected. Use **Download library to
this device** after signing in on another device. The app shell becomes available
offline after the initial online visit and service-worker installation. Browser
storage eviction or clearing site data removes downloads. Signing out clears this
device's copies; the server library is retained for the next sign-in.

This is a personal device feature: someone with access to the unlocked browser
profile can read its offline copies. A server cannot remotely revoke data from a
disconnected device. Online authorization failures clear local copies. On shared
devices, sign out before leaving.

## Local build and use

Requires Go 1.27.3+, Node and npm. Dependencies are pinned in go.sum and package-lock.json.

```powershell
npm --prefix web ci
npm --prefix web run build
go test ./...
go build -o MyNews.exe ./cmd/my-news
./MyNews.exe
```

First standalone launch opens a one-time account setup link on this PC. The app
listens on loopback port 9030 and keeps state in the user's configuration directory
under `OneWorkspace/MyNews`, separately from its executable. An explicit `--data`
folder must be private NTFS on Windows. Managed installations use MyAccounts instead.

For Workspace, use the normal installer/catalog path; copying an executable into
an existing installation is not a managed install. Source integration is in the
separate `workspace-my-news` checkout. See [delivery status](docs/ACCEPTANCE.md),
[security](docs/SECURITY.md) and [source notes](docs/RESEARCH.md).

Browser acceptance uses synthetic publishers and isolated accounts:

```powershell
$env:NEWS_BROWSER_TEST='1'
go test ./internal/news -run TestBrowserAcceptance -v -count=1
```

It requires installed Chrome. No fake headlines are shipped in the app.
