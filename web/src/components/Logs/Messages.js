import React from "react";
import Message from "./Message";

const Messages = ({messages}) => {
  return (
    <table style={{width: '100%'}}>
      <thead>
      <tr>
        <th></th>
        <th></th>
        <th></th>
        <th colSpan={2}>Content</th>
      </tr>
      <tr>
        <th>Type</th>
        <th>Scraper / Action</th>
        <th>URL</th>
        <th>status</th>
        <th>Message</th>
      </tr>
      </thead>
      <tbody>
      {messages.map((message, index) => <Message key={index} message={message}/>)}
      </tbody>
    </table>
  );
}

export default Messages;