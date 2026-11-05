package news

var defaults = []Source{
	{ID: "bbc-world", Name: "BBC News", URL: "https://feeds.bbci.co.uk/news/world/rss.xml", Region: "World", Category: "Headlines"},
	{ID: "guardian-world", Name: "The Guardian", URL: "https://www.theguardian.com/world/rss", Region: "World", Category: "Headlines"},
	{ID: "aljazeera", Name: "Al Jazeera", URL: "https://www.aljazeera.com/xml/rss/all.xml", Region: "World", Category: "Headlines"},
	{ID: "bbc-uk", Name: "BBC News · UK", URL: "https://feeds.bbci.co.uk/news/uk/rss.xml", Region: "United Kingdom", Category: "Headlines"},
	{ID: "bbc-us", Name: "BBC News · US & Canada", URL: "https://feeds.bbci.co.uk/news/world/us_and_canada/rss.xml", Region: "North America", Category: "Headlines"},
	{ID: "bbc-europe", Name: "BBC News · Europe", URL: "https://feeds.bbci.co.uk/news/world/europe/rss.xml", Region: "Europe", Category: "Headlines"},
	{ID: "bbc-asia", Name: "BBC News · Asia", URL: "https://feeds.bbci.co.uk/news/world/asia/rss.xml", Region: "Asia", Category: "Headlines"},
	{ID: "bbc-africa", Name: "BBC News · Africa", URL: "https://feeds.bbci.co.uk/news/world/africa/rss.xml", Region: "Africa", Category: "Headlines"},
	{ID: "bbc-middleeast", Name: "BBC News · Middle East", URL: "https://feeds.bbci.co.uk/news/world/middle_east/rss.xml", Region: "Middle East", Category: "Headlines"},
	{ID: "bbc-latin", Name: "BBC News · Latin America", URL: "https://feeds.bbci.co.uk/news/world/latin_america/rss.xml", Region: "Latin America", Category: "Headlines"},
	{ID: "guardian-australia", Name: "The Guardian · Australia", URL: "https://www.theguardian.com/australia-news/rss", Region: "Oceania", Category: "Headlines"},
	{ID: "bbc-business", Name: "BBC Business", URL: "https://feeds.bbci.co.uk/news/business/rss.xml", Region: "World", Category: "Business"},
	{ID: "guardian-tech", Name: "The Guardian · Technology", URL: "https://www.theguardian.com/uk/technology/rss", Region: "World", Category: "Technology"},
	{ID: "bbc-science", Name: "BBC Science", URL: "https://feeds.bbci.co.uk/news/science_and_environment/rss.xml", Region: "World", Category: "Science"},
	{ID: "bbc-health", Name: "BBC Health", URL: "https://feeds.bbci.co.uk/news/health/rss.xml", Region: "World", Category: "Health"},
	{ID: "guardian-sport", Name: "The Guardian · Sport", URL: "https://www.theguardian.com/sport/rss", Region: "World", Category: "Sports"},
	{ID: "guardian-culture", Name: "The Guardian · Culture", URL: "https://www.theguardian.com/culture/rss", Region: "World", Category: "Entertainment"},
	{ID: "guardian-climate", Name: "The Guardian · Environment", URL: "https://www.theguardian.com/environment/rss", Region: "World", Category: "Environment"},
}
