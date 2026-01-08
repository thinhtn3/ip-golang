import { useParams } from "react-router-dom"
import Chatbox from "../components/Chat/Chatbox"

export default function Chat() {
    const { sessionId } = useParams()
    return (
        <div>
            <Chatbox sessionId={sessionId} />
        </div>
    )
}