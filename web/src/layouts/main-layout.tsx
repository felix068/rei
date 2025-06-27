import { Outlet } from "react-router";

export default function MainLayout() {
	return (
		<div className="flex flex-col justify-center items-center">
			<div className="flex gap-x-4">
				<h1>Rei</h1>
				<button className="font-medium">
					Unread <span className="font-bold">(1)</span>
				</button>
				<button>Feeds</button>
				<button>Settings</button>
			</div>

			<div>
				<Outlet />
			</div>
		</div>
	);
}
