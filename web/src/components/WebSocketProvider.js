import React, {createContext, useEffect} from "react";
import {useDispatch} from 'react-redux';
import {connect} from "../api/ws";
import {addMessage} from "../actions/wsSlice";

export const WebSocketContext = createContext(null)

let socket, ws;
const WebSocketProvider = ({children}) => {
  const dispatch = useDispatch();

  const handleMessage = (payload) => {
    try {
      const data = JSON.parse(payload.data)
      dispatch(addMessage(data));
    } catch {
    }
  }

  useEffect(() => {
    if (!socket || socket.readyState === WebSocket.CLOSED) {
      socket = connect(handleMessage);

      ws = {
        socket: socket
      }
    }
  }, [handleMessage]);

  return (
    <WebSocketContext.Provider value={ws}>
      {children}
    </WebSocketContext.Provider>);
}

export default WebSocketProvider;