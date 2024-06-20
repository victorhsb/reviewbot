import { Message, Product, User } from "./models";

export interface API {
  loadMessages: (id: string) => Promise<Message[]>
  sendMessage: (id: string, message: string) => Promise<void>

  loadProducts: () => Promise<Product[]>
  loadUser: (id: string) => Promise<User>
  loadUsers: () => Promise<User[]>
}

export default (baseurl: string): API => ({
  loadMessages: async (id: string): Promise<Message[]> => {
    const response = await fetch(`${baseurl}/v1/users/${id}/messages`);

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  },
  sendMessage: async (id: string, message: string): Promise<void> => {
    const response = await fetch(`${baseurl}/v1/users/${id}/message`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ message })
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
  },
  loadUser: async (id: string): Promise<User> => {
    const response = await fetch(`${baseurl}/v1/users/${id}`);

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  },
  loadUsers: async (): Promise<User[]> => {
    const response = await fetch(`${baseurl}/v1/users`);

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  }
})
