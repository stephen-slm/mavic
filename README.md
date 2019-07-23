<h1 align="center">
    <a href="https://github.com/tehstun/Mavic">
      <img src="./docs/img/logo.png" alt="mavic-logo" width="200">
    </a>
    <br/>
    <a href="https://github.com/tehstun/mavic">
      <img src="https://img.shields.io/badge/Mavic-v0.0.1-blue.svg" alt="mavic-Version">
    </a>
</h1>

<h4 align="center">Mavic is a CLI application designed to download direct images found on selected reddit subreddits.</h4>

<p align="center">
  <a>
    <img src="https://img.shields.io/badge/CommandLineParser-2.50.0-brightgreen.svg">
    <img src="https://img.shields.io/badge/Newtonsoft.Json-12.0.2-brightgreen.svg">
  </a>
</p>

<p align="center">
  <a href="#how-to-use">How To Use</a> •
  <a href="#releases">Releases</a> •
  <a href="#license">License</a>
</p>

# How to Use

Display basic help related information about the application for when you quickly need to understand possible options.

```
.\Mavic.exe --help

  -o, --output        (Default: ./) The output directory to store the images.

  -l, --limit         (Default: 50) The total number of posts max per sub-reddit.

  -f, --frontpage     (Default: false) If the front page should be scrapped or not.

  -t, --type          (Default: hot) What kind of page type reddit should be scraping, e.g hot, new, top, rising.

  -s, --subreddits    Required. What subreddits are going to be scrapped for downloading images.

  --help              Display this help screen.

  --version           Display version information.
```

Downloading all images from the last 50 r/cute currently on hot.

`.\Mavic.exe --subreddits cute`

Downloading all images from the top 25 r/cute, r/cats, r/aww into a picture folder.

`.\Mavic.exe --subreddits cute cats aww -l 25 --output ./pictures`

Downloading cat pictures and the front page images of the last 100 items.

`.\Mavic.exe -s cute -f --limit 100`

Downloading all top gifs from the top 100 r/gifs posts of all time.

`.\Mavic.exe -s gifs -l 100 --type top`

Downloading all cute and frontpage images of the hot 100 posts and ouputting to a pictures folder.

`.\Mavic.exe --subreddits cute -l 100 --output ./pictures -f`

<div align="center">
    <img src="./docs/img/home.gif" width="650" />
</div>

# Releases

Release information can be found here: https://github.com/tehstun/Mavic/releases

# License

Mavic is licensed with a MIT License.
