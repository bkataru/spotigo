package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/bkataru/spotigo/internal/config"
	"github.com/bkataru/spotigo/internal/ollama"
	"github.com/bkataru/spotigo/internal/tools"
)

const (
	// MaxInputLength is the maximum allowed length for user input
	MaxInputLength = 4096
	// MaxConversationHistory limits the number of messages kept in context
	MaxConversationHistory = 50
)

var (
	chatModel    string
	chatContext  int
	enableTools  bool
	musicDataDir string
)

func init() {
	chatCmd.Flags().StringVar(&chatModel, "model", "", "override the default chat model")
	chatCmd.Flags().IntVar(&chatContext, "context", 4096, "context window size")
	chatCmd.Flags().BoolVar(&enableTools, "tools", true, "enable tool calling for music queries")
	chatCmd.Flags().StringVar(&musicDataDir, "data-dir", "./data", "directory containing music data files")
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start an AI chat session about your music",
	Long: `Start an interactive AI chat session to discuss your music library.

You can ask questions like:
  - "What are my most played genres?"
  - "Find songs similar to [track name]"
  - "When did I start listening to [artist]?"
  - "What's my musical taste evolution over time?"
  - "Recommend something based on my recent listening"

The AI runs entirely locally using Ollama. No data leaves your machine.

Exit the chat with 'exit', 'quit', Ctrl+C, or Ctrl+D.`,
	Run: func(cmd *cobra.Command, args []string) {
		runChat()
	},
}

