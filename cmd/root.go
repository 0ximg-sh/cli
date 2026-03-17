package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"0ximg.sh/cli/internal/api"
	"0ximg.sh/cli/internal/models"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

const buyMeACoffeeURL = "https://buymeacoffee.com/levinhne"

var (
	fromClipboard    bool
	outputPath       string
	shareOutput      bool
	language         string
	title            string
	theme            string
	background       string
	backgroundImage  string
	codePadRight     int
	font             string
	highlightLines   string
	lineOffset       int
	linePad          int
	noLineNumber     bool
	noRoundCorner    bool
	noWindowControls bool
	padHoriz         int
	padVert          int
	shadowBlurRadius int
	shadowColor      string
	shadowOffsetX    int
	shadowOffsetY    int
	tabWidth         int
	windowTitle      string
)

var rootCmd = &cobra.Command{
	Use:          "0ximg [file]",
	Short:        "0ximg CLI",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	Example:      "  0ximg main.go --theme Dracula --output main.go.png\n  cat main.go | 0ximg --language go --theme Dracula --output main.go.png\n  0ximg --from-clipboard --language go --share",
	RunE: func(cmd *cobra.Command, args []string) error {
		if shouldShowHelp(args) {
			return cmd.Help()
		}

		code, sourcePath, err := readCode(args)
		if err != nil {
			return err
		}

		req := models.RenderRequest{
			Code:             code,
			Language:         language,
			Title:            title,
			Theme:            theme,
			Background:       background,
			BackgroundImage:  backgroundImage,
			Font:             font,
			HighlightLines:   highlightLines,
			NoLineNumber:     noLineNumber,
			NoRoundCorner:    noRoundCorner,
			NoWindowControls: noWindowControls,
			ShadowColor:      shadowColor,
			WindowTitle:      windowTitle,
		}

		if req.Language == "" && sourcePath != "" {
			req.Language = strings.TrimPrefix(filepath.Ext(sourcePath), ".")
		}
		if req.WindowTitle == "" && sourcePath != "" {
			req.WindowTitle = filepath.Base(sourcePath)
		}

		req.HighlightLines = normalizeHighlightLines(req.HighlightLines)

		if err := validateHighlightLines(req.Code, req.HighlightLines); err != nil {
			return err
		}

		bindOptionalIntFlag(cmd, &req.CodePadRight, "code-pad-right", codePadRight)
		bindOptionalIntFlag(cmd, &req.LineOffset, "line-offset", lineOffset)
		bindOptionalIntFlag(cmd, &req.LinePad, "line-pad", linePad)
		bindOptionalIntFlag(cmd, &req.PadHoriz, "pad-horiz", padHoriz)
		bindOptionalIntFlag(cmd, &req.PadVert, "pad-vert", padVert)
		bindOptionalIntFlag(cmd, &req.ShadowBlurRadius, "shadow-blur-radius", shadowBlurRadius)
		bindOptionalIntFlag(cmd, &req.ShadowOffsetX, "shadow-offset-x", shadowOffsetX)
		bindOptionalIntFlag(cmd, &req.ShadowOffsetY, "shadow-offset-y", shadowOffsetY)
		bindOptionalIntFlag(cmd, &req.TabWidth, "tab-width", tabWidth)

		imageURL, previewURL, err := api.RenderImage(req)
		if err != nil {
			return err
		}

		if shareOutput {
			if err := clipboard.WriteAll(previewURL); err != nil {
				return fmt.Errorf("copy preview link to clipboard: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "✨ Preview link copied to clipboard: %s\n", previewURL)
			if previewURL != imageURL {
				fmt.Fprintf(cmd.OutOrStdout(), "🖼 Image URL: %s\n", imageURL)
			}
			return nil
		}

		resolvedOutputPath := resolveOutputPath(sourcePath)
		if err := downloadImage(imageURL, resolvedOutputPath); err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), "🎨 Image saved successfully!")
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetHelpTemplate(rootCmd.HelpTemplate() + "\nSupport: Buy me a coffee at " + buyMeACoffeeURL + "\n")
	rootCmd.Version = versionString()
	rootCmd.AddCommand(versionCmd)

	rootCmd.Flags().BoolVar(&fromClipboard, "from-clipboard", false, "Read code from clipboard")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output image path")
	rootCmd.Flags().BoolVarP(&shareOutput, "share", "s", false, "Copy the shareable link to clipboard")
	rootCmd.Flags().StringVarP(&language, "language", "l", "", "Language for syntax highlighting")
	rootCmd.Flags().StringVar(&title, "title", "", "Title metadata for the render")
	rootCmd.Flags().StringVarP(&theme, "theme", "t", "", "Theme used by the renderer")
	rootCmd.Flags().StringVarP(&background, "background", "b", "", "Background color")
	rootCmd.Flags().StringVar(&backgroundImage, "background-image", "", "Background image path or URL")
	rootCmd.Flags().IntVar(&codePadRight, "code-pad-right", 0, "Right code padding")
	rootCmd.Flags().StringVarP(&font, "font", "f", "", "Font family")
	rootCmd.Flags().StringVar(&highlightLines, "highlight-lines", "", "Line ranges to highlight, e.g. 1;5-12")
	rootCmd.Flags().IntVar(&lineOffset, "line-offset", 0, "Starting line number")
	rootCmd.Flags().IntVar(&linePad, "line-pad", 0, "Padding between lines")
	rootCmd.Flags().BoolVar(&noLineNumber, "no-line-number", false, "Hide line numbers")
	rootCmd.Flags().BoolVar(&noRoundCorner, "no-round-corner", false, "Disable rounded corners")
	rootCmd.Flags().BoolVar(&noWindowControls, "no-window-controls", false, "Hide window controls")
	rootCmd.Flags().IntVar(&padHoriz, "pad-horiz", 0, "Horizontal padding")
	rootCmd.Flags().IntVar(&padVert, "pad-vert", 0, "Vertical padding")
	rootCmd.Flags().IntVar(&shadowBlurRadius, "shadow-blur-radius", 0, "Shadow blur radius")
	rootCmd.Flags().StringVar(&shadowColor, "shadow-color", "", "Shadow color")
	rootCmd.Flags().IntVar(&shadowOffsetX, "shadow-offset-x", 0, "Shadow offset X")
	rootCmd.Flags().IntVar(&shadowOffsetY, "shadow-offset-y", 0, "Shadow offset Y")
	rootCmd.Flags().IntVar(&tabWidth, "tab-width", 0, "Tab width")
	rootCmd.Flags().StringVar(&windowTitle, "window-title", "", "Window title")
}

func readCode(args []string) (string, string, error) {
	if fromClipboard {
		code, err := clipboard.ReadAll()
		if err != nil {
			return "", "", fmt.Errorf("read code from clipboard: %w", err)
		}
		if strings.TrimSpace(code) == "" {
			return "", "", fmt.Errorf("clipboard is empty")
		}
		return code, "", nil
	}

	if len(args) == 0 {
		stdinInfo, err := os.Stdin.Stat()
		if err != nil {
			return "", "", fmt.Errorf("inspect stdin: %w", err)
		}

		if stdinInfo.Mode()&os.ModeCharDevice == 0 {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return "", "", fmt.Errorf("read code from stdin: %w", err)
			}
			if strings.TrimSpace(string(data)) == "" {
				return "", "", fmt.Errorf("stdin is empty")
			}
			return string(data), "", nil
		}

		return "", "", fmt.Errorf("missing input file: provide a file path, pipe code via stdin, or use --from-clipboard")
	}

	sourcePath := args[0]
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return "", "", fmt.Errorf("read source file %q: %w", sourcePath, err)
	}

	return string(data), sourcePath, nil
}

