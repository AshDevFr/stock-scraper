import React, {useEffect, useRef} from "react";
import styled from 'styled-components';
import {useSelector} from "react-redux";

const TerminalContainer = styled.div`
`;

const LogDiv = styled.div`
font-size: small;
font-family: monospace;
`;

const LogSpan = styled.span`
font-weight: bold;
`;

const ErrorSpan = styled(LogSpan)`
  color: red;
`;

const InfoSpan = styled(LogSpan)`
color: yellow;
`;

const DefaultLogSpan = styled(LogSpan)`
color: white;
`;

const renderLevel = (level) => {
  switch (level) {
    case "error": {
      return (<ErrorSpan>[Error] </ErrorSpan>)
    }
    case "warn": {
      return (<InfoSpan>[Warn] &nbsp;</InfoSpan>)
    }
    default: {
      return (<DefaultLogSpan>[Info] &nbsp;</DefaultLogSpan>)
    }
  }
}

const addLeadZero = (n) => {
  return n > 9 ? String(n) : `0${n}`
}

const renderTime = (time) => {
  const date = new Date(time * 1000);
  const year = date.getFullYear();
  const month = addLeadZero(date.getMonth());
  const day = addLeadZero(date.getDay());

  const hours = addLeadZero(date.getHours());
  const minutes = addLeadZero(date.getMinutes());
  const seconds = addLeadZero(date.getSeconds());

  return `[${year}-${month}-${day} ${hours}:${minutes}:${seconds}]`
}

const renderLog = (log) => {
  const level = log.status;
  const time = log.time;
  const itemUrl = log.item ? (log.item.trackedUrl || log.item.url) : '';
  const key = `${time}${level}${itemUrl}`

  return (
    <LogDiv key={key}>
      {renderLevel(level)}{renderTime(time)} {log.message} {itemUrl}
    </LogDiv>
  );
}

const Terminal = () => {
  const messagesEndRef = useRef(null)
  const {messages} = useSelector((state) => state.ws);
  const logs = messages.filter(m => m.type === 'update').slice(-50);

  const scrollToBottom = () => {
    messagesEndRef.current.scrollIntoView({behavior: "smooth"})
  }

  useEffect(scrollToBottom, [logs]);

  return (
    <TerminalContainer>
      {logs.map(renderLog)}
      <div ref={messagesEndRef}/>
    </TerminalContainer>
  );
}

export default Terminal;