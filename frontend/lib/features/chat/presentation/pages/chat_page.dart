import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../state/chat_provider.dart';
import '../widgets/message_bubble.dart';

class ChatPage extends StatefulWidget {
  final String chatId;
  const ChatPage({super.key, required this.chatId});

  @override
  State<ChatPage> createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final _controller = TextEditingController();
  final _scrollCtrl = ScrollController();

  @override
  void dispose() {
    _controller.dispose();
    _scrollCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final provider = context.watch<ChatProvider>();
    final chat = provider.chats.firstWhere((c) => c.id == widget.chatId);
    final messages = provider.messagesFor(widget.chatId);

    return Scaffold(
      appBar: AppBar(title: Text(chat.title)),
      body: SafeArea(
        child: Column(
          children: [
            Expanded(
              child: ListView.builder(
                controller: _scrollCtrl,
                reverse: true,
                padding: const EdgeInsets.symmetric(
                  horizontal: 12,
                  vertical: 8,
                ),
                itemCount: messages.length,
                itemBuilder: (context, index) {
                  final m = messages[messages.length - 1 - index];
                  return MessageBubble(message: m);
                },
              ),
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(12, 6, 12, 12),
              child: Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _controller,
                      textInputAction: TextInputAction.send,
                      onSubmitted: (_) => _send(provider),
                      decoration: const InputDecoration(
                        hintText: 'Message…',
                        border: OutlineInputBorder(),
                        isDense: true,
                      ),
                    ),
                  ),
                  const SizedBox(width: 8),
                  FilledButton(
                    onPressed: () => _send(provider),
                    child: const Icon(Icons.send),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _send(ChatProvider provider) {
    final text = _controller.text;
    _controller.clear();
    provider.send(widget.chatId, text);
    _scrollCtrl.animateTo(
      0,
      duration: const Duration(milliseconds: 250),
      curve: Curves.easeOut,
    );
  }
}
