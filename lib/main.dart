import 'package:flutter/material.dart';
import 'package:acer_red/pages/index.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: '红枫',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.white),
      ),
      home: const Index(),
      builder: (context, child) {
        if (Theme.of(context).platform == TargetPlatform.macOS) {
          return SizedBox(height: 730, child: child);
        }
        return child!;
      },
    );
  }
}
