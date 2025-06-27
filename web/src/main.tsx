import { createRoot } from "react-dom/client";
import { RouterProvider, createBrowserRouter } from "react-router";
import "./index.css";
import Home from "./Home.tsx";
import MainLayout from "./layouts/main-layout.tsx";
import Feeds from "./pages/Feeds.tsx";

const router = createBrowserRouter([
	// {
	// 	path: "/",
	// 	Component: App,
	// },
	{
		Component: MainLayout,
		children: [
			{ index: true, Component: Home },
			{ path: "feeds", Component: Feeds },
		],
	},
]);

const root = document.getElementById("root")!;

createRoot(root).render(<RouterProvider router={router} />);
