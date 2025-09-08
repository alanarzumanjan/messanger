import 'dart:async';
import '../domain/chat.dart';
import '../domain/message.dart';

abstract class ChatRepository {
  Future<List<Chat>> fetchChats();
  Future<List<Message>> fetchMessages(String chatId);
  Future<Message> sendMessage(String chatId, String text);
}

/// Демонстрационная реализация в памяти.
/// Позже заменишь на работу с REST/WebSocket.
class DemoChatRepository implements ChatRepository {
  final List<Chat> _chats = [
    Chat(
      id: 'c1',
      title: 'General',
      lastMessage: 'Welcome to Alanogram 👋',
      updatedAt: DateTime.now().subtract(const Duration(minutes: 5)),
    ),
    Chat(
      id: 'c2',
      title: 'Design',
      lastMessage: 'Figma sucks? 😅',
      updatedAt: DateTime.now().subtract(const Duration(hours: 1)),
    ),
  ];

  final Map<String, List<Message>> _messages = {
    'c1': [
      Message(
        id: 'm1',
        chatId: 'c1',
        senderId: 'u2',
        text: 'Welcome to Alanogram 👋',
        createdAt: DateTime.now().subtract(const Duration(minutes: 5)),
        isMine: false,
      ),
    ],
    'c2': [
      Message(
        id: 'm2',
        chatId: 'c2',
        senderId: 'u1',
        text: 'Figma sucks? 😅',
        createdAt: DateTime.now().subtract(const Duration(hours: 1)),
        isMine: true,
      ),
    ],
  };

  @override
  Future<List<Chat>> fetchChats() async {
    await Future<void>.delayed(const Duration(milliseconds: 150));
    _chats.sort((a, b) => b.updatedAt.compareTo(a.updatedAt));
    return List<Chat>.from(_chats);
  }

  @override
  Future<List<Message>> fetchMessages(String chatId) async {
    await Future<void>.delayed(const Duration(milliseconds: 100));
    return List<Message>.from(_messages[chatId] ?? const []);
  }

  @override
  Future<Message> sendMessage(String chatId, String text) async {
    final msg = Message(
      id: DateTime.now().microsecondsSinceEpoch.toString(),
      chatId: chatId,
      senderId: 'me',
      text: text.trim(),
      createdAt: DateTime.now(),
      isMine: true,
    );
    _messages.putIfAbsent(chatId, () => []);
    _messages[chatId]!.add(msg);

    final chat = _chats.firstWhere((c) => c.id == chatId);
    chat.lastMessage = msg.text;
    chat.updatedAt = msg.createdAt;

    // Демо-эхо
    unawaited(
      Future<void>.delayed(const Duration(milliseconds: 700), () {
        final echo = Message(
          id: DateTime.now().microsecondsSinceEpoch.toString(),
          chatId: chatId,
          senderId: 'bot',
          text: 'echo: ${msg.text}',
          createdAt: DateTime.now(),
          isMine: false,
        );
        _messages[chatId]!.add(echo);
        chat.lastMessage = echo.text;
        chat.updatedAt = echo.createdAt;
      }),
    );

    return msg;
  }
}
