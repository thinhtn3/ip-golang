import { useParams } from "react-router-dom"
import Chatbox from "../components/Chat/Chatbox"
import axios from "axios";
import { useState, useEffect } from "react";
import { supabase } from "../config/supabase";

export default function Chat() {
    const apiUrl = import.meta.env.VITE_DEVELOPMENT_URL;
    const [message, setMessage] = useState("");
    const [messages, setMessages] = useState([]);
    const [role, setRole] = useState("user");
    const [authorized, setAuthorized] = useState(false);
    const { sessionId } = useParams();
  
    const getAccessToken = async () => {
      const { data } = await supabase.auth.getSession()
      return data.session?.access_token
    }
  
    useEffect(() => {
      const fetchMessages = async () => {
        const accessToken = await getAccessToken()
        const response = await axios.get(`${apiUrl}/chat/sessions/${sessionId}/messages`, {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          }
        })
        console.log(response)
        if (response.status === 200) {
          setMessages(response.data.messages)
          setAuthorized(true)
        } else if (response.status === 403) {
        //Blur page and show error message
          console.error('Forbidden: User does not own session')
        } else {
          console.error('Error fetching messages:', response.data.message)
        }
      }
      
      fetchMessages()
    }, [])
  
    const handleSubmit = async () => {
      const accessToken = await getAccessToken()
      const response = await axios.post(`${apiUrl}/chat/sessions/${sessionId}/messages`, {
        message: message,
        role: role,
      }, {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        }
      })
      if (response.status === 200) {
        setMessages([...messages, response.data.chat])
      }
      else if (response.status === 403) {
        console.error('Forbidden: User does not own session')
      }
      else {
        console.error('Error sending message:', response.data.message)
      }
    }

    return (
        <div>
            {authorized ? <Chatbox sessionId={sessionId} messages={messages} setMessages={setMessages} handleSubmit={handleSubmit} /> : <div className="flex flex-col items-center justify-center h-screen">
                <h1 className="text-2xl font-bold">Forbidden: User does not own session</h1>
                <button className="mt-4 bg-blue-500 text-white px-4 py-2 rounded-md" onClick={() => navigate("/")}>Go to Home</button>
            </div>}
        </div>
    )
}