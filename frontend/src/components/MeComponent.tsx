import { useQuery } from "@connectrpc/connect-query"
import { getMe } from "../gen/account/v1/account-AccountService_connectquery"

export function MeComponent() {
	const { isPending, error, data } = useQuery(getMe)

	if (isPending) {
		return (<h1>Pending Me ...</h1>)
	}

	if (error) {
		return (<h1>error Me ${error.message}</h1>)
	}
	return <h1>user id:{data.account?.id}</h1>
}


