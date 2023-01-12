
<p align="center">
    <img src="assets/go-sound.png" height="300px" style="border-radius:8px" alt="goandsoundcloud">
</p>

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://go.dev/dl/)

A simple GO project to download sound tracks from soundcloud.


# Installation

There are multiple ways to install, the easiest is :
```
go install github.com/AYehia0/soundcloud-dl@latest
```
other way is to grab is the source code and build.

# How to use ?

<p align="center">
    <img src="assets/soundcloud-dl-github.gif" height="300px" style="border-radius:8px" alt="goandsoundcloud">
</p>


```
Usage:
  sc <url> [flags]

Flags:
  -b, --best                   Download with the best available quality.
  -p, --download-path string   The download path where tracks are stored. (default "/home/<username>")
  -h, --help                   help for sc
  -s, --search                 Check if the track exists or not.
```
Note `search` flag isn't implemented yet!

# Features

- Download track with multiple qualities.
- Nice UI
- [Blazingly Fast](https://youtu.be/Z0GX2mTUtfo)
- More to be added.
