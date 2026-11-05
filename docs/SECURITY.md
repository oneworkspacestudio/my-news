# Security and storage

My News uses the shared Go apphost for MyAccounts OIDC, owner binding, standalone
setup, session expiry, CSRF/origin validation, private NTFS state, logging and
encrypted backup/restore. API rows are scoped to the authenticated immutable user
ID. Managed mode does not permit an independent local login. Health reveals only
application/version status. All news APIs require authentication.

RSS/Atom text is untrusted. Scripts/styles/frames/SVG are discarded; the renderer
only renders text nodes. Links accept HTTP(S), reject credentials and unsafe schemes.
Feeds and article pages go through a public-address-only HTTP transport. DNS results
are checked and the chosen IP is pinned for each connection, including redirects.
Loopback, private, link-local, multicast, reserved/documentation and transition
addresses are denied. Ports are limited to 80/443. System proxies are deliberately
not used. Redirects, response sizes, refresh concurrency and request times are bounded.
No session cookies or OIDC credentials are sent to publishers.

The server stores feeds, stories, saved copies and preferences in apphost's SQLite
database in the app's private data root. There are no sidecars, startup downloads,
or original user files in the replaceable program directory. Apphost backups cover
the database. Refresh does not remove saved copies. Publisher failures preserve
the last successful feed. Users can remove sources independently of saved stories.

The service worker caches only known static app assets. It never caches API or
authentication responses. IndexedDB holds an account identifier, feed snapshot and
explicitly saved reading copies, never passwords, CSRF tokens or OIDC tokens.
Account changes and successful sign-out clear local reading data. Offline copies
are intentionally readable without contacting the server; see the README's shared
device limitation. Browser quotas/eviction apply. Saves fail visibly on quota errors.

Reading extraction is best effort. Only article/main elements from the ordinary
public HTML response are considered. No scripts are run and no authentication,
paywall circumvention, special archives or reader proxies are used.
