import { Button, Grid, TextField } from '@mui/material';
import { useState } from 'react';

type InputMessageProps = {
  onSend: (message: string) => void;
  disabled?: boolean;
}

function InputMessage({ onSend, disabled=false }: InputMessageProps) {
  const [ val, setval ] = useState<string>('')

  const handleSend = () => {
    onSend(val)
    setval('')
  }

  return (
    <Grid container>
      <Grid item xs={10}>
        <TextField value={val} disabled={disabled} onChange={(e) => setval(e.target.value)} fullWidth />
      </Grid>
      <Grid item xs={2}>
        <Button disabled={disabled} onClick={() => handleSend()} variant="outlined" size="large" fullWidth>Send</Button>
      </Grid>
    </Grid>
  );
}

export default InputMessage
