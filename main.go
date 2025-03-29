package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	cmds := flag.String("cmds", "", "comma separated list of shell commands to use")

	flag.Parse()

	if cmds == nil || *cmds == "" {
		env := os.Getenv("CMDS")
		if env != "" {
			*cmds = env
		} else {
			fmt.Println("usage: mpcsh --cmds ls,cat,head")
			os.Exit(42)
		}
	}

	list := strings.Split(*cmds, ",")

	// Create MCP server
	s := server.NewMCPServer(
		"MPCsh 🚀",
		"0.1.0",
	)

	for _, cmdName := range list {
		bb, err := exec.Command("man", cmdName).CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}

		tool := mcp.NewTool(cmdName,
			mcp.WithDescription(string(bb)),
			mcp.WithString("args",
				mcp.Required(),
				mcp.Description("args to execute command with"),
			),
		)

		s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args, ok := request.Params.Arguments["args"].(string)
			if !ok {
				return nil, errors.New("name must be a string")
			}

			bb, err := exec.Command(cmdName, args).CombinedOutput()
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(bb)), nil
		})
	}

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
