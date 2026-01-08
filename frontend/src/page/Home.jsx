import { useNavigate } from 'react-router-dom'
import { supabase } from '../config/supabase'
import React, { useState, useEffect } from 'react';
import axios from 'axios';


export default function Home() {
    //if not authenticated, redirect to login
    const navigate = useNavigate()
    const [questions, setQuestions] = useState([])

    const getAccessToken = async () => {
        const { data } = await supabase.auth.getSession()
        return data.session?.access_token
    }

    useEffect(() => {
        const fetchQuestions = async () => {
            const { data, error } = await supabase.from('question_bank').select('*')
            if (error) {
                console.error('Error fetching questions:', error)
            } else {
                setQuestions(data)
            }
        }
        fetchQuestions()
    }, [])

    const sendQuestion = async (questionId) => {
        const accessToken = await getAccessToken()
        const { data, error } = await axios.post('http://localhost:8080/chat/create', {
            question_id: questionId
        }, {
            headers: {
                Authorization: `Bearer ${accessToken}`,
            },
        })
        if (error) {
            console.error('Error sending question:', error)
        }
        const sessionId = data.session.id
        navigate(`/chat/${sessionId}`)
    }

    const logout = async () => {
        const { error } = await supabase.auth.signOut()
        if (error) {
            console.error('Error logging out:', error)
        }
    }

    return (
        <div>
            {/* create logout button */}
            <button onClick={logout}>Logout</button>
            <div>
                {questions.map((question) => (
                    <div key={question.id}>{question.title}, {question.id} <button onClick={() => sendQuestion(question.id)}>Send</button></div>
                ))}
            </div>
        </div>
    )
}