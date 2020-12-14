import React from "react";
import "./Logs.scss";
import Messages from "./Logs/Messages";

const Logs = ({logs}) => {
  const messages = logs.map(log => {
    try {
      const data = JSON.parse(log.data)
      return data
    } catch {
    }
    return null
  })
  return (
    <div className="Logs">
      <h2>Logs</h2>
      <Messages messages={messages}/>
    </div>
  );
}

export default Logs;