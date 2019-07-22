using System.Collections.Generic;
using CommandLine;

namespace Mavic
{
    internal class Program
    {
        /// <summary>
        ///     Process the property parsed commandline options.
        /// </summary>
        /// <param name="commandLineOptions">The options parsed.</param>
        private static void ProcessParsedArguments(CommandLineOptions commandLineOptions)
        {
            var scraper = new RedditScraper(commandLineOptions);
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
            Parser.Default.ParseArguments<CommandLineOptions>(args)
                .WithParsed(ProcessParsedArguments)
                .WithNotParsed(ProcessParseErrors);
        }
    }
}