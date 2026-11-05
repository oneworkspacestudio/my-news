// Generated from screens/App.aui; edit the AUI source.
import { Alert, Badge, Button, Card, Divider, Grid, Heading, Icon, Image, Input, Link, Page, Row, Section, Select, Spacer, Stack, Text } from "@/components/ui"

export function App({ data, onAction }: { data: any; onAction: (name: string) => void }) {
  return (
    <Page>
      <Row variant="topbar">
        <Link href="#briefing" variant="brand">
          <Image src="./news.svg" alt="" />
          <Text>My News</Text>
        </Link>
        <Row variant="search">
          <Icon name="Search" />
          <Input name="query" placeholder="Search stories, topics and publishers" label="Search news" />
        </Row>
        <Spacer />
        <Button onClick={() => onAction("theme")} variant="icon" label="Switch light or dark theme"><Icon name="SunMoon" /></Button>
        <Link href={data.workspace} variant="apps" label="Workspace apps"><Icon name="Grid2X2" /></Link>
        {data.loggedIn && (
          <Button onClick={() => onAction("logout")} variant="avatar" label="Sign out and remove offline data">{data.initials}</Button>
        )}
      </Row>
      {data.message && (
        <Alert variant={data.messageKind}>{data.message}</Alert>
      )}
      {data.loggedIn ? (
        <Grid variant="layout">
          <Stack variant="sidebar">
            <Text variant="nav-label">YOUR NEWSROOM</Text>
            {data.navigation.map((item, i) => (
              <Button key={i} onClick={() => onAction("navigate")} item={item.id} variant={item.variant} pressed={item.active}>
                <Icon name={item.icon} />
                <Text>{item.label}</Text>
              </Button>
            ))}
            <Divider />
            <Text variant="nav-label">EXPLORE TOPICS</Text>
            {data.topics.map((item, i) => (
              <Button key={i} onClick={() => onAction("topic")} item={item.id} variant={item.variant} pressed={item.active}>
                <Icon name={item.icon} />
                <Text>{item.id}</Text>
              </Button>
            ))}
            <Spacer />
            <Stack variant="sidebar-note">
              <Icon name="Globe2" size={28} />
              <Text>A world of perspectives.</Text>
              <Text variant="muted">Independent sources. Your own reading list.</Text>
            </Stack>
          </Stack>
          <Stack variant="main">
            {data.offline && (
              <Alert variant="offline">You’re offline. Your saved reading copies are ready to read.</Alert>
            )}
            {data.reader ? (
              <Stack variant="reader">
                <Row variant="reader-actions">
                  <Button onClick={() => onAction("closeReader")} variant="secondary">
                    <Icon name="ArrowLeft" />
                    <Text>Back to stories</Text>
                  </Button>
                  <Spacer />
                  <Button onClick={() => onAction("fontSmaller")} variant="icon" label="Smaller reading text"><Icon name="Minus" /></Button>
                  <Button onClick={() => onAction("fontLarger")} variant="icon" label="Larger reading text"><Icon name="Plus" /></Button>
                </Row>
                <Row variant="story-meta">
                  <Badge>{data.reader.category}</Badge>
                  <Text>{data.reader.region}</Text>
                  <Text>{data.reader.dateLabel}</Text>
                </Row>
                <Heading level={1}>{data.reader.title}</Heading>
                <Row variant="byline">
                  <Icon name="Newspaper" />
                  <Text>{data.reader.source}</Text>
                  <Spacer />
                  <Button onClick={() => onAction("saveReader")} variant="primary" disabled={data.busy}>
                    <Icon name="Download" />
                    <Text>{data.reader.saveLabel}</Text>
                  </Button>
                </Row>
                <Alert variant="reading-note">{data.reader.readingNote}</Alert>
                <Section variant="article-body">
                  {data.reader.paragraphs.map((item, i) => (
                    <Text key={i} variant="body">{item}</Text>
                  ))}
                </Section>
                <Link href={data.reader.url} external={true} variant="publisher">
                  <Text>Read original at publisher</Text>
                  <Icon name="ArrowUpRight" />
                </Link>
              </Stack>
            ) : (
              <>
                <Row variant="page-heading">
                  <Stack>
                    <Text variant="eyebrow">{data.date}</Text>
                    <Heading level={1}>{data.title}</Heading>
                    <Text variant="subtitle">{data.subtitle}</Text>
                  </Stack>
                  <Spacer />
                  <Button onClick={() => onAction("refresh")} variant="secondary" disabled={data.refreshDisabled} label="Refresh news">
                    <Icon name="RefreshCw" />
                    <Text>{data.refreshLabel}</Text>
                  </Button>
                </Row>
                {data.showSources ? (
                  <>
                    <Card variant="source-form">
                      <Heading level={2}>Bring your own perspective</Heading>
                      <Text variant="muted">Add a publisher’s RSS or Atom feed. Local, international, or something wonderfully specific.</Text>
                      <Grid variant="form">
                        <Input name="sourceName" placeholder="Publisher name" label="Publisher name" />
                        <Input name="sourceURL" placeholder="https://example.com/feed.xml" label="RSS or Atom feed URL" />
                        <Input name="sourceRegion" placeholder="City, country or region" label="Source location" />
                        <Select name="sourceCategory" options={data.categoryOptions} label="Source category" />
                      </Grid>
                      <Button onClick={() => onAction("addSource")} variant="primary" disabled={data.busy}>
                        <Icon name="Plus" />
                        <Text>Add source</Text>
                      </Button>
                    </Card>
                    <Stack variant="sources">
                      {data.sources.map((item, i) => (
                        <Card key={i} variant="source">
                          <Row>
                            <Icon name="Rss" />
                            <Stack>
                              <Heading level={3}>{item.name}</Heading>
                              <Text variant="muted">{item.description}</Text>
                            </Stack>
                            <Spacer />
                            <Button onClick={() => onAction("removeSource")} item={item.id} variant="icon" label="Remove source" disabled={data.offline}><Icon name="Trash2" /></Button>
                          </Row>
                          {item.error ? (
                            <Text variant="error">{item.error}</Text>
                          ) : (
                            <Text variant="caption">{item.status}</Text>
                          )}
                        </Card>
                      ))}
                    </Stack>
                  </>
                ) : (
                  <>
                    {data.showBriefing && (
                      <Card variant="welcome">
                        <Stack>
                          <Badge variant="welcome">THE BIGGER PICTURE</Badge>
                          <Heading level={2}>Your world, in perspective.</Heading>
                          <Text>Explore the stories shaping our world. Make room for a different point of view.</Text>
                          <Row variant="welcome-stats">
                            <Text>{data.sourceCount}</Text>
                            <Text>•</Text>
                            <Text>{data.regionCount}</Text>
                          </Row>
                        </Stack>
                        <Section variant="globe"><Icon name="Globe2" size={150} /></Section>
                      </Card>
                    )}
                    <Row variant="filters">
                      <Row variant="location">
                        <Icon name="MapPin" />
                        <Select name="region" options={data.regionOptions} label="News location" />
                      </Row>
                      <Select name="sort" options={data.sortOptions} label="Sort stories" />
                      <Spacer />
                      {data.showSaved ? (
                        <Button onClick={() => onAction("downloadSaved")} variant="secondary" disabled={data.refreshDisabled}>
                          <Icon name="Download" />
                          <Text>Download library to this device</Text>
                        </Button>
                      ) : (
                        <Button onClick={() => onAction("followTopic")} variant="secondary">
                          <Icon name="Star" />
                          <Text>{data.followLabel}</Text>
                        </Button>
                      )}
                    </Row>
                    <Row variant="section-heading">
                      <Heading level={2}>{data.sectionTitle}</Heading>
                      <Text variant="muted">{data.resultLabel}</Text>
                    </Row>
                    {data.empty ? (
                      <Card variant="empty">
                        <Icon name={data.emptyIcon} size={42} />
                        <Heading level={2}>{data.emptyTitle}</Heading>
                        <Text>{data.emptyText}</Text>
                        {data.canRefresh && (
                          <Button onClick={() => onAction("refresh")} variant="primary" disabled={data.busy}>Refresh news</Button>
                        )}
                      </Card>
                    ) : (
                      <>
                        <Grid variant="stories">
                          {data.stories.map((item, i) => (
                            <Card key={i} variant={item.cardVariant}>
                              <Button onClick={() => onAction("read")} item={item.id} variant="story-open" label={item.readLabel}>
                                <Row variant="story-meta">
                                  <Text variant="source">{item.source}</Text>
                                  <Spacer />
                                  <Icon name={item.icon} />
                                </Row>
                                <Badge variant="category">{item.category}</Badge>
                                <Heading level={2}>{item.title}</Heading>
                                <Text variant="summary">{item.summary}</Text>
                              </Button>
                              <Row variant="story-footer">
                                <Text variant="caption">{item.region}</Text>
                                <Text variant="caption">{item.dateLabel}</Text>
                                <Spacer />
                                <Button onClick={() => onAction("save")} item={item.id} variant={item.saveVariant} label={item.saveLabel} disabled={data.busy}><Icon name={item.saveIcon} /></Button>
                              </Row>
                              {item.offlineLabel && (
                                <Text variant="offline-label">{item.offlineLabel}</Text>
                              )}
                            </Card>
                          ))}
                        </Grid>
                        {data.hasMore && (
                          <Button onClick={() => onAction("more")} variant="secondary">Show more stories</Button>
                        )}
                      </>
                    )}
                  </>
                )}
                <Row variant="footer">
                  <Text>MY NEWS</Text>
                  <Text>{data.footer}</Text>
                  <Spacer />
                  <Text>A little more informed. A little more connected.</Text>
                </Row>
              </>
            )}
          </Stack>
        </Grid>
      ) : (
        <Stack variant="login">
          <Image src="./news.svg" alt="My News" />
          <Text variant="eyebrow">YOUR WORLD, IN PERSPECTIVE</Text>
          <Heading level={1}>{data.loginTitle}</Heading>
          <Text>A personal newsroom. Stories from around the world, ready whenever you are.</Text>
          {data.oidc ? (
            <Link href="./auth/login" variant="primary">Continue with MyAccounts</Link>
          ) : (
            <>
              <Input name="email" type="email" placeholder="Email address" label="Email address" />
              <Input name="password" type="password" placeholder="Password" label="Password" onClick={() => onAction("login")} />
              <Button onClick={() => onAction("login")} variant="primary" disabled={data.busy}>{data.loginLabel}</Button>
            </>
          )}
          <Text variant="muted">Saved reading copies stay on this device until you sign out or clear them.</Text>
        </Stack>
      )}
    </Page>
  )
}
