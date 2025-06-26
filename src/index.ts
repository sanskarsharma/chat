import { Container, getContainer } from '@cloudflare/containers';

export class ChatContainer extends Container {
  // Configure default port for the container
  defaultPort = 8080;
  sleepAfter = "5m";

  override onStart() {
    console.log('Chat container successfully started');
  }

  override onStop() {
    console.log('Chat container successfully shut down');
  }

  override onError(error: unknown) {
    console.log('Chat container error:', error);
  }
}

export default {
  async fetch(request, env) {
    console.log("fetching");
    return getContainer(env.CHAT_CONTAINER).fetch(request);
  },
};