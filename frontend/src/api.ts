import { Message } from "./models";

export interface API {
  loadMessages: (id: string) => Promise<Message[]>
  sendMessage: (id: string, message: string) => Promise<void>
}

export const stubapi = (): API => {
  const messages = [
    { message: "Hello, how are you?", author: "John Doe" },
    { message: "I'm good, thanks! How about you?", author: "Bot" },
    { message: "I'm doing well, thank you.", author: "John Doe" },
    { message: "Great to hear!", author: "Bot" },
    { message: "Have a good day!", author: "John Doe" },
    { message: "You too!", author: "Bot" }
  ]
  return {
    loadMessages: async (): Promise<Message[]> => {
      return messages;
    },
    sendMessage: async (_: string, message: string): Promise<void> => {
      messages.push({ message, author: "John Doe" });
      return
    }
  }
}

export default (baseurl: string): API => ({
  loadMessages: async (id: string): Promise<Message[]> => {
    const response = await fetch(`${baseurl}/v1/messages/${id}`);

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    console.log(data)
    return data;
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
})
