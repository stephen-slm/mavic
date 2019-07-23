namespace Mavic.Types
{
    /// <summary>
    ///     A Author of a given reddit post.
    /// </summary>
    public class Author
    {
        /// <summary>
        ///     The name of the reddit user, commonly the reddit username.
        /// </summary>
        public string Name { get; set; }

        /// <summary>
        ///     The link directly to the reddit users profile page. u/username.
        /// </summary>
        public string Link { get; set; }
    }

    /// <summary>
    ///     A given image and basic meta data about said image posting on reddit.
    /// </summary>
    public class Image
    {
        /// <summary>
        ///     The id of the given image post.
        /// </summary>
        public string Id { get; set; }

        /// <summary>
        ///     THe image id, the ending part of the link.
        /// </summary>
        public string ImageId { get; set; }

        /// <summary>
        ///     The author of the given post.
        /// </summary>
        public Author Author { get; set; }

        /// <summary>
        ///     The source post link directory to reddit.
        /// </summary>
        public string PostLink { get; set; }

        /// <summary>
        ///     The link to the source image (e.g i.imgur.com)
        /// </summary>
        public string Link { get; set; }

        /// <summary>
        ///     The title of the given post.
        /// </summary>
        public string Title { get; set; }

        /// <summary>
        ///     The sub reddit that that the image was posted onto.
        /// </summary>
        public string Subreddit { get; set; }

        /// <summary>
        ///     The source in which the image is hosted, e.g imgur, reddit.
        /// </summary>
        public string Source { get; set; }
    }
}