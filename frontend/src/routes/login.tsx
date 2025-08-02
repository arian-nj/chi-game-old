import { useQuery, type QueryFunctionContext } from '@tanstack/react-query'
import { createFileRoute, useRouter, useSearch } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { SetJwtToken } from '../lib/auth'

export const Route = createFileRoute('/login')({
	component: RouteComponent,
})

function RouteComponent() {
	const [userID, setUserId] = useState<number>(0)
	const [userInput, setUserInput] = useState("0")
	const router = useRouter()
	const search: any = useSearch({ strict: false })

	const { isPending, error, data } = useQuery({
		enabled: userID !== 0,
		queryKey: ["auth", userID],
		queryFn: ValidateDummy
	})

	useEffect(() => {
		if (data?.token) {
			SetJwtToken(data.token)
			router.history.push(search.redirect)
		}
	}, [data])

	if (isPending) {
		return <>
			<div>Finding out who are you ...</div>
			<input
				className='bg-gray-200 rounded-sm'
				placeholder='what is your id'
				onChange={(e) => setUserInput(e.target.value)}
				value={userInput}
			/>
			<button type='submit' onClick={(e) => {
				e.preventDefault()
				setUserId(Number(userInput))
				setUserInput("0")
			}
			}>Send</button>
		</>
	}

	if (error) {
		return <div>Error: {error.message}</div>
	}
}


type AuthResponse = { token: string };

async function ValidateDummy({ queryKey }: QueryFunctionContext<[string, number]>): Promise<AuthResponse> {
	const [, userId] = queryKey
	const response = await fetch("http://localhost:8383/api/auth/validate-dummy/", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ user_id: userId }),
	});
	if (!response.ok) throw new Error("Login failed");

	return response.json();
}


