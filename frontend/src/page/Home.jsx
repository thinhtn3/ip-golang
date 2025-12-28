import { useAuth } from '../context/authProvider'
import { useNavigate } from 'react-router-dom'
import { supabase } from '../config/supabase'
import React, { useState, useEffect } from 'react';
import axios from 'axios';


export default function Home() {
    //if not authenticated, redirect to login
    const navigate = useNavigate()
    const [questions, setQuestions] = useState([])

    let accessToken;
    const getAccessToken = async () => {
        const { data, error } = await supabase.auth.getSession()
        if (error) {
            console.error('Error getting session:', error)
            navigate('/login')
        }
        accessToken = data.session?.access_token
        return accessToken
    }
    accessToken = getAccessToken()

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
        console.log(accessToken)
        const { data, error } = await axios.post('http://localhost:8080/chat/create', {
            questionId: questionId
        }, {
            headers: {
                Authorization: `Bearer ${accessToken}`,
            },
        })
        if (error) {
            console.error('Error sending question:', error)
        }
    }
    const logout = async () => {
        const { error } = await supabase.auth.signOut()
        if (error) {
            console.error('Error logging out:', error)
        }
        //create local storage item to store access token
        localStorage.removeItem('accessToken')
        //remove user from context
        setUser(null)
        //remove session from context
        setSession(null)
        //remove loading from context
        setLoading(false)
        //navigate to login page
        navigate('/login')
    }

    return (
        <div>
            <h1>Home</h1>
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