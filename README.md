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
    {route = "/", content = "./index.html", contentType = "text/html"},
    -- Add more routes here
}

-- You don't have to make any changes here unless you need to process something before registering.
for i, r in ipairs(routes) do
    local route = route.new(r.route, r.content, r.contentType)
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
> This lua code supports multiple routes and can be greatly simplified if you only intend on serrving one route.

## Planned features

- Support for custom endpoint logic with the lua configuration :trollface:
