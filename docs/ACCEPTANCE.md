# Delivery status

My News **0.1.0-preview.1** is implemented and packaged locally. The final Windows
executable has passed native browser and managed-installer acceptance. Linux x64
compiles; Linux runtime, physical phones, reboot and clean-host tests are not claimed.
No source commit/push, public release/channel promotion or owner-PC installation
was performed. The public app repository has not been created by this task.

## Verified

- Unit/native API checks: RSS and Atom, unsafe markup/URL rejection, private-network
  and redirect blocking, CSRF, anonymous denial, account isolation and saved storage.
- Chrome at1440px and390px: refresh, category/location/search, followed topics,
  custom feeds, public reading extraction, offline reload/reading, account switching,
  and logout clearing downloads. No browser script errors. Synthetic publishers.
- All18 default live RSS sources returned valid stories in the recorded check.
- Normal complete Workspace candidate **0.3.0-preview.114**, installer integration9:
  **38 apps /462 packages**. All11 News-owned payloads match their catalog hashes.
- Actual installer on isolated NTFS data: fresh setup and explicit News omission;
  real MyAccounts sign-in/reopen, shared picker/category/icon, anonymous API denial.
- News-only update to a synthetic2-test binary and rollback to1: saved reading copy
  and settings preserved; manager/other app contracts unchanged; exact installed
  executable version/hash and rollback bytes checked.
- Supported113-to114 protected upgrade: owner identity, Drive contract and an
  original synthetic Files sentinel retained, new default News installed and usable.

Evidence: [browser](browser-acceptance.json), [live sources](live-source-check.json),
[installer paths](evidence/managed-acceptance.json), [catalog](evidence/final-catalog.json).
Screenshots contain synthetic news, not bundled production headlines.

## Artifacts and boundaries

Final app package: `I:/Projects/tooling/my-news-20261105/build3/my-news-windows-amd64.zip`.
The companion suite candidate is under `suite3` in the same task artifact directory.
Windows executable SHA256:
`e2fb5a4b45cca9a0cd0a66a4b901acad2a0bed2bad56be20071625e00afca7a0`.

Managed contract: news/My News, `/news`, MyAccounts, canonical9030/rebased+40,
private `news` data, General/default-on, single executable, news-* ownership and
app recovery scope. Publication requires a real myNews component pin and the
mandatory signed-Android suite publisher. Reconcile current public main and
other pending app work before choosing a public version;114 is a local candidate.

Offline copies contain feed text or best-effort text from the ordinary public page.
Paywalled/unavailable/script-rendered pages may remain excerpts; media is not
saved. Offline data is readable by someone with access to the browser profile.
Browser storage eviction/clearing removes downloads. See README and SECURITY.

All tests used synthetic accounts and dedicated test roots. Owner applications,
data, credentials, files, autostart and recovery copies were not test fixtures.
No Actions or public assets were changed. Test services are stopped before cleanup.
