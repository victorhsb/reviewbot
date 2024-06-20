import Container from "@mui/material/Container";
import Chat from "./components/chat"
import ProductList from "./components/products";
import { Grid } from "@mui/material";

const stubId = "0f1f2e9d-66c4-442b-b020-93a9b1e863e6"

function App() {
  return (
    <Container>
      <Grid container spacing={2}>
        <Grid item xs={12} md={8}>
          <Chat userId={stubId} sx={{ flexGrow: 1 }} />
        </Grid>
        <Grid item xs={12} md={4}>
          <ProductList />
        </Grid>
      </Grid>
    </Container>
  )
}

export default App
