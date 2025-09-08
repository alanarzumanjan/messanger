import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../core/config/theme.dart';
import '../features/chat/data/chat_repository.dart';
import '../features/chat/presentation/pages/chats_page.dart';
import '../features/chat/presentation/pages/chats_page.dart';
import '../features/chat/presentation/state/chat_provider.dart';

class AlanogramApp extends StatelessWidget {
  const AlanogramApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(
          create: (_) => ChatProvider(DemoChatRepository())..load(),
        ),
      ],
      child: MaterialApp(
        debugShowCheckedModeBanner: false,
        title: 'Alanogram',
        theme: lightTheme,
        darkTheme: darkTheme,
        themeMode: ThemeMode.system,
        initialRoute: '/',
        onGenerateRoute: (settings) {
          if (settings.name == '/') {
            return MaterialPageRoute(builder: (_) => const ChatsPage());
          }
          if (settings.name == '/chat') {
            final chatId = settings.arguments as String;
            return MaterialPageRoute(builder: (_) => ChatPage(chatId: chatId));
          }
          return null;
        },
      ),
    );
  }
}
