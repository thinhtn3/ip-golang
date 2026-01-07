import { useParams } from "react-router-dom"

export default function Chat() {
    const { sessionId } = useParams()
    return (
        <div>
            <h1>Chat {sessionId}</h1>
        </div>
    )
}