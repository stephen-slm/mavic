using System.Collections.Generic;
using CommandLine;

namespace Mavic
{
    internal class Program
    {
        /// <summary>
        /// </summary>
        /// <param name="commandLineOptions"></param>
        private static void ProcessParsedArguments(CommandLineOptions commandLineOptions)
        {
            var scraper = new Scraper(commandLineOptions);
            scraper.ProcessSubreddits().Wait();
        }

        /// <summary>
        /// </summary>
        /// <param name="errors"></param>
        private static void ProcessParseErrors(IEnumerable<Error> errors)
        {
        }

        private static void Main(string[] args)
        {
            Parser.Default.ParseArguments<CommandLineOptions>(args)
                .WithParsed(ProcessParsedArguments)
                .WithNotParsed(ProcessParseErrors);
        }
    }
}