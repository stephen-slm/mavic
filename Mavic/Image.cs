namespace Mavic
{
    public class Author
    {
        /// <summary>
        /// </summary>
        public string Name { get; set; }

        /// <summary>
        /// </summary>
        public string Link { get; set; }
    }

    public class Image
    {
        /// <summary>
        /// </summary>
        public string Id { get; set; }

        /// <summary>
        /// </summary>
        public Author Author { get; set; }

        /// <summary>
        /// </summary>
        public string PostLink { get; set; }

        /// <summary>
        /// </summary>
        public string Link { get; set; }

        /// <summary>
        /// </summary>
        public string Title { get; set; }

        /// <summary>
        /// </summary>
        public string Category { get; set; }
    }
}