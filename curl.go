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
	boxWidth     = 56
)

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
			Edges []struct {
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

func fetchAvatarASCII(username string) ([]string, error) {
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

	var lines []string
	bounds := resized.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		var b strings.Builder
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, bv, a := resized.At(x, y).RGBA()
			if a == 0 {
				b.WriteByte(' ')
				continue
			}
			gray := 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(bv/257)
			idx := int(gray * float64(len(asciiChars)-1) / 255)
			if idx >= len(asciiChars) {
				idx = len(asciiChars) - 1
			}
			if idx < 0 {
				idx = 0
			}
			b.WriteByte(asciiChars[idx])
		}
		lines = append(lines, b.String())
	}
	return lines, nil
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

func boxTop() string {
	return "┌" + strings.Repeat("─", boxWidth-2) + "┐"
}

func boxDivider() string {
	return "├" + strings.Repeat("─", boxWidth-2) + "┤"
}

func boxBottom() string {
	return "└" + strings.Repeat("─", boxWidth-2) + "┘"
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
			lines []string
			err   error
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
			lines, err := fetchAvatarASCII(username)
			avCh <- avResult{lines, err}
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
		writeLine(fmtLine(centerText("github.com/"+username, boxWidth-4)))
		writeLine(boxDivider())

		if avRes.err != nil {
			writeLine(fmtLine(" avatar: " + avRes.err.Error()))
		} else {
			padding := (boxWidth - 4 - avatarWidth) / 2
			padStr := strings.Repeat(" ", padding)
			for _, line := range avRes.lines {
				writeLine(fmtLine(padStr + line))
			}
		}
		writeLine(boxDivider())

		writeLine(fmtLine(" pinned repos"))
		if prRes.err != nil {
			writeLine(fmtLine(" error fetching repos"))
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

		writeLine(fmtLine(" languages"))
		if lsRes.err != nil {
			writeLine(fmtLine(" error fetching languages"))
		} else if len(lsRes.stats) == 0 {
			writeLine(fmtLine(" no language data"))
		} else {
			for _, lang := range lsRes.stats {
				bar := langBar(lang.Pct, 10)
				s := fmt.Sprintf("  %-16s %s %5.1f%%", lang.Name, bar, lang.Pct)
				writeLine(fmtLine(s))
				time.Sleep(15 * time.Millisecond)
			}
		}
		writeLine(boxDivider())

		writeLine(fmtLine(centerText("generated by github-stats-backend", boxWidth-4)))
		writeLine(boxBottom())
	}
}
