import React, {useEffect, useState} from 'react';
import './App.css';
import Header from "./components/Header.js";
import Logs from "./components/Logs";
import {connect} from "./api/ws";

let socket;
const App = () => {
  const [messages, setMessages] = useState([]);
  const addMessage = (msg) => setMessages(previousMessages => [...previousMessages, msg]);

  useEffect(() => {
    if (!socket || socket.readyState === WebSocket.CLOSED) {
      socket = connect(addMessage);
    }
  }, [addMessage]);

  return (
    <div>
      <Header/>
      <Logs logs={messages}/>
    </div>
  );
}

export default App;
