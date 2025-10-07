import { useEffect, useState } from "react";
import { api, type Post } from "./utils/api";

type SortType = "date" | "feed" | "alphabetical";

export default function Home() {
	const [posts, setPosts] = useState<Post[]>([]);
	const [sortedPosts, setSortedPosts] = useState<Post[]>([]);
	const [sortBy, setSortBy] = useState<SortType>("date");
	const [dark, setDark] = useState(() => {
		const saved = localStorage.getItem("theme");
		return saved === "dark";
	});

	async function fetchUnreadPosts() {
		try {
			const unreadPosts = await api.listUnreadPosts();
			setPosts(unreadPosts || []);
		} catch (error) {
			console.error("Error fetching unread posts:", error);
			setPosts([]);
		}
	}

	async function handleMarkAsRead(postId: string, postLink: string) {
		try {
			await api.markPostAsRead(postId);
			setPosts(posts.filter((p) => p.id !== postId));
			window.open(postLink, "_blank");
		} catch (error) {
			console.error("Error marking post as read:", error);
		}
	}

	useEffect(() => {
		fetchUnreadPosts();
		const interval = setInterval(fetchUnreadPosts, 10000);
		return () => clearInterval(interval);
	}, []);

	useEffect(() => {
		const checkTheme = () => {
			const saved = localStorage.getItem("theme");
			setDark(saved === "dark");
		};
		const interval = setInterval(checkTheme, 100);
		return () => clearInterval(interval);
	}, []);

	useEffect(() => {
		let sorted = [...posts];
		switch (sortBy) {
			case "feed":
				sorted.sort((a, b) => a.feedName.localeCompare(b.feedName));
				break;
			case "alphabetical":
				sorted.sort((a, b) => a.name.localeCompare(b.name));
				break;
			case "date":
			default:
				sorted.sort(
					(a, b) =>
						new Date(b.createdAt).getTime() -
						new Date(a.createdAt).getTime()
				);
				break;
		}
		setSortedPosts(sorted);
	}, [posts, sortBy]);

	const cardBg = dark ? "#2d2d2d" : "#ffffff";
	const textColor = dark ? "#ffffff" : "#000000";
	const textSecondary = dark ? "#e5e5e5" : "#6b7280";
	const borderColor = dark ? "#4d4d4d" : "#e5e7eb";

	return (
		<div className="flex flex-col items-center p-4 md:p-8">
			<h1
				className="text-2xl font-bold mb-6"
				style={{ color: textColor }}
			>
				Unread ({posts.length})
			</h1>

			<div className="w-full max-w-5xl mb-4 flex flex-col sm:flex-row gap-2 items-center">
				<span style={{ color: textColor }} className="font-medium">
					Sort by:
				</span>
				<select
					value={sortBy}
					onChange={(e) => setSortBy(e.target.value as SortType)}
					className="px-4 py-2 border rounded"
					style={{
						backgroundColor: cardBg,
						borderColor,
						color: textColor,
					}}
				>
					<option value="date">Date (newest first)</option>
					<option value="feed">Feed name</option>
					<option value="alphabetical">Title (A-Z)</option>
				</select>
			</div>

			<div className="w-full max-w-5xl space-y-4">
				{sortedPosts.length === 0 ? (
					<div
						className="text-center p-8 border rounded"
						style={{ borderColor }}
					>
						<p style={{ color: textSecondary }}>No unread posts</p>
					</div>
				) : (
					sortedPosts.map((post) => (
						<div
							key={post.id}
							className="border rounded p-4 md:p-6"
							style={{
								backgroundColor: cardBg,
								borderColor,
							}}
						>
							<div className="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-2 mb-2">
								<h3
									className="font-bold text-lg"
									style={{ color: textColor }}
								>
									{post.name}
								</h3>
								<span
									className="text-xs px-2 py-1 rounded"
									style={{
										backgroundColor: dark
											? "#3d3d3d"
											: "#f5f5f5",
										color: "rgb(237, 108, 61)",
										whiteSpace: "nowrap",
									}}
								>
									{post.feedName}
								</span>
							</div>
							<p
								className="text-sm mb-3"
								style={{ color: textSecondary }}
							>
								{post.description}
							</p>
							<div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-2">
								<span
									className="text-xs"
									style={{ color: textSecondary }}
								>
									{new Date(post.createdAt).toLocaleDateString(
										"en-US",
										{
											year: "numeric",
											month: "short",
											day: "numeric",
											hour: "2-digit",
											minute: "2-digit",
										}
									)}
								</span>
								<button
									onClick={() =>
										handleMarkAsRead(post.id, post.link)
									}
									className="px-4 py-2 text-white rounded text-sm w-full sm:w-auto"
									style={{
										backgroundColor: "rgb(237, 108, 61)",
									}}
									onMouseEnter={(e) => {
										e.currentTarget.style.backgroundColor =
											"rgb(245, 166, 128)";
									}}
									onMouseLeave={(e) => {
										e.currentTarget.style.backgroundColor =
											"rgb(237, 108, 61)";
									}}
								>
									Read
								</button>
							</div>
						</div>
					))
				)}
			</div>
		</div>
	);
}
