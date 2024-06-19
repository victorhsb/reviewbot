import { SxProps } from "@mui/material";
import { Message } from "../models";
import ChatMessage from "./message";
import Stack from "@mui/material/Stack";

type ChatProps = {
  messages: Message[];
  sx?: SxProps;
}

function Chat({ messages, sx }: ChatProps) {
  return (
    <Stack spacing={1} sx={{ ...sx }}>
      {messages.map((m, i) =>
        <ChatMessage key={i} content={m.message} author={m.author} />
      )}
    </Stack>
  )
}

export default Chat
