package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

const (
	asciiChars   = "@%#*+=-:. "
	avatarWidth  = 40
	avatarHeight = 20
	boxWidth     = 60

	ansiReset    = "\033[0m"
	ansiCyan     = "\033[36m"
	ansiBoldCyan = "\033[1;36m"
	// ansiYellow     = "\033[33m"
	ansiBoldYellow = "\033[1;33m"
	ansiGray       = "\033[90m"
	ansiRed        = "\033[31m"
)

type avatarPixel struct {
	Ch      byte
	R, G, B uint8
}

type PinnedRepoData struct {
	Name        string
	URL         string
	Description string
	Stars       int
	Forks       int
	Language    string
	LangColor   string
}

type LanguageStatData struct {
	Name  string
	Color string
	Pct   float64
}

type pinnedReposResponse struct {
	User struct {
		PinnedItems struct {
			TotalCount int `json:"totalCount"`
			Edges      []struct {
				Node struct {
					Name        string `json:"name"`
					URL         string `json:"url"`
					Description string `json:"description"`
					Stargazers  struct {
						TotalCount int `json:"totalCount"`
					} `json:"stargazers"`
					ForkCount       int `json:"forkCount"`
					PrimaryLanguage *struct {
						Name  string `json:"name"`
						Color string `json:"color"`
					} `json:"primaryLanguage"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"pinnedItems"`
	} `json:"user"`
}

type langStatsResponse struct {
	User struct {
		Repositories struct {
			Nodes []struct {
				Languages *struct {
					TotalSize int `json:"totalSize"`
					Edges     []struct {
						Size int `json:"size"`
						Node struct {
							Name  string `json:"name"`
							Color string `json:"color"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"languages"`
			} `json:"nodes"`
		} `json:"repositories"`
	} `json:"user"`
}

func (c *GitHubClient) QueryInto(query string, variables map[string]interface{}, target interface{}) error {
	resp, err := c.Query(query, variables)
	if err != nil {
		return err
	}
	if len(resp.Errors) > 0 {
		return fmt.Errorf("graphql: %s", resp.Errors[0].Message)
	}
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func hexToANSI(hex string) string {
	if len(hex) < 7 || hex[0] != '#' {
		return ""
	}
	var r, g, b uint8
	if _, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b); err != nil {
		return ""
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

func fetchAvatarASCII(username string) ([][]avatarPixel, error) {
	resp, err := http.Get(fmt.Sprintf("https://github.com/%s.png", username))
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	resized := resize.Resize(avatarWidth, avatarHeight, img, resize.Bilinear)

	n := len(asciiChars) - 1
	var rows [][]avatarPixel
	bounds := resized.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		var row []avatarPixel
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := resized.At(x, y).RGBA()
			if a == 0 {
				row = append(row, avatarPixel{Ch: ' '})
				continue
			}
			r8, g8, b8 := uint8(r/257), uint8(g/257), uint8(b/257)
			gray := 0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8)
			idx := int(gray * float64(n) / 255)
			if idx > n {
				idx = n
			}
			if idx < 0 {
				idx = 0
			}
			row = append(row, avatarPixel{Ch: asciiChars[idx], R: r8, G: g8, B: b8})
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func getPinnedRepos(client *GitHubClient, username string) ([]PinnedRepoData, error) {
	query := `
		query PinnedRepos($username: String!) {
			user(login: $username) {
				pinnedItems(first: 6, types: [REPOSITORY]) {
					totalCount
					edges {
						node {
							... on Repository {
								name
								url
								description
								stargazers { totalCount }
								forkCount
								primaryLanguage {
									name
									color
								}
							}
						}
					}
				}
			}
		}
	`
	variables := map[string]interface{}{"username": username}

	var result pinnedReposResponse
	if err := client.QueryInto(query, variables, &result); err != nil {
		return nil, err
	}

	var repos []PinnedRepoData
	for _, edge := range result.User.PinnedItems.Edges {
		r := PinnedRepoData{
			Name:  edge.Node.Name,
			URL:   edge.Node.URL,
			Stars: edge.Node.Stargazers.TotalCount,
			Forks: edge.Node.ForkCount,
		}
		if edge.Node.PrimaryLanguage != nil {
			r.Language = edge.Node.PrimaryLanguage.Name
			r.LangColor = edge.Node.PrimaryLanguage.Color
		}
		repos = append(repos, r)
	}
	return repos, nil
}

func getLanguageStats(client *GitHubClient, username string) ([]LanguageStatData, error) {
	query := `
		query LangStats($username: String!) {
			user(login: $username) {
				repositories(first: 100, ownerAffiliations: OWNER, isFork: false) {
					nodes {
						languages(first: 10, orderBy: {field: SIZE, direction: DESC}) {
							edges {
								size
								node {
									name
									color
								}
							}
						}
					}
				}
			}
		}
	`
	variables := map[string]interface{}{"username": username}

	var result langStatsResponse
	if err := client.QueryInto(query, variables, &result); err != nil {
		return nil, err
	}

	type langEntry struct {
		color string
		size  int64
	}
	langMap := make(map[string]*langEntry)
	var totalSize int64

	for _, node := range result.User.Repositories.Nodes {
		if node.Languages == nil {
			continue
		}
		for _, edge := range node.Languages.Edges {
			name := edge.Node.Name
			if name == "Jupyter Notebook" {
				continue
			}
			if entry, ok := langMap[name]; ok {
				entry.size += int64(edge.Size)
			} else {
				langMap[name] = &langEntry{
					color: edge.Node.Color,
					size:  int64(edge.Size),
				}
			}
			totalSize += int64(edge.Size)
		}
	}

	if totalSize == 0 {
		return nil, nil
	}

	var stats []LanguageStatData
	for name, entry := range langMap {
		stats = append(stats, LanguageStatData{
			Name:  name,
			Color: entry.color,
			Pct:   float64(entry.size) / float64(totalSize) * 100,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Pct > stats[j].Pct
	})

	if len(stats) > 8 {
		stats = stats[:8]
	}

	return stats, nil
}

func fmtLine(content string) string {
	inner := boxWidth - 4
	runes := []rune(content)
	if len(runes) > inner {
		runes = runes[:inner]
	}
	return "│ " + string(runes) + strings.Repeat(" ", inner-len(runes)) + " │"
}

func coloredLine(content, code string) string {
	inner := boxWidth - 4
	runes := []rune(content)
	if len(runes) > inner {
		runes = runes[:inner]
	}
	padded := string(runes) + strings.Repeat(" ", inner-len(runes))
	return "│ " + code + padded + ansiReset + " │"
}

func boxTop() string {
	return ansiCyan + "┌" + strings.Repeat("─", boxWidth-2) + "┐" + ansiReset
}

func boxDivider() string {
	return ansiCyan + "├" + strings.Repeat("─", boxWidth-2) + "┤" + ansiReset
}

func boxBottom() string {
	return ansiCyan + "└" + strings.Repeat("─", boxWidth-2) + "┘" + ansiReset
}

func langBar(pct float64, width int) string {
	filled := int(pct * float64(width) / 100.0)
	if filled > width {
		filled = width
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

func centerText(s string, width int) string {
	runes := []rune(s)
	if len(runes) >= width {
		return string(runes[:width])
	}
	padding := width - len(runes)
	left := padding / 2
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", padding-left)
}

func renderAvatarRow(row []avatarPixel) string {
	var b strings.Builder
	lastR, lastG, lastB := uint8(0), uint8(0), uint8(0)
	for _, p := range row {
		if p.Ch == ' ' {
			b.WriteByte(' ')
			continue
		}
		if p.R != lastR || p.G != lastG || p.B != lastB {
			fmt.Fprintf(&b, "\033[38;2;%d;%d;%dm", p.R, p.G, p.B)
			lastR, lastG, lastB = p.R, p.G, p.B
		}
		b.WriteByte(p.Ch)
	}
	b.WriteString(ansiReset)
	return b.String()
}

func curlHandler(client *GitHubClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		w := c.Writer

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")

		flusher, ok := w.(http.Flusher)
		if !ok {
			c.String(http.StatusInternalServerError, "streaming not supported")
			return
		}

		type avResult struct {
			pixels [][]avatarPixel
			err    error
		}
		type prResult struct {
			repos []PinnedRepoData
			err   error
		}
		type lsResult struct {
			stats []LanguageStatData
			err   error
		}

		avCh := make(chan avResult, 1)
		prCh := make(chan prResult, 1)
		lsCh := make(chan lsResult, 1)

		go func() {
			pixels, err := fetchAvatarASCII(username)
			avCh <- avResult{pixels, err}
		}()
		go func() {
			repos, err := getPinnedRepos(client, username)
			prCh <- prResult{repos, err}
		}()
		go func() {
			stats, err := getLanguageStats(client, username)
			lsCh <- lsResult{stats, err}
		}()

		avRes := <-avCh
		prRes := <-prCh
		lsRes := <-lsCh

		writeLine := func(line string) {
			fmt.Fprintln(w, line)
			flusher.Flush()
			time.Sleep(20 * time.Millisecond)
		}

		writeLine(boxTop())
		writeLine(coloredLine(centerText("github.com/"+username, boxWidth-4), ansiBoldCyan))
		writeLine(boxDivider())

		if avRes.err != nil {
			writeLine(coloredLine(" avatar: "+avRes.err.Error(), ansiRed))
		} else {
			padding := (boxWidth - 4 - avatarWidth) / 2
			padStr := strings.Repeat(" ", padding)
			for _, row := range avRes.pixels {
				colored := renderAvatarRow(row)
				line := "│ " + padStr + colored + padStr + " │"
				fmt.Fprintln(w, line)
				flusher.Flush()
				time.Sleep(20 * time.Millisecond)
			}
		}
		writeLine(boxDivider())

		writeLine(coloredLine(" pinned repos", ansiBoldYellow))
		if prRes.err != nil {
			writeLine(coloredLine(" error fetching repos", ansiRed))
		} else if len(prRes.repos) == 0 {
			writeLine(fmtLine(" no pinned repos"))
		} else {
			for _, repo := range prRes.repos {
				s := "  " + repo.Name
				if repo.Language != "" {
					s += "  [" + repo.Language + "]"
				}
				s += fmt.Sprintf("  ★%d", repo.Stars)
				writeLine(fmtLine(s))
				time.Sleep(15 * time.Millisecond)
			}
		}
		writeLine(boxDivider())

		writeLine(coloredLine(" languages", ansiBoldYellow))
		if lsRes.err != nil {
			writeLine(coloredLine(" error fetching languages", ansiRed))
		} else if len(lsRes.stats) == 0 {
			writeLine(fmtLine(" no language data"))
		} else {
			for _, lang := range lsRes.stats {
				bar := langBar(lang.Pct, 10)
				s := fmt.Sprintf("  %-16s %s %5.1f%%", lang.Name, bar, lang.Pct)
				c := hexToANSI(lang.Color)
				writeLine(coloredLine(s, c))
				time.Sleep(15 * time.Millisecond)
			}
		}
		writeLine(boxDivider())

		writeLine(coloredLine(centerText("generated by github-stats-backend", boxWidth-4), ansiGray))
		writeLine(coloredLine(centerText("https://github.com/Mystery-Coder/github-stats-backend", boxWidth-4), ansiGray))

		writeLine(boxBottom())
	}
}
