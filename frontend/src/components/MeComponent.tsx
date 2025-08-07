import { useQuery } from "@tanstack/react-query"
import { GetMe } from "../lib/auth"

export function MeComponent() {
	const { isPending, error, data } = useQuery({
		queryKey: ["getMe"],
		queryFn: GetMe
	})

	if (isPending) {
		return (<h1>Pending Me ...</h1>)
	}

	if (error) {
		return (<h1>error Me ${error.message}</h1>)
	}
	return <h1>user id:{data.id}</h1>
}