func runChat() {
	fmt.Println("Spotigo AI Chat")
	fmt.Println("===============")
	fmt.Println("Ask me anything about your music library!")
	fmt.Println("Type 'exit', 'quit', or press Ctrl+C / Ctrl+D to end the session.")
	fmt.Println()

	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded. Cannot start chat.")
		return
	}

	// Initialize Ollama client
	ollamaClient := ollama.NewClient(cfg.Ollama.Host, time.Duration(cfg.Ollama.Timeout)*time.Second)

	// Test Ollama connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := ollamaClient.Ping(ctx); err != nil {
		fmt.Printf("❌ Cannot connect to Ollama: %v\n", err)
		fmt.Println("Make sure Ollama is running:")
		fmt.Println("  ollama serve")
		fmt.Println()
		fmt.Println("Or check your Ollama configuration:")
		fmt.Printf("  Host: %s\n", cfg.Ollama.Host)
		return
	}

	// Load model configuration
	modelCfg, err := config.LoadModelConfig("./config")
	if err != nil {
		fmt.Printf("Warning: Could not load model config: %v\n", err)
		fmt.Println("Using default model: granite4:1b")
	}

	// Determine which model to use
	modelName := chatModel
	if modelName == "" {
		if modelCfg != nil {
			var err error
			modelName, err = modelCfg.GetModelForRole("chat")
			if err != nil {
				fmt.Printf("Warning: failed to get chat model: %v\n", err)
			}
		}
		if modelName == "" {
			modelName = "granite4:1b"
		}
	}

	fmt.Printf("✅ Connected to Ollama using model: %s\n", modelName)
	fmt.Println()

	// Load system prompt
	var systemPrompt string
	if modelCfg != nil {
		if agentCfg, exists := modelCfg.Agents["chat_agent"]; exists {
			systemPrompt = agentCfg.SystemPrompt
		}
	}
	if systemPrompt == "" {
		systemPrompt = `You are Spotigo, a friendly music intelligence assistant. You help users explore their Spotify library, discover insights about their listening habits, and find new music. Be conversational, knowledgeable about music, and helpful. Always be concise but informative.`

		if enableTools {
			systemPrompt += `

You have access to tools that can query the user's music library. Use these tools to answer questions about their music:
- get_library_stats: Get overall statistics
- search_tracks: Search for tracks by text
- get_tracks_by_artist: Get tracks by a specific artist
- get_recently_added_tracks: Get recently added tracks
- get_all_artists: Get all unique artists
- get_playlist_by_name: Find a playlist
- query_music_data: Execute custom queries with filters, sorting, aggregation

When you need specific information from the library, use the appropriate tool. After getting results, summarize them in a natural, conversational way.`
		}
	}

	// Initialize music tools if enabled
	var musicTools *tools.MusicTools
	var toolDefs []ollama.Tool
	if enableTools {
		musicTools = tools.NewMusicTools(musicDataDir)
		toolDefs = musicTools.GetToolDefinitions()
		fmt.Println("🔧 Tool calling enabled - I can query your music library!")
		fmt.Println()
	}

	// Set up signal handling for graceful exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Channel to signal chat loop to exit
	exitChan := make(chan struct{})

	// Handle signals in background
	go func() {
		<-sigChan
		fmt.Println("\n\nReceived interrupt signal. Goodbye!")
		close(exitChan)
	}()

	// Start chat loop
	reader := bufio.NewReader(os.Stdin)

	// Initialize conversation with system message
	messages := []ollama.Message{
		{Role: "system", Content: systemPrompt},
	}

	for {
		// Check if we should exit before prompting
		select {
		case <-exitChan:
			return
		default:
		}

		fmt.Print("You: ")
		input, err := reader.ReadString('\n')

		// Handle EOF (Ctrl+D)
		if errors.Is(err, io.EOF) {
			fmt.Println("\n\nEnd of input (Ctrl+D). Goodbye!")
			return
		}

		// Handle other read errors
		if err != nil {
			// Check if it's due to signal interrupt
			select {
			case <-exitChan:
				return
			default:
				fmt.Printf("\nError reading input: %v\n", err)
				return
			}
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Handle exit commands (case-insensitive)
		inputLower := strings.ToLower(input)
		if inputLower == "exit" || inputLower == "quit" || inputLower == "bye" || inputLower == "q" {
			fmt.Println("Goodbye!")
			return
		}

		// Handle help command
		if inputLower == "help" || inputLower == "?" {
			printChatHelp()
			continue
		}

		// Handle clear command
		if inputLower == "clear" || inputLower == "reset" {
			messages = []ollama.Message{
				{Role: "system", Content: systemPrompt},
			}
			fmt.Println("Conversation cleared. Starting fresh!")
			fmt.Println()
			continue
		}

		// Validate input
		if valErr := validateChatInput(input); valErr != nil {
			fmt.Printf("Invalid input: %v\n", valErr)
			continue
		}

		// Trim conversation history if too long
		if len(messages) > MaxConversationHistory {
			// Keep system message and last N-1 messages
			messages = append(messages[:1], messages[len(messages)-MaxConversationHistory+1:]...)
		}

		// Add user message to conversation
		messages = append(messages, ollama.Message{
			Role:    "user",
			Content: input,
		})

		fmt.Print("Spotigo: ")

		// Chat loop with tool calling support
		maxToolCalls := 5
		for iteration := 0; iteration < maxToolCalls; iteration++ {
			// Create cancellable context for the chat request
			chatCtx, chatCancel := context.WithTimeout(context.Background(), 60*time.Second)

			// Run chat request with signal awareness
			respChan := make(chan *ollama.ChatResponse, 1)
			errChan := make(chan error, 1)

			go func() {
				req := ollama.ChatRequest{
					Model:    modelName,
					Messages: messages,
					Options: &ollama.Options{
						Temperature: 0.7,
						NumPredict:  chatContext,
					},
				}

				// Add tools if enabled
				if enableTools && toolDefs != nil {
					req.Tools = toolDefs
				}

				resp, err := ollamaClient.Chat(chatCtx, req)
				if err != nil {
					errChan <- err
					return
				}
				respChan <- resp
			}()

			// Wait for response or exit signal
			select {
			case <-exitChan:
				chatCancel()
				fmt.Println("\n\nChat interrupted. Goodbye!")
				return

			case err := <-errChan:
				chatCancel()
				fmt.Printf("❌ Error: %v\n", err)
				fmt.Println("Retrying with fallback model...")

				// Try fallback model
				fallbackModel := "qwen3:0.6b"
				if modelCfg != nil {
					fallbackModel, _ = modelCfg.GetFallbackForRole("chat") //nolint:errcheck // Fallback model is optional
				}

				fallbackCtx, fallbackCancel := context.WithTimeout(context.Background(), 60*time.Second)
				req := ollama.ChatRequest{
					Model:    fallbackModel,
					Messages: messages,
					Options: &ollama.Options{
						Temperature: 0.7,
						NumPredict:  chatContext,
					},
				}
				resp, fallbackErr := ollamaClient.Chat(fallbackCtx, req)
				fallbackCancel()

				if fallbackErr != nil {
					fmt.Printf("❌ Fallback also failed: %v\n", fallbackErr)
					// Remove user message from history if chat failed
					messages = messages[:len(messages)-1]
					break
				}

				// Add assistant response to conversation
				messages = append(messages, resp.Message)
				fmt.Println(resp.Message.Content)
				fmt.Println()
				break

			case resp := <-respChan:
				chatCancel()

				// Check if response contains tool calls
				if enableTools && len(resp.Message.ToolCalls) > 0 {
					// Add assistant message with tool calls to history
					messages = append(messages, resp.Message)

					// Execute tool calls
					for _, toolCall := range resp.Message.ToolCalls {
						fmt.Printf("🔧 Calling tool: %s\n", toolCall.Function.Name)

						// Show arguments for debugging
						var args map[string]interface{}
						if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err == nil {
							argsJSON, _ := json.MarshalIndent(args, "", "  ")
							fmt.Printf("   Arguments: %s\n", string(argsJSON))
						}

						result, err := musicTools.ExecuteToolCall(toolCall)
						if err != nil {
							result = fmt.Sprintf("Error executing tool: %v", err)
						}

						// Add tool result to conversation
						messages = append(messages, ollama.Message{
							Role:    "tool",
							Content: result,
						})
					}

					// Continue to next iteration to get final response
					continue
				}

				// No tool calls, this is the final response
				messages = append(messages, resp.Message)
				fmt.Println(resp.Message.Content)
				fmt.Println()
				break
			}

			// Break out of tool calling loop if we got here
			break
		}
	}
}

