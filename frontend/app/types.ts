export interface LanguageEdge {
  size: number;
  node: {
    name: string;
    color: string;
  };
}

export interface RepositoryNode {
  name: string;
  stargazers: {
    totalCount: number;
  };
  forkCount: number;
  isPrivate: boolean;
  isFork: boolean;
  languages?: {
    edges: LanguageEdge[];
    totalSize: number;
  };
  description?: string;
  primaryLanguage?: {
    name: string;
    color: string;
  };
  url?: string;
}

export interface ContributionDay {
  contributionCount: number;
  date: string;
  weekday: number;
}

export interface Week {
  contributionDays: ContributionDay[];
}

export interface UserStats {
  user: {
    name: string;
    login: string;
    bio: string;
    avatarUrl: string;
    company: string;
    location: string;
    email: string;
    websiteUrl: string;
    twitterUsername: string;
    createdAt: string;
    pronouns: string;
    status?: {
      emoji: string;
      message: string;
    };
    followers: { totalCount: number };
    following: { totalCount: number };
    repositories: {
      totalCount: number;
      nodes: RepositoryNode[];
    };
    repositoriesContributedTo: {
      totalCount: number;
      nodes: {
        name: string;
        owner: { login: string };
        isPrivate: boolean;
      }[];
    };
    starredRepositories: { totalCount: number };
    organizations: {
      totalCount: number;
      nodes: {
        name: string;
        login: string;
        avatarUrl: string;
        description: string;
      }[];
    };
    contributionsCollection: {
      totalCommitContributions: number;
      totalPullRequestContributions: number;
      totalIssueContributions: number;
      totalRepositoryContributions: number;
      totalPullRequestReviewContributions: number;
      restrictedContributionsCount: number;
      contributionYears: number[];
      contributionCalendar: {
        totalContributions: number;
        weeks: Week[];
      };
    };
    gists: { totalCount: number };
    watching: { totalCount: number };
    sponsors: { totalCount: number };
    sponsoring: { totalCount: number };
    packages: { totalCount: number };
    pullRequests: { totalCount: number };
    issues: { totalCount: number };
  };
}

export interface PinnedRepos {
  user: {
    pinnedItems: {
      totalCount: number;
      edges: {
        node: RepositoryNode;
      }[];
    };
  };
}
