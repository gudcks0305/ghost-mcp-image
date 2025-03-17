package main

import (
	"context"
	"fmt"
	"ghost-images/tools"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Ghost Images",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	s.AddTool(mcp.Tool{
		Name:        "upload_image_base64",
		Description: "Upload images to Ghost",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"base64_image": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		base64Image := request.Params.Arguments["base64_image"].(string)
		result, err := tools.UploadImageBase64Image(base64Image)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%v", result)), nil
	})

	s.AddTool(mcp.Tool{
		Name:        "upload_image_local_path",
		Description: "Upload images to Ghost",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"local_path": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		localPath := request.Params.Arguments["local_path"].(string)
		result, err := tools.UploadImageLocalPath(localPath)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%v", result)), nil
	})
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
