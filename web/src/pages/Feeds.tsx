import { useEffect, useState } from "react";
import { api, type Feed } from "../utils/api";

export default function Feeds() {
	const [feeds, setFeeds] = useState<Feed[]>([]);
	const [newFeedUrl, setNewFeedUrl] = useState("");
	const [error, setError] = useState("");
	const [dark, setDark] = useState(() => {
		const saved = localStorage.getItem("theme");
		return saved === "dark";
	});

	async function fetchFeeds() {
		try {
			const feedsList = await api.listFeeds();
			setFeeds(feedsList || []);
		} catch (error) {
			console.error("Error fetching feeds:", error);
			setFeeds([]);
		}
	}

	async function handleAddFeed(e: React.FormEvent) {
		e.preventDefault();
		if (!newFeedUrl.trim()) return;

		try {
			setError("");
			await api.addFeed(newFeedUrl);
			setNewFeedUrl("");
			await fetchFeeds();
		} catch (err: any) {
			const errorMsg =
				err.message || "Failed to add feed. Check the URL and try again.";
			setError(errorMsg);
		}
	}

	async function handleDeleteFeed(feedId: string) {
		try {
			await api.deleteFeed(feedId);
			await fetchFeeds();
		} catch (error) {
			console.error("Error deleting feed:", error);
			alert("Failed to delete feed");
		}
	}

	useEffect(() => {
		fetchFeeds();
	}, []);

	useEffect(() => {
		const checkTheme = () => {
			const saved = localStorage.getItem("theme");
			setDark(saved === "dark");
		};
		const interval = setInterval(checkTheme, 100);
		return () => clearInterval(interval);
	}, []);

	const cardBg = dark ? "#2d2d2d" : "#ffffff";
	const textColor = dark ? "#ffffff" : "#000000";
	const textSecondary = dark ? "#e5e5e5" : "#6b7280";
	const borderColor = dark ? "#4d4d4d" : "#e5e7eb";
	const inputBg = dark ? "#2d2d2d" : "#ffffff";

	return (
		<div className="flex flex-col items-center p-8">
			<h1
				className="text-2xl font-bold mb-6"
				style={{ color: textColor }}
			>
				Feeds
			</h1>

			<form onSubmit={handleAddFeed} className="mb-4 flex gap-2">
				<input
					type="url"
					value={newFeedUrl}
					onChange={(e) => setNewFeedUrl(e.target.value)}
					placeholder="RSS feed URL"
					className="px-4 py-2 border rounded w-96"
					style={{
						backgroundColor: inputBg,
						borderColor,
						color: textColor,
					}}
				/>
				<button
					type="submit"
					className="px-4 py-2 text-white rounded"
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
					Add
				</button>
			</form>

			{error && (
				<div className="mb-4 p-3 border rounded w-96"
					style={{
						backgroundColor: dark ? "rgb(127, 29, 29)" : "rgb(254, 226, 226)",
						borderColor: dark ? "rgb(220, 38, 38)" : "rgb(248, 113, 113)",
						color: dark ? "rgb(254, 202, 202)" : "rgb(153, 27, 27)",
					}}
				>
					{error}
				</div>
			)}

			<div className="w-full max-w-2xl space-y-4">
				{feeds.map((feed) => (
					<div
						key={feed.id}
						className="border rounded p-4 flex justify-between items-start"
						style={{
							backgroundColor: cardBg,
							borderColor,
						}}
					>
						<div>
							<h3
								className="font-bold"
								style={{ color: textColor }}
							>
								{feed.name}
							</h3>
							<p
								className="text-sm"
								style={{ color: textSecondary }}
							>
								{feed.description}
							</p>
							<a
								href={feed.link}
								target="_blank"
								rel="noopener noreferrer"
								className="text-xs"
								style={{
									color: "rgb(237, 108, 61)",
								}}
								onMouseEnter={(e) => {
									e.currentTarget.style.color =
										"rgb(245, 166, 128)";
								}}
								onMouseLeave={(e) => {
									e.currentTarget.style.color =
										"rgb(237, 108, 61)";
								}}
							>
								{feed.link}
							</a>
						</div>
						<button
							onClick={() => handleDeleteFeed(feed.id)}
							className="text-sm"
							style={{
								color: dark ? "rgb(248, 113, 113)" : "rgb(239, 68, 68)",
							}}
							onMouseEnter={(e) => {
								e.currentTarget.style.color = dark
									? "rgb(254, 202, 202)"
									: "rgb(220, 38, 38)";
							}}
							onMouseLeave={(e) => {
								e.currentTarget.style.color = dark
									? "rgb(248, 113, 113)"
									: "rgb(239, 68, 68)";
							}}
						>
							Delete
						</button>
					</div>
				))}
			</div>
		</div>
	);
}
