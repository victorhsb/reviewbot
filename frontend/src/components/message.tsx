import "./message.css"
import moment from "moment"
import { Message, User } from "../models"

type MessageProps = {
  content: Message;
  user: User;
}

function MessageBubble({content, user}: MessageProps) {
  const { message, direction, timestamp } = content

  const sent = direction == 'sent'

  const auth = <b>{sent ? user.username : "BOT"}:</b>
  const align = !sent ? 'align-left' : 'align-right'

  const date = timestamp && <><br /><small>{moment(timestamp).format('MM/DD/YYYY HH:mm')}</small></>

  return <div className={`message-card ${align}`}>
    {auth}<br /><p>{message}{date}</p>
  </div>
}

export default MessageBubble
