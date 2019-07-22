using System.Collections.Generic;
using CommandLine;

namespace Mavic
{
    internal class Program
    {
        /// <summary>
        ///     Process the property parsed commandline options.
        /// </summary>
        /// <param name="scrapingOptions">The options parsed.</param>
        private static void ProcessParsedArguments(ScrapingOptions scrapingOptions)
        {
            var scraper = new RedditScraper(scrapingOptions);
            scraper.ProcessSubreddits().Wait();
        }

        /// <summary>
        /// </summary>
        /// <param name="errors"></param>
        private static void ProcessParseErrors(IEnumerable<Error> errors)
        {
            // ignore, generic help message output.
        }

        /// <summary>
        ///     The main entry point of the Mavic application.
        /// </summary>
        /// <param name="args">Standard arguments to be parsed.</param>
        private static void Main(string[] args)
        {
            Parser.Default.ParseArguments<ScrapingOptions>(args)
                .WithParsed(ProcessParsedArguments)
                .WithNotParsed(ProcessParseErrors);
        }
    }
}