func shouldShowHelp(args []string) bool {
	if fromClipboard || len(args) > 0 {
		return false
	}

	stdinInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return stdinInfo.Mode()&os.ModeCharDevice != 0
}

func bindOptionalIntFlag(cmd *cobra.Command, target **int, flagName string, value int) {
	if cmd.Flags().Changed(flagName) {
		valueCopy := value
		*target = &valueCopy
	}
}

func validateHighlightLines(code string, raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	lineCount := strings.Count(code, "\n") + 1
	for _, segment := range strings.Split(raw, ";") {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			return fmt.Errorf("invalid --highlight-lines %q: empty segment", raw)
		}

		if strings.Contains(segment, "-") {
			bounds := strings.Split(segment, "-")
			if len(bounds) != 2 {
				return fmt.Errorf("invalid --highlight-lines segment %q: expected start-end", segment)
			}

			start, err := parsePositiveLineNumber(bounds[0], segment)
			if err != nil {
				return err
			}
			end, err := parsePositiveLineNumber(bounds[1], segment)
			if err != nil {
				return err
			}

			if start > end {
				return fmt.Errorf("invalid --highlight-lines segment %q: start must be <= end", segment)
			}
			if end > lineCount {
				return fmt.Errorf("invalid --highlight-lines segment %q: file only has %d lines", segment, lineCount)
			}

			continue
		}

		line, err := parsePositiveLineNumber(segment, segment)
		if err != nil {
			return err
		}
		if line > lineCount {
			return fmt.Errorf("invalid --highlight-lines segment %q: file only has %d lines", segment, lineCount)
		}
	}

	return nil
}

func normalizeHighlightLines(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}

	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ';' || r == ','
	})
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	return strings.Join(parts, ";")
}

func parsePositiveLineNumber(value string, segment string) (int, error) {
	value = strings.TrimSpace(value)
	line, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid --highlight-lines segment %q: %q is not a number", segment, value)
	}
	if line <= 0 {
		return 0, fmt.Errorf("invalid --highlight-lines segment %q: line numbers must be greater than 0", segment)
	}

	return line, nil
}

func resolveOutputPath(sourcePath string) string {
	if strings.TrimSpace(outputPath) != "" {
		return outputPath
	}

	filename := "snippet.png"
	if sourcePath != "" {
		filename = filepath.Base(sourcePath) + ".png"
	}

	return filepath.Join(".", filename)
}

func downloadImage(url string, destination string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download rendered image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		payload, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("download rendered image failed with status %s and unreadable body: %w", resp.Status, readErr)
		}
		message := strings.TrimSpace(string(payload))
		if message == "" {
			message = "empty response body"
		}
		return fmt.Errorf("download rendered image failed with status %s: %s", resp.Status, message)
	}

	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return fmt.Errorf("create output directory for %q: %w", destination, err)
	}

	file, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("create output file %q: %w", destination, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("write image to %q: %w", destination, err)
	}

	return nil
}
