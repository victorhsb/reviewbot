import Card from "@mui/material/Card";
import { Box } from "@mui/system";
import "./message.css"
import { Message } from "../models"

type MessageProps = {
  content: Message;
  isReceiver?: boolean;
}

function MessageBubble({content, isReceiver}: MessageProps) {
  const { senderName, sender, message } = content

  const auth = senderName && <b>{senderName}:</b>
  const align = isReceiver ? 'align-left' : 'align-right'

  return <Box component={Card} className={`message-card ${align}`}>
    {auth}<br />{message}
  </Box>
}

export default MessageBubble
