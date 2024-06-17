# E-commerce Review Chatbot
This project implements a review chatbot, designed to streamline the collection and processing of customer reviews after product delivery. The chatbot interacts with customers through a conversation flow, gathering ratings and feedback to enhance the company's understanding of customer satisfaction and product performance.

## Features
**Initiation**: The chatbot initiates a conversation with the customer upon receiving a message through the queue confirming the product delivery.
**Conversation Flow**: Guides the customer through a structured conversation flow to collect ratings and feedback on the delivered product.
**Rating**: Collects a star rating (1-5) from the customer for the delivered product.
**Feedback**: Gathers additional feedback from the customer through open-ended questions or multiple-choice options.
**Storage**: Stores review data in a structured format, including the customer's ID, product ID, star rating, and feedback comments.
**Queue/Stream-Based Architecture**: Employs RabbitMQ for message queuing, ensuring reliable message delivery and facilitating scalable processing of reviews.

## Structure
The structure used on the backend of the project is based on a 2018 gophercon talk that can be found [here](https://www.youtube.com/watch?v=oL6JBUk6tj0&t=245s).
The code reference can be found [here](https://github.com/katzien/go-structure-examples/tree/master)
