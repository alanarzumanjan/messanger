import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import '../state/chat_provider.dart';

class ChatsPage extends StatelessWidget {
  const ChatsPage({super.key});

  @override
  Widget build(BuildContext context) {
    final provider = context.watch<ChatProvider>();
    final chats = provider.chats;

    return Scaffold(
      appBar: AppBar(title: const Text('Alanogram')),
      body: provider.loading
          ? const Center(child: CircularProgressIndicator())
          : ListView.separated(
              itemCount: chats.length,
              separatorBuilder: (_, __) => const Divider(height: 1),
              itemBuilder: (context, i) {
                final c = chats[i];
                final time = DateFormat.Hm().format(c.updatedAt);
                return ListTile(
                  title: Text(c.title),
                  subtitle: Text(
                    c.lastMessage,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  trailing: Text(
                    time,
                    style: Theme.of(context).textTheme.labelSmall,
                  ),
                  onTap: () =>
                      Navigator.pushNamed(context, '/chat', arguments: c.id),
                );
              },
            ),
    );
  }
}
