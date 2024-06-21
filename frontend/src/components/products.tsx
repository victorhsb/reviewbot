import { Button, Card, CardContent, CardHeader, Stack } from "@mui/material";
import { Product } from "../models";
import api from "../api";
import { useEffect, useState } from "react";

const _api = api("http://localhost:8080")

type ProductProps = {
  product: Product
};

function Review({ review }: { review: { rating: number, username: string, review: string, timestamp: Date } }) {
  return <Card>
    <CardHeader title={`${ review.username } - ${"★".repeat(review.rating)}`}/>
    <CardContent>
      {review.review}
    </CardContent>
  </Card>
}

function ProductCard({ product }: ProductProps) {
  return <Card variant='outlined'>
    <CardHeader title={product.title} />
    <CardContent>
      {product.id && <Button onClick={() => navigator.clipboard.writeText(product.id)}>{product.id}</Button> }
      <Stack spacing={1}>
        {product.reviews.map((review, i) => <Review key={i} review={review} />)}
      </Stack>
    </CardContent>
  </Card>
}

function ProductList() {
  const [products, setProducts] = useState<Product[]>([])

  useEffect(() => {
    const intervalId = setInterval(() => {
      _api.loadProducts().then(setProducts);
    }, 2000);

    // Clear the interval when the component is unmounted
    return () => clearInterval(intervalId);
  }, [])

  return (
    <Stack spacing={1}>
      {products.map((product, index) => <ProductCard key={index} product={product} />)}
    </Stack>
  )
}

export default ProductList;
