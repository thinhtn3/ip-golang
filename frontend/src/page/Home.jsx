import { supabase } from '../config/supabase'
import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'


export default function Home() {
    const navigate = useNavigate()
    const handleLogout = () => {
        supabase.auth.signOut()
        localStorage.clear()
    }

    useEffect(() => {
        const user = localStorage.getItem('user')
        if (!user) {
            navigate('/login')
        }
    }, [])
    return (
        <div>
            <h1>Home</h1>
            <button onClick={handleLogout}>Logout</button>
        </div>
    )
}