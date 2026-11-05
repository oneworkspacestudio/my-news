# My News implementation

User requested a Google News-style aggregator with categories, locations and
offline reading. Creation is authorized; public repository creation, commits,
pushes/releases/channel changes and owner-PC installation are not authorized.
No agents or Actions. Original My News branding, not Google assets.

App source is the new independent `I:/Projects/my-news` repository. The outer
OneWorkspace repo is unborn and must not be staged. The workspace-folder suite
checkout is older/dirty. Integration uses isolated `codex/my-news` worktree
`I:/Projects/workspace-my-news`, based on Books delivery HEAD 055d59c. Preserve
parent checkout's unrelated Books sandbox fixes, other tasks and owner installs.
Any future publication must reconcile current public main and other pending apps;
114/integration9 are local candidate contracts, not claimed public reservations.

App: Go apphost/SQLite, MyAccounts and standalone first-use setup, 18 default RSS
feeds across nine regions, category/location/search/followed topics, custom RSS/
Atom, safe plain-text reader, best-effort public article extraction on save,
excerpt/content/copy labels, IndexedDB offline library plus static-only service
worker. Exact URLs deduplicated after filtering; attribution retained. Offline
copies are per browser/profile and remain readable disconnected; sign-out and
online account change clear them. No paywall bypass or offline media downloading.

All external HTTP uses DNS-pinned public IP transport, no proxies/cookies, bounded
redirects/body/time/concurrency. APIs and SQL are account-scoped. AUI screens are
source; generated React is checked in, dist/node_modules excluded. Dependencies
are pinned. Actual public feeds all18 returned valid items in the recorded live
check, separate from deterministic synthetic native/browser fixtures.

Native tests cover RSS/Atom, script stripping, private-address/redirect denial,
CSRF, anonymous denial, account isolation, persisted saved text and feed removal.
Desktop1440/mobile390 Chrome test covers real native host+SQLite, refresh/search/
region/follow/custom feed, public reading extraction, offline reload/read, account
switch and logout. Screenshots use synthetic news only. Windows package and Linux
binary compile. Linux runtime/physical phone/reboot not tested.

Artifacts `I:/Projects/tooling/my-news-20261105`; dedicated synthetic NTFS temp/
fixture root `C:/Projects/.cache/my-news-20261105`. Stop task processes and verify
paths/reparse points before recursive cleanup. No owner installation is a fixture.
Final managed installer acceptance and evidence are recorded below when finished.

## Final local acceptance

Final build3 app Windows SHA256e2fb5a4b45cca9a0cd0a66a4b901acad2a0bed2bad56be20071625e00afca7a0;
archive dbe659af474e15d59daba4628f645d879bed94f3d3a8eb2074ef42616a84ff4e.
Final native Chrome checks pass, including account switching and correct extracted-
copy/excerpt labelling. All18 live sources returned valid items. Linux compiles only.

Normal suite3 builder passes38 apps/462 packages; all11 News-owned payloads rehashed
and final embedded EXE matches build3. Actual isolated installer acceptance passes
fresh, explicit omission, News-only2-test update/rollback and113→114 protected
upgrade. Real MyAccounts owner/issuer binding, shared picker/icon/category, app
reopening, anonymous API401, exact installed version/hash, saved text/preferences,
unchanged manager/other app contracts on app-only update and Files sentinel retained.
See docs/evidence (app) or docs/evidence/my-news (suite). No owner fixture/deployment.

Build incident: full installer link exhausted C temp space. Rerun used named I
GOTMPDIR and succeeded; NTFS runtime fixtures stayed on C. A broad vet collided
with bundle rewriting, then contended with installation disk I/O; changed-package
vet is the final scoped check. Full suite Go tests passed; app tests/vet passed.
No code guard was bypassed. The normal updater and rollback were used.

No commit/push/release/channel/owner install is authorized or performed. Public
release remains a later task: reconcile current main/other new apps, commit reviewed
source, create/pin real app release, use the signed Android publication entrypoint,
verify public channel and every explicitly requested PC. Never deploy these local
candidate/test-update artifacts as a published release. Cleanup result follows.

Cleanup complete: native fixture services were stopped, process references and
reparse points checked, and dedicated C test/private/temp roots plus I synthetic
program copies and disposable Go temp removed. No task service remains. Final
changed-package suite vet, News component/publisher tests and diff check pass.
Owner state and newest recovery remain untouched. Source/build3/suite3 and sanitised
evidence are retained. Acceptance docs were finalized after immutable local package
creation; their source copies are the current delivery record.

## Authorized publication and two-PC delivery

The user subsequently explicitly requested release and installation on both PCs.
Publication, channel promotion and normal protected owner updates are now authorized.
Suite rebased onto public d84c798, retaining Finance5, Voice8, MyAI24, Tube7 and Slicer8.
Workspace114 is already public; News integration9 now requires Workspace115.
Earlier candidate114 evidence remains historical; final public verification follows.

## Publication verified

My News1 is public from0ab0d4c; Windows archive SHA256
41d341da2f9316b26066672ddeef73cfebd74262c6ca1758b2b7394417e64abd,
with the exact native-tested executable e2fb5a4b45cc. Workspace115 is public
from08e175d after final fresh/omit/News-only-update/rollback/114-to115 acceptance.
All480 suite assets and permanent signed Android13 companion verified. An upload
duplicate was recovered by retaining matching draft assets and uploading only
missing files, then rerunning the mandatory publisher. No public assets replaced.

Concurrent Workspace116 retains all11 tested News payloads and115 source. Shared
channel merge af70e7b preserves116, News1, Commute1 and Finance6. Anonymous full
latest-bundle verification and both owner installations remain in progress; record
their actual completion before claiming both PCs delivered. All development and
acceptance tests used isolated synthetic identities/data, never owner fixtures.

Primary owner delivery completed115 + News1 at17:02 UTC: all11 payload hashes and
manager match public release; five native-served static assets match the tested
dist;37 services run; unselected apps/packages, identities/credentials, locations
and latest verified recovery preserved. Secondary verified116 repair is active;
its tunnel has not reconnected yet, and News installation remains pending there.
Do not repeat the active repair. The owner's browser requires manual dismissal of
a Windows permission prompt and MyAccounts sign-in before visual launch verification.

Public channel/download verification completed using the matching116 installer:
its successful no-token check was reused after comparing the checked binary to
the public hash and confirming News1 in the output. The older115 check correctly
refused the116 integration format. No published asset or release tag was changed.
