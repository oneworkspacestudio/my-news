# Source and implementation notes

Publisher documentation checked during implementation:

- BBC feed directory: https://support.bbc.co.uk/platform/feeds/NewsFeeds.htm
- Guardian RSS guidance: https://www.theguardian.com/help/feeds
- Le Monde's section-feed examples: https://www.lemonde.fr/en/rss/article/2022/06/23/the-rss-feeds-of-monde-fr_5987755_165.html

The starting catalog is ordinary RSS, not an undocumented Google News scraping API.
Feed availability and coverage can change. No API key is needed. The app's Sources
page reports per-feed refresh failures and retains older stories. The acceptance
record distinguishes deterministic synthetic tests from live source checks.

Categories/locations come from the publisher feed's configured section, not from
an unreliable inferred geolocation or AI classifier. A global Technology feed is
tagged World; a UK news feed is tagged United Kingdom. Custom feeds can use any
city/country/region label. Exact duplicate article URLs are suppressed in the view
after filtering so region/section membership is preserved.

Screens are authored in AUI and compiled to React. Original branding and CSS are
used; Google logos, names and proprietary assets are not copied. Dependencies:
React, Lucide, Vite/AUI compiler and Go's shared Workspace apphost. No third-party
executable is downloaded at application startup.
