import { useState } from "react";
import viteLogo from "/vite.svg";
import reactLogo from "./assets/react.svg";

function App() {
	const [count, setCount] = useState(0);
	const [data, setData] = useState({});

	async function fetchData() {
		await fetch("http://localhost:1323/list_feeds")
			.then((a) => a.json())
			.then((data) => setData(data))
			.catch((e) => console.error(e));
	}

	return (
		<>
			<h1>Vite + React</h1>
			<div>
				<button onClick={fetchData}>fetch</button>

				{data && <div>{data.name}</div>}
			</div>
		</>
	);
}

export default App;
