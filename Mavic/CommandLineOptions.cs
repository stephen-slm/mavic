using System.Collections.Generic;
using CommandLine;

namespace Mavic
{
    public class CommandLineOptions
    {
        /// <summary>
        ///     The directory in which we will be downloading all the images into, based on the folder name of the given
        ///     sub-reddit.
        /// </summary>
        [Option('o', "output", Required = false, HelpText = "The output directory to store the images.",
            Default = "./")]
        public string OutputDirectory { get; set; }

        /// <summary>
        ///     The total number of images to download max per sub-reddit before we continue to the next one.
        /// </summary>
        [Option('l', "limit", Required = false, HelpText = "The total number of posts max per sub-reddit",
            Default = 50)]
        public int ImageLimit { get; set; }

        /// <summary>
        ///     If set to true, the tool will scrape the front page of reddit for the current most active sub-reddits and
        ///     then scrape all the imgur links directly from them sub-reddits.
        /// </summary>
        [Option('f', "front", Required = false, HelpText = "If the front page should be scrapped or not.",
            Default = true)]
        public bool FrontPage { get; set; }

        /// <summary>
        ///     You can change this to adjust on what kind of images you get from Reddits filtering options (hot, new, rising,
        ///     controversial, top), hot is the default by reddit while also the default in the tool.
        /// </summary>
        [Option('t', "type", Required = false, HelpText = "What kind of page type reddit should be scraping, e.g hot",
            Default = "hot")]
        public string PageType { get; set; }

        /// <summary>
        ///     What subreddits are going to be scrapped for downloading of sad images.
        /// </summary>
        [Option('s', "subreddits", Required = true,
            HelpText = "What subreddits are going to be scrapped for downloading imgur images")]
        public IEnumerable<string> Subreddits { get; set; }
    }
}