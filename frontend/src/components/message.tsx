import "./message.css"
import moment from "moment"
import { Message, User } from "../models"

type MessageProps = {
  content: Message;
  user?: User;
}

function MessageBubble({content, user}: MessageProps) {
  const { message, direction, timestamp } = content

  const auth =  user && <b>{user.username}:</b>
  const align = direction == 'sender' ? 'align-left' : 'align-right'

  const date = timestamp && <><br /><small className="message-timestamp">{moment(timestamp).format('MM/DD/YYYY HH:mm')}</small></>

  return <div className={`message-card ${align}`}>
    {auth}<br /><p>{message}{date}</p>
  </div>
}

export default MessageBubble
