import React from "react";
import "./Logs.scss";
import Messages from "./Logs/Messages";
import {useSelector} from "react-redux";

const Logs = () => {
  const {messages} = useSelector((state) => state.ws);

  return (
    <div className="Logs">
      <h2>Logs</h2>
      <Messages messages={messages}/>
    </div>
  );
}

export default Logs;