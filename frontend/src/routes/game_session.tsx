import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/game_session')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/game_session"!</div>
}
