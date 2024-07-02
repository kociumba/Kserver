package luaintegration

var Simplify = `
function Register_routes(routes)
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
end
`

// Testing for the custom endpoint logic
var SimplifyFUCK = `
function Register_routes(routes)
    if not routes then
        error("Routes parameter is nil")
    end

    for i, r in ipairs(routes) do
        if not r then
            error(string.format("Route at index %d is nil", i))
        end

        if not r.route then
            error(string.format("Route at index %d has no route", i))
        end

        if not r.content then
            error(string.format("Route at index %d has no content", i))
        end

        if not r.contentType then
            error(string.format("Route at index %d has no contentType", i))
        end

        if r.contentType == "application/json" and not r.handler then
            error(string.format("Route at index %d has contentType 'application/json' but no handler", i))
        elseif r.contentType ~= "application/json" and r.handler then
            error(string.format("Route at index %d has contentType '%s' but has handler", i, r.contentType))
        end

        local route
        if r.contentType == "application/json" then
            route = route.new(r.route, r.content, r.contentType, r.handler)
        else
            route = route.new(r.route, r.content, r.contentType, nil)
        end

        local success, err = pcall(function()
            local result = registerRoutes(route)
            if not result then
                error(string.format("Failed to register route '%s' at index %d", r.route, i))
            end
        end)

        if not success then
            print(string.format("Error registering route '%s' at index %d: %s", r.route, i, err))
        else
        --    print(string.format("Registered route '%s' at index %d", r.route, i))
        end
    end
end
`
