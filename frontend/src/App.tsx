import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import XoGame from "./XoGame";
import Chat from "./Chat";

const queryClient = new QueryClient()

export default function App() {
	return (
		<QueryClientProvider client={queryClient}>
			<Auth />
			<XoGame />
			<Chat />
		</QueryClientProvider>
	)
};

function Auth() {
	return <>
	</>
}

// function Example() {
// 	const { isPending, error, data } = useQuery({
// 		queryKey: ['repoData'],
// 		queryFn: async () => {
// 			const response = await fetch('http://api.github.com/repos/TanStack/query')
// 			if (!response.ok) {
// 				throw new Error('Network response was not ok')
// 			}
// 			return response.json()
// 		},
// 	})
//
//
// 	if (isPending) return 'Loading...'
//
// 	if (error) return <span>Error: {error.message}</span>
//
// 	return (
// 		<div>
// 			<h1>{data.name}</h1>
// 			<p>{data.description}</p>
// 			<strong>👀 {data.subscribers_count}</strong>{' '}
// 			<strong>✨ {data.stargazers_count}</strong>{' '}
// 			<strong>🍴 {data.forks_count}</strong>
// 		</div>
// 	)
//
// }
