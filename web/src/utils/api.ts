// Use /api in production (nginx proxy) or localhost:1323 in development
const API_BASE =
	import.meta.env.MODE === "production" ? "/api" : "http://localhost:1323";

export interface Feed {
	id: string;
	name: string;
	description: string;
	link: string;
	rssLink: string;
	createdAt: string;
	updatedAt: string;
}

export interface Post {
	id: string;
	name: string;
	description: string;
	link: string;
	feedId: string;
	feedName: string;
	isRead: boolean;
	createdAt: string;
	updatedAt: string;
}

export const api = {
	async addFeed(rssLink: string): Promise<string> {
		const response = await fetch(`${API_BASE}/add_feed`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ link: rssLink }),
		});

		if (!response.ok) {
			const error = await response.text();
			throw new Error(error);
		}

		return response.json();
	},

	async listFeeds(): Promise<Feed[]> {
		const response = await fetch(`${API_BASE}/list_feeds`);
		if (!response.ok) throw new Error("Failed to fetch feeds");
		return response.json();
	},

	async listPosts(): Promise<Post[]> {
		const response = await fetch(`${API_BASE}/list_posts`);
		if (!response.ok) throw new Error("Failed to fetch posts");
		return response.json();
	},

	async listUnreadPosts(): Promise<Post[]> {
		const response = await fetch(`${API_BASE}/list_unread_posts`);
		if (!response.ok) throw new Error("Failed to fetch unread posts");
		return response.json();
	},

	async markPostAsRead(postId: string): Promise<void> {
		const response = await fetch(`${API_BASE}/posts/${postId}/read`, {
			method: "PUT",
		});
		if (!response.ok) throw new Error("Failed to mark post as read");
	},

	async deleteFeed(feedId: string): Promise<void> {
		const response = await fetch(`${API_BASE}/feeds/${feedId}`, {
			method: "DELETE",
		});
		if (!response.ok) throw new Error("Failed to delete feed");
	},
};
