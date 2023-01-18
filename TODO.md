# TODO
- [X] Add a flag `--highest-quality` : Get the best quality.
- [ ] ~~Add a flag `--get-sizes` : Get all the sizes of a track.~~
- [X] Add a flag `--download-path` : path to the download location.
- [X] Add support for downloading `hls`.
- [X] Add the metadata to the track after downloading.
- [X] Add Download track through search.
- [ ] Add auto-completion for the flags.
- [X] Use the soundcloud api call ~~instead of using `goquery`.~~, `goquery` is still used to fetch the `client_id`.
- [ ] Change the architecture, and re-org the structure.
- [X] Download a playlist :
    - [X] Check if the url is a playlist, ~~prompt the user Y/N to continue if `--playlist` isn't passed.~~ Automatically detect if the url is a playlist.
- [ ] Config file to save settings.

# Maybe not going TODO
- [ ] Load the waveform json data and stream directly into the terminal.
