import 'package:flutter/material.dart';
import 'package:acer_red/pages/index.dart';
import 'package:acer_red/env/ui.dart';
import 'package:acer_red/env/config.dart';

void main() {
  Settings().init();
  runApp(const MyApp());
}


class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: '红枫',
      navigatorKey: navigatorKey,
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.light(
          primary: const Color(0xFF191923),
          secondary: const Color(0xFFf2efea),
        ),
      ),
      home: const Home(),
      builder: (context, child) {
        if (Theme.of(context).platform == TargetPlatform.macOS) {
          return SizedBox(height: 730, child: child);
        }
        return child!;
      },
    );
  }
}
