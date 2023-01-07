# TODO

- Before Starting:

    - [x] Learn the basics of go lang.
    - Learn working with streams and files
    - Learn how to send requests
    - Learn how to structure the code in go
    - Learn some OOP
    - Learn Async and multi-threading

- Get the client ID ?
    - Make a GET request to : `https://soundcloud.com`
    - Find the last asset : `https://a-v2.sndcdn.com/assets/*`
    - Make a get request to that url and extract the `client_id`

- Validate client_id ?
    - Make a get request to `/me` and check the status response.

- Get auth_token ?
    - Find your OAuth token by visiting SoundCloud after logging in and watching any of the browsers requests to the SoundCloud API, the token will be under the `Authorization` header of any of these requests

    - soundcloud uses OAuth so the Authorization header should be as `OAuth {auth_token}`
- Download a track ?
