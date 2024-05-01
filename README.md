# Websocket - Golang

[![CodeQL](https://github.com/gitnoober/simple-ws-golang/actions/workflows/codeql.yml/badge.svg)](https://github.com/gitnoober/simple-ws-golang/actions/workflows/codeql.yml)
[![Go](https://github.com/gitnoober/simple-ws-golang/actions/workflows/go.yml/badge.svg)](https://github.com/gitnoober/simple-ws-golang/actions/workflows/go.yml)
[![Dependency review](https://github.com/gitnoober/simple-ws-golang/actions/workflows/dependency-review.yml/badge.svg)](https://github.com/gitnoober/simple-ws-golang/actions/workflows/dependency-review.yml)
<hr>

## Features

- Listens for requests to the home page and renders the HTML page.
- Added handler for WS connection upgrade along with a route.
- Added JS code for handling event listeners for sending messages to all users.
- WebSocket listener constantly listens for payloads and acts on them based on the action type:
  - **Broadcast Message:** Broadcast message to all online users.
  - **Remove User:** Delete user from online users when the user leaves the webpage.
  - **Add User:** Add user to the list of currently online users.

<hr>

![Image Alt Text](static/image.png)

###### Note- Listens to port: 8080
