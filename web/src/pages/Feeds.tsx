import { useEffect, useState } from "react";

export default function Feeds() {
	const [count, setCount] = useState(0);
	const [data, setData] = useState({});

	async function fetchData() {
		await fetch("http://localhost:1323/list_feeds")
			.then((a) => a.json())
			.then((data) => setData(data))
			.catch((e) => console.error(e));
	}

	useEffect(() => {
		fetchData();
	}, []);

	return (
		<div className="flex flex-col justify-center items-center">
			<button onClick={fetchData}>Refresh</button>

			<div className="mt-4">{data && <div>{data.name}</div>}</div>
		</div>
	);
}
