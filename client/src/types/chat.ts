export interface MessageWrapper {
  senderId: string;
  content: Message;
}
interface Message {
  message: string;
  name: string;
}
