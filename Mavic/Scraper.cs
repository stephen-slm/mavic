using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading.Tasks;
using System.Xml.Linq;

namespace Mavic
{
    public class Scraper
    {
        /// <summary>
        ///     The options to be used throughout the parsing and scraping of the reddit site.
        /// </summary>
        private readonly CommandLineOptions _options;

        /// <summary>
        ///     The list of supported image file types that can be downloaded from reddit.
        /// </summary>
        private readonly List<string> _supportedFileTypes = new List<string>
            {"jpeg", "png", "gif", "apng", "tiff", "pdf", "xcf"};

        /// <summary>
        ///     The list of supported page types on reddit to be used.
        /// </summary>
        private readonly List<string> _supportedPageTypes = new List<string>
            {"hot", "new", "rising", "controversial", "top"};

        /// <summary>
        ///     Creates a new instance of the scraper with the command line options.
        /// </summary>
        /// <param name="options">The command line options</param>
        public Scraper(CommandLineOptions options)
        {
            _options = options;

            // if the limit goes outside the bounds of the upper and lower scopes, reset back to the 50 limit.
            if (_options.ImageLimit <= 0 || _options.ImageLimit > 500) _options.ImageLimit = 50;
        }

        /// <summary>
        ///     Process each subreddits nd start downloading all the images.
        /// </summary>
        public async Task ProcessSubreddits()
        {
            foreach (var subreddit in _options.Subreddits)
            {
                var feed = await GatherRedditRssFeed(subreddit);
                var links = ParseImgurLinksFromFeed(feed);
            }
        }

        /// <summary>
        ///     Downloads and parses the reddit XML rss feed into a XDocument based on the sub reddit and the limit.
        /// </summary>
        /// <param name="subreddit">The sub reddit being downloaded</param>
        /// <returns></returns>
        private async Task<XDocument> GatherRedditRssFeed(string subreddit)
        {
            if (string.IsNullOrEmpty(subreddit))
                throw new ArgumentException("sub reddit is required for downloading", "subreddit");

            using var httpClient = new HttpClient();

            var url = string.IsNullOrEmpty(_options.PageType) || _options.PageType == "hot" ||
                      !_supportedPageTypes.Contains(_options.PageType)
                ? $"https://www.reddit.com/r/{subreddit}/.rss?limit={_options.ImageLimit}&after=0"
                : $"https://www.reddit.com/r/{subreddit}/{_options.PageType}.rss?limit={_options.ImageLimit}&after=0";

            var source = await httpClient.GetAsync(url);
            var stringContent = await source.Content.ReadAsStringAsync();

            return XDocument.Parse(stringContent);
        }

        /// <summary>
        ///     Parses all the imgur links from the given reddit RSS feed.
        /// </summary>
        /// <param name="feed">The xml parsed from the rss feed</param>
        /// <returns></returns>
        private List<string> ParseImgurLinksFromFeed(XDocument feed)
        {
            return new List<string>();
        }
    }
}