package handlers

const defCfg = `# yaml-language-server: $schema=https://raw.githubusercontent.com/kociumba/kserver/main/.kserver

port: 8080

handlers:
- route: /
  content: ./index.html
  contentType: text/html

`
