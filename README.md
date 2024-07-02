# Kserver
 
A proof of concept server that can be configured with YAML ⚙️ or LUA 🌕

## Usage

For YAML create a `kserver.yml` file in the directory you want to serve, 
the file should look like this:

```yaml copy
# yaml-language-server: $schema=https://raw.githubusercontent.com/kociumba/kserver/main/.kserver

port: 8080

handlers:
- route: /
  content: ./index.html
  contentType: text/html
# Add more routes here, in VSCode you can use ctrl+space to autocomplete the array
```

>[!NOTE]
> There is a YAML schema availible if you include the top comment, which allows for autocomplete in VSCode.

> [!TIP]
>**To use lua configuration, pass the `-lua` flag to `kserver`**
>
>*There is no way to set the port in lua so you can use `-port INT` to configure it or use the default `8000`*

For LUA create a `kserver.lua` file in the directory you want to serve,
the file should look like this:

```lua copy
local routes = {
    -- It's important to keep the order of these parameteres otherwise the server will confuse them
    {route = "/", content = "./index.html", contentType = "text/html"},
    -- Add more routes here
}

Register_routes(routes)
```

If you need custom logic before registering you can add it before the call to Register_routes(routes) or replace it with this:

```lua

for i, r in ipairs(routes) do
    local route = route.new(r.route, r.content, r.contentType, r.handler)
    if not route then
        error("Failed to create route")
    end
    local success, err = pcall(function()
        local result = registerRoutes(route)
        if not result then
            error("Failed to register route: " .. r.route)
        end
    end)
    if not success then
        error("Error registering route: " .. err)
    end
end
```

>[!NOTE]
> Register_routes and registerRoutes will show up as undefined as they are created at runtime when the server starts.

## Planned features

- Support for custom endpoint logic with the lua configuration :trollface:
