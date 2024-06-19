import Container from "@mui/material/Container";
import Chat from "./components/chat"
import { Paper } from "@mui/material";
import InputMessage from "./components/input";
import { useEffect, useState } from "react";
import api from "./api";
import { Message } from "./models";

const _api = api("http://localhost:8080")
const stubID = "4ecb6555-2fe5-42b2-8451-be3c24fbb1c8"

function App() {
  const [messages, setMessages] = useState<Message[]>([])

  useEffect(() => {
    const intervalId = setInterval(() => {
      _api.loadMessages(stubID).then(setMessages);
    }, 1000);

    // Clear the interval when the component is unmounted
    return () => clearInterval(intervalId);
  }, [])

  return (
    <Container>
      <Paper sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
        <Chat messages={messages} sx={{ flexGrow: 1 }} />
        <InputMessage onSend={(msg) => _api.sendMessage(stubID, msg)} />
      </Paper>
    </Container>
  )
}

export default App
