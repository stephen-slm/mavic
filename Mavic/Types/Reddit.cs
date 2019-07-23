using System;
using System.Collections.Generic;
using Newtonsoft.Json;

namespace Mavic.Types
{
    /// <summary>
    /// </summary>
    public class RedditListing
    {
        /// <summary>
        /// </summary>
        [JsonProperty("data")]
        public RedditImageData Data { get; set; }
    }

    /// <summary>
    /// </summary>
    public class RedditImageData
    {
        /// <summary>
        /// </summary>
        [JsonProperty("children")]
        public List<RedditImageChild> Children { get; set; }
    }

    /// <summary>
    /// </summary>
    public class RedditImageChild
    {
        /// <summary>
        /// </summary>
        [JsonProperty("data")]
        public RedditImageChildData Data { get; set; }
    }

    /// <summary>
    /// </summary>
    public class RedditImageChildData
    {
        /// <summary>
        /// </summary>
        [JsonProperty("title")]
        public string Title { get; set; }

        /// <summary>
        /// </summary>
        [JsonProperty("domain")]
        public string Domain { get; set; }

        /// <summary>
        /// </summary>
        [JsonProperty("id")]
        public string Id { get; set; }

        /// <summary>
        /// </summary>
        [JsonProperty("author")]
        public string Author { get; set; }

        /// <summary>
        /// </summary>
        [JsonProperty("permalink")]
        public string Permalink { get; set; }

        /// <summary>
        /// </summary>
        [JsonProperty("post_hint")]
        public string PostHint { get; set; }

        /// <summary>
        /// </summary>
        [JsonProperty("url")]
        public Uri Url { get; set; }

        /// <summary>
        /// </summary>
        [JsonProperty("subreddit")]
        public string Subreddit { get; set; }
    }
}