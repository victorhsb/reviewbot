import Card from "@mui/material/Card";
import { Box } from "@mui/system";
import "./message.css"

type MessageProps = {
  author?: string;
  content: string;
}

function Message({ content, author }: MessageProps) {
  const auth = author && <b>{author}:</b>
  const align = author === 'Bot' ? 'align-left' : 'align-right'

  return <Box component={Card} className={`message-card ${align}`}>
    {auth}<br />{content}
  </Box>
}

export default Message
