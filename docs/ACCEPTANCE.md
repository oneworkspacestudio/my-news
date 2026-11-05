# Delivery status

My News **0.1.0-preview.1** is [publicly released](https://github.com/oneworkspacestudio/my-news/releases/tag/v0.1.0-preview.1). The final Windows
executable has passed native browser and managed-installer acceptance. Linux x64
compiles; Linux runtime, physical phones, reboot and clean-host tests are not claimed.
The app source is published at commit `0ab0d4c`. Workspace **0.3.0-preview.115**
has completed final installer acceptance and is
[publicly released](https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.115)
through the required signed-Android entry point. All480 public asset hashes match.
The channel also includes the concurrent Workspace116 release, whose eleven
News payloads exactly match the tested115 catalog. Owner-PC delivery is recorded
separately below when verified.

The matching public116 installer completed an anonymous channel/bundle check,
including News1. The115 installer correctly refused116's newer integration format.
The first owner PC completed the normal protected115 + News1 update: all11 News
payloads and the manager match the release, five native-served static assets match
the tested app, and37 services run without new failures. Existing app contracts,
identities/credentials, paths and the verified recovery copy are preserved.
The second PC reconnected after its protected reinstall request, with all36 existing
services healthy. Its final updater receipt/version and My News installation still
need the authenticated Updates UI; a Windows permission prompt and MyAccounts
sign-in require the user's input. Owner PCs were deployment targets, never fixtures.

## Verified

- Unit/native API checks: RSS and Atom, unsafe markup/URL rejection, private-network
  and redirect blocking, CSRF, anonymous denial, account isolation and saved storage.
- Chrome at1440px and390px: refresh, category/location/search, followed topics,
  custom feeds, public reading extraction, offline reload/reading, account switching,
  and logout clearing downloads. No browser script errors. Synthetic publishers.
- All18 default live RSS sources returned valid stories in the recorded check.
- Normal complete Workspace **0.3.0-preview.115**, installer integration9:
  **38 apps /462 packages**. All11 News-owned payloads match their catalog hashes.
- Actual installer on isolated NTFS data: fresh setup and explicit News omission;
  real MyAccounts sign-in/reopen, shared picker/category/icon, anonymous API denial.
- News-only update to a synthetic2-test binary and rollback to1: saved reading copy
  and settings preserved; manager/other app contracts unchanged; exact installed
  executable version/hash and rollback bytes checked.
- Supported114-to115 protected upgrade: synthetic owner identity, Drive contract and an
  original synthetic Files sentinel retained, new default News installed and usable.

Evidence: [browser](browser-acceptance.json), [live sources](live-source-check.json),
[final public115 installer paths](evidence/public115-managed-acceptance.json),
[earlier candidate paths](evidence/managed-acceptance.json), [catalog](evidence/final-catalog.json).
Screenshots contain synthetic news, not bundled production headlines.

## Artifacts and boundaries

The public Windows archive is `my-news-0.1.0-preview.1-windows-amd64.zip`.
Windows executable SHA256:
`e2fb5a4b45cca9a0cd0a66a4b901acad2a0bed2bad56be20071625e00afca7a0`.

Managed contract: news/My News, `/news`, MyAccounts, canonical9030/rebased+40,
private `news` data, General/default-on, single executable, news-* ownership and
app recovery scope. Workspace115 pins the verified public My News1 archive.
Earlier candidate114 evidence is retained as historical local acceptance;
the final public115 acceptance uses the Windows progress-write fix and exact
release installer bytes.

Offline copies contain feed text or best-effort text from the ordinary public page.
Paywalled/unavailable/script-rendered pages may remain excerpts; media is not
saved. Offline data is readable by someone with access to the browser profile.
Browser storage eviction/clearing removes downloads. See README and SECURITY.

All tests used synthetic accounts and dedicated test roots. Owner applications,
data, credentials, files, autostart and recovery copies were not test fixtures.
GitHub Actions remain disabled. Published assets are immutable. Test services
and dedicated synthetic installation roots were cleaned after acceptance.
