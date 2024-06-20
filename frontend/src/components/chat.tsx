import { Paper, SxProps } from "@mui/material";
import { Message } from "../models";
import ChatMessage from "./message";
import Stack from "@mui/material/Stack";
import { useEffect, useState } from "react";
import api from "../api";
import InputMessage from "./input";

type ChatProps = {
  userId: string;
  sx?: SxProps;
}

const _api = api("http://localhost:8080")

function Chat({ userId, sx }: ChatProps) {
  const [messages, setMessages] = useState<Message[]>([])

  useEffect(() => {
    const intervalId = setInterval(() => {
      _api.loadMessages(userId).then(setMessages);
    }, 2000);

    // Clear the interval when the component is unmounted
    return () => clearInterval(intervalId);
  }, [userId])

  return (
    <Paper sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Stack spacing={1} sx={{ ...sx }}>
        {messages.map((m, i) =>
          <ChatMessage key={i} content={m} isReceiver={m.sender != userId} />
        )}
      </Stack>
      <InputMessage onSend={(msg: string) => _api.sendMessage(userId, msg)} />
    </Paper>
  )
}

export default Chat
