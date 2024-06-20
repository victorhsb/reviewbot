import Container from "@mui/material/Container";
import Chat from "./components/chat"
import ProductList from "./components/products";
import { Grid } from "@mui/material";

function App() {
  return (
    <Container>
      <Grid container spacing={2}>
        <Grid item xs={12} md={8}>
          <Chat />
        </Grid>
        <Grid item xs={12} md={4}>
          <ProductList />
        </Grid>
      </Grid>
    </Container>
  )
}

export default App
