package types

type Project struct {
    Title       string
    Description string
    Tags        []string
    ImageURL    string
    GitHubURL   string
    LiveURL     string
}

// GetProjects returns a list of featured projects
func GetProjects() []Project {
    return []Project{
        {
            Title:       "AI-Powered Portfolio",
            Description: "An interactive portfolio site with AI chat capabilities, built with Go and HTMX. Features real-time interactions and modern UI.",
            Tags:        []string{"Go", "HTMX", "TailwindCSS", "OpenAI"},
            GitHubURL:   "https://github.com/sub0xdai/portfolio-site",
            ImageURL:    "/static/images/portfolio.png",
        },
        {
            Title:       "Blockchain Explorer",
            Description: "A comprehensive blockchain explorer for analyzing transactions, smart contracts, and network statistics in real-time.",
            Tags:        []string{"React", "TypeScript", "Web3.js", "GraphQL"},
            GitHubURL:   "https://github.com/sub0xdai/blockchain-explorer",
            ImageURL:    "/static/images/blockchain.png",
        },
        {
            Title:       "DeFi Dashboard",
            Description: "A decentralized finance dashboard for tracking investments, yields, and portfolio performance across multiple chains.",
            Tags:        []string{"Next.js", "Solidity", "Ethers.js", "TailwindCSS"},
            GitHubURL:   "https://github.com/sub0xdai/defi-dashboard",
            LiveURL:     "https://defi.yourdomain.com",
            ImageURL:    "/static/images/defi.png",
        },
    }
}
