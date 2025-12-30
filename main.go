package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// GitHub API client
type GitHubClient struct {
	token   string
	baseURL string
}

func (c *GitHubClient) Query(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var gqlResp GraphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return nil, err
	}

	return &gqlResp, nil
}

func NewGitHubClient(token string) *GitHubClient {
	return &GitHubClient{
		token:   token,
		baseURL: "https://api.github.com/graphql",
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("GITHUB_TOKEN environment variable is required")
		os.Exit(1)
	}

	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	ghClient := NewGitHubClient(token)

	// Get pinned repositories
	r.GET("/api/pinned/:username", func(c *gin.Context) {
		username := c.Param("username")

		query := `
			query PinnedRepos($username: String!) {
				user(login: $username) {
					pinnedItems(first: 6, types: [REPOSITORY]) {
						totalCount
						edges {
							node {
								... on Repository {
									name
									id
									url
									description
									stargazers {
										totalCount
									}
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

		variables := map[string]interface{}{
			"username": username,
		}

		resp, err := ghClient.Query(query, variables)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(resp.Errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": resp.Errors})
			return
		}

		c.JSON(http.StatusOK, resp.Data)
	})

	// Get user stats
	r.GET("/api/stats/:username", func(c *gin.Context) {
		username := c.Param("username")

		query := `
			query UserStats($username: String!) {
				user(login: $username) {
					name
					login
					bio
					avatarUrl
					company
					location
					email
					websiteUrl
					twitterUsername
					createdAt
					pronouns
					followers {
						totalCount
					}
					following {
						totalCount
					}
					repositories(first: 100, ownerAffiliations: OWNER, privacy: null) {
						totalCount
						nodes {
							name
							stargazers {
								totalCount
							}
							forkCount
							isPrivate
							isFork
							languages(first: 10, orderBy: {field: SIZE, direction: DESC}) {
								edges {
									size
									node {
										name
										color
									}
								}
								totalSize
							}
						}
					}
					repositoriesContributedTo(first: 100, contributionTypes: [COMMIT, ISSUE, PULL_REQUEST, REPOSITORY]) {
						totalCount
						nodes {
							name
							owner {
								login
							}
							isPrivate
						}
					}
					starredRepositories {
						totalCount
					}
					contributionsCollection {
						totalCommitContributions
						totalPullRequestContributions
						totalIssueContributions
						totalRepositoryContributions
						totalPullRequestReviewContributions
						restrictedContributionsCount
						contributionYears
					}
					pullRequests {
						totalCount
					}
					issues {
						totalCount
					}
				}
			}
		`

		variables := map[string]interface{}{
			"username": username,
		}

		resp, err := ghClient.Query(query, variables)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(resp.Errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": resp.Errors})
			return
		}

		c.JSON(http.StatusOK, resp.Data)
	})

	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	r.Run(":" + port)
}
