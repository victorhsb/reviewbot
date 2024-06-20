import { Message, Product } from "./models";

export interface API {
  loadMessages: (id: string) => Promise<Message[]>
  sendMessage: (id: string, message: string) => Promise<void>

  loadProducts: () => Promise<Product[]>
}

export default (baseurl: string): API => ({
  loadMessages: async (id: string): Promise<Message[]> => {
    const response = await fetch(`${baseurl}/v1/messages/${id}`);

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  },
  sendMessage: async (id: string, message: string): Promise<void> => {
    const response = await fetch(`${baseurl}/v1/messages`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ message, sender: id })
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
  },
  loadProducts: async (): Promise<Product[]> =>{
    const response = await fetch(`${baseurl}/v1/products`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }
})