// printChatHelp prints available chat commands
func printChatHelp() {
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  help, ?     - Show this help message")
	fmt.Println("  clear, reset - Clear conversation history")
	fmt.Println("  exit, quit, q, bye - Exit the chat")
	fmt.Println()
	fmt.Println("Keyboard shortcuts:")
	fmt.Println("  Ctrl+C      - Exit immediately")
	fmt.Println("  Ctrl+D      - Exit (end of input)")
	fmt.Println()
	fmt.Println("Tips:")
	fmt.Println("  - Ask about your music library, listening habits, or recommendations")
	fmt.Println("  - The conversation history is preserved for context")
	fmt.Println("  - Use 'clear' to start a fresh conversation")
	fmt.Println()
}

// validateChatInput validates user input for the chat
func validateChatInput(input string) error {
	// Check length
	if len(input) > MaxInputLength {
		return fmt.Errorf("input too long (max %d characters, got %d)", MaxInputLength, len(input))
	}

	// Check for valid UTF-8
	if !utf8.ValidString(input) {
		return fmt.Errorf("input contains invalid UTF-8 characters")
	}

	// Check for null bytes or other control characters that could cause issues
	for i, r := range input {
		if r == 0 {
			return fmt.Errorf("input contains null byte at position %d", i)
		}
		// Allow common control characters (tab, newline) but reject others
		if r < 32 && r != '\t' && r != '\n' && r != '\r' {
			return fmt.Errorf("input contains invalid control character at position %d", i)
		}
	}

	return nil
}
