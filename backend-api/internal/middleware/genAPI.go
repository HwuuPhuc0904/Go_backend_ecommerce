package middleware

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "os"
    "sort"
    "strings"
)

func GenerateAPIDocumentation(r *gin.Engine) {
    routes := r.Routes()
    
    // Sort routes by path for better organization
    sort.Slice(routes, func(i, j int) bool {
        return routes[i].Path < routes[j].Path
    })
    
    // Open file for writing
    f, err := os.Create("api_documentation.md")
    if err != nil {
        global.Logger.Error("Failed to create API documentation file", zap.Error(err))
        return
    }
    defer f.Close()
    
    // Write header
    f.WriteString("# API Documentation\n\n")
    f.WriteString("| Method | Path | Handler |\n")
    f.WriteString("|--------|------|--------|\n")
    
    // Group routes by their base path
    routeGroups := make(map[string][]gin.RouteInfo)
    for _, route := range routes {
        // Skip the /ping route
        if route.Path == "/ping" {
            continue
        }
        
        basePath := strings.Split(route.Path, "/")[1]
        if basePath == "" {
            basePath = "root"
        }
        routeGroups[basePath] = append(routeGroups[basePath], route)
    }
    
    // Write routes by group
    for group, routes := range routeGroups {
        f.WriteString("\n## " + strings.ToUpper(group) + " Routes\n\n")
        f.WriteString("| Method | Path | Handler |\n")
        f.WriteString("|--------|------|--------|\n")
        
        for _, route := range routes {
            f.WriteString("| " + route.Method + " | " + route.Path + " | " + route.Handler + " |\n")
        }
    }
    
    global.Logger.Info("API documentation generated successfully")
}
