import { Button, Grid, TextField } from '@mui/material';
import { useState } from 'react';

type InputMessageProps = {
  onSend: (message: string) => void;
}

function InputMessage({ onSend }: InputMessageProps) {
  const [ val, setval ] = useState<string>('')

  const handleSend = () => {
    onSend(val)
    setval('')
  }

  return (
    <Grid container>
      <Grid item xs={10}>
        <TextField value={val} onChange={(e) => setval(e.target.value)} fullWidth />
      </Grid>
      <Grid item xs={2}>
        <Button onClick={() => handleSend()} variant="outlined" size="large" fullWidth>Send</Button>
      </Grid>
    </Grid>
  );
}

export default InputMessage
