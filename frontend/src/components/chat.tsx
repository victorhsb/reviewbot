import { Box, Button, Divider, Grid, Paper } from "@mui/material";
import { Message, User } from "../models";
import ChatMessage from "./message";
import Stack from "@mui/material/Stack";
import { useEffect, useState } from "react";
import api from "../api";
import InputMessage from "./input";

const _api = api("http://localhost:8080")

function UserList({ users, selectedId, onSelect }: { users: User[], selectedId?: string, onSelect: (u: User) => void }) {
  return <Stack spacing={1}>
    {users.map((u, i) => <Button key={i} disabled={selectedId == u.id} variant="contained" onClick={() => onSelect(u)} > {u.username} </Button>)}
  </Stack>
}

function Chat() {
  const [messages, setMessages] = useState<Message[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [selectedUser, setSelectedUser] = useState<User | undefined>()

  useEffect(() => {
    if (selectedUser) {
      const intervalId = setInterval(() => {
        _api.loadMessages(selectedUser.id)
          .then(setMessages)
      }, 2000);

      return () => clearInterval(intervalId);
    }
  }, [selectedUser])
  useEffect(() => { _api.loadUsers().then(setUsers) }, [])

  const selectUser = (u: User) => {
    setSelectedUser(u)
    setMessages([])
    _api.loadMessages(u.id).then(setMessages)
  }

  return (
    <Grid container spacing={1} sx={{ backgroundColor: "#cecece" }} paddingX={2}>
      <Grid item xs={2}>
        <UserList users={users} onSelect={u => selectUser(u)} selectedId={selectedUser && selectedUser.id!} />
      </Grid>
      <Grid item xs={10}>
        <Paper sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
          <h1 style={{ paddingLeft: "1em" }}>{selectedUser ? selectedUser.username : 'Select a user'}</h1>
          <Divider sx={{ marginBottom: 2 }} />
          <Stack paddingX={2} spacing={1} sx={{ flex: 1, overflowY: 'auto', maxHeight: 'calc(90vh)' }}>
            {messages.slice(-8).map((m, i) =>
                <ChatMessage key={i} content={m} user={selectedUser} />
            )}
          </Stack>
          <InputMessage disabled={!selectedUser} onSend={(msg: string) => selectedUser && _api.sendMessage(selectedUser.id, msg)} />
        </Paper>
      </Grid>
    </Grid >
  )
}

export default Chat
