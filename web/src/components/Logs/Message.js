import React from "react";

const renderAction = message => (
  <tr>
    <th>Action</th>
    <td>{message.action}</td>
    <td style={{overflowWrap: 'break-word', wordBreak: 'break-word'}}>{message.item && message.item.url}</td>
    <td colSpan={3}>{message.content}</td>
  </tr>
)

const renderUpdate = message => (
  <tr>
    <th>Update</th>
    <td>{message.scraper}</td>
    <td style={{overflowWrap: 'break-word', wordBreak: 'break-word'}}>{message.item && message.item.url}</td>
    <td>{message.status}</td>
    <td>{message.message}</td>
  </tr>
)


const Message = ({message}) => message.type === 'action' ? renderAction(message) : renderUpdate(message)

export default Message;