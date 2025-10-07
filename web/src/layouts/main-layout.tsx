import { useEffect, useState } from "react";
import { Link, Outlet } from "react-router";
import { api } from "../utils/api";

export default function MainLayout() {
	const [unreadCount, setUnreadCount] = useState(0);
	const [dark, setDark] = useState(() => {
		const saved = localStorage.getItem("theme");
		return saved === "dark";
	});

	async function fetchUnreadCount() {
		try {
			const posts = await api.listUnreadPosts();
			setUnreadCount(posts?.length || 0);
		} catch (error) {
			console.error("Error fetching unread count:", error);
		}
	}

	useEffect(() => {
		fetchUnreadCount();
		const interval = setInterval(fetchUnreadCount, 10000);
		return () => clearInterval(interval);
	}, []);

	useEffect(() => {
		const html = document.documentElement;
		if (dark) {
			html.classList.add("dark");
			html.classList.remove("light");
			localStorage.setItem("theme", "dark");
		} else {
			html.classList.add("light");
			html.classList.remove("dark");
			localStorage.setItem("theme", "light");
		}
	}, [dark]);

	const bgColor = dark ? "#1a1a1a" : "#ffffff";
	const textColor = dark ? "#ffffff" : "#000000";
	const buttonBg = dark ? "#2d2d2d" : "#f5f5f5";
	const buttonHoverBg = dark ? "#3d3d3d" : "#e5e5e5";
	const buttonBorder = dark ? "#4d4d4d" : "#d1d5db";
	const logoSrc = dark
		? "/logo/white/REI-white-transparent-low.png"
		: "/logo/black/REI-black-transparent-low.png";

	return (
		<div
			data-theme={dark ? "dark" : "light"}
			style={{
				minHeight: "100vh",
				backgroundColor: `${bgColor} !important`,
				color: `${textColor} !important`,
			}}
		>
			<div className="flex flex-col items-center">
				<div className="flex flex-wrap gap-2 md:gap-x-4 p-4 md:p-6 items-center justify-center">
					<img src={logoSrc} height={72} alt="REI" style={{ height: "56px" }} className="md:h-[72px]" />
					<Link
						to="/"
						style={{
							backgroundColor: buttonBg,
							color: textColor,
							borderColor: buttonBorder,
						}}
						className="font-medium px-3 md:px-4 py-2 rounded border text-sm md:text-base"
						onMouseEnter={(e) => {
							e.currentTarget.style.backgroundColor =
								buttonHoverBg;
						}}
						onMouseLeave={(e) => {
							e.currentTarget.style.backgroundColor = buttonBg;
						}}
					>
						Unread <span className="font-bold">({unreadCount})</span>
					</Link>
					<Link
						to="/feeds"
						style={{
							backgroundColor: buttonBg,
							color: textColor,
							borderColor: buttonBorder,
						}}
						className="px-3 md:px-4 py-2 rounded border text-sm md:text-base"
						onMouseEnter={(e) => {
							e.currentTarget.style.backgroundColor =
								buttonHoverBg;
						}}
						onMouseLeave={(e) => {
							e.currentTarget.style.backgroundColor = buttonBg;
						}}
					>
						Feeds
					</Link>
					<button
						onClick={() => setDark(!dark)}
						style={{
							backgroundColor: buttonBg,
							color: textColor,
							borderColor: buttonBorder,
						}}
						className="px-3 md:px-4 py-2 rounded border text-sm md:text-base"
						onMouseEnter={(e) => {
							e.currentTarget.style.backgroundColor =
								buttonHoverBg;
						}}
						onMouseLeave={(e) => {
							e.currentTarget.style.backgroundColor = buttonBg;
						}}
					>
						{dark ? "☀️" : "🌙"}
					</button>
				</div>

				<div className="w-full">
					<Outlet />
				</div>
			</div>
		</div>
	);
}
