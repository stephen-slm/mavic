using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Net;
using System.Net.Http;
using System.Threading.Tasks;
using Mavic.Types;
using Newtonsoft.Json;

namespace Mavic
{
    public class RedditScraper
    {
        /// <summary>
        ///     After is used for paging when the user selects more than 100 items to search for.
        /// </summary>
        private readonly int _after = 0;

        /// <summary>
        ///     The options to be used throughout the parsing and scraping of the reddit site.
        /// </summary>
        private readonly ScrapingOptions _scrapingOptions;

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
        ///     A unique list of image ids to stop duplicate images being downloaded for a given subreddit.
        /// </summary>
        private readonly Dictionary<string, HashSet<string>> _uniqueImages = new Dictionary<string, HashSet<string>>();

        /// <summary>
        ///     Creates a new instance of the scraper with the command line options.
        /// </summary>
        /// <param name="scrapingOptions">The command line options</param>
        public RedditScraper(ScrapingOptions scrapingOptions)
        {
            this._scrapingOptions = scrapingOptions;

            if (this._scrapingOptions.ImageLimit > 100)
            {
                Console.Out.WriteLine("Option 'limit' is currently enforced to 100 or less due to a on going problem");
                this._scrapingOptions.ImageLimit = 100;
            }

            // if the limit goes outside the bounds of the upper and lower scopes, reset back to the 50 limit.
            if (this._scrapingOptions.ImageLimit <= 0 || this._scrapingOptions.ImageLimit > 500)
                this._scrapingOptions.ImageLimit = 50;

            if (this._scrapingOptions.FrontPage)
                this._scrapingOptions.Subreddits = this._scrapingOptions.Subreddits.Prepend("frontpage");
        }

        /// <summary>
        ///     Process each subreddits nd start downloading all the images.
        /// </summary>
        public async Task ProcessSubreddits()
        {
            foreach (var subreddit in this._scrapingOptions.Subreddits)
            {
                // if we have not already done the subreddit before, then create a new unique entry into the unique
                // images list to keep track of all the already downloaded images by imgur image id.
                if (!this._uniqueImages.ContainsKey(subreddit))
                    this._uniqueImages.Add(subreddit, new HashSet<string>());

                try
                {
                    var feed = await this.GatherRedditFeed(subreddit);
                    var links = ParseImgurLinksFromFeed(feed);

                    var directory = Path.Combine(this._scrapingOptions.OutputDirectory, subreddit);

                    if (!Directory.Exists(directory)) Directory.CreateDirectory(directory);

                    Console.Out.WriteLine();
                    Console.Out.WriteLine($"Downloading {links.Count()} images from /r/{subreddit}");

                    foreach (var image in links)
                    {
                        if (string.IsNullOrEmpty(image.ImageId) ||
                            this._uniqueImages[subreddit].Contains(image.ImageId)) continue;

                        this._uniqueImages[subreddit].Add(image.ImageId);

                        Console.Out.WriteLine($"Downloading {image.ImageId} from /r/{image.Subreddit}");

                        await DownloadImage(directory, image);
                    }
                }
                catch (Exception)
                {
                    // ignored
                }
            }
        }

        /// <summary>
        ///     Downloads a given image based on its properties.
        /// </summary>
        /// <param name="outputDirectory">The output directory for the image to be stored</param>
        /// <param name="image">The image being downloaded</param>
        /// <returns></returns>
        private static async Task DownloadImage(string outputDirectory, Image image)
        {
            // replace gifv with mp4 for a preferred download as gifv files do not work really well on windows/desktop
            // machines but require additional processing, while mp4 will be file.
            if (image.Link.EndsWith("gifv")) image.Link = image.Link.Substring(0, image.Link.Length - 4) + "mp4";
            var imageImgurId = image.Link.Split("/").Last();

            var imageFullPath = Path.Combine(outputDirectory, imageImgurId);
            if (File.Exists(imageFullPath)) return;

            using var webClient = new WebClient();

            try
            {
                // the image is probably a collection of images, which cannot be downloaded as of yet.
                if (string.IsNullOrEmpty(Path.GetExtension(imageFullPath))) return;
                await webClient.DownloadFileTaskAsync(new Uri(image.Link), imageFullPath);
            }
            catch (Exception)
            {
                // ignored
            }
        }

        /// <summary>
        ///     Downloads and parses the reddit XML rss feed into a XDocument based on the sub reddit and the limit.
        /// </summary>
        /// <param name="subreddit">The sub reddit being downloaded</param>
        /// <returns></returns>
        private async Task<RedditListing> GatherRedditFeed(string subreddit)
        {
            Debug.Assert(!string.IsNullOrEmpty(subreddit));

            if (string.IsNullOrEmpty(subreddit))
                throw new ArgumentException("sub reddit is required for downloading", "subreddit");

            using var httpClient = new HttpClient();

            var url = string.IsNullOrEmpty(this._scrapingOptions.PageType) || this._scrapingOptions.PageType == "hot" ||
                      !this._supportedPageTypes.Contains(this._scrapingOptions.PageType)
                ? $"https://www.reddit.com/r/{subreddit}/.json?limit={this._scrapingOptions.ImageLimit}&after={this._after}"
                : $"https://www.reddit.com/r/{subreddit}/{this._scrapingOptions.PageType}.json?limit={this._scrapingOptions.ImageLimit}&after={this._after}";

            if (subreddit == "frontpage")
                url = $"https://www.reddit.com/.json?limit={this._scrapingOptions.ImageLimit}&after={this._after}";

            var source = await httpClient.GetAsync(url);

            var stringContent = await source.Content.ReadAsStringAsync();
            return JsonConvert.DeserializeObject<RedditListing>(stringContent);
        }

        /// <summary>
        ///     Parses all the imgur links from the given reddit RSS feed.
        /// </summary>
        /// <param name="redditListing">The xml parsed from the rss feed</param>
        /// <returns></returns>
        private static IEnumerable<Image> ParseImgurLinksFromFeed(RedditListing redditListing)
        {
            // ensure that the feed is not null, returning a empty list if is.
            if (redditListing == null) return new List<Image>();
            var possibleDataImages = redditListing.Data.Children.Where(e => e.Data.Domain.Contains("imgur")).ToList();

            var linkImages = new List<Image>();
            foreach (var possibleDataImage in possibleDataImages)
            {
                if (possibleDataImage.Data.Url == null || !possibleDataImage.Data.Url.Host.Contains("imgur")) continue;
                if (!possibleDataImage.Data.Url.LocalPath.Contains(".")) continue;

                linkImages.Add(new Image
                {
                    Author = new Author
                    {
                        Link = $"https://www.reddit.com/user/{possibleDataImage.Data.Author}",
                        Name = possibleDataImage.Data.Author
                    },
                    Id = possibleDataImage.Data.Id,
                    ImageId = possibleDataImage.Data.Url.AbsoluteUri.Split("/").Last().Split(".").First(),
                    Link = possibleDataImage.Data.Url.AbsoluteUri,
                    Subreddit = possibleDataImage.Data.Subreddit,
                    PostLink = possibleDataImage.Data.Permalink,
                    Title = possibleDataImage.Data.Title
                });
            }

            return linkImages;
        }
    }
